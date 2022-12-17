package mongo

import (
	"context"
	"errors"

	"github.com/go-seidon/hippo/internal/repository"
	db_mongo "github.com/go-seidon/provider/mongo"
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

func (p *provider) GetAuthRepo() repository.Auth {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.File {
	return p.fileRepo
}
