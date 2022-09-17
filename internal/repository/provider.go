package repository

import "context"

const (
	DB_PROVIDER_MYSQL = "mysql"
	DB_PROVIDER_MONGO = "mongo"
)

type Provider interface {
	Init(ctx context.Context) error
	GetAuthRepo() AuthRepository
	GetFileRepo() FileRepository
}
