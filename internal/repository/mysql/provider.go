package mysql

import (
	"context"

	"github.com/go-seidon/hippo/internal/repository"
	db_mysql "github.com/go-seidon/provider/mysql"
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

func (p *provider) GetAuthRepo() repository.Auth {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.File {
	return p.fileRepo
}
