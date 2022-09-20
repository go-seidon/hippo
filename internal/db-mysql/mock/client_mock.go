// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db-mysql/client.go

// Package mock_dbmysql is a generated GoMock package.
package mock_dbmysql

import (
	context "context"
	sql "database/sql"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockClient) Begin() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockClientMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockClient)(nil).Begin))
}

// BeginTx mocks base method.
func (m *MockClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, opts)
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockClientMockRecorder) BeginTx(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockClient)(nil).BeginTx), ctx, opts)
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// Exec mocks base method.
func (m *MockClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockClientMockRecorder) Exec(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockClient)(nil).Exec), varargs...)
}

// ExecContext mocks base method.
func (m *MockClient) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockClientMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockClient)(nil).ExecContext), varargs...)
}

// Ping mocks base method.
func (m *MockClient) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockClientMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockClient)(nil).Ping))
}

// PingContext mocks base method.
func (m *MockClient) PingContext(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingContext", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingContext indicates an expected call of PingContext.
func (mr *MockClientMockRecorder) PingContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingContext", reflect.TypeOf((*MockClient)(nil).PingContext), ctx)
}

// Prepare mocks base method.
func (m *MockClient) Prepare(query string) (*sql.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(*sql.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockClientMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockClient)(nil).Prepare), query)
}

// PrepareContext mocks base method.
func (m *MockClient) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContext", ctx, query)
	ret0, _ := ret[0].(*sql.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContext indicates an expected call of PrepareContext.
func (mr *MockClientMockRecorder) PrepareContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContext", reflect.TypeOf((*MockClient)(nil).PrepareContext), ctx, query)
}

// Query mocks base method.
func (m *MockClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockClientMockRecorder) Query(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockClient)(nil).Query), varargs...)
}

// QueryContext mocks base method.
func (m *MockClient) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockClientMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockClient)(nil).QueryContext), varargs...)
}

// QueryRow mocks base method.
func (m *MockClient) QueryRow(query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockClientMockRecorder) QueryRow(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockClient)(nil).QueryRow), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockClient) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockClientMockRecorder) QueryRowContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockClient)(nil).QueryRowContext), varargs...)
}

// SetConnMaxIdleTime mocks base method.
func (m *MockClient) SetConnMaxIdleTime(d time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConnMaxIdleTime", d)
}

// SetConnMaxIdleTime indicates an expected call of SetConnMaxIdleTime.
func (mr *MockClientMockRecorder) SetConnMaxIdleTime(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnMaxIdleTime", reflect.TypeOf((*MockClient)(nil).SetConnMaxIdleTime), d)
}

// SetConnMaxLifetime mocks base method.
func (m *MockClient) SetConnMaxLifetime(d time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConnMaxLifetime", d)
}

// SetConnMaxLifetime indicates an expected call of SetConnMaxLifetime.
func (mr *MockClientMockRecorder) SetConnMaxLifetime(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnMaxLifetime", reflect.TypeOf((*MockClient)(nil).SetConnMaxLifetime), d)
}

// SetMaxIdleConns mocks base method.
func (m *MockClient) SetMaxIdleConns(n int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxIdleConns", n)
}

// SetMaxIdleConns indicates an expected call of SetMaxIdleConns.
func (mr *MockClientMockRecorder) SetMaxIdleConns(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxIdleConns", reflect.TypeOf((*MockClient)(nil).SetMaxIdleConns), n)
}

// SetMaxOpenConns mocks base method.
func (m *MockClient) SetMaxOpenConns(n int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxOpenConns", n)
}

// SetMaxOpenConns indicates an expected call of SetMaxOpenConns.
func (mr *MockClientMockRecorder) SetMaxOpenConns(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxOpenConns", reflect.TypeOf((*MockClient)(nil).SetMaxOpenConns), n)
}

// MockPingable is a mock of Pingable interface.
type MockPingable struct {
	ctrl     *gomock.Controller
	recorder *MockPingableMockRecorder
}

// MockPingableMockRecorder is the mock recorder for MockPingable.
type MockPingableMockRecorder struct {
	mock *MockPingable
}

// NewMockPingable creates a new mock instance.
func NewMockPingable(ctrl *gomock.Controller) *MockPingable {
	mock := &MockPingable{ctrl: ctrl}
	mock.recorder = &MockPingableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPingable) EXPECT() *MockPingableMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockPingable) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPingableMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPingable)(nil).Ping))
}

// PingContext mocks base method.
func (m *MockPingable) PingContext(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingContext", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingContext indicates an expected call of PingContext.
func (mr *MockPingableMockRecorder) PingContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingContext", reflect.TypeOf((*MockPingable)(nil).PingContext), ctx)
}

// MockManageable is a mock of Manageable interface.
type MockManageable struct {
	ctrl     *gomock.Controller
	recorder *MockManageableMockRecorder
}

// MockManageableMockRecorder is the mock recorder for MockManageable.
type MockManageableMockRecorder struct {
	mock *MockManageable
}

// NewMockManageable creates a new mock instance.
func NewMockManageable(ctrl *gomock.Controller) *MockManageable {
	mock := &MockManageable{ctrl: ctrl}
	mock.recorder = &MockManageableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManageable) EXPECT() *MockManageableMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockManageable) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockManageableMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockManageable)(nil).Close))
}

// SetConnMaxIdleTime mocks base method.
func (m *MockManageable) SetConnMaxIdleTime(d time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConnMaxIdleTime", d)
}

// SetConnMaxIdleTime indicates an expected call of SetConnMaxIdleTime.
func (mr *MockManageableMockRecorder) SetConnMaxIdleTime(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnMaxIdleTime", reflect.TypeOf((*MockManageable)(nil).SetConnMaxIdleTime), d)
}

// SetConnMaxLifetime mocks base method.
func (m *MockManageable) SetConnMaxLifetime(d time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConnMaxLifetime", d)
}

// SetConnMaxLifetime indicates an expected call of SetConnMaxLifetime.
func (mr *MockManageableMockRecorder) SetConnMaxLifetime(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnMaxLifetime", reflect.TypeOf((*MockManageable)(nil).SetConnMaxLifetime), d)
}

// SetMaxIdleConns mocks base method.
func (m *MockManageable) SetMaxIdleConns(n int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxIdleConns", n)
}

// SetMaxIdleConns indicates an expected call of SetMaxIdleConns.
func (mr *MockManageableMockRecorder) SetMaxIdleConns(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxIdleConns", reflect.TypeOf((*MockManageable)(nil).SetMaxIdleConns), n)
}

// SetMaxOpenConns mocks base method.
func (m *MockManageable) SetMaxOpenConns(n int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxOpenConns", n)
}

// SetMaxOpenConns indicates an expected call of SetMaxOpenConns.
func (mr *MockManageableMockRecorder) SetMaxOpenConns(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxOpenConns", reflect.TypeOf((*MockManageable)(nil).SetMaxOpenConns), n)
}

// MockBeginable is a mock of Beginable interface.
type MockBeginable struct {
	ctrl     *gomock.Controller
	recorder *MockBeginableMockRecorder
}

// MockBeginableMockRecorder is the mock recorder for MockBeginable.
type MockBeginableMockRecorder struct {
	mock *MockBeginable
}

// NewMockBeginable creates a new mock instance.
func NewMockBeginable(ctrl *gomock.Controller) *MockBeginable {
	mock := &MockBeginable{ctrl: ctrl}
	mock.recorder = &MockBeginableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBeginable) EXPECT() *MockBeginableMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockBeginable) Begin() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockBeginableMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockBeginable)(nil).Begin))
}

// BeginTx mocks base method.
func (m *MockBeginable) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, opts)
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockBeginableMockRecorder) BeginTx(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockBeginable)(nil).BeginTx), ctx, opts)
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockTransaction) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTransactionMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTransaction)(nil).Commit))
}

// Exec mocks base method.
func (m *MockTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockTransactionMockRecorder) Exec(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockTransaction)(nil).Exec), varargs...)
}

// ExecContext mocks base method.
func (m *MockTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockTransactionMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockTransaction)(nil).ExecContext), varargs...)
}

// Prepare mocks base method.
func (m *MockTransaction) Prepare(query string) (*sql.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(*sql.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockTransactionMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockTransaction)(nil).Prepare), query)
}

// PrepareContext mocks base method.
func (m *MockTransaction) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContext", ctx, query)
	ret0, _ := ret[0].(*sql.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContext indicates an expected call of PrepareContext.
func (mr *MockTransactionMockRecorder) PrepareContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContext", reflect.TypeOf((*MockTransaction)(nil).PrepareContext), ctx, query)
}

// Query mocks base method.
func (m *MockTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockTransactionMockRecorder) Query(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockTransaction)(nil).Query), varargs...)
}

// QueryContext mocks base method.
func (m *MockTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockTransactionMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockTransaction)(nil).QueryContext), varargs...)
}

// QueryRow mocks base method.
func (m *MockTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockTransactionMockRecorder) QueryRow(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockTransaction)(nil).QueryRow), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockTransactionMockRecorder) QueryRowContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockTransaction)(nil).QueryRowContext), varargs...)
}

// Rollback mocks base method.
func (m *MockTransaction) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTransactionMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTransaction)(nil).Rollback))
}

// MockPrepareable is a mock of Prepareable interface.
type MockPrepareable struct {
	ctrl     *gomock.Controller
	recorder *MockPrepareableMockRecorder
}

// MockPrepareableMockRecorder is the mock recorder for MockPrepareable.
type MockPrepareableMockRecorder struct {
	mock *MockPrepareable
}

// NewMockPrepareable creates a new mock instance.
func NewMockPrepareable(ctrl *gomock.Controller) *MockPrepareable {
	mock := &MockPrepareable{ctrl: ctrl}
	mock.recorder = &MockPrepareableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrepareable) EXPECT() *MockPrepareableMockRecorder {
	return m.recorder
}

// Prepare mocks base method.
func (m *MockPrepareable) Prepare(query string) (*sql.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(*sql.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockPrepareableMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockPrepareable)(nil).Prepare), query)
}

// PrepareContext mocks base method.
func (m *MockPrepareable) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContext", ctx, query)
	ret0, _ := ret[0].(*sql.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContext indicates an expected call of PrepareContext.
func (mr *MockPrepareableMockRecorder) PrepareContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContext", reflect.TypeOf((*MockPrepareable)(nil).PrepareContext), ctx, query)
}

// MockQueryable is a mock of Queryable interface.
type MockQueryable struct {
	ctrl     *gomock.Controller
	recorder *MockQueryableMockRecorder
}

// MockQueryableMockRecorder is the mock recorder for MockQueryable.
type MockQueryableMockRecorder struct {
	mock *MockQueryable
}

// NewMockQueryable creates a new mock instance.
func NewMockQueryable(ctrl *gomock.Controller) *MockQueryable {
	mock := &MockQueryable{ctrl: ctrl}
	mock.recorder = &MockQueryableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryable) EXPECT() *MockQueryableMockRecorder {
	return m.recorder
}

// Query mocks base method.
func (m *MockQueryable) Query(query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockQueryableMockRecorder) Query(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockQueryable)(nil).Query), varargs...)
}

// QueryContext mocks base method.
func (m *MockQueryable) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockQueryableMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockQueryable)(nil).QueryContext), varargs...)
}

// QueryRow mocks base method.
func (m *MockQueryable) QueryRow(query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockQueryableMockRecorder) QueryRow(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockQueryable)(nil).QueryRow), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockQueryable) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockQueryableMockRecorder) QueryRowContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockQueryable)(nil).QueryRowContext), varargs...)
}

// MockExecutable is a mock of Executable interface.
type MockExecutable struct {
	ctrl     *gomock.Controller
	recorder *MockExecutableMockRecorder
}

// MockExecutableMockRecorder is the mock recorder for MockExecutable.
type MockExecutableMockRecorder struct {
	mock *MockExecutable
}

// NewMockExecutable creates a new mock instance.
func NewMockExecutable(ctrl *gomock.Controller) *MockExecutable {
	mock := &MockExecutable{ctrl: ctrl}
	mock.recorder = &MockExecutableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecutable) EXPECT() *MockExecutableMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockExecutable) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockExecutableMockRecorder) Exec(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockExecutable)(nil).Exec), varargs...)
}

// ExecContext mocks base method.
func (m *MockExecutable) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockExecutableMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockExecutable)(nil).ExecContext), varargs...)
}
