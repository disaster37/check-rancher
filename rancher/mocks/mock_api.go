// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/disaster37/check-rancher/v2/rancher/api (interfaces: API,ClusterAPI,ETCDBackupAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	rancherapi "github.com/disaster37/check-rancher/v2/rancher/api"
	gomock "github.com/golang/mock/gomock"
)

// MockAPI is a mock of API interface.
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI.
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance.
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// Cluster mocks base method.
func (m *MockAPI) Cluster() rancherapi.ClusterAPI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cluster")
	ret0, _ := ret[0].(rancherapi.ClusterAPI)
	return ret0
}

// Cluster indicates an expected call of Cluster.
func (mr *MockAPIMockRecorder) Cluster() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cluster", reflect.TypeOf((*MockAPI)(nil).Cluster))
}

// ETCDBackup mocks base method.
func (m *MockAPI) ETCDBackup() rancherapi.ETCDBackupAPI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETCDBackup")
	ret0, _ := ret[0].(rancherapi.ETCDBackupAPI)
	return ret0
}

// ETCDBackup indicates an expected call of ETCDBackup.
func (mr *MockAPIMockRecorder) ETCDBackup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETCDBackup", reflect.TypeOf((*MockAPI)(nil).ETCDBackup))
}

// MockClusterAPI is a mock of ClusterAPI interface.
type MockClusterAPI struct {
	ctrl     *gomock.Controller
	recorder *MockClusterAPIMockRecorder
}

// MockClusterAPIMockRecorder is the mock recorder for MockClusterAPI.
type MockClusterAPIMockRecorder struct {
	mock *MockClusterAPI
}

// NewMockClusterAPI creates a new mock instance.
func NewMockClusterAPI(ctrl *gomock.Controller) *MockClusterAPI {
	mock := &MockClusterAPI{ctrl: ctrl}
	mock.recorder = &MockClusterAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterAPI) EXPECT() *MockClusterAPIMockRecorder {
	return m.recorder
}

// GetByName mocks base method.
func (m *MockClusterAPI) GetByName(arg0 string) (*rancherapi.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0)
	ret0, _ := ret[0].(*rancherapi.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockClusterAPIMockRecorder) GetByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockClusterAPI)(nil).GetByName), arg0)
}

// List mocks base method.
func (m *MockClusterAPI) List() ([]*rancherapi.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*rancherapi.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockClusterAPIMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockClusterAPI)(nil).List))
}

// MockETCDBackupAPI is a mock of ETCDBackupAPI interface.
type MockETCDBackupAPI struct {
	ctrl     *gomock.Controller
	recorder *MockETCDBackupAPIMockRecorder
}

// MockETCDBackupAPIMockRecorder is the mock recorder for MockETCDBackupAPI.
type MockETCDBackupAPIMockRecorder struct {
	mock *MockETCDBackupAPI
}

// NewMockETCDBackupAPI creates a new mock instance.
func NewMockETCDBackupAPI(ctrl *gomock.Controller) *MockETCDBackupAPI {
	mock := &MockETCDBackupAPI{ctrl: ctrl}
	mock.recorder = &MockETCDBackupAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockETCDBackupAPI) EXPECT() *MockETCDBackupAPIMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockETCDBackupAPI) List() ([]*rancherapi.ETCDBackup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*rancherapi.ETCDBackup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockETCDBackupAPIMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockETCDBackupAPI)(nil).List))
}

// ListByClusterID mocks base method.
func (m *MockETCDBackupAPI) ListByClusterID(arg0 string) ([]*rancherapi.ETCDBackup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByClusterID", arg0)
	ret0, _ := ret[0].([]*rancherapi.ETCDBackup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByClusterID indicates an expected call of ListByClusterID.
func (mr *MockETCDBackupAPIMockRecorder) ListByClusterID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByClusterID", reflect.TypeOf((*MockETCDBackupAPI)(nil).ListByClusterID), arg0)
}