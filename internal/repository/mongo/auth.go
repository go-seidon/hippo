package mongo

import (
	"context"

	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/datetime"
	db_mongo "github.com/go-seidon/provider/mongo"
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
		return nil, repository.ErrNotFound
	}
	return nil, err
}

func NewAuth(opts ...RepoOption) *authRepository {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	clock := p.clock
	if clock == nil {
		clock = datetime.NewClock()
	}

	return &authRepository{
		dbClient: p.dbClient,
		dbConfig: p.dbConfig,
		clock:    clock,
	}
}
