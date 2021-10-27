// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
)

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// GetNetworkGraphConfig mocks base method.
func (m *MockDataStore) GetNetworkGraphConfig(ctx context.Context) (*storage.NetworkGraphConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkGraphConfig", ctx)
	ret0, _ := ret[0].(*storage.NetworkGraphConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkGraphConfig indicates an expected call of GetNetworkGraphConfig.
func (mr *MockDataStoreMockRecorder) GetNetworkGraphConfig(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkGraphConfig", reflect.TypeOf((*MockDataStore)(nil).GetNetworkGraphConfig), ctx)
}

// UpdateNetworkGraphConfig mocks base method.
func (m *MockDataStore) UpdateNetworkGraphConfig(ctx context.Context, config *storage.NetworkGraphConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNetworkGraphConfig", ctx, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNetworkGraphConfig indicates an expected call of UpdateNetworkGraphConfig.
func (mr *MockDataStoreMockRecorder) UpdateNetworkGraphConfig(ctx, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNetworkGraphConfig", reflect.TypeOf((*MockDataStore)(nil).UpdateNetworkGraphConfig), ctx, config)
}
