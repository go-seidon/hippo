package repository_mysql

import (
	"context"

	db_mysql "github.com/go-seidon/hippo/internal/db-mysql"
	"github.com/go-seidon/hippo/internal/repository"
)

type provider struct {
	mClient  db_mysql.Client
	rClient  db_mysql.Client
	authRepo *authRepository
	fileRepo *fileRepository
}

func (p *provider) Init(ctx context.Context) error {
	return nil
}

func (p *provider) Ping(ctx context.Context) error {
	return p.rClient.PingContext(ctx)
}

func (p *provider) GetAuthRepo() repository.AuthRepository {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.FileRepository {
	return p.fileRepo
}
