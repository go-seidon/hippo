// Code generated by MockGen. DO NOT EDIT.
// Source: internal/healthcheck/health.go

// Package mock_healthcheck is a generated GoMock package.
package mock_healthcheck

import (
	context "context"
	reflect "reflect"

	healthcheck "github.com/go-seidon/hippo/internal/healthcheck"
	system "github.com/go-seidon/provider/system"
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
func (m *MockHealthCheck) Check(ctx context.Context) (*healthcheck.CheckResult, *system.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx)
	ret0, _ := ret[0].(*healthcheck.CheckResult)
	ret1, _ := ret[1].(*system.Error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockHealthCheckMockRecorder) Check(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockHealthCheck)(nil).Check), ctx)
}
