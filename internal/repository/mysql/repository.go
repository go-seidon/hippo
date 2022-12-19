package mysql

import (
	"context"
	"fmt"

	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/mysql"
)

type mysqlRepository struct {
	dbClient mysql.Pingable
	authRepo *auth
	fileRepo *file
}

func (p *mysqlRepository) Init(ctx context.Context) error {
	return nil
}

func (p *mysqlRepository) Ping(ctx context.Context) error {
	return p.dbClient.PingContext(ctx)
}

func (p *mysqlRepository) GetAuth() repository.Auth {
	return p.authRepo
}

func (p *mysqlRepository) GetFile() repository.File {
	return p.fileRepo
}

func NewRepository(opts ...RepoOption) (*mysqlRepository, error) {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.dbClient == nil && p.gormClient == nil {
		return nil, fmt.Errorf("invalid db client")
	}

	var err error
	dbClient := p.dbClient
	if dbClient == nil {
		dbClient, err = p.gormClient.DB()
		if err != nil {
			return nil, err
		}
	}

	authRepo := &auth{
		gormClient: p.gormClient,
	}
	fileRepo := &file{
		gormClient: p.gormClient,
	}

	repo := &mysqlRepository{
		dbClient: dbClient,
		authRepo: authRepo,
		fileRepo: fileRepo,
	}
	return repo, nil
}
