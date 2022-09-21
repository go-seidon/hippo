package rest_app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/auth"
	"github.com/go-seidon/local/internal/deleting"
	"github.com/go-seidon/local/internal/encoding"
	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/hashing"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/repository"
	"github.com/go-seidon/local/internal/retrieving"
	"github.com/go-seidon/local/internal/serialization"
	"github.com/go-seidon/local/internal/text"
	"github.com/go-seidon/local/internal/uploading"

	"github.com/gorilla/mux"
)

type ContextKey string

const CorrelationIdKey = "correlationId"
const CorrelationIdCtxKey ContextKey = CorrelationIdKey

type restApp struct {
	config     *RestAppConfig
	server     Server
	logger     logging.Logger
	repository repository.Provider

	healthService healthcheck.HealthCheck
}

func (a *restApp) Run(ctx context.Context) error {
	a.logger.Infof("Running %s:%s", a.config.GetAppName(), a.config.GetAppVersion())

	err := a.healthService.Start()
	if err != nil {
		return err
	}

	err = a.repository.Init(ctx)
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

func (a *restApp) Stop(ctx context.Context) error {
	a.logger.Infof("Stopping %s on: %s", a.config.GetAppName(), a.config.GetAddress())

	err := a.healthService.Stop()
	if err != nil {
		a.logger.Errorf("Failed stopping healthcheck, err: %s", err.Error())
	}

	return a.server.Shutdown(ctx)
}

func NewRestApp(opts ...RestAppOption) (*restApp, error) {
	p := RestAppParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.Config == nil {
		return nil, fmt.Errorf("invalid config")
	}

	var err error
	logger := p.Logger
	if logger == nil {
		logger, err = app.NewDefaultLog(p.Config)
		if err != nil {
			return nil, err
		}
	}

	repo := p.Repository
	if repo == nil {
		repo, err = app.NewDefaultRepository(p.Config)
		if err != nil {
			return nil, err
		}
	}

	healthService := p.HealthService
	if healthService == nil {
		healthService, err = app.NewDefaultHealthCheck(logger, repo)
		if err != nil {
			return nil, err
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

	app := &restApp{
		server:        server,
		config:        raCfg,
		logger:        logger,
		healthService: healthService,
		repository:    repo,
	}
	return app, nil
}
