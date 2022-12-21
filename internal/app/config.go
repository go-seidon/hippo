package app

import (
	"github.com/go-seidon/provider/config/viper"
)

type Config struct {
	AppName    string `env:"APP_NAME"`
	AppEnv     string `env:"APP_ENV"`
	AppVersion string `env:"APP_VERSION"`
	AppDebug   bool   `env:"APP_DEBUG"`

	RESTAppHost string `env:"REST_APP_HOST"`
	RESTAppPort int    `env:"REST_APP_PORT"`

	GRPCAppHost string `env:"GRPC_APP_HOST"`
	GRPCAppPort int    `env:"GRPC_APP_PORT"`

	RepositoryProvider string `env:"REPOSITORY_PROVIDER"`

	MySQLPrimaryHost     string `env:"MYSQL_PRIMARY_HOST"`
	MySQLPrimaryPort     int    `env:"MYSQL_PRIMARY_PORT"`
	MySQLPrimaryUser     string `env:"MYSQL_PRIMARY_USER"`
	MySQLPrimaryPassword string `env:"MYSQL_PRIMARY_PASSWORD"`
	MySQLPrimaryDBName   string `env:"MYSQL_PRIMARY_DB_NAME"`

	MySQLSecondaryHost     string `env:"MYSQL_SECONDARY_HOST"`
	MySQLSecondaryPort     int    `env:"MYSQL_SECONDARY_PORT"`
	MySQLSecondaryUser     string `env:"MYSQL_SECONDARY_USER"`
	MySQLSecondaryPassword string `env:"MYSQL_SECONDARY_PASSWORD"`
	MySQLSecondaryDBName   string `env:"MYSQL_SECONDARY_DB_NAME"`

	MongoMode           string   `env:"MONGO_MODE"`
	MongoDBName         string   `env:"MONGO_DB_NAME"`
	MongoAuthMode       string   `env:"MONGO_AUTH_MODE"`
	MongoAuthUser       string   `env:"MONGO_AUTH_USER"`
	MongoAuthPassword   string   `env:"MONGO_AUTH_PASSWORD"`
	MongoAuthSource     string   `env:"MONGO_AUTH_SOURCE"`
	MongoStandaloneHost string   `env:"MONGO_STANDALONE_HOST"`
	MongoStandalonePort int      `env:"MONGO_STANDALONE_PORT"`
	MongoReplicaName    string   `env:"MONGO_REPLICA_NAME"`
	MongoReplicaHosts   []string `env:"MONGO_REPLICA_HOSTS"`

	UploadFormSize  int64  `env:"UPLOAD_FORM_SIZE"`
	UploadDirectory string `env:"UPLOAD_DIRECTORY"`
}

func NewDefaultConfig() (*Config, error) {
	tomlConfig, err := viper.NewConfig(
		viper.WithFileName("config/default.toml"),
	)
	if err != nil {
		return nil, err
	}

	err = tomlConfig.LoadConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = tomlConfig.ParseConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
