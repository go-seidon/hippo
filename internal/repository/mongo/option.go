package mongo

import (
	db_mongo "github.com/go-seidon/provider/mongo"
)

type RepositoryParam struct {
	dbClient db_mongo.Client
	dbConfig *DbConfig
}

type DbConfig struct {
	DbName string
}

type RepoOption = func(*RepositoryParam)

func WithDbClient(dbClient db_mongo.Client) RepoOption {
	return func(ro *RepositoryParam) {
		ro.dbClient = dbClient
	}
}

func WithDbConfig(dbConfig *DbConfig) RepoOption {
	return func(ro *RepositoryParam) {
		ro.dbConfig = dbConfig
	}
}
