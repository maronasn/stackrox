package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/detection/alertmanager"
	"github.com/stackrox/rox/central/detection/deploytime"
	"github.com/stackrox/rox/central/detection/lifecycle/metrics"
	"github.com/stackrox/rox/central/detection/runtime"
	"github.com/stackrox/rox/central/enrichment"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	processIndicatorDatastore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/central/processwhitelist"
	whitelistDataStore "github.com/stackrox/rox/central/processwhitelist/datastore"
	"github.com/stackrox/rox/central/reprocessor"
	riskManager "github.com/stackrox/rox/central/risk/manager"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	deployTimePkg "github.com/stackrox/rox/pkg/detection/deploytime"
	"github.com/stackrox/rox/pkg/enforcers"
	"github.com/stackrox/rox/pkg/expiringcache"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/images/enricher"
	"github.com/stackrox/rox/pkg/policies"
	"github.com/stackrox/rox/pkg/process/filter"
	processWhitelistPkg "github.com/stackrox/rox/pkg/processwhitelist"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
	"golang.org/x/time/rate"
)

var (
	lifecycleMgrCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Alert, resources.Deployment, resources.Image, resources.Indicator, resources.Policy, resources.ProcessWhitelist)))
)

type processWhitelistKey struct {
	deploymentID  string
	containerName string
	clusterID     string
	namespace     string
}

type indicatorWithInjector struct {
	indicator           *storage.ProcessIndicator
	msgToSensorInjector common.MessageInjector
}

type managerImpl struct {
	reprocessor        reprocessor.Loop
	enricher           enrichment.Enricher
	riskManager        riskManager.Manager
	runtimeDetector    runtime.Detector
	deploytimeDetector deploytime.Detector
	alertManager       alertmanager.AlertManager

	deploymentDataStore     deploymentDatastore.DataStore
	processesDataStore      processIndicatorDatastore.DataStore
	whitelists              whitelistDataStore.DataStore
	imageDataStore          imageDatastore.DataStore
	deletedDeploymentsCache expiringcache.Cache
	processFilter           filter.Filter

	queuedIndicators map[string]indicatorWithInjector

	indicatorQueueLock   sync.Mutex
	flushProcessingLock  concurrency.TransparentMutex
	indicatorRateLimiter *rate.Limiter
	indicatorFlushTicker *time.Ticker

	deploymentsPendingEnrichment *deploymentsPendingEnrichment

	policyAlertsLock *concurrency.KeyedMutex
}

func (m *managerImpl) copyAndResetIndicatorQueue() map[string]indicatorWithInjector {
	m.indicatorQueueLock.Lock()
	defer m.indicatorQueueLock.Unlock()
	if len(m.queuedIndicators) == 0 {
		return nil
	}
	copiedMap := m.queuedIndicators
	m.queuedIndicators = make(map[string]indicatorWithInjector)

	return copiedMap
}

func (m *managerImpl) buildIndicatorFilter() {
	ctx := sac.WithAllAccess(context.Background())
	var processesToRemove []string

	deploymentIDs, err := m.deploymentDataStore.GetDeploymentIDs()
	if err != nil {
		utils.Should(errors.Wrap(err, "error getting deployment IDs"))
		return
	}

	deploymentIDSet := set.NewStringSet(deploymentIDs...)

	err = m.processesDataStore.WalkAll(ctx, func(pi *storage.ProcessIndicator) error {
		if !deploymentIDSet.Contains(pi.GetDeploymentId()) {
			// Don't remove as these processes will be removed by GC
			// but don't add to the filter
			return nil
		}
		if !m.processFilter.Add(pi) {
			processesToRemove = append(processesToRemove, pi.GetId())
		}
		return nil
	})
	if err != nil {
		utils.Should(errors.Wrap(err, "error building indicator filter"))
	}

	log.Infof("Cleaning up %d processes as a part of building process filter", len(processesToRemove))
	if err := m.processesDataStore.RemoveProcessIndicators(ctx, processesToRemove); err != nil {
		utils.Should(errors.Wrap(err, "error removing process indicators"))
	}
	log.Infof("Successfully cleaned up those %d processes", len(processesToRemove))
}

func (m *managerImpl) flushQueuePeriodically() {
	defer m.indicatorFlushTicker.Stop()
	for range m.indicatorFlushTicker.C {
		m.flushIndicatorQueue()
	}
}

func (m *managerImpl) generateProcessAlertsAndEnforcement(indicators map[string]indicatorWithInjector, whitelistMap map[processWhitelistKey][]*storage.ProcessIndicator) {
	deploymentIDs := uniqueDeploymentIDs(indicators)

	newAlerts, err := m.runtimeDetector.AlertsForDeployments(deploymentIDs...)
	if err != nil {
		log.Errorf("Failed to compute runtime alerts: %s", err)
		return
	}

	deploymentToMatchingIndicators := make(map[string][]*storage.ProcessIndicator)
	for key, indicators := range whitelistMap {
		userWhitelist, err := m.checkAndUpdateWhitelist(key, indicators)
		if err != nil {
			log.Errorf("error checking and updating whitelist for %+v: %v", key, err)
			continue
		}
		if userWhitelist {
			deploymentToMatchingIndicators[key.deploymentID] = indicators
		}
	}
	if len(deploymentToMatchingIndicators) > 0 {
		// Compute whitelist alerts here
		whitelistAlerts, err := m.getWhitelistAlerts(deploymentToMatchingIndicators)
		if err != nil {
			log.Errorf("failed to get whitelist alerts: %v", err)
			return
		}
		newAlerts = append(newAlerts, whitelistAlerts...)
	}

	modifiedDeployments, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, newAlerts, alertmanager.WithLifecycleStage(storage.LifecycleStage_RUNTIME), alertmanager.WithDeploymentIDs(deploymentIDs...))
	if err != nil {
		log.Errorf("Couldn't alert and notify: %s", err)
	} else if modifiedDeployments.Cardinality() > 0 {
		defer m.reprocessor.ReprocessRiskForDeployments(modifiedDeployments.AsSlice()...)
	}

	// Create enforcement actions for the new alerts and send them with the stored injectors.
	m.generateAndSendEnforcements(newAlerts, indicators)
}

func indicatorToWhitelistKey(indicator *storage.ProcessIndicator) processWhitelistKey {
	return processWhitelistKey{
		deploymentID:  indicator.GetDeploymentId(),
		containerName: indicator.GetContainerName(),
		clusterID:     indicator.GetClusterId(),
		namespace:     indicator.GetNamespace(),
	}
}

func (m *managerImpl) flushIndicatorQueue() {
	// This is a potentially long-running operation, and we don't want to have a pile of goroutines queueing up on
	// this lock.
	if !m.flushProcessingLock.MaybeLock() {
		return
	}
	defer m.flushProcessingLock.Unlock()

	copiedQueue := m.copyAndResetIndicatorQueue()
	if len(copiedQueue) == 0 {
		return
	}

	// Map copiedQueue to slice
	indicatorSlice := make([]*storage.ProcessIndicator, 0, len(copiedQueue))
	for id, i := range copiedQueue {
		if deleted, _ := m.deletedDeploymentsCache.Get(i.indicator.GetDeploymentId()).(bool); deleted {
			delete(copiedQueue, id)
			continue
		}
		if !m.processFilter.Add(i.indicator) {
			delete(copiedQueue, id)
			metrics.ProcessFilterCounterInc("NotAdded")
			continue
		}
		metrics.ProcessFilterCounterInc("Added")
		indicatorSlice = append(indicatorSlice, i.indicator)
	}

	// Index the process indicators in batch
	if err := m.processesDataStore.AddProcessIndicators(lifecycleMgrCtx, indicatorSlice...); err != nil {
		log.Errorf("Error adding process indicators: %v", err)
	}

	// Group the processes into particular whitelist segments
	whitelistMap := make(map[processWhitelistKey][]*storage.ProcessIndicator)
	for _, indicator := range indicatorSlice {
		key := indicatorToWhitelistKey(indicator)
		whitelistMap[key] = append(whitelistMap[key], indicator)
	}

	if !features.SensorBasedDetection.Enabled() {
		m.generateProcessAlertsAndEnforcement(copiedQueue, whitelistMap)
	} else {
		for key, indicators := range whitelistMap {
			if _, err := m.checkAndUpdateWhitelist(key, indicators); err != nil {
				log.Errorf("error checking and updating whitelist for %+v: %v", key, err)
			}
		}
	}
}

func (m *managerImpl) addToQueue(indicator *storage.ProcessIndicator, injector common.MessageInjector) {
	m.indicatorQueueLock.Lock()
	defer m.indicatorQueueLock.Unlock()

	m.queuedIndicators[indicator.GetId()] = indicatorWithInjector{
		indicator:           indicator,
		msgToSensorInjector: injector,
	}
}

func (m *managerImpl) getWhitelistAlerts(deploymentsToIndicators map[string][]*storage.ProcessIndicator) ([]*storage.Alert, error) {
	whitelistExecutor := newWhitelistExecutor(lifecycleMgrCtx, m.deploymentDataStore, deploymentsToIndicators)
	if err := m.runtimeDetector.PolicySet().ForEach(whitelistExecutor); err != nil {
		return nil, err
	}
	return whitelistExecutor.alerts, nil
}

func (m *managerImpl) checkAndUpdateWhitelist(whitelistKey processWhitelistKey, indicators []*storage.ProcessIndicator) (bool, error) {
	key := &storage.ProcessWhitelistKey{
		DeploymentId:  whitelistKey.deploymentID,
		ContainerName: whitelistKey.containerName,
		ClusterId:     whitelistKey.clusterID,
		Namespace:     whitelistKey.namespace,
	}

	// TODO joseph what to do if whitelist doesn't exist?  Always create for now?
	whitelist, exists, err := m.whitelists.GetProcessWhitelist(lifecycleMgrCtx, key)
	if err != nil {
		return false, err
	}

	existingProcess := set.NewStringSet()
	for _, element := range whitelist.GetElements() {
		existingProcess.Add(element.GetElement().GetProcessName())
	}

	var elements []*storage.WhitelistItem
	var hasNonStartupProcess bool
	for _, indicator := range indicators {
		if !processwhitelist.IsStartupProcess(indicator) {
			hasNonStartupProcess = true
		}
		whitelistItem := processWhitelistPkg.WhitelistItemFromProcess(indicator)
		if !existingProcess.Add(whitelistItem) {
			continue
		}
		insertableElement := &storage.WhitelistItem{Item: &storage.WhitelistItem_ProcessName{ProcessName: whitelistItem}}
		elements = append(elements, insertableElement)
	}
	if len(elements) == 0 {
		return false, nil
	}
	if !exists {
		_, err = m.whitelists.UpsertProcessWhitelist(lifecycleMgrCtx, key, elements, true)
		return false, err
	}

	userWhitelist := processwhitelist.IsUserLocked(whitelist)
	roxWhitelist := processwhitelist.IsRoxLocked(whitelist) && hasNonStartupProcess
	if userWhitelist || roxWhitelist {
		// We already checked if it's in the whitelist and it is not, so reprocess risk to mark the results are suspicious if necessary
		m.reprocessor.ReprocessRiskForDeployments(whitelistKey.deploymentID)
		return userWhitelist, nil
	}
	_, err = m.whitelists.UpdateProcessWhitelistElements(lifecycleMgrCtx, key, elements, nil, true)
	return userWhitelist, err
}

func (m *managerImpl) IndicatorAdded(indicator *storage.ProcessIndicator, injector common.MessageInjector) error {
	if indicator.GetId() == "" {
		return fmt.Errorf("invalid indicator received: %s, id was empty", proto.MarshalTextString(indicator))
	}

	m.addToQueue(indicator, injector)

	if m.indicatorRateLimiter.Allow() {
		go m.flushIndicatorQueue()
	}
	return nil
}

func (m *managerImpl) generateRuntimeAlertsOnCreate(deployment *storage.Deployment, injector common.MessageInjector) error {
	indicatorQuery := search.NewQueryBuilder().AddExactMatches(search.DeploymentID, deployment.GetId()).ProtoQuery()
	indicators, err := m.processesDataStore.SearchRawProcessIndicators(lifecycleMgrCtx, indicatorQuery)
	if err != nil {
		return err
	}
	if len(indicators) == 0 {
		return nil
	}

	indicatorMap := make(map[string]indicatorWithInjector, len(indicators))
	for _, ind := range indicators {
		indicatorMap[ind.GetId()] = indicatorWithInjector{
			indicator:           ind,
			msgToSensorInjector: injector,
		}
	}

	m.generateProcessAlertsAndEnforcement(indicatorMap, nil)
	return nil
}

func (m *managerImpl) DeploymentUpdated(ctx enricher.EnrichmentContext, deployment *storage.Deployment, create bool, injector common.MessageInjector) error {
	retrievedInjector := m.deploymentsPendingEnrichment.removeAndRetrieveInjector(deployment.GetId())
	// Enforcement-related: IF the pending deployment had an injector, that means we have an enforcement decision pending.
	// If we deleted it, we would lose the opportunity to perform that enforcement forever.
	// So, delete it, but keep the enforcement injector.
	// Doing it this way ensures that we are performing detection on the most up-to-date version of the deployment
	// (which is the argument passed to this function); the one in the pendingCache might be stale.
	// This also ensures that we're more likely to persist the image (since we may not get the image ID until an update)
	if injector == nil && retrievedInjector != nil {
		injector = retrievedInjector
	}
	enrichmentPending, err := m.processDeploymentUpdate(ctx, deployment, injector)
	if err != nil {
		return err
	}
	if enrichmentPending {
		m.deploymentsPendingEnrichment.add(ctx, deployment, injector)
	}
	if create {
		if err := m.generateRuntimeAlertsOnCreate(deployment, injector); err != nil {
			log.Errorf("Could not generate runtime alerts for deployment %s: %v", deployment.GetId(), err)
		}
	}
	return nil
}

func (m *managerImpl) HandleAlerts(deploymentID string, alerts []*storage.Alert, stage storage.LifecycleStage) error {
	defer m.reprocessor.ReprocessRiskForDeployments(deploymentID)

	if _, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, alerts,
		alertmanager.WithLifecycleStage(stage), alertmanager.WithDeploymentIDs(deploymentID)); err != nil {
		return err
	}

	return nil
}

func (m *managerImpl) processDeploymentUpdate(ctx enricher.EnrichmentContext, deployment *storage.Deployment, injector common.MessageInjector) (bool, error) {
	// Attempt to enrich the image before detection.
	images, updatedIndices, pendingEnrichment, err := m.enricher.EnrichDeployment(ctx, deployment)
	if err != nil {
		log.Errorf("Error enriching deployment %s: %s", deployment.GetName(), err)
	}
	if len(updatedIndices) > 0 {
		for _, idx := range updatedIndices {
			img := images[idx]
			if img.GetId() == "" {
				continue
			}
			if err := m.imageDataStore.UpsertImage(lifecycleMgrCtx, img); err != nil {
				log.Errorf("Error persisting image %s: %s", img.GetName().GetFullName(), err)
			}
		}
	}

	// Update risk after processing and save the deployment.
	// There is no need to save the deployment in this function as it will be saved post reprocessing risk
	defer m.riskManager.ReprocessDeploymentRiskWithImages(deployment, images)

	presentAlerts, err := m.deploytimeDetector.Detect(deployTimePkg.DetectionContext{}, deployment, images)
	if err != nil {
		return false, errors.Wrap(err, "fetching deploy time alerts")
	}
	if _, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, presentAlerts,
		alertmanager.WithLifecycleStage(storage.LifecycleStage_DEPLOY), alertmanager.WithDeploymentIDs(deployment.GetId())); err != nil {
		return false, err
	}
	m.maybeInjectEnforcement(presentAlerts, deployment, injector)

	return pendingEnrichment, nil
}

func (m *managerImpl) maybeInjectEnforcement(presentAlerts []*storage.Alert, deployment *storage.Deployment, injector common.MessageInjector) {
	// If we're not passed an injector, that's our signal that nobody cares about any enforcement action. (example: on deployment updates)
	if injector == nil {
		return
	}

	// Generate enforcement actions based on the currently generated alerts.
	resp := determineEnforcementForDeployment(presentAlerts, deployment)
	// No enforcement, all good!
	if resp == nil {
		return
	}

	// Log if we are not enforcing because of an annotation.
	if !enforcers.ShouldEnforce(deployment.GetAnnotations()) {
		log.Warnf("Did not inject enforcement because deployment %s contained Enforcement Bypass annotations", deployment.GetName())
		return
	}

	err := injector.InjectMessage(context.Background(), &central.MsgToSensor{
		Msg: &central.MsgToSensor_Enforcement{
			Enforcement: resp,
		},
	})
	if err != nil {
		log.Errorf("Failed to inject enforcement action %s: %v", proto.MarshalTextString(resp), err)
	}
}

func (m *managerImpl) UpsertPolicy(policy *storage.Policy) error {
	var presentAlerts []*storage.Alert

	m.policyAlertsLock.Lock(policy.GetId())
	defer m.policyAlertsLock.Unlock(policy.GetId())
	// Add policy to set.
	if policies.AppliesAtDeployTime(policy) {
		if err := m.deploytimeDetector.PolicySet().UpsertPolicy(policy); err != nil {
			return errors.Wrapf(err, "adding policy %s to deploy time detector", policy.GetName())
		}
		deployTimeAlerts, err := m.deploytimeDetector.AlertsForPolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "error generating deploy-time alerts for policy %s", policy.GetName())
		}
		presentAlerts = append(presentAlerts, deployTimeAlerts...)
	} else {
		err := m.deploytimeDetector.PolicySet().RemovePolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "removing policy %s from deploy time detector", policy.GetName())
		}
	}

	if policies.AppliesAtRunTime(policy) {
		if err := m.runtimeDetector.PolicySet().UpsertPolicy(policy); err != nil {
			return errors.Wrapf(err, "adding policy %s to runtime detector", policy.GetName())
		}
		runTimeAlerts, err := m.runtimeDetector.AlertsForPolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "error generating runtime alerts for policy %s", policy.GetName())
		}
		presentAlerts = append(presentAlerts, runTimeAlerts...)
	} else {
		err := m.runtimeDetector.PolicySet().RemovePolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "removing policy %s from runtime detector", policy.GetName())
		}
	}

	// Perform notifications and update DB.
	modifiedDeployments, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, presentAlerts, alertmanager.WithPolicyID(policy.GetId()))
	if err != nil {
		return err
	}
	if modifiedDeployments.Cardinality() > 0 {
		defer m.reprocessor.ReprocessRiskForDeployments(modifiedDeployments.AsSlice()...)
	}
	return nil
}

func (m *managerImpl) DeploymentRemoved(deployment *storage.Deployment) error {
	m.deploymentsPendingEnrichment.remove(deployment.GetId())
	_, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, nil, alertmanager.WithDeploymentIDs(deployment.GetId()))
	return err
}

func (m *managerImpl) RemovePolicy(policyID string) error {
	m.policyAlertsLock.Lock(policyID)
	defer m.policyAlertsLock.Unlock(policyID)
	if err := m.deploytimeDetector.PolicySet().RemovePolicy(policyID); err != nil {
		return err
	}
	if err := m.runtimeDetector.PolicySet().RemovePolicy(policyID); err != nil {
		return err
	}
	modifiedDeployments, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, nil, alertmanager.WithPolicyID(policyID))
	if err != nil {
		return err
	}
	if modifiedDeployments.Cardinality() > 0 {
		m.reprocessor.ReprocessRiskForDeployments(modifiedDeployments.AsSlice()...)
	}
	return nil
}

func deploymentAndAlertToEnforcementProto(deployment *storage.Deployment, alert *storage.Alert) *central.SensorEnforcement {
	return &central.SensorEnforcement{
		Enforcement: alert.GetEnforcement().GetAction(),
		Resource: &central.SensorEnforcement_Deployment{
			Deployment: &central.DeploymentEnforcement{
				DeploymentId:   deployment.GetId(),
				DeploymentName: deployment.GetName(),
				DeploymentType: deployment.GetType(),
				Namespace:      deployment.GetNamespace(),
				AlertId:        alert.GetId(),
				PolicyName:     alert.GetPolicy().GetName(),
			},
		},
	}
}

// determineEnforcement returns the alert and its enforcement action to use from the input list (if any have enforcement).
func determineEnforcementForDeployment(alerts []*storage.Alert, deployment *storage.Deployment) *central.SensorEnforcement {
	// Only form and return the response if there is an enforcement action to be taken.
	var candidate *central.SensorEnforcement
	for _, alert := range alerts {
		// Prioritize scale to zero, so return immediately.
		if alert.GetEnforcement().GetAction() == storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT {
			return deploymentAndAlertToEnforcementProto(deployment, alert)
		}
		if candidate == nil && alert.GetEnforcement().GetAction() != storage.EnforcementAction_UNSET_ENFORCEMENT {
			candidate = deploymentAndAlertToEnforcementProto(deployment, alert)
		}
	}
	return candidate
}

func uniqueDeploymentIDs(indicatorsToInfo map[string]indicatorWithInjector) []string {
	m := set.NewStringSet()
	for _, infoWithInjector := range indicatorsToInfo {
		deploymentID := infoWithInjector.indicator.GetDeploymentId()
		if deploymentID == "" {
			continue
		}
		m.Add(deploymentID)
	}
	return m.AsSlice()
}

func (m *managerImpl) generateAndSendEnforcements(alerts []*storage.Alert, indicatorsToInfo map[string]indicatorWithInjector) {
	for _, alert := range alerts {
		// Skip alerts without runtime enforcement.
		if alert.GetEnforcement().GetAction() != storage.EnforcementAction_KILL_POD_ENFORCEMENT {
			continue
		}

		// If the alert has enforcement, we want to generate a list of enforcement and injector pairs.
		for _, singleIndicator := range alert.GetProcessViolation().GetProcesses() {
			if infoWithInjector, ok := indicatorsToInfo[singleIndicator.GetId()]; ok {
				// Generate the enforcement action.
				enforcement := createEnforcementAction(alert, infoWithInjector.indicator.GetPodId())
				// Attempt to send the enforcement with the injector.
				err := infoWithInjector.msgToSensorInjector.InjectMessage(context.Background(), &central.MsgToSensor{
					Msg: &central.MsgToSensor_Enforcement{
						Enforcement: enforcement,
					},
				})
				if err != nil {
					log.Errorf("Failed to inject enforcement action %s: %v", proto.MarshalTextString(enforcement), err)
				}
			}
		}
	}
}

func createEnforcementAction(alert *storage.Alert, podID string) *central.SensorEnforcement {
	resource := &central.SensorEnforcement_ContainerInstance{
		ContainerInstance: &central.ContainerInstanceEnforcement{
			PodId: podID,
			DeploymentEnforcement: &central.DeploymentEnforcement{
				DeploymentId:   alert.GetDeployment().GetId(),
				DeploymentName: alert.GetDeployment().GetName(),
				Namespace:      alert.GetDeployment().GetNamespace(),
				PolicyName:     alert.GetPolicy().GetName(),
			},
		},
	}
	return &central.SensorEnforcement{
		Enforcement: storage.EnforcementAction_KILL_POD_ENFORCEMENT,
		Resource:    resource,
	}
}
