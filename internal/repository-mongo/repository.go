package repository_mongo

import (
	"github.com/go-seidon/local/internal/datetime"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryOption struct {
	dbClient *mongo.Client
	dbConfig *DbConfig
	clock    datetime.Clock
}

type DbConfig struct {
	DbName string
}

type RepoOption = func(*RepositoryOption)

func WithDbClient(dbClient *mongo.Client) RepoOption {
	return func(ro *RepositoryOption) {
		ro.dbClient = dbClient
	}
}

func WithDbConfig(dbConfig *DbConfig) RepoOption {
	return func(ro *RepositoryOption) {
		ro.dbConfig = dbConfig
	}
}

func WithClock(clock datetime.Clock) RepoOption {
	return func(ro *RepositoryOption) {
		ro.clock = clock
	}
}
