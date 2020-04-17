// Code generated by MockGen. DO NOT EDIT.
// Source: dao/datastore.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	dao "github.com/rest_api_creator/backend-sls/dao"
	errors "github.com/rest_api_creator/backend-sls/errors"
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

// CreateUser mocks base method
func (m *MockDataStore) CreateUser(arg0, arg1, arg2 string) errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(errors.ApiError)
	return ret0
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockDataStoreMockRecorder) CreateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockDataStore)(nil).CreateUser), arg0, arg1, arg2)
}

// GetUser mocks base method
func (m *MockDataStore) GetUser(arg0 string) (dao.User, errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(errors.ApiError)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockDataStoreMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockDataStore)(nil).GetUser), arg0)
}

// GetProject mocks base method
func (m *MockDataStore) GetProject(arg0, arg1 string) (dao.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", arg0, arg1)
	ret0, _ := ret[0].(dao.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject
func (mr *MockDataStoreMockRecorder) GetProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockDataStore)(nil).GetProject), arg0, arg1)
}

// UpdateUserToken mocks base method
func (m *MockDataStore) UpdateUserToken(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserToken indicates an expected call of UpdateUserToken
func (mr *MockDataStoreMockRecorder) UpdateUserToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserToken", reflect.TypeOf((*MockDataStore)(nil).UpdateUserToken), arg0, arg1)
}
