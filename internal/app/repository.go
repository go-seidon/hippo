package app

import (
	"fmt"

	"github.com/go-seidon/hippo/internal/repository"
	repository_mongo "github.com/go-seidon/hippo/internal/repository/mongo"
	repository_mysql "github.com/go-seidon/hippo/internal/repository/mysql"
	db_mongo "github.com/go-seidon/provider/mongo"
	db_mysql "github.com/go-seidon/provider/mysql"
)

func NewDefaultRepository(config *Config) (repository.Provider, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid config")
	}

	if config.RepositoryProvider != repository.PROVIDER_MYSQL &&
		config.RepositoryProvider != repository.PROVIDER_MONGO {
		return nil, fmt.Errorf("invalid repository provider")
	}

	var repo repository.Provider
	if config.RepositoryProvider == repository.PROVIDER_MYSQL {
		dbMaster, err := db_mysql.NewClient(
			db_mysql.WithAuth(config.MySQLMasterUser, config.MySQLMasterPassword),
			db_mysql.WithConfig(db_mysql.ClientConfig{DbName: config.MySQLMasterDBName}),
			db_mysql.WithLocation(config.MySQLMasterHost, config.MySQLMasterPort),
			db_mysql.ParseTime(),
		)
		if err != nil {
			return nil, err
		}

		dbReplica, err := db_mysql.NewClient(
			db_mysql.WithAuth(config.MySQLReplicaUser, config.MySQLReplicaPassword),
			db_mysql.WithConfig(db_mysql.ClientConfig{DbName: config.MySQLReplicaDBName}),
			db_mysql.WithLocation(config.MySQLReplicaHost, config.MySQLReplicaPort),
			db_mysql.ParseTime(),
		)
		if err != nil {
			return nil, err
		}

		repo, err = repository_mysql.NewRepository(
			repository_mysql.WithDbMaster(dbMaster),
			repository_mysql.WithDbReplica(dbReplica),
		)
		if err != nil {
			return nil, err
		}
	} else if config.RepositoryProvider == repository.PROVIDER_MONGO {
		opts := []db_mongo.ClientOption{}
		if config.MongoAuthMode == db_mongo.AUTH_BASIC {
			opts = append(opts, db_mongo.WithBasicAuth(
				config.MongoAuthUser,
				config.MongoAuthPassword,
				config.MongoAuthSource,
			))
		}

		if config.MongoMode == db_mongo.MODE_STANDALONE {
			opts = append(opts, db_mongo.UsingStandalone(
				config.MongoStandaloneHost, config.MongoStandalonePort,
			))
		} else if config.MongoMode == db_mongo.MODE_REPLICATION {
			opts = append(opts, db_mongo.UsingReplication(
				config.MongoReplicaName,
				config.MongoReplicaHosts,
			))
		}

		opts = append(opts, db_mongo.WithConfig(db_mongo.ClientConfig{
			DbName: config.MongoDBName,
		}))

		dbClient, err := db_mongo.NewClient(opts...)
		if err != nil {
			return nil, err
		}

		repo, err = repository_mongo.NewRepository(
			repository_mongo.WithDbClient(dbClient),
			repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
				DbName: config.MongoDBName,
			}),
		)
		if err != nil {
			return nil, err
		}
	}
	return repo, nil
}
