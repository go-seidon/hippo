package mongo

import (
	"context"
	"time"

	"github.com/go-seidon/hippo/internal/repository"
	db_mongo "github.com/go-seidon/provider/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type file struct {
	dbConfig *DbConfig
	dbClient db_mongo.Client
}

func (r *file) CreateFile(ctx context.Context, p repository.CreateFileParam) (*repository.CreateFileResult, error) {
	err := p.CreateFn(ctx, repository.CreateFnParam{
		FilePath: p.Path,
	})
	if err != nil {
		return nil, err
	}

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

	res := &repository.CreateFileResult{
		UniqueId:  p.UniqueId,
		Name:      p.Name,
		Path:      p.Path,
		Mimetype:  p.Mimetype,
		Extension: p.Extension,
		Size:      p.Size,
		CreatedAt: p.CreatedAt,
	}
	return res, nil
}

func (r *file) RetrieveFile(ctx context.Context, p repository.RetrieveFileParam) (*repository.RetrieveFileResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("file")
	findFilter := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
	}
	file := struct {
		Id        string     `bson:"_id"`
		Name      string     `bson:"name"`
		Path      string     `bson:"path"`
		Mimetype  string     `bson:"mimetype"`
		Extension string     `bson:"extension"`
		Size      int64      `bson:"size"`
		CreatedAt time.Time  `bson:"created_at"`
		DeletedAt *time.Time `bson:"deleted_at"`
	}{}
	err := cl.FindOne(ctx, findFilter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	res := &repository.RetrieveFileResult{
		UniqueId:  file.Id,
		Name:      file.Name,
		Path:      file.Path,
		Mimetype:  file.Mimetype,
		Extension: file.Extension,
		Size:      file.Size,
		CreatedAt: file.CreatedAt,
		DeletedAt: file.DeletedAt,
	}
	return res, nil
}

func (r *file) DeleteFile(ctx context.Context, p repository.DeleteFileParam) (*repository.DeleteFileResult, error) {
	cl := r.dbClient.Database(r.dbConfig.DbName).Collection("file")
	findFilter := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
	}
	file := struct {
		Id        string     `bson:"_id"`
		Name      string     `bson:"name"`
		Path      string     `bson:"path"`
		DeletedAt *time.Time `bson:"deleted_at"`
	}{}
	err := cl.FindOne(ctx, findFilter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	if file.DeletedAt != nil {
		return nil, repository.ErrDeleted
	}

	err = p.DeleteFn(ctx, repository.DeleteFnParam{
		FilePath: file.Path,
	})
	if err != nil {
		return nil, err
	}

	updateFilter := bson.D{
		{
			Key:   "_id",
			Value: p.UniqueId,
		},
	}
	data := bson.M{
		"$set": bson.M{
			"deleted_at": p.DeletedAt,
		},
	}
	_, err = cl.UpdateOne(ctx, updateFilter, data)
	if err != nil {
		return nil, err
	}

	res := &repository.DeleteFileResult{
		DeletedAt: p.DeletedAt,
	}
	return res, nil
}

func NewFile(opts ...RepoOption) *file {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	return &file{
		dbClient: p.dbClient,
		dbConfig: p.dbConfig,
	}
}
