package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-seidon/hippo/internal/repository"
	db_mongo "github.com/go-seidon/provider/mongo"
	"github.com/go-seidon/provider/typeconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type auth struct {
	dbConfig *DbConfig
	dbClient db_mongo.Client
}

// @note: return `ErrExists` if client_id is already created
func (r *auth) CreateClient(ctx context.Context, p repository.CreateClientParam) (*repository.CreateClientResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("auth_client")

	currentClient := struct {
		Id       string `bson:"_id"`
		ClientId string `bson:"client_id"`
	}{}
	err := cl.FindOne(ctx, bson.D{
		{
			Key:   "client_id",
			Value: p.ClientId,
		},
	}).Decode(&currentClient)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		if err == nil {
			return nil, repository.ErrExists
		}
		return nil, err
	}

	data := bson.D{
		{
			Key:   "_id",
			Value: p.Id,
		},
		{
			Key:   "client_id",
			Value: p.ClientId,
		},
		{
			Key:   "client_secret",
			Value: p.ClientSecret,
		},
		{
			Key:   "name",
			Value: p.Name,
		},
		{
			Key:   "type",
			Value: p.Type,
		},
		{
			Key:   "status",
			Value: p.Status,
		},
		{
			Key:   "created_at",
			Value: p.CreatedAt,
		},
		{
			Key:   "updated_at",
			Value: p.CreatedAt,
		},
	}
	_, err = cl.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	client := struct {
		Id           string    `bson:"_id"`
		Name         string    `bson:"name"`
		Type         string    `bson:"type"`
		Status       string    `bson:"status"`
		ClientId     string    `bson:"client_id"`
		ClientSecret string    `bson:"client_secret"`
		CreatedAt    time.Time `bson:"created_at"`
	}{}
	err = cl.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: p.Id,
		},
	}).Decode(&client)
	if err != nil {
		return nil, err
	}

	res := &repository.CreateClientResult{
		Id:           client.Id,
		Name:         client.Name,
		Type:         client.Type,
		Status:       client.Status,
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
		CreatedAt:    client.CreatedAt.UTC(),
	}
	return res, nil
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

func (r *auth) UpdateClient(ctx context.Context, p repository.UpdateClientParam) (*repository.UpdateClientResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("auth_client")

	currentClient := struct {
		Id       string `bson:"_id"`
		ClientId string `bson:"client_id"`
	}{}
	err := cl.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: p.Id,
		},
	}).Decode(&currentClient)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	updateFilter := bson.D{
		{
			Key:   "_id",
			Value: p.Id,
		},
	}
	data := bson.M{
		"$set": bson.M{
			"name":       p.Name,
			"type":       p.Type,
			"status":     p.Status,
			"updated_at": p.UpdatedAt,
		},
	}
	_, err = cl.UpdateOne(ctx, updateFilter, data)
	if err != nil {
		return nil, err
	}

	client := struct {
		Id           string    `bson:"_id"`
		Name         string    `bson:"name"`
		Type         string    `bson:"type"`
		Status       string    `bson:"status"`
		ClientId     string    `bson:"client_id"`
		ClientSecret string    `bson:"client_secret"`
		CreatedAt    time.Time `bson:"created_at"`
		UpdatedAt    time.Time `bson:"updated_at"`
	}{}
	err = cl.FindOne(ctx, bson.D{
		{
			Key:   "_id",
			Value: p.Id,
		},
	}).Decode(&client)
	if err != nil {
		return nil, err
	}

	res := &repository.UpdateClientResult{
		Id:           client.Id,
		Name:         client.Name,
		Type:         client.Type,
		Status:       client.Status,
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
		CreatedAt:    client.CreatedAt.UTC(),
		UpdatedAt:    client.UpdatedAt.UTC(),
	}
	return res, nil
}

func (r *auth) SearchClient(ctx context.Context, p repository.SearchClientParam) (*repository.SearchClientResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("auth_client")

	filter := bson.D{}

	if len(p.Statuses) > 0 {
		filter = append(filter, primitive.E{
			Key: "status",
			Value: bson.D{
				{
					Key:   "$in",
					Value: p.Statuses,
				},
			},
		})
	}

	if p.Keyword != "" {
		filter = append(filter, primitive.E{
			Key: "$or",
			Value: bson.A{
				bson.D{
					{
						Key: "name",
						Value: bson.D{
							{
								Key:   "$regex",
								Value: fmt.Sprintf(".*%s.*", p.Keyword),
							},
						},
					},
				},
				bson.D{
					{
						Key: "client_id",
						Value: bson.D{
							{
								Key:   "$regex",
								Value: fmt.Sprintf(".*%s.*", p.Keyword),
							},
						},
					},
				},
			},
		})
	}

	options := options.Find()
	if p.Limit > 0 {
		options.SetLimit(int64(p.Limit))
	}
	if p.Offset > 0 {
		options.SetSkip(p.Offset)
	}

	findRes, err := cl.Find(ctx, filter, options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	total, err := cl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	clients := []struct {
		Id           string     `bson:"_id"`
		Name         string     `bson:"name"`
		Type         string     `bson:"type"`
		Status       string     `bson:"status"`
		ClientId     string     `bson:"client_id"`
		ClientSecret string     `bson:"client_secret"`
		CreatedAt    time.Time  `bson:"created_at"`
		UpdatedAt    *time.Time `bson:"updated_at"`
	}{}
	err = findRes.All(ctx, &clients)
	if err != nil {
		return nil, err
	}

	items := []repository.SearchClientItem{}
	for _, client := range clients {
		var updatedAt *time.Time
		if client.UpdatedAt != nil {
			updatedAt = typeconv.Time(client.UpdatedAt.UTC())
		}

		items = append(items, repository.SearchClientItem{
			Id:           client.Id,
			ClientId:     client.ClientId,
			ClientSecret: client.ClientSecret,
			Name:         client.Name,
			Type:         client.Type,
			Status:       client.Status,
			CreatedAt:    client.CreatedAt.UTC(),
			UpdatedAt:    updatedAt,
		})
	}

	res := &repository.SearchClientResult{
		Summary: repository.SearchClientSummary{
			TotalItems: total,
		},
		Items: items,
	}
	return res, nil
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
