package repository_mongo

import (
	"context"
	"fmt"

	"github.com/go-seidon/local/internal/datetime"
	db_mongo "github.com/go-seidon/local/internal/db-mongo"
	"github.com/go-seidon/local/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type authRepository struct {
	dbConfig *DbConfig
	dbClient db_mongo.Client
	clock    datetime.Clock
}

func (r *authRepository) FindClient(ctx context.Context, p repository.FindClientParam) (*repository.FindClientResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("auth_client")
	filter := bson.D{
		{
			Key:   "client_id",
			Value: p.ClientId,
		},
	}
	projection := options.FindOne().SetProjection(bson.D{
		{
			Key:   "client_id",
			Value: 1,
		},
		{
			Key:   "client_secret",
			Value: 1,
		},
	})

	authClient := struct {
		ClientId     string `bson:"client_id"`
		ClientSecret string `bson:"client_secret"`
	}{}
	err := cl.FindOne(ctx, filter, projection).Decode(&authClient)
	if err == nil {
		res := repository.FindClientResult{
			ClientId:     authClient.ClientId,
			ClientSecret: authClient.ClientSecret,
		}
		return &res, nil
	}

	if err == mongo.ErrNoDocuments {
		return nil, repository.ErrorRecordNotFound
	}
	return nil, err
}

func NewAuthRepository(opts ...RepoOption) (*authRepository, error) {
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

	r := &authRepository{
		dbClient: p.dbClient,
		dbConfig: p.dbConfig,
		clock:    clock,
	}
	return r, nil
}
