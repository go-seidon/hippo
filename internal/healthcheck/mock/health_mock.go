// Code generated by MockGen. DO NOT EDIT.
// Source: internal/healthcheck/health.go

// Package mock_healthcheck is a generated GoMock package.
package mock_healthcheck

import (
	reflect "reflect"

	healthcheck "github.com/go-seidon/local/internal/healthcheck"
	gomock "github.com/golang/mock/gomock"
)

// MockHealthCheck is a mock of HealthCheck interface.
type MockHealthCheck struct {
	ctrl     *gomock.Controller
	recorder *MockHealthCheckMockRecorder
}

// MockHealthCheckMockRecorder is the mock recorder for MockHealthCheck.
type MockHealthCheckMockRecorder struct {
	mock *MockHealthCheck
}

// NewMockHealthCheck creates a new mock instance.
func NewMockHealthCheck(ctrl *gomock.Controller) *MockHealthCheck {
	mock := &MockHealthCheck{ctrl: ctrl}
	mock.recorder = &MockHealthCheckMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthCheck) EXPECT() *MockHealthCheckMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockHealthCheck) Check() (*healthcheck.CheckResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check")
	ret0, _ := ret[0].(*healthcheck.CheckResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockHealthCheckMockRecorder) Check() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockHealthCheck)(nil).Check))
}

// Start mocks base method.
func (m *MockHealthCheck) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockHealthCheckMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockHealthCheck)(nil).Start))
}

// Stop mocks base method.
func (m *MockHealthCheck) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockHealthCheckMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockHealthCheck)(nil).Stop))
}

// MockChecker is a mock of Checker interface.
type MockChecker struct {
	ctrl     *gomock.Controller
	recorder *MockCheckerMockRecorder
}

// MockCheckerMockRecorder is the mock recorder for MockChecker.
type MockCheckerMockRecorder struct {
	mock *MockChecker
}

// NewMockChecker creates a new mock instance.
func NewMockChecker(ctrl *gomock.Controller) *MockChecker {
	mock := &MockChecker{ctrl: ctrl}
	mock.recorder = &MockCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChecker) EXPECT() *MockCheckerMockRecorder {
	return m.recorder
}

// Status mocks base method.
func (m *MockChecker) Status() (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *MockCheckerMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockChecker)(nil).Status))
}