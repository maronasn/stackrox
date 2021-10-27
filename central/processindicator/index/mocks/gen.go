// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/processindicator/index (interfaces: Indexer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// AddProcessIndicator mocks base method.
func (m *MockIndexer) AddProcessIndicator(arg0 *storage.ProcessIndicator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProcessIndicator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProcessIndicator indicates an expected call of AddProcessIndicator.
func (mr *MockIndexerMockRecorder) AddProcessIndicator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProcessIndicator", reflect.TypeOf((*MockIndexer)(nil).AddProcessIndicator), arg0)
}

// AddProcessIndicators mocks base method.
func (m *MockIndexer) AddProcessIndicators(arg0 []*storage.ProcessIndicator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProcessIndicators", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProcessIndicators indicates an expected call of AddProcessIndicators.
func (mr *MockIndexerMockRecorder) AddProcessIndicators(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProcessIndicators", reflect.TypeOf((*MockIndexer)(nil).AddProcessIndicators), arg0)
}

// Count mocks base method.
func (m *MockIndexer) Count(arg0 *v1.Query, arg1 ...blevesearch.SearchOption) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Count", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIndexerMockRecorder) Count(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIndexer)(nil).Count), varargs...)
}

// DeleteProcessIndicator mocks base method.
func (m *MockIndexer) DeleteProcessIndicator(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProcessIndicator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProcessIndicator indicates an expected call of DeleteProcessIndicator.
func (mr *MockIndexerMockRecorder) DeleteProcessIndicator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProcessIndicator", reflect.TypeOf((*MockIndexer)(nil).DeleteProcessIndicator), arg0)
}

// DeleteProcessIndicators mocks base method.
func (m *MockIndexer) DeleteProcessIndicators(arg0 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProcessIndicators", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProcessIndicators indicates an expected call of DeleteProcessIndicators.
func (mr *MockIndexerMockRecorder) DeleteProcessIndicators(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProcessIndicators", reflect.TypeOf((*MockIndexer)(nil).DeleteProcessIndicators), arg0)
}

// MarkInitialIndexingComplete mocks base method.
func (m *MockIndexer) MarkInitialIndexingComplete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkInitialIndexingComplete")
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkInitialIndexingComplete indicates an expected call of MarkInitialIndexingComplete.
func (mr *MockIndexerMockRecorder) MarkInitialIndexingComplete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkInitialIndexingComplete", reflect.TypeOf((*MockIndexer)(nil).MarkInitialIndexingComplete))
}

// NeedsInitialIndexing mocks base method.
func (m *MockIndexer) NeedsInitialIndexing() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NeedsInitialIndexing")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NeedsInitialIndexing indicates an expected call of NeedsInitialIndexing.
func (mr *MockIndexerMockRecorder) NeedsInitialIndexing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NeedsInitialIndexing", reflect.TypeOf((*MockIndexer)(nil).NeedsInitialIndexing))
}

// Search mocks base method.
func (m *MockIndexer) Search(arg0 *v1.Query, arg1 ...blevesearch.SearchOption) ([]search.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Search", varargs...)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockIndexerMockRecorder) Search(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexer)(nil).Search), varargs...)
}
