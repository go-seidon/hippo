package rest_app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/auth"
	db_mongo "github.com/go-seidon/local/internal/db-mongo"
	db_mysql "github.com/go-seidon/local/internal/db-mysql"
	"github.com/go-seidon/local/internal/deleting"
	"github.com/go-seidon/local/internal/encoding"
	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/hashing"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/repository"
	repository_mongo "github.com/go-seidon/local/internal/repository-mongo"
	repository_mysql "github.com/go-seidon/local/internal/repository-mysql"
	"github.com/go-seidon/local/internal/retrieving"
	"github.com/go-seidon/local/internal/serialization"
	"github.com/go-seidon/local/internal/text"
	"github.com/go-seidon/local/internal/uploading"

	"github.com/gorilla/mux"
)

type ContextKey string

const CorrelationIdKey = "correlationId"
const CorrelationIdCtxKey ContextKey = CorrelationIdKey

type RestApp struct {
	config *RestAppConfig
	server app.Server
	logger logging.Logger
	repo   repository.Provider

	healthService healthcheck.HealthCheck
}

func (a *RestApp) Run() error {
	a.logger.Infof("Running %s:%s", a.config.GetAppName(), a.config.GetAppVersion())

	err := a.healthService.Start()
	if err != nil {
		return err
	}

	err = a.repo.Init(context.Background())
	if err != nil {
		return err
	}

	a.logger.Infof("Listening on: %s", a.config.GetAddress())
	err = a.server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *RestApp) Stop() error {
	a.logger.Infof("Stopping %s on: %s", a.config.GetAppName(), a.config.GetAddress())
	return a.server.Shutdown(context.Background())
}

func NewRestApp(opts ...Option) (*RestApp, error) {
	p := RestAppParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.Config == nil {
		return nil, fmt.Errorf("invalid rest app config")
	}
	cfgValid := p.Config.DBProvider == repository.DB_PROVIDER_MYSQL ||
		p.Config.DBProvider == repository.DB_PROVIDER_MONGO
	if !cfgValid {
		return nil, fmt.Errorf("invalid config")
	}

	logger := p.Logger
	if logger == nil {
		opts := []logging.Option{}

		appOpt := logging.WithAppContext(p.Config.AppName, p.Config.AppVersion)
		opts = append(opts, appOpt)

		if p.Config.AppDebug {
			debugOpt := logging.EnableDebugging()
			opts = append(opts, debugOpt)
		}

		if p.Config.AppEnv == app.ENV_LOCAL || p.Config.AppEnv == app.ENV_TEST {
			prettyOpt := logging.EnablePrettyPrint()
			opts = append(opts, prettyOpt)
		}

		stackSkipOpt := logging.AddStackSkip("github.com/go-seidon/local/internal/logging")
		opts = append(opts, stackSkipOpt)

		logger = logging.NewLogrusLog(opts...)
	}

	inetPingJob, err := healthcheck.NewHttpPingJob(healthcheck.NewHttpPingJobParam{
		Name:     "internet-connection",
		Interval: 30 * time.Second,
		Url:      "https://google.com",
	})
	if err != nil {
		return nil, err
	}

	appDiskJob, err := healthcheck.NewDiskUsageJob(healthcheck.NewDiskUsageJobParam{
		Name:      "app-disk",
		Interval:  60 * time.Second,
		Directory: "/",
	})
	if err != nil {
		return nil, err
	}

	healthService := p.HealthService
	if healthService == nil {
		healthService, err = healthcheck.NewGoHealthCheck(
			healthcheck.WithLogger(logger),
			healthcheck.AddJob(inetPingJob),
			healthcheck.AddJob(appDiskJob),
		)
		if err != nil {
			return nil, err
		}
	}

	repo := p.Repository
	if repo == nil {
		if p.Config.DBProvider == repository.DB_PROVIDER_MYSQL {
			dbMaster, err := db_mysql.NewClient(
				db_mysql.WithAuth(p.Config.MySQLMasterUser, p.Config.MySQLMasterPassword),
				db_mysql.WithConfig(db_mysql.ClientConfig{DbName: p.Config.MySQLMasterDBName}),
				db_mysql.WithLocation(p.Config.MySQLMasterHost, p.Config.MySQLMasterPort),
				db_mysql.ParseTime(),
			)
			if err != nil {
				return nil, err
			}

			dbReplica, err := db_mysql.NewClient(
				db_mysql.WithAuth(p.Config.MySQLReplicaUser, p.Config.MySQLReplicaPassword),
				db_mysql.WithConfig(db_mysql.ClientConfig{DbName: p.Config.MySQLReplicaDBName}),
				db_mysql.WithLocation(p.Config.MySQLReplicaHost, p.Config.MySQLReplicaPort),
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
		} else if p.Config.DBProvider == repository.DB_PROVIDER_MONGO {
			opts := []db_mongo.ClientOption{}
			if p.Config.MongoAuthMode == db_mongo.AUTH_BASIC {
				opts = append(opts, db_mongo.WithBasicAuth(
					p.Config.MongoAuthUser,
					p.Config.MongoAuthPassword,
					p.Config.MongoAuthSource,
				))
			}
			if p.Config.MongoMode == db_mongo.MODE_STANDALONE {
				opts = append(opts, db_mongo.UsingStandalone(
					p.Config.MongoStandaloneHost, p.Config.MongoStandalonePort,
				))
			} else if p.Config.MongoMode == db_mongo.MODE_REPLICATION {
				opts = append(opts, db_mongo.UsingReplication(
					p.Config.MongoReplicaName,
					p.Config.MongoReplicaHosts,
				))
			}

			opts = append(opts, db_mongo.WithConfig(db_mongo.ClientConfig{
				DbName: p.Config.MongoDBName,
			}))

			dbClient, err := db_mongo.NewClient(opts...)
			if err != nil {
				return nil, err
			}

			repo, err = repository_mongo.NewRepository(
				repository_mongo.WithDbClient(dbClient),
				repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
					DbName: p.Config.MongoDBName,
				}),
			)
			if err != nil {
				return nil, err
			}
		}
	}

	fileManager := filesystem.NewFileManager()
	dirManager := filesystem.NewDirectoryManager()
	identifier := text.NewKsuid()

	deleteService, err := deleting.NewDeleter(deleting.NewDeleterParam{
		FileRepo:    repo.GetFileRepo(),
		Logger:      logger,
		FileManager: fileManager,
	})
	if err != nil {
		return nil, err
	}

	retrieveService, err := retrieving.NewRetriever(retrieving.NewRetrieverParam{
		FileRepo:    repo.GetFileRepo(),
		Logger:      logger,
		FileManager: fileManager,
	})
	if err != nil {
		return nil, err
	}

	uploadService, err := uploading.NewUploader(uploading.NewUploaderParam{
		FileRepo:    repo.GetFileRepo(),
		FileManager: fileManager,
		Logger:      logger,
		Identifier:  identifier,
		DirManager:  dirManager,
	})
	if err != nil {
		return nil, err
	}

	raCfg := &RestAppConfig{
		AppName:        p.Config.AppName,
		AppVersion:     p.Config.AppVersion,
		AppHost:        p.Config.RESTAppHost,
		AppPort:        p.Config.RESTAppPort,
		UploadFormSize: p.Config.UploadFormSize,
		UploadDir:      p.Config.UploadDirectory,
	}
	locator := uploading.NewDailyRotate(uploading.NewDailyRotateParam{})
	serializer := serialization.NewJsonSerializer()
	encoder := encoding.NewBase64Encoder()
	hasher := hashing.NewBcryptHasher()

	router := mux.NewRouter()
	generalRouter := router.NewRoute().Subrouter()
	fileRouter := router.NewRoute().Subrouter()

	RequestLogMiddleware, err := NewRequestLogMiddleware(RequestLogMiddlewareParam{
		Logger: logger,
		IngoreURI: map[string]bool{
			"/health": true,
		},
		Header: map[string]string{
			"X-Correlation-Id": CorrelationIdKey,
		},
	})
	if err != nil {
		return nil, err
	}

	router.Use(RequestLogMiddleware)
	router.Use(NewDefaultMiddleware(DefaultMiddlewareParam{
		CorrelationIdHeaderKey: "X-Correlation-Id",
		CorrelationIdCtxKey:    CorrelationIdCtxKey,
	}))
	router.HandleFunc(
		"/",
		NewRootHandler(logger, serializer, raCfg),
	)
	generalRouter.HandleFunc(
		"/health",
		NewHealthCheckHandler(logger, serializer, healthService),
	).Methods(http.MethodGet)
	fileRouter.HandleFunc(
		"/file/{id}",
		NewDeleteFileHandler(logger, serializer, deleteService),
	).Methods(http.MethodDelete)
	fileRouter.HandleFunc(
		"/file/{id}",
		NewRetrieveFileHandler(logger, serializer, retrieveService),
	).Methods(http.MethodGet)
	fileRouter.HandleFunc(
		"/file",
		NewUploadFileHandler(logger, serializer, uploadService, locator, raCfg),
	).Methods(http.MethodPost)

	router.NotFoundHandler = NewNotFoundHandler(logger, serializer)
	router.MethodNotAllowedHandler = NewMethodNotAllowedHandler(logger, serializer)

	basicAuth, err := auth.NewBasicAuth(auth.NewBasicAuthParam{
		Encoder:  encoder,
		Hasher:   hasher,
		AuthRepo: repo.GetAuthRepo(),
	})
	if err != nil {
		return nil, err
	}
	BasicAuthMiddleware := NewBasicAuthMiddleware(basicAuth, serializer)
	generalRouter.Use(BasicAuthMiddleware)
	fileRouter.Use(BasicAuthMiddleware)

	server := p.Server
	if p.Server == nil {
		server = &http.Server{
			Addr:     raCfg.GetAddress(),
			Handler:  router,
			ErrorLog: log.New(logger.WriterLevel("error"), "", 0),
		}
	}

	app := &RestApp{
		server:        server,
		config:        raCfg,
		logger:        logger,
		healthService: healthService,
		repo:          repo,
	}
	return app, nil
}
