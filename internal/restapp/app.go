package restapp

import (
	"context"
	"fmt"
	net_http "net/http"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/auth"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/hippo/internal/resthandler"
	"github.com/go-seidon/hippo/internal/restmiddleware"
	"github.com/go-seidon/hippo/internal/storage/multipart"
	"github.com/go-seidon/provider/encoding/base64"
	"github.com/go-seidon/provider/hashing/bcrypt"
	"github.com/go-seidon/provider/identity/ksuid"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/serialization/json"
	"github.com/go-seidon/provider/validation/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type restApp struct {
	config     *RestAppConfig
	server     Server
	logger     logging.Logger
	repository repository.Provider

	healthClient healthcheck.HealthCheck
}

func (a *restApp) Run(ctx context.Context) error {
	a.logger.Infof("Running %s:%s", a.config.GetAppName(), a.config.GetAppVersion())

	err := a.healthClient.Start(ctx)
	if err != nil {
		return err
	}

	err = a.repository.Init(ctx)
	if err != nil {
		return err
	}

	a.logger.Infof("Listening on: %s", a.config.GetAddress())
	err = a.server.Start(a.config.GetAddress())
	if err != net_http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *restApp) Stop(ctx context.Context) error {
	a.logger.Infof("Stopping %s on: %s", a.config.GetAppName(), a.config.GetAddress())

	err := a.healthClient.Stop(ctx)
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

	healthClient := p.HealthClient
	if healthClient == nil {
		healthClient, err = app.NewDefaultHealthCheck(logger, repo)
		if err != nil {
			return nil, err
		}
	}

	server := p.Server
	if p.Server == nil {
		e := echo.New()
		e.Use(middleware.Recover())
		server = &echoServer{e}

		jsonSerializer := json.NewSerializer()
		base64Encoder := base64.NewEncoder()
		bcryptHasher := bcrypt.NewHasher()
		govalidator := govalidator.NewValidator()
		ksuIdentifier := ksuid.NewIdentifier()
		fileManager := filesystem.NewFileManager()
		dirManager := filesystem.NewDirectoryManager()
		locator := file.NewDailyRotate(file.DailyRotateParam{})

		basicClient := auth.NewBasicAuth(auth.NewBasicAuthParam{
			Encoder:  base64Encoder,
			Hasher:   bcryptHasher,
			AuthRepo: repo.GetAuthRepo(),
		})

		fileClient := file.NewFile(file.FileParam{
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

		basicHandler := resthandler.NewBasic(resthandler.BasicParam{
			Config: &resthandler.BasicConfig{
				AppName:    config.AppName,
				AppVersion: config.AppVersion,
			},
		})
		healthHandler := resthandler.NewHealth(resthandler.HealthParam{
			HealthClient: healthClient,
		})
		fileHandler := resthandler.NewFile(resthandler.FileParam{
			FileClient: fileClient,
			FileParser: multipart.FileParser,
		})

		requestLog := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
			Logger: logger,
			IgnoreURI: map[string]bool{
				"/health": true,
			},
			Header: map[string]string{
				"X-Correlation-Id": "correlationId",
			},
		})
		logMiddleware := echo.WrapMiddleware(requestLog.Handle)

		basicAuth := restmiddleware.NewBasicAuth(restmiddleware.BasicAuthParam{
			Serializer:  jsonSerializer,
			BasicClient: basicClient,
		})
		basicAuthMiddleware := echo.WrapMiddleware(basicAuth.Handle)

		basicGroup := e.Group("", logMiddleware)
		basicGroup.GET("/", basicHandler.GetAppInfo)

		basicAuthGroup := e.Group("", basicAuthMiddleware)
		basicAuthGroup.GET("/health", healthHandler.CheckHealth)
		basicAuthGroup.POST("/v1/file", fileHandler.UploadFile)
		basicAuthGroup.GET("/v1/file/:id", fileHandler.RetrieveFileById)
		basicAuthGroup.DELETE("/v1/file/:id", fileHandler.DeleteFileById)
	}

	app := &restApp{
		server:       server,
		config:       config,
		logger:       logger,
		healthClient: healthClient,
		repository:   repo,
	}
	return app, nil
}
