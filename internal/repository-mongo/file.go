package repository_mongo

import (
	"context"
	"fmt"

	db_mongo "github.com/go-seidon/hippo/internal/db-mongo"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/datetime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileRepository struct {
	dbConfig *DbConfig
	dbClient db_mongo.Client
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
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("file")
	findFilter := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
	}
	file := struct {
		Id        string `bson:"_id"`
		Name      string `bson:"name"`
		Path      string `bson:"path"`
		Mimetype  string `bson:"mimetype"`
		Extension string `bson:"extension"`
		Size      int64  `bson:"size"`
	}{}
	err := cl.FindOne(ctx, findFilter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrorRecordNotFound
		}
		return nil, err
	}

	res := &repository.RetrieveFileResult{
		UniqueId:  file.Id,
		Name:      file.Name,
		Path:      file.Path,
		MimeType:  file.Mimetype,
		Extension: file.Extension,
		Size:      file.Size,
	}
	return res, nil
}

func (r *fileRepository) CreateFile(ctx context.Context, p repository.CreateFileParam) (*repository.CreateFileResult, error) {
	err := p.CreateFn(ctx, repository.CreateFnParam{
		FilePath: p.Path,
	})
	if err != nil {
		return nil, err
	}

	currentTimestamp := r.clock.Now()
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("file")
	data := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
		{
			Key:   "name",
			Value: p.Name,
		},
		{
			Key:   "path",
			Value: p.Path,
		},
		{
			Key:   "mimetype",
			Value: p.Mimetype,
		},
		{
			Key:   "extension",
			Value: p.Extension,
		},
		{
			Key:   "size",
			Value: p.Size,
		},
		{
			Key:   "created_at",
			Value: currentTimestamp,
		},
		{
			Key:   "updated_at",
			Value: currentTimestamp,
		},
	}
	_, err = cl.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	res := &repository.CreateFileResult{
		UniqueId:  p.UniqueId,
		Name:      p.Name,
		Path:      p.Path,
		Mimetype:  p.Mimetype,
		Extension: p.Extension,
		Size:      p.Size,
		CreatedAt: currentTimestamp,
	}
	return res, nil
}

func NewFileRepository(opts ...RepoOption) (*fileRepository, error) {
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

	r := &fileRepository{
		dbClient: p.dbClient,
		dbConfig: p.dbConfig,
		clock:    clock,
	}
	return r, nil
}
