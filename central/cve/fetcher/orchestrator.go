package fetcher

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	"github.com/stackrox/rox/central/cve/converter"
	cveDataStore "github.com/stackrox/rox/central/cve/datastore"
	cveMatcher "github.com/stackrox/rox/central/cve/matcher"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/sac"
	pkgScanners "github.com/stackrox/rox/pkg/scanners"
	"github.com/stackrox/rox/pkg/scanners/types"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	readCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
			sac.ResourceScopeKeys(resources.Cluster),
		))
)

type orchestratorCVEManager struct {
	embeddedCVEIdToClusters map[converter.CVEType]map[string][]*storage.Cluster

	clusterDataStore clusterDataStore.DataStore
	cveDataStore     cveDataStore.DataStore
	cveMatcher       *cveMatcher.CVEMatcher

	creators map[string]pkgScanners.OrchestratorScannerCreator
	scanners map[string]types.OrchestratorScanner

	mutex sync.Mutex
}

func (m *orchestratorCVEManager) initialize() {
	m.Reconcile()
}

// Reconcile fetches new CVEs from scanner and reconciles them
func (m *orchestratorCVEManager) Reconcile() {
	clusters, err := m.clusterDataStore.GetClusters(readCtx)
	if err != nil {
		log.Errorf("failed to get clusters %v", err)
		return
	}
	log.Infof("Found %d clusters to scan for orchestrator vulnerabilities.", len(clusters))

	err = m.reconcileCVEs(clusters, converter.K8s)
	if err != nil {
		log.Errorf("failed to reconcile orchestrator Kubernetes CVEs: %v", err)
	}
	err = m.reconcileCVEs(clusters, converter.OpenShift)
	if err != nil {
		log.Errorf("failed to reconcile orchestrator OpenShift CVEs: %v", err)
	}
}

func (m *orchestratorCVEManager) Scan(version string, cveType converter.CVEType) ([]*storage.EmbeddedVulnerability, error) {
	scanners := map[string]types.OrchestratorScanner{}

	m.mutex.Lock()
	for k, v := range m.scanners {
		scanners[k] = v
	}
	m.mutex.Unlock()

	if len(scanners) == 0 {
		return nil, errors.New("no orchestrator scanners are integrated")
	}
	switch cveType {
	case converter.K8s:
		return k8sScan(version, scanners)
	case converter.OpenShift:
		return openShiftScan(version, scanners)
	}
	return nil, errors.Errorf("unexpected kind %s", cveType)
}

func (m *orchestratorCVEManager) updateCVEs(embeddedCVEs []*storage.EmbeddedVulnerability, embeddedCVEToClusters map[string][]*storage.Cluster, cveType converter.CVEType) error {
	newCVEIDs := set.NewStringSet()
	var newCVEs []converter.ClusterCVEParts
	for _, embeddedCVE := range embeddedCVEs {
		cve := converter.EmbeddedCVEToProtoCVE("", embeddedCVE)
		newCVEIDs.Add(cve.GetId())
		newCVEs = append(newCVEs, converter.NewClusterCVEParts(cve, embeddedCVEToClusters[embeddedCVE.GetCve()], embeddedCVE.GetFixedBy()))
	}

	m.embeddedCVEIdToClusters[cveType] = embeddedCVEToClusters
	return m.updateCVEsInDB(newCVEIDs, newCVEs, cveType)
}

func (m *orchestratorCVEManager) updateCVEsInDB(cveIds set.StringSet, cves []converter.ClusterCVEParts, cveType converter.CVEType) error {
	if err := m.cveDataStore.UpsertClusterCVEs(cveElevatedCtx, cves...); err != nil {
		return err
	}
	return reconcileCVEsInDB(m.cveDataStore, cveType.ToStorageCVEType(), cveIds)
}

// createOrchestratorScanner creates a types.OrchestratorScanner out of the given storage.OrchestratorIntegration.
func (m *orchestratorCVEManager) createOrchestratorScanner(source *storage.OrchestratorIntegration) (types.OrchestratorScanner, error) {
	creator, exists := m.creators[source.GetType()]
	if !exists {
		return nil, fmt.Errorf("scanner with type %q does not exist", source.GetType())
	}
	scanner, err := creator(source)
	if err != nil {
		return nil, err
	}
	return scanner, nil
}

func (m *orchestratorCVEManager) UpsertOrchestratorScanner(integration *storage.OrchestratorIntegration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	scanner, err := m.createOrchestratorScanner(integration)
	if err != nil {
		return errors.Wrap(err, "Failed to create orchestrator scanner")
	}
	m.scanners[integration.GetId()] = scanner
	return nil
}

func (m *orchestratorCVEManager) RemoveIntegration(integrationID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.scanners, integrationID)
}

func k8sScan(version string, scanners map[string]types.OrchestratorScanner) ([]*storage.EmbeddedVulnerability, error) {
	errorList := errorhelpers.NewErrorList(fmt.Sprintf("error scanning orchestrator for Kubernetes:%s", version))

	var allVulns []*storage.EmbeddedVulnerability
	for _, scanner := range scanners {
		result, err := scanner.KubernetesScan(version)
		if err != nil {
			errorList.AddError(err)
			continue
		}
		vulnIDsSet := set.NewStringSet()
		for _, v := range result {
			for _, vuln := range v {
				if vulnIDsSet.Add(vuln.GetCve()) {
					allVulns = append(allVulns, vuln)
				}
			}
		}
		return allVulns, nil
	}

	return nil, errorList.ToError()
}

func openShiftScan(version string, scanners map[string]types.OrchestratorScanner) ([]*storage.EmbeddedVulnerability, error) {
	errorList := errorhelpers.NewErrorList(fmt.Sprintf("error scanning orchestrator for OpenShift:%s", version))
	for _, scanner := range scanners {
		result, err := scanner.OpenShiftScan(version)
		if err != nil {
			errorList.AddError(err)
			continue
		}
		return result, nil
	}

	return nil, errorList.ToError()
}

func (m *orchestratorCVEManager) reconcileCVEs(clusters []*storage.Cluster, cveType converter.CVEType) error {
	versionToClusters := make(map[string][]*storage.Cluster)
	for _, cluster := range clusters {
		var version string
		switch cveType {
		case converter.K8s:
			version = cluster.GetStatus().GetOrchestratorMetadata().GetVersion()
		case converter.OpenShift:
			version = cluster.GetStatus().GetOrchestratorMetadata().GetOpenshiftVersion()
		}

		if version == "" {
			continue
		}
		versionToClusters[version] = append(versionToClusters[version], cluster)
	}

	embeddedCVEIDToClusters := make(map[string][]*storage.Cluster)
	var allEmbeddedCVEs []*storage.EmbeddedVulnerability
	for version := range versionToClusters {
		vulns, err := m.Scan(version, cveType)
		if err != nil {
			return err
		}
		for _, vuln := range vulns {
			if _, ok := embeddedCVEIDToClusters[vuln.GetCve()]; !ok {
				allEmbeddedCVEs = append(allEmbeddedCVEs, vuln)
			}
			embeddedCVEIDToClusters[vuln.GetCve()] = append(embeddedCVEIDToClusters[vuln.GetCve()], versionToClusters[version]...)
		}
	}

	err := m.updateCVEs(allEmbeddedCVEs, embeddedCVEIDToClusters, cveType)
	if err != nil {
		return err
	}
	log.Infof("Successfully fetched %d %s CVEs", len(m.embeddedCVEIdToClusters[cveType]), cveType)
	return nil
}

func (m *orchestratorCVEManager) getAffectedClusters(cveID string, cveType converter.CVEType) ([]*storage.Cluster, error) {
	if cveToClusters, ok := m.embeddedCVEIdToClusters[cveType]; ok {
		if clusters, ok := cveToClusters[cveID]; ok {
			return clusters, nil
		}
	}
	return nil, errors.Errorf("Cannot find cve with type %v id %s", cveType, cveID)
}
