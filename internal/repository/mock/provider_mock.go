// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/provider.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	repository "github.com/go-seidon/local/internal/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// GetAuthRepo mocks base method.
func (m *MockProvider) GetAuthRepo() repository.AuthRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthRepo")
	ret0, _ := ret[0].(repository.AuthRepository)
	return ret0
}

// GetAuthRepo indicates an expected call of GetAuthRepo.
func (mr *MockProviderMockRecorder) GetAuthRepo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthRepo", reflect.TypeOf((*MockProvider)(nil).GetAuthRepo))
}

// GetFileRepo mocks base method.
func (m *MockProvider) GetFileRepo() repository.FileRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileRepo")
	ret0, _ := ret[0].(repository.FileRepository)
	return ret0
}

// GetFileRepo indicates an expected call of GetFileRepo.
func (mr *MockProviderMockRecorder) GetFileRepo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileRepo", reflect.TypeOf((*MockProvider)(nil).GetFileRepo))
}

// Init mocks base method.
func (m *MockProvider) Init(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockProviderMockRecorder) Init(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockProvider)(nil).Init), ctx)
}

// Ping mocks base method.
func (m *MockProvider) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockProviderMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockProvider)(nil).Ping), ctx)
}
