package db_mysql

import (
	"database/sql"
	"fmt"

	"context"

	_ "github.com/go-sql-driver/mysql"
)

type Client interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Query
}

type Transaction interface {
	Commit() error
	Rollback() error
	Query
}

type Query interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type ClientOption = func(*ClientParam)

type ClientParam struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DbName          string
	ShouldParseTime bool
}

type ClientConfig struct {
	DbName string
}

func ParseTime() ClientOption {
	return func(cp *ClientParam) {
		cp.ShouldParseTime = true
	}
}

func WithLocation(host string, port int) ClientOption {
	return func(cp *ClientParam) {
		cp.Host = host
		cp.Port = port
	}
}

func WithAuth(username, password string) ClientOption {
	return func(cp *ClientParam) {
		cp.Username = username
		cp.Password = password
	}
}

func WithConfig(cfg ClientConfig) ClientOption {
	return func(cp *ClientParam) {
		cp.DbName = cfg.DbName
	}
}

func NewClient(opts ...ClientOption) (*sql.DB, error) {
	p := ClientParam{}
	for _, opt := range opts {
		opt(&p)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		p.Username, p.Password,
		p.Host, p.Port, p.DbName,
	)

	if p.ShouldParseTime {
		dsn = fmt.Sprintf("%s?parseTime=true", dsn)
	}
	return sql.Open("mysql", dsn)
}
