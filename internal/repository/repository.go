package repository

import "context"

const (
	PROVIDER_MYSQL = "mysql"
	PROVIDER_MONGO = "mongo"
)

type Repository interface {
	Init(ctx context.Context) error
	Ping(ctx context.Context) error
	GetAuth() Auth
	GetFile() File
}
