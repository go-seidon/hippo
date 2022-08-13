package app

type Config struct {
	AppName    string `env:"APP_NAME"`
	AppEnv     string `env:"APP_ENV"`
	AppVersion string `env:"APP_VERSION"`
	AppDebug   bool   `env:"APP_DEBUG"`

	RESTAppHost string `env:"REST_APP_HOST"`
	RESTAppPort int    `env:"REST_APP_PORT"`

	RPCAppHost string `env:"RPC_APP_HOST"`
	RPCAppPort int    `env:"RPC_APP_PORT"`

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

	UploadFormSize  int64  `env:"UPLOAD_FORM_SIZE"`
	UploadDirectory string `env:"UPLOAD_DIRECTORY"`
}
