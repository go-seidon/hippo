package repository_mongo

import (
	"context"
	"fmt"

	"github.com/go-seidon/local/internal/datetime"
	db_mongo "github.com/go-seidon/local/internal/db-mongo"
	"github.com/go-seidon/local/internal/repository"
)

type RepositoryParam struct {
	dbClient db_mongo.Client
	dbConfig *DbConfig
	clock    datetime.Clock
}

type DbConfig struct {
	DbName string
}

type RepoOption = func(*RepositoryParam)

func WithDbClient(dbClient db_mongo.Client) RepoOption {
	return func(ro *RepositoryParam) {
		ro.dbClient = dbClient
	}
}

func WithDbConfig(dbConfig *DbConfig) RepoOption {
	return func(ro *RepositoryParam) {
		ro.dbConfig = dbConfig
	}
}

func WithClock(clock datetime.Clock) RepoOption {
	return func(ro *RepositoryParam) {
		ro.clock = clock
	}
}

type provider struct {
	dbClient db_mongo.Client
	authRepo *authRepository
	fileRepo *fileRepository
}

func (p *provider) Init(ctx context.Context) error {
	err := p.dbClient.Connect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *provider) GetAuthRepo() repository.AuthRepository {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.FileRepository {
	return p.fileRepo
}

func NewRepository(opts ...RepoOption) (*provider, error) {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.dbClient == nil {
		return nil, fmt.Errorf("invalid db client specified")
	}
	if p.dbConfig == nil {
		return nil, fmt.Errorf("invalid db config specified")
	}

	var clock datetime.Clock
	if p.clock == nil {
		clock = datetime.NewClock()
	} else {
		clock = p.clock
	}

	authRepo := &authRepository{
		dbConfig: p.dbConfig,
		dbClient: p.dbClient,
		clock:    clock,
	}
	fileRepo := &fileRepository{
		dbConfig: p.dbConfig,
		dbClient: p.dbClient,
		clock:    clock,
	}

	repo := &provider{
		dbClient: p.dbClient,
		authRepo: authRepo,
		fileRepo: fileRepo,
	}
	return repo, nil
}
