package repository_mongo

import (
	"context"
	"errors"

	db_mongo "github.com/go-seidon/local/internal/db-mongo"
	"github.com/go-seidon/local/internal/repository"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

type provider struct {
	dbClient db_mongo.Client
	authRepo *authRepository
	fileRepo *fileRepository
}

func (p *provider) Init(ctx context.Context) error {
	err := p.dbClient.Connect(ctx)
	if err == nil {
		return nil
	}
	if errors.Is(err, topology.ErrTopologyConnected) {
		return nil
	}
	return err
}

func (p *provider) Ping(ctx context.Context) error {
	return p.dbClient.Ping(ctx, readpref.Secondary())
}

func (p *provider) GetAuthRepo() repository.AuthRepository {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.FileRepository {
	return p.fileRepo
}
