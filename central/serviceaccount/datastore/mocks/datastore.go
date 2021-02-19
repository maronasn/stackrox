// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	reflect "reflect"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// Search mocks base method
func (m *MockDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// Count mocks base method
func (m *MockDataStore) Count(ctx context.Context, q *v1.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, q)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockDataStoreMockRecorder) Count(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDataStore)(nil).Count), ctx, q)
}

// SearchRawServiceAccounts mocks base method
func (m *MockDataStore) SearchRawServiceAccounts(ctx context.Context, q *v1.Query) ([]*storage.ServiceAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawServiceAccounts", ctx, q)
	ret0, _ := ret[0].([]*storage.ServiceAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawServiceAccounts indicates an expected call of SearchRawServiceAccounts
func (mr *MockDataStoreMockRecorder) SearchRawServiceAccounts(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawServiceAccounts", reflect.TypeOf((*MockDataStore)(nil).SearchRawServiceAccounts), ctx, q)
}

// SearchServiceAccounts mocks base method
func (m *MockDataStore) SearchServiceAccounts(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchServiceAccounts", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchServiceAccounts indicates an expected call of SearchServiceAccounts
func (mr *MockDataStoreMockRecorder) SearchServiceAccounts(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchServiceAccounts", reflect.TypeOf((*MockDataStore)(nil).SearchServiceAccounts), ctx, q)
}

// GetServiceAccount mocks base method
func (m *MockDataStore) GetServiceAccount(ctx context.Context, id string) (*storage.ServiceAccount, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceAccount", ctx, id)
	ret0, _ := ret[0].(*storage.ServiceAccount)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetServiceAccount indicates an expected call of GetServiceAccount
func (mr *MockDataStoreMockRecorder) GetServiceAccount(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceAccount", reflect.TypeOf((*MockDataStore)(nil).GetServiceAccount), ctx, id)
}

// UpsertServiceAccount mocks base method
func (m *MockDataStore) UpsertServiceAccount(ctx context.Context, request *storage.ServiceAccount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertServiceAccount", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertServiceAccount indicates an expected call of UpsertServiceAccount
func (mr *MockDataStoreMockRecorder) UpsertServiceAccount(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertServiceAccount", reflect.TypeOf((*MockDataStore)(nil).UpsertServiceAccount), ctx, request)
}

// RemoveServiceAccount mocks base method
func (m *MockDataStore) RemoveServiceAccount(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveServiceAccount", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveServiceAccount indicates an expected call of RemoveServiceAccount
func (mr *MockDataStoreMockRecorder) RemoveServiceAccount(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServiceAccount", reflect.TypeOf((*MockDataStore)(nil).RemoveServiceAccount), ctx, id)
}
