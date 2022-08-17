package app

import (
	"database/sql"
	"fmt"

	"github.com/go-seidon/local/internal/repository"
	repository_mysql "github.com/go-seidon/local/internal/repository-mysql"
	_ "github.com/go-sql-driver/mysql"
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

	authRepo, err := repository_mysql.NewOAuthRepository(
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
