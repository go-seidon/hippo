package repository_mongo_test

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func OpenDb(dsn string) (*mongo.Client, error) {
	if dsn == "" {
		dsn = "mongodb://admin:123456@localhost:27018/goseidon_local_test"
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
		"file://../../migration/mongo",
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
	}

	_, err := cl.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
