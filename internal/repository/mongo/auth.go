package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-seidon/hippo/internal/repository"
	db_mongo "github.com/go-seidon/provider/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type auth struct {
	dbConfig *DbConfig
	dbClient db_mongo.Client
}

func (r *auth) FindClient(ctx context.Context, p repository.FindClientParam) (*repository.FindClientResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("auth_client")

	var filter bson.D
	if p.ClientId != "" {
		filter = bson.D{
			{
				Key:   "client_id",
				Value: p.ClientId,
			},
		}
	} else {
		filter = bson.D{
			{
				Key:   "_id",
				Value: p.Id,
			},
		}
	}

	projection := options.FindOne().SetProjection(bson.D{
		{
			Key:   "_id",
			Value: 1,
		},
		{
			Key:   "name",
			Value: 1,
		},
		{
			Key:   "type",
			Value: 1,
		},
		{
			Key:   "status",
			Value: 1,
		},
		{
			Key:   "client_id",
			Value: 1,
		},
		{
			Key:   "client_secret",
			Value: 1,
		},
		{
			Key:   "created_at",
			Value: 1,
		},
		{
			Key:   "updated_at",
			Value: 1,
		},
	})

	client := struct {
		Id           string     `bson:"_id"`
		Name         string     `bson:"name"`
		Type         string     `bson:"type"`
		Status       string     `bson:"status"`
		ClientId     string     `bson:"client_id"`
		ClientSecret string     `bson:"client_secret"`
		CreatedAt    time.Time  `bson:"created_at"`
		UpdatedAt    *time.Time `bson:"updated_at"`
	}{}
	err := cl.FindOne(ctx, filter, projection).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	res := &repository.FindClientResult{
		Id:           client.Id,
		Name:         client.Name,
		Type:         client.Type,
		Status:       client.Status,
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}
	return res, nil
}

// @todo: add implementation
func (r *auth) CreateClient(ctx context.Context, p repository.CreateClientParam) (*repository.CreateClientResult, error) {
	return nil, fmt.Errorf("unimplemented")
}

// @todo: add implementation
func (r *auth) UpdateClient(ctx context.Context, p repository.UpdateClientParam) (*repository.UpdateClientResult, error) {
	return nil, fmt.Errorf("unimplemented")
}

// @todo: add implementation
func (r *auth) SearchClient(ctx context.Context, p repository.SearchClientParam) (*repository.SearchClientResult, error) {
	return nil, fmt.Errorf("unimplemented")
}

func NewAuth(opts ...RepoOption) *auth {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	return &auth{
		dbClient: p.dbClient,
		dbConfig: p.dbConfig,
	}
}
