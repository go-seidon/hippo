// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/auth.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	repository "github.com/go-seidon/hippo/internal/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// CreateClient mocks base method.
func (m *MockAuth) CreateClient(ctx context.Context, p repository.CreateClientParam) (*repository.CreateClientResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClient", ctx, p)
	ret0, _ := ret[0].(*repository.CreateClientResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateClient indicates an expected call of CreateClient.
func (mr *MockAuthMockRecorder) CreateClient(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClient", reflect.TypeOf((*MockAuth)(nil).CreateClient), ctx, p)
}

// FindClient mocks base method.
func (m *MockAuth) FindClient(ctx context.Context, p repository.FindClientParam) (*repository.FindClientResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindClient", ctx, p)
	ret0, _ := ret[0].(*repository.FindClientResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindClient indicates an expected call of FindClient.
func (mr *MockAuthMockRecorder) FindClient(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindClient", reflect.TypeOf((*MockAuth)(nil).FindClient), ctx, p)
}

// SearchClient mocks base method.
func (m *MockAuth) SearchClient(ctx context.Context, p repository.SearchClientParam) (*repository.SearchClientResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchClient", ctx, p)
	ret0, _ := ret[0].(*repository.SearchClientResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchClient indicates an expected call of SearchClient.
func (mr *MockAuthMockRecorder) SearchClient(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchClient", reflect.TypeOf((*MockAuth)(nil).SearchClient), ctx, p)
}

// UpdateClient mocks base method.
func (m *MockAuth) UpdateClient(ctx context.Context, p repository.UpdateClientParam) (*repository.UpdateClientResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClient", ctx, p)
	ret0, _ := ret[0].(*repository.UpdateClientResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateClient indicates an expected call of UpdateClient.
func (mr *MockAuthMockRecorder) UpdateClient(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClient", reflect.TypeOf((*MockAuth)(nil).UpdateClient), ctx, p)
}
