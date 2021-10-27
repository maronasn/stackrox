// Code generated by MockGen. DO NOT EDIT.
// Source: pipeline.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/stackrox/rox/central/sensor/service/common"
	pipeline "github.com/stackrox/rox/central/sensor/service/pipeline"
	reconciliation "github.com/stackrox/rox/central/sensor/service/pipeline/reconciliation"
	central "github.com/stackrox/rox/generated/internalapi/central"
)

// MockBasePipeline is a mock of BasePipeline interface.
type MockBasePipeline struct {
	ctrl     *gomock.Controller
	recorder *MockBasePipelineMockRecorder
}

// MockBasePipelineMockRecorder is the mock recorder for MockBasePipeline.
type MockBasePipelineMockRecorder struct {
	mock *MockBasePipeline
}

// NewMockBasePipeline creates a new mock instance.
func NewMockBasePipeline(ctrl *gomock.Controller) *MockBasePipeline {
	mock := &MockBasePipeline{ctrl: ctrl}
	mock.recorder = &MockBasePipelineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBasePipeline) EXPECT() *MockBasePipelineMockRecorder {
	return m.recorder
}

// OnFinish mocks base method.
func (m *MockBasePipeline) OnFinish(clusterID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnFinish", clusterID)
}

// OnFinish indicates an expected call of OnFinish.
func (mr *MockBasePipelineMockRecorder) OnFinish(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnFinish", reflect.TypeOf((*MockBasePipeline)(nil).OnFinish), clusterID)
}

// MockClusterPipeline is a mock of ClusterPipeline interface.
type MockClusterPipeline struct {
	ctrl     *gomock.Controller
	recorder *MockClusterPipelineMockRecorder
}

// MockClusterPipelineMockRecorder is the mock recorder for MockClusterPipeline.
type MockClusterPipelineMockRecorder struct {
	mock *MockClusterPipeline
}

// NewMockClusterPipeline creates a new mock instance.
func NewMockClusterPipeline(ctrl *gomock.Controller) *MockClusterPipeline {
	mock := &MockClusterPipeline{ctrl: ctrl}
	mock.recorder = &MockClusterPipelineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterPipeline) EXPECT() *MockClusterPipelineMockRecorder {
	return m.recorder
}

// OnFinish mocks base method.
func (m *MockClusterPipeline) OnFinish(clusterID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnFinish", clusterID)
}

// OnFinish indicates an expected call of OnFinish.
func (mr *MockClusterPipelineMockRecorder) OnFinish(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnFinish", reflect.TypeOf((*MockClusterPipeline)(nil).OnFinish), clusterID)
}

// Reconcile mocks base method.
func (m *MockClusterPipeline) Reconcile(ctx context.Context, reconciliationStore *reconciliation.StoreMap) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reconcile", ctx, reconciliationStore)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reconcile indicates an expected call of Reconcile.
func (mr *MockClusterPipelineMockRecorder) Reconcile(ctx, reconciliationStore interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reconcile", reflect.TypeOf((*MockClusterPipeline)(nil).Reconcile), ctx, reconciliationStore)
}

// Run mocks base method.
func (m *MockClusterPipeline) Run(ctx context.Context, msg *central.MsgFromSensor, injector common.MessageInjector) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx, msg, injector)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockClusterPipelineMockRecorder) Run(ctx, msg, injector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockClusterPipeline)(nil).Run), ctx, msg, injector)
}

// MockFactory is a mock of Factory interface.
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory.
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance.
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// PipelineForCluster mocks base method.
func (m *MockFactory) PipelineForCluster(ctx context.Context, clusterID string) (pipeline.ClusterPipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PipelineForCluster", ctx, clusterID)
	ret0, _ := ret[0].(pipeline.ClusterPipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PipelineForCluster indicates an expected call of PipelineForCluster.
func (mr *MockFactoryMockRecorder) PipelineForCluster(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PipelineForCluster", reflect.TypeOf((*MockFactory)(nil).PipelineForCluster), ctx, clusterID)
}

// MockFragment is a mock of Fragment interface.
type MockFragment struct {
	ctrl     *gomock.Controller
	recorder *MockFragmentMockRecorder
}

// MockFragmentMockRecorder is the mock recorder for MockFragment.
type MockFragmentMockRecorder struct {
	mock *MockFragment
}

// NewMockFragment creates a new mock instance.
func NewMockFragment(ctrl *gomock.Controller) *MockFragment {
	mock := &MockFragment{ctrl: ctrl}
	mock.recorder = &MockFragmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFragment) EXPECT() *MockFragmentMockRecorder {
	return m.recorder
}

// Match mocks base method.
func (m *MockFragment) Match(msg *central.MsgFromSensor) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Match", msg)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Match indicates an expected call of Match.
func (mr *MockFragmentMockRecorder) Match(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Match", reflect.TypeOf((*MockFragment)(nil).Match), msg)
}

// OnFinish mocks base method.
func (m *MockFragment) OnFinish(clusterID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnFinish", clusterID)
}

// OnFinish indicates an expected call of OnFinish.
func (mr *MockFragmentMockRecorder) OnFinish(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnFinish", reflect.TypeOf((*MockFragment)(nil).OnFinish), clusterID)
}

// Reconcile mocks base method.
func (m *MockFragment) Reconcile(ctx context.Context, clusterID string, reconciliationStore *reconciliation.StoreMap) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reconcile", ctx, clusterID, reconciliationStore)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reconcile indicates an expected call of Reconcile.
func (mr *MockFragmentMockRecorder) Reconcile(ctx, clusterID, reconciliationStore interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reconcile", reflect.TypeOf((*MockFragment)(nil).Reconcile), ctx, clusterID, reconciliationStore)
}

// Run mocks base method.
func (m *MockFragment) Run(ctx context.Context, clusterID string, msg *central.MsgFromSensor, injector common.MessageInjector) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx, clusterID, msg, injector)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockFragmentMockRecorder) Run(ctx, clusterID, msg, injector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockFragment)(nil).Run), ctx, clusterID, msg, injector)
}

// MockFragmentFactory is a mock of FragmentFactory interface.
type MockFragmentFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFragmentFactoryMockRecorder
}

// MockFragmentFactoryMockRecorder is the mock recorder for MockFragmentFactory.
type MockFragmentFactoryMockRecorder struct {
	mock *MockFragmentFactory
}

// NewMockFragmentFactory creates a new mock instance.
func NewMockFragmentFactory(ctrl *gomock.Controller) *MockFragmentFactory {
	mock := &MockFragmentFactory{ctrl: ctrl}
	mock.recorder = &MockFragmentFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFragmentFactory) EXPECT() *MockFragmentFactoryMockRecorder {
	return m.recorder
}

// GetFragment mocks base method.
func (m *MockFragmentFactory) GetFragment(ctx context.Context, clusterID string) (pipeline.Fragment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFragment", ctx, clusterID)
	ret0, _ := ret[0].(pipeline.Fragment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFragment indicates an expected call of GetFragment.
func (mr *MockFragmentFactoryMockRecorder) GetFragment(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFragment", reflect.TypeOf((*MockFragmentFactory)(nil).GetFragment), ctx, clusterID)
}
