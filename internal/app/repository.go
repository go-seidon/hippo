package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-seidon/local/internal/repository"
	repository_mongo "github.com/go-seidon/local/internal/repository-mongo"
	repository_mysql "github.com/go-seidon/local/internal/repository-mysql"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_PROVIDER_MYSQL = "mysql"
	DB_PROVIDER_MONGO = "mongo"
)

func NewRepository(o ...RepositoryOption) (*NewRepositoryResult, error) {
	if len(o) == 0 {
		return nil, fmt.Errorf("invalid repository option")
	}

	var p NewRepositoryOption
	for _, opt := range o {
		opt.Apply(&p)
	}

	if p.Provider == DB_PROVIDER_MYSQL {
		return newMySQLRepository(p)
	} else if p.Provider == DB_PROVIDER_MONGO {
		return newMongoRepository(p)
	}

	return nil, fmt.Errorf("db provider is not supported")
}

func newMySQLRepository(p NewRepositoryOption) (*NewRepositoryResult, error) {
	masterDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		p.MySQLMasterUser, p.MySQLMasterPassword,
		p.MySQLMasterHost, p.MySQLMasterPort, p.MySQLMasterDBName,
	)
	masterClient, err := sql.Open("mysql", masterDsn)
	if err != nil {
		return nil, err
	}

	replicaDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		p.MySQLReplicaUser, p.MySQLReplicaPassword,
		p.MySQLReplicaHost, p.MySQLReplicaPort, p.MySQLReplicaDBName,
	)
	replicaClient, err := sql.Open("mysql", replicaDsn)
	if err != nil {
		return nil, err
	}

	fileRepo, err := repository_mysql.NewFileRepository(
		repository_mysql.WithDbMaster(masterClient),
		repository_mysql.WithDbReplica(replicaClient),
	)
	if err != nil {
		return nil, err
	}

	authRepo, err := repository_mysql.NewAuthRepository(
		repository_mysql.WithDbMaster(masterClient),
		repository_mysql.WithDbReplica(replicaClient),
	)
	if err != nil {
		return nil, err
	}

	r := &NewRepositoryResult{
		FileRepo:  fileRepo,
		OAuthRepo: authRepo,
	}
	return r, nil
}

func newMongoRepository(p NewRepositoryOption) (*NewRepositoryResult, error) {
	dbDsn := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/?authSource=%s",
		p.MongoUser, p.MongoPassword, p.MongoHost, p.MongoPort, p.MongoDBName,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbDsn))
	if err != nil {
		return nil, err
	}

	dbConfig := &repository_mongo.DbConfig{DbName: p.MongoDBName}
	authRepo, err := repository_mongo.NewAuthRepository(
		repository_mongo.WithDbClient(dbClient),
		repository_mongo.WithDbConfig(dbConfig),
	)
	if err != nil {
		return nil, err
	}

	fileRepo, err := repository_mongo.NewFileRepository(
		repository_mongo.WithDbClient(dbClient),
		repository_mongo.WithDbConfig(dbConfig),
	)
	if err != nil {
		return nil, err
	}

	r := &NewRepositoryResult{
		FileRepo:  fileRepo,
		OAuthRepo: authRepo,
	}
	return r, nil
}

type RepositoryOption interface {
	Apply(*NewRepositoryOption)
}

type NewRepositoryOption struct {
	Provider string

	MySQLMasterHost     string
	MySQLMasterPort     int
	MySQLMasterUser     string
	MySQLMasterPassword string
	MySQLMasterDBName   string

	MySQLReplicaHost     string
	MySQLReplicaPort     int
	MySQLReplicaUser     string
	MySQLReplicaPassword string
	MySQLReplicaDBName   string

	MongoHost     string
	MongoPort     int
	MongoUser     string
	MongoPassword string
	MongoDBName   string
}

type NewRepositoryResult struct {
	FileRepo  repository.FileRepository
	OAuthRepo repository.AuthRepository
}

type MySQLConn struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type mysqlOption struct {
	masterHost     string
	masterPort     int
	masterUser     string
	masterPassword string
	masterDbName   string

	replicaHost     string
	replicaPort     int
	replicaUser     string
	replicaPassword string
	replicaDbName   string
}

func (o *mysqlOption) Apply(p *NewRepositoryOption) {
	p.MySQLMasterHost = o.masterHost
	p.MySQLMasterPort = o.masterPort
	p.MySQLMasterDBName = o.masterDbName
	p.MySQLMasterUser = o.masterUser
	p.MySQLMasterPassword = o.masterPassword
	p.MySQLReplicaHost = o.replicaHost
	p.MySQLReplicaPort = o.replicaPort
	p.MySQLReplicaDBName = o.replicaDbName
	p.MySQLReplicaUser = o.replicaUser
	p.MySQLReplicaPassword = o.replicaPassword
	p.Provider = DB_PROVIDER_MYSQL
}

func WithMySQL(mConn MySQLConn, rConn MySQLConn) *mysqlOption {
	return &mysqlOption{
		masterUser:     mConn.User,
		masterPassword: mConn.Password,
		masterDbName:   mConn.DbName,
		masterHost:     mConn.Host,
		masterPort:     mConn.Port,

		replicaUser:     rConn.User,
		replicaPassword: rConn.Password,
		replicaDbName:   rConn.DbName,
		replicaHost:     rConn.Host,
		replicaPort:     rConn.Port,
	}
}

type MongoConn struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type mongoOption struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func (o *mongoOption) Apply(p *NewRepositoryOption) {
	p.MongoDBName = o.DbName
	p.MongoHost = o.Host
	p.MongoPort = o.Port
	p.MongoUser = o.User
	p.MongoPassword = o.Password
	p.Provider = DB_PROVIDER_MONGO
}

func WithMongo(mConn MongoConn) *mongoOption {
	return &mongoOption{
		Host:     mConn.Host,
		Port:     mConn.Port,
		User:     mConn.User,
		Password: mConn.Password,
		DbName:   mConn.DbName,
	}
}
