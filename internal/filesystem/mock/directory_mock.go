// Code generated by MockGen. DO NOT EDIT.
// Source: internal/filesystem/directory.go

// Package mock_filesystem is a generated GoMock package.
package mock_filesystem

import (
	context "context"
	reflect "reflect"

	filesystem "github.com/go-seidon/hippo/internal/filesystem"
	gomock "github.com/golang/mock/gomock"
)

// MockDirectoryManager is a mock of DirectoryManager interface.
type MockDirectoryManager struct {
	ctrl     *gomock.Controller
	recorder *MockDirectoryManagerMockRecorder
}

// MockDirectoryManagerMockRecorder is the mock recorder for MockDirectoryManager.
type MockDirectoryManagerMockRecorder struct {
	mock *MockDirectoryManager
}

// NewMockDirectoryManager creates a new mock instance.
func NewMockDirectoryManager(ctrl *gomock.Controller) *MockDirectoryManager {
	mock := &MockDirectoryManager{ctrl: ctrl}
	mock.recorder = &MockDirectoryManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDirectoryManager) EXPECT() *MockDirectoryManagerMockRecorder {
	return m.recorder
}

// CreateDir mocks base method.
func (m *MockDirectoryManager) CreateDir(ctx context.Context, p filesystem.CreateDirParam) (*filesystem.CreateDirResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDir", ctx, p)
	ret0, _ := ret[0].(*filesystem.CreateDirResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDir indicates an expected call of CreateDir.
func (mr *MockDirectoryManagerMockRecorder) CreateDir(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDir", reflect.TypeOf((*MockDirectoryManager)(nil).CreateDir), ctx, p)
}

// IsDirectoryExists mocks base method.
func (m *MockDirectoryManager) IsDirectoryExists(ctx context.Context, p filesystem.IsDirectoryExistsParam) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDirectoryExists", ctx, p)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsDirectoryExists indicates an expected call of IsDirectoryExists.
func (mr *MockDirectoryManagerMockRecorder) IsDirectoryExists(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDirectoryExists", reflect.TypeOf((*MockDirectoryManager)(nil).IsDirectoryExists), ctx, p)
}
