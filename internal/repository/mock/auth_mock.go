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

// MockAuthRepository is a mock of AuthRepository interface.
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository.
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance.
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// FindClient mocks base method.
func (m *MockAuthRepository) FindClient(ctx context.Context, p repository.FindClientParam) (*repository.FindClientResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindClient", ctx, p)
	ret0, _ := ret[0].(*repository.FindClientResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindClient indicates an expected call of FindClient.
func (mr *MockAuthRepositoryMockRecorder) FindClient(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindClient", reflect.TypeOf((*MockAuthRepository)(nil).FindClient), ctx, p)
}
