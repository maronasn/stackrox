// Code generated by MockGen. DO NOT EDIT.
// Source: scope_checker.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	permissions "github.com/stackrox/rox/pkg/auth/permissions"
	sac "github.com/stackrox/rox/pkg/sac"
	effectiveaccessscope "github.com/stackrox/rox/pkg/sac/effectiveaccessscope"
)

// MockScopeChecker is a mock of ScopeChecker interface.
type MockScopeChecker struct {
	ctrl     *gomock.Controller
	recorder *MockScopeCheckerMockRecorder
}

// MockScopeCheckerMockRecorder is the mock recorder for MockScopeChecker.
type MockScopeCheckerMockRecorder struct {
	mock *MockScopeChecker
}

// NewMockScopeChecker creates a new mock instance.
func NewMockScopeChecker(ctrl *gomock.Controller) *MockScopeChecker {
	mock := &MockScopeChecker{ctrl: ctrl}
	mock.recorder = &MockScopeCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScopeChecker) EXPECT() *MockScopeCheckerMockRecorder {
	return m.recorder
}

// AccessMode mocks base method.
func (m *MockScopeChecker) AccessMode(am storage.Access) sac.ScopeChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccessMode", am)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// AccessMode indicates an expected call of AccessMode.
func (mr *MockScopeCheckerMockRecorder) AccessMode(am interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccessMode", reflect.TypeOf((*MockScopeChecker)(nil).AccessMode), am)
}

// AllAllowed mocks base method.
func (m *MockScopeChecker) AllAllowed(ctx context.Context, subScopeKeyss [][]sac.ScopeKey) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllAllowed", ctx, subScopeKeyss)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllAllowed indicates an expected call of AllAllowed.
func (mr *MockScopeCheckerMockRecorder) AllAllowed(ctx, subScopeKeyss interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllAllowed", reflect.TypeOf((*MockScopeChecker)(nil).AllAllowed), ctx, subScopeKeyss)
}

// Allowed mocks base method.
func (m *MockScopeChecker) Allowed(subScopeKeys ...sac.ScopeKey) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range subScopeKeys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Allowed", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Allowed indicates an expected call of Allowed.
func (mr *MockScopeCheckerMockRecorder) Allowed(subScopeKeys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Allowed", reflect.TypeOf((*MockScopeChecker)(nil).Allowed), subScopeKeys...)
}

// Check mocks base method.
func (m *MockScopeChecker) Check(ctx context.Context, pred sac.ScopePredicate) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, pred)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockScopeCheckerMockRecorder) Check(ctx, pred interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockScopeChecker)(nil).Check), ctx, pred)
}

// ClusterID mocks base method.
func (m *MockScopeChecker) ClusterID(clusterID string) sac.ScopeChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClusterID", clusterID)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// ClusterID indicates an expected call of ClusterID.
func (mr *MockScopeCheckerMockRecorder) ClusterID(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClusterID", reflect.TypeOf((*MockScopeChecker)(nil).ClusterID), clusterID)
}

// EffectiveAccessScope mocks base method.
func (m *MockScopeChecker) EffectiveAccessScope(resource permissions.ResourceWithAccess) (*effectiveaccessscope.ScopeTree, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EffectiveAccessScope", resource)
	ret0, _ := ret[0].(*effectiveaccessscope.ScopeTree)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EffectiveAccessScope indicates an expected call of EffectiveAccessScope.
func (mr *MockScopeCheckerMockRecorder) EffectiveAccessScope(resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EffectiveAccessScope", reflect.TypeOf((*MockScopeChecker)(nil).EffectiveAccessScope), resource)
}

// ForClusterScopedObject mocks base method.
func (m *MockScopeChecker) ForClusterScopedObject(obj sac.ClusterScopedObject) sac.ScopeChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForClusterScopedObject", obj)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// ForClusterScopedObject indicates an expected call of ForClusterScopedObject.
func (mr *MockScopeCheckerMockRecorder) ForClusterScopedObject(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForClusterScopedObject", reflect.TypeOf((*MockScopeChecker)(nil).ForClusterScopedObject), obj)
}

// ForNamespaceScopedObject mocks base method.
func (m *MockScopeChecker) ForNamespaceScopedObject(obj sac.NamespaceScopedObject) sac.ScopeChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForNamespaceScopedObject", obj)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// ForNamespaceScopedObject indicates an expected call of ForNamespaceScopedObject.
func (mr *MockScopeCheckerMockRecorder) ForNamespaceScopedObject(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForNamespaceScopedObject", reflect.TypeOf((*MockScopeChecker)(nil).ForNamespaceScopedObject), obj)
}

// Namespace mocks base method.
func (m *MockScopeChecker) Namespace(namespace string) sac.ScopeChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Namespace", namespace)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// Namespace indicates an expected call of Namespace.
func (mr *MockScopeCheckerMockRecorder) Namespace(namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Namespace", reflect.TypeOf((*MockScopeChecker)(nil).Namespace), namespace)
}

// Resource mocks base method.
func (m *MockScopeChecker) Resource(resource permissions.ResourceHandle) sac.ScopeChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resource", resource)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// Resource indicates an expected call of Resource.
func (mr *MockScopeCheckerMockRecorder) Resource(resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resource", reflect.TypeOf((*MockScopeChecker)(nil).Resource), resource)
}

// SubScopeChecker mocks base method.
func (m *MockScopeChecker) SubScopeChecker(keys ...sac.ScopeKey) sac.ScopeChecker {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubScopeChecker", varargs...)
	ret0, _ := ret[0].(sac.ScopeChecker)
	return ret0
}

// SubScopeChecker indicates an expected call of SubScopeChecker.
func (mr *MockScopeCheckerMockRecorder) SubScopeChecker(keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubScopeChecker", reflect.TypeOf((*MockScopeChecker)(nil).SubScopeChecker), keys...)
}
