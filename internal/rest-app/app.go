package rest_app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/auth"
	"github.com/go-seidon/local/internal/encoding"
	"github.com/go-seidon/local/internal/file"
	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/hashing"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/repository"
	"github.com/go-seidon/local/internal/serialization"
	"github.com/go-seidon/local/internal/text"
	"github.com/go-seidon/local/internal/validation"

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
	locator := file.NewDailyRotate(file.NewDailyRotateParam{})
	validator := validation.NewGoValidator()

	fileService, err := file.NewFile(file.NewFileParam{
		FileRepo:    repo.GetFileRepo(),
		FileManager: fileManager,
		Logger:      logger,
		Identifier:  identifier,
		DirManager:  dirManager,
		Locator:     locator,
		Validator:   validator,
		Config: &file.FileConfig{
			UploadDir: p.Config.UploadDirectory,
		},
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
	}
	serializer := serialization.NewJsonSerializer()
	encoder := encoding.NewBase64Encoder()
	hasher := hashing.NewBcryptHasher()

	basicHandler := NewBasicHandler(BasicHandlerParam{
		Logger:     p.Logger,
		Serializer: serializer,
		Config:     raCfg,
	})
	healthHandler := NewHealthHandler(HealthHandlerParam{
		Logger:        p.Logger,
		Serializer:    serializer,
		HealthService: healthService,
	})
	fileHandler := NewFileHandler(FileHandlerParam{
		Logger:      p.Logger,
		Serializer:  serializer,
		Config:      raCfg,
		FileService: fileService,
	})

	RequestLogMiddleware, err := NewRequestLogMiddleware(RequestLogMiddlewareParam{
		Logger: logger,
		IgnoreURI: map[string]bool{
			"/health": true,
		},
		Header: map[string]string{
			"X-Correlation-Id": CorrelationIdKey,
		},
	})
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	generalRouter := router.NewRoute().Subrouter()
	fileRouter := router.NewRoute().Subrouter()

	router.Use(RequestLogMiddleware)
	router.Use(NewDefaultMiddleware(DefaultMiddlewareParam{
		CorrelationIdHeaderKey: "X-Correlation-Id",
		CorrelationIdCtxKey:    CorrelationIdCtxKey,
	}))
	router.HandleFunc("/", basicHandler.GetAppInfo)
	router.NotFoundHandler = http.HandlerFunc(basicHandler.GetAppInfo)
	router.MethodNotAllowedHandler = http.HandlerFunc(basicHandler.MethodNotAllowed)

	generalRouter.HandleFunc("/health", healthHandler.CheckHealth).Methods(http.MethodGet)
	fileRouter.HandleFunc("/file/{id}", fileHandler.DeleteFileById).Methods(http.MethodDelete)
	fileRouter.HandleFunc("/file/{id}", fileHandler.RetrieveFileById).Methods(http.MethodGet)
	fileRouter.HandleFunc("/file", fileHandler.UploadFile).Methods(http.MethodPost)

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
