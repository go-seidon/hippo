package repository_mongo

import (
	"fmt"

	db_mongo "github.com/go-seidon/hippo/internal/db-mongo"
	"github.com/go-seidon/provider/datetime"
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

	clock := p.clock
	if clock == nil {
		clock = datetime.NewClock()
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
