// Code generated by MockGen. DO NOT EDIT.
// Source: jolokia/jolokia.go

// Package mock_jolokia is a generated GoMock package.
package mock_jolokia

import (
	reflect "reflect"

	jolokia "github.com/artemiscloud/activemq-artemis-management/jolokia"
	gomock "github.com/golang/mock/gomock"
)

// MockIData is a mock of IData interface.
type MockIData struct {
	ctrl     *gomock.Controller
	recorder *MockIDataMockRecorder
}

// MockIDataMockRecorder is the mock recorder for MockIData.
type MockIDataMockRecorder struct {
	mock *MockIData
}

// NewMockIData creates a new mock instance.
func NewMockIData(ctrl *gomock.Controller) *MockIData {
	mock := &MockIData{ctrl: ctrl}
	mock.recorder = &MockIDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIData) EXPECT() *MockIDataMockRecorder {
	return m.recorder
}

// Print mocks base method.
func (m *MockIData) Print() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Print")
}

// Print indicates an expected call of Print.
func (mr *MockIDataMockRecorder) Print() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Print", reflect.TypeOf((*MockIData)(nil).Print))
}

// MockIJolokia is a mock of IJolokia interface.
type MockIJolokia struct {
	ctrl     *gomock.Controller
	recorder *MockIJolokiaMockRecorder
}

// MockIJolokiaMockRecorder is the mock recorder for MockIJolokia.
type MockIJolokiaMockRecorder struct {
	mock *MockIJolokia
}

// NewMockIJolokia creates a new mock instance.
func NewMockIJolokia(ctrl *gomock.Controller) *MockIJolokia {
	mock := &MockIJolokia{ctrl: ctrl}
	mock.recorder = &MockIJolokiaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIJolokia) EXPECT() *MockIJolokiaMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockIJolokia) Exec(path, postJsonString string) (*jolokia.ResponseData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", path, postJsonString)
	ret0, _ := ret[0].(*jolokia.ResponseData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockIJolokiaMockRecorder) Exec(path, postJsonString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockIJolokia)(nil).Exec), path, postJsonString)
}

// Read mocks base method.
func (m *MockIJolokia) Read(path string) (*jolokia.ResponseData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", path)
	ret0, _ := ret[0].(*jolokia.ResponseData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIJolokiaMockRecorder) Read(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIJolokia)(nil).Read), path)
}
