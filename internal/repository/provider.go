package repository

import "context"

const (
	PROVIDER_MYSQL = "mysql"
	PROVIDER_MONGO = "mongo"
)

type Provider interface {
	Init(ctx context.Context) error
	Ping(ctx context.Context) error
	GetAuthRepo() Auth
	GetFileRepo() File
}
