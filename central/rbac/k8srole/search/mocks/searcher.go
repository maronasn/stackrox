// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/rbac/k8srole/search (interfaces: Searcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
)

// MockSearcher is a mock of Searcher interface.
type MockSearcher struct {
	ctrl     *gomock.Controller
	recorder *MockSearcherMockRecorder
}

// MockSearcherMockRecorder is the mock recorder for MockSearcher.
type MockSearcherMockRecorder struct {
	mock *MockSearcher
}

// NewMockSearcher creates a new mock instance.
func NewMockSearcher(ctrl *gomock.Controller) *MockSearcher {
	mock := &MockSearcher{ctrl: ctrl}
	mock.recorder = &MockSearcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearcher) EXPECT() *MockSearcherMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockSearcher) Count(arg0 context.Context, arg1 *v1.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockSearcherMockRecorder) Count(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockSearcher)(nil).Count), arg0, arg1)
}

// Search mocks base method.
func (m *MockSearcher) Search(arg0 context.Context, arg1 *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0, arg1)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockSearcherMockRecorder) Search(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockSearcher)(nil).Search), arg0, arg1)
}

// SearchRawRoles mocks base method.
func (m *MockSearcher) SearchRawRoles(arg0 context.Context, arg1 *v1.Query) ([]*storage.K8SRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawRoles", arg0, arg1)
	ret0, _ := ret[0].([]*storage.K8SRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawRoles indicates an expected call of SearchRawRoles.
func (mr *MockSearcherMockRecorder) SearchRawRoles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawRoles", reflect.TypeOf((*MockSearcher)(nil).SearchRawRoles), arg0, arg1)
}

// SearchRoles mocks base method.
func (m *MockSearcher) SearchRoles(arg0 context.Context, arg1 *v1.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRoles", arg0, arg1)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRoles indicates an expected call of SearchRoles.
func (mr *MockSearcherMockRecorder) SearchRoles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRoles", reflect.TypeOf((*MockSearcher)(nil).SearchRoles), arg0, arg1)
}
