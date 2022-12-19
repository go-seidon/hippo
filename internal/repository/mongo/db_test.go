package mongo_test

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func OpenDb(dsn string) (*mongo.Client, error) {
	if dsn == "" {
		dsn = "mongodb://admin:123456@localhost:27020/hippo_test"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}

type RunDbMigrationParam struct {
	DbName string
}

func RunDbMigration(dbClient *mongo.Client, p RunDbMigrationParam) error {
	driver, err := mongodb.WithInstance(dbClient, &mongodb.Config{
		DatabaseName: p.DbName,
	})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://../../../migration/mongo",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err == nil {
		return nil
	}

	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

type InsertAuthClientParam struct {
	Id           string
	Name         string
	ClientId     string
	ClientSecret string
	Type         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DbName       string
}

func InsertAuthClient(dbClient *mongo.Client, p InsertAuthClientParam) error {
	cl := dbClient.Database(p.DbName).Collection("auth_client")
	ctx := context.Background()
	data := bson.D{
		{
			Key:   "_id",
			Value: p.Id,
		},
		{
			Key:   "name",
			Value: p.Name,
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
			Key:   "type",
			Value: p.Type,
		},
		{
			Key:   "status",
			Value: p.Status,
		},
	}

	if !p.CreatedAt.IsZero() {
		data = append(data, primitive.E{
			Key:   "created_at",
			Value: p.CreatedAt,
		})
	}

	if !p.UpdatedAt.IsZero() {
		data = append(data, primitive.E{
			Key:   "updated_at",
			Value: p.UpdatedAt,
		})
	}

	_, err := cl.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

type InsertFileParam struct {
	Id        string
	Name      string
	Path      string
	Mimetype  string
	Extension string
	Size      int64
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
	DbName    string
}

func InsertFile(dbClient *mongo.Client, p InsertFileParam) error {
	cl := dbClient.Database(p.DbName).Collection("file")
	ctx := context.Background()
	createdAt := time.UnixMilli(p.CreatedAt).UTC()
	updatedAt := time.UnixMilli(p.UpdatedAt).UTC()
	deletedAt := time.UnixMilli(p.DeletedAt).UTC()

	data := bson.D{
		{
			Key:   "_id",
			Value: p.Id,
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
			Value: createdAt,
		},
		{
			Key:   "updated_at",
			Value: updatedAt,
		},
	}

	if p.DeletedAt != 0 {
		data = append(data, primitive.E{
			Key:   "deleted_at",
			Value: deletedAt,
		})
	}

	_, err := cl.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
