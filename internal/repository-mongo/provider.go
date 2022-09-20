package repository_mongo

import (
	"context"

	db_mongo "github.com/go-seidon/local/internal/db-mongo"
	"github.com/go-seidon/local/internal/repository"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

func (p *provider) Ping(ctx context.Context) error {
	return p.dbClient.Ping(ctx, readpref.Secondary())
}

func (p *provider) GetAuthRepo() repository.AuthRepository {
	return p.authRepo
}

func (p *provider) GetFileRepo() repository.FileRepository {
	return p.fileRepo
}
