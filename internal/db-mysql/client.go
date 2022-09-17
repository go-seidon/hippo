package db_mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ClientOption = func(*ClientParam)

type ClientParam struct {
	Host      string
	Port      int
	Username  string
	Password  string
	DbName    string
	ParseTime bool
}

type ClientConfig struct {
	DbName string
}

func ParseTime() ClientOption {
	return func(cp *ClientParam) {
		cp.ParseTime = true
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
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		p.Username, p.Password,
		p.Host, p.Port, p.DbName,
	)
	return sql.Open("mysql", dsn)
}
