package repository_mongo

import (
	"context"
	"fmt"

	"github.com/go-seidon/local/internal/datetime"
	"github.com/go-seidon/local/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileRepository struct {
	dbConfig *DbConfig
	dbClient *mongo.Client
	clock    datetime.Clock
}

func (r *fileRepository) DeleteFile(ctx context.Context, p repository.DeleteFileParam) (*repository.DeleteFileResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("file")
	findFilter := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
	}
	file := struct {
		Id   string `bson:"_id"`
		Name string `bson:"name"`
		Path string `bson:"path"`
	}{}
	err := cl.FindOne(ctx, findFilter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrorRecordNotFound
		}
		return nil, err
	}

	err = p.DeleteFn(ctx, repository.DeleteFnParam{
		FilePath: file.Path,
	})
	if err != nil {
		return nil, err
	}

	currentTimestamp := r.clock.Now()
	updateFilter := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
	}
	data := bson.M{
		"$set": bson.M{
			"deleted_at": currentTimestamp,
		},
	}
	deleteRes, err := cl.UpdateOne(ctx, updateFilter, data)
	if err != nil {
		return nil, err
	}

	if deleteRes.ModifiedCount != 1 {
		return nil, fmt.Errorf("record is not updated")
	}

	res := &repository.DeleteFileResult{
		DeletedAt: currentTimestamp,
	}
	return res, nil
}

func (r *fileRepository) RetrieveFile(ctx context.Context, p repository.RetrieveFileParam) (*repository.RetrieveFileResult, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (r *fileRepository) CreateFile(ctx context.Context, p repository.CreateFileParam) (*repository.CreateFileResult, error) {
	return nil, fmt.Errorf("unimplemented")
}

func NewFileRepository(opts ...RepoOption) (*fileRepository, error) {
	option := RepositoryOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.dbClient == nil {
		return nil, fmt.Errorf("invalid db client specified")
	}
	if option.dbConfig == nil {
		return nil, fmt.Errorf("invalid db config specified")
	}

	var clock datetime.Clock
	if option.clock == nil {
		clock = datetime.NewClock()
	} else {
		clock = option.clock
	}

	r := &fileRepository{
		dbClient: option.dbClient,
		dbConfig: option.dbConfig,
		clock:    clock,
	}
	return r, nil
}
