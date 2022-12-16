package restapp

import (
	"context"
	"fmt"
	"log"
	net_http "net/http"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/auth"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/encoding/base64"
	"github.com/go-seidon/provider/hashing/bcrypt"
	"github.com/go-seidon/provider/http"
	"github.com/go-seidon/provider/identifier/ksuid"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/serialization/json"
	"github.com/go-seidon/provider/validation/govalidator"
	"github.com/gorilla/mux"
)

type ContextKey string

const CorrelationIdKey = "correlationId"
const CorrelationIdCtxKey ContextKey = CorrelationIdKey

type restApp struct {
	config     *RestAppConfig
	server     http.Server
	logger     logging.Logger
	repository repository.Provider

	healthService healthcheck.HealthCheck
}

func (a *restApp) Run(ctx context.Context) error {
	a.logger.Infof("Running %s:%s", a.config.GetAppName(), a.config.GetAppVersion())

	err := a.healthService.Start(ctx)
	if err != nil {
		return err
	}

	err = a.repository.Init(ctx)
	if err != nil {
		return err
	}

	a.logger.Infof("Listening on: %s", a.config.GetAddress())
	err = a.server.ListenAndServe()
	if err != net_http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *restApp) Stop(ctx context.Context) error {
	a.logger.Infof("Stopping %s on: %s", a.config.GetAppName(), a.config.GetAddress())

	err := a.healthService.Stop(ctx)
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

	config := &RestAppConfig{
		AppName:        fmt.Sprintf("%s-rest", p.Config.AppName),
		AppVersion:     p.Config.AppVersion,
		AppHost:        p.Config.RESTAppHost,
		AppPort:        p.Config.RESTAppPort,
		UploadFormSize: p.Config.UploadFormSize,
	}

	var err error
	logger := p.Logger
	if logger == nil {
		logger, err = app.NewDefaultLog(p.Config, config.AppName)
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

	govalidator := govalidator.NewValidator()
	ksuIdentifier := ksuid.NewIdentifier()
	fileManager := filesystem.NewFileManager()
	dirManager := filesystem.NewDirectoryManager()
	locator := file.NewDailyRotate(file.NewDailyRotateParam{})

	fileService, err := file.NewFile(file.NewFileParam{
		FileRepo:    repo.GetFileRepo(),
		FileManager: fileManager,
		Logger:      logger,
		Identifier:  ksuIdentifier,
		DirManager:  dirManager,
		Locator:     locator,
		Validator:   govalidator,
		Config: &file.FileConfig{
			UploadDir: p.Config.UploadDirectory,
		},
	})
	if err != nil {
		return nil, err
	}

	jsonSerializer := json.NewSerializer()
	base64Encoder := base64.NewEncoder()
	bcryptHasher := bcrypt.NewHasher()

	basicHandler := NewBasicHandler(BasicHandlerParam{
		Logger:     p.Logger,
		Serializer: jsonSerializer,
		Config:     config,
	})
	healthHandler := NewHealthHandler(HealthHandlerParam{
		Logger:        p.Logger,
		Serializer:    jsonSerializer,
		HealthService: healthService,
	})
	fileHandler := NewFileHandler(FileHandlerParam{
		Logger:      p.Logger,
		Serializer:  jsonSerializer,
		Config:      config,
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
	router.NotFoundHandler = net_http.HandlerFunc(basicHandler.NotFound)
	router.MethodNotAllowedHandler = net_http.HandlerFunc(basicHandler.MethodNotAllowed)

	generalRouter.HandleFunc("/health", healthHandler.CheckHealth).Methods(net_http.MethodGet)
	fileRouter.HandleFunc("/v1/file/{id}", fileHandler.DeleteFileById).Methods(net_http.MethodDelete)
	fileRouter.HandleFunc("/v1/file/{id}", fileHandler.RetrieveFileById).Methods(net_http.MethodGet)
	fileRouter.HandleFunc("/v1/file", fileHandler.UploadFile).Methods(net_http.MethodPost)

	basicAuth, err := auth.NewBasicAuth(auth.NewBasicAuthParam{
		Encoder:  base64Encoder,
		Hasher:   bcryptHasher,
		AuthRepo: repo.GetAuthRepo(),
	})
	if err != nil {
		return nil, err
	}
	BasicAuthMiddleware := NewBasicAuthMiddleware(basicAuth, jsonSerializer)
	generalRouter.Use(BasicAuthMiddleware)
	fileRouter.Use(BasicAuthMiddleware)

	server := p.Server
	if p.Server == nil {
		server = &net_http.Server{
			Addr:     config.GetAddress(),
			Handler:  router,
			ErrorLog: log.New(logger.WriterLevel("error"), "", 0),
		}
	}

	app := &restApp{
		server:        server,
		config:        config,
		logger:        logger,
		healthService: healthService,
		repository:    repo,
	}
	return app, nil
}
