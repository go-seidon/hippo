package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-seidon/hippo/internal/repository"
	db_mongo "github.com/go-seidon/provider/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

type mongoRepository struct {
	dbClient db_mongo.Client
	authRepo *auth
	fileRepo *file
}

func (p *mongoRepository) Init(ctx context.Context) error {
	err := p.dbClient.Connect(ctx)
	if err == nil {
		return nil
	}
	if errors.Is(err, topology.ErrTopologyConnected) {
		return nil
	}
	return err
}

func (p *mongoRepository) Ping(ctx context.Context) error {
	return p.dbClient.Ping(ctx, readpref.Secondary())
}

func (p *mongoRepository) GetAuth() repository.Auth {
	return p.authRepo
}

func (p *mongoRepository) GetFile() repository.File {
	return p.fileRepo
}

func NewRepository(opts ...RepoOption) (*mongoRepository, error) {
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

	authRepo := &auth{
		dbConfig: p.dbConfig,
		dbClient: p.dbClient,
	}
	fileRepo := &file{
		dbConfig: p.dbConfig,
		dbClient: p.dbClient,
	}

	repo := &mongoRepository{
		dbClient: p.dbClient,
		authRepo: authRepo,
		fileRepo: fileRepo,
	}
	return repo, nil
}
