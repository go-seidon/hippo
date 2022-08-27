package rest_app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/auth"
	"github.com/go-seidon/local/internal/deleting"
	"github.com/go-seidon/local/internal/encoding"
	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/hashing"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
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

	healthService healthcheck.HealthCheck
}

func (a *RestApp) Run() error {
	a.logger.Infof("Running %s:%s", a.config.GetAppName(), a.config.GetAppVersion())

	err := a.healthService.Start()
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
	option := RestAppOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.Config == nil {
		return nil, fmt.Errorf("invalid rest app config")
	}
	if option.Config.DBProvider != app.DB_PROVIDER_MYSQL &&
		option.Config.DBProvider != app.DB_PROVIDER_MONGO {
		return nil, fmt.Errorf("unsupported db provider")
	}

	logger := option.Logger
	if option.Logger == nil {
		opts := []logging.Option{}

		appOpt := logging.WithAppContext(option.Config.AppName, option.Config.AppVersion)
		opts = append(opts, appOpt)

		if option.Config.AppDebug {
			debugOpt := logging.EnableDebugging()
			opts = append(opts, debugOpt)
		}

		if option.Config.AppEnv == app.ENV_LOCAL || option.Config.AppEnv == app.ENV_TEST {
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

	healthService := option.HealthService
	if option.HealthService == nil {
		healthCheck, err := healthcheck.NewGoHealthCheck(
			healthcheck.WithLogger(logger),
			healthcheck.AddJob(inetPingJob),
			healthcheck.AddJob(appDiskJob),
		)
		if err != nil {
			return nil, err
		}
		healthService = healthCheck
	}

	var repoOpt app.RepositoryOption
	if option.Config.DBProvider == app.DB_PROVIDER_MYSQL {
		repoOpt = app.WithMySQL(app.MySQLConn{
			Host:     option.Config.MySQLMasterHost,
			Port:     option.Config.MySQLMasterPort,
			User:     option.Config.MySQLMasterUser,
			Password: option.Config.MySQLMasterPassword,
			DbName:   option.Config.MySQLMasterDBName,
		}, app.MySQLConn{
			Host:     option.Config.MySQLReplicaHost,
			Port:     option.Config.MySQLReplicaPort,
			User:     option.Config.MySQLReplicaUser,
			Password: option.Config.MySQLReplicaPassword,
			DbName:   option.Config.MySQLReplicaDBName,
		})
	} else if option.Config.DBProvider == app.DB_PROVIDER_MONGO {
		repoOpt = app.WithMongo(app.MongoConn{
			Host:     option.Config.MongoMasterHost,
			Port:     option.Config.MongoMasterPort,
			User:     option.Config.MongoMasterUser,
			Password: option.Config.MongoMasterPassword,
			DbName:   option.Config.MongoMasterDBName,
		})
	}
	repo, err := app.NewRepository(repoOpt)
	if err != nil {
		return nil, err
	}

	fileManager := filesystem.NewFileManager()
	dirManager := filesystem.NewDirectoryManager()
	identifier := text.NewKsuid()

	deleteService, err := deleting.NewDeleter(deleting.NewDeleterParam{
		FileRepo:    repo.FileRepo,
		Logger:      logger,
		FileManager: fileManager,
	})
	if err != nil {
		return nil, err
	}

	retrieveService, err := retrieving.NewRetriever(retrieving.NewRetrieverParam{
		FileRepo:    repo.FileRepo,
		Logger:      logger,
		FileManager: fileManager,
	})
	if err != nil {
		return nil, err
	}

	uploadService, err := uploading.NewUploader(uploading.NewUploaderParam{
		FileRepo:    repo.FileRepo,
		FileManager: fileManager,
		Logger:      logger,
		Identifier:  identifier,
		DirManager:  dirManager,
	})
	if err != nil {
		return nil, err
	}

	raCfg := &RestAppConfig{
		AppName:        option.Config.AppName,
		AppVersion:     option.Config.AppVersion,
		AppHost:        option.Config.RESTAppHost,
		AppPort:        option.Config.RESTAppPort,
		UploadFormSize: option.Config.UploadFormSize,
		UploadDir:      option.Config.UploadDirectory,
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
		Encoder:   encoder,
		Hasher:    hasher,
		OAuthRepo: repo.OAuthRepo,
	})
	if err != nil {
		return nil, err
	}
	BasicAuthMiddleware := NewBasicAuthMiddleware(basicAuth, serializer)
	generalRouter.Use(BasicAuthMiddleware)
	fileRouter.Use(BasicAuthMiddleware)

	server := option.Server
	if option.Server == nil {
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
	}
	return app, nil
}
