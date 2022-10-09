package app

type Config struct {
	AppName    string `env:"APP_NAME"`
	AppEnv     string `env:"APP_ENV"`
	AppVersion string `env:"APP_VERSION"`
	AppDebug   bool   `env:"APP_DEBUG"`

	RESTAppHost string `env:"REST_APP_HOST"`
	RESTAppPort int    `env:"REST_APP_PORT"`

	GRPCAppHost string `env:"GRPC_APP_HOST"`
	GRPCAppPort int    `env:"GRPC_APP_PORT"`

	DBProvider string `env:"DB_PROVIDER"`

	MySQLMasterHost     string `env:"MYSQL_MASTER_HOST"`
	MySQLMasterPort     int    `env:"MYSQL_MASTER_PORT"`
	MySQLMasterUser     string `env:"MYSQL_MASTER_USER"`
	MySQLMasterPassword string `env:"MYSQL_MASTER_PASSWORD"`
	MySQLMasterDBName   string `env:"MYSQL_MASTER_DB_NAME"`

	MySQLReplicaHost     string `env:"MYSQL_REPLICA_HOST"`
	MySQLReplicaPort     int    `env:"MYSQL_REPLICA_PORT"`
	MySQLReplicaUser     string `env:"MYSQL_REPLICA_USER"`
	MySQLReplicaPassword string `env:"MYSQL_REPLICA_PASSWORD"`
	MySQLReplicaDBName   string `env:"MYSQL_REPLICA_DB_NAME"`

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
