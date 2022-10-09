package grpc_app

import (
	"context"
	"fmt"

	grpc_v1 "github.com/go-seidon/hippo/generated/grpc-v1"
	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/auth"
	"github.com/go-seidon/hippo/internal/encoding"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/filesystem"
	grpc_auth "github.com/go-seidon/hippo/internal/grpc-auth"
	grpc_log "github.com/go-seidon/hippo/internal/grpc-log"
	"github.com/go-seidon/hippo/internal/hashing"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/logging"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/hippo/internal/text"
	"github.com/go-seidon/hippo/internal/validation"
	"google.golang.org/grpc"
)

type grpcApp struct {
	server        Server
	config        *GrpcAppConfig
	logger        logging.Logger
	repository    repository.Provider
	healthService healthcheck.HealthCheck
}

func (a *grpcApp) Run(ctx context.Context) error {
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
	if err != grpc.ErrServerStopped {
		return err
	}

	return nil
}

func (a *grpcApp) Stop(ctx context.Context) error {
	a.logger.Infof("Stopping %s on: %s", a.config.GetAppName(), a.config.GetAddress())

	err := a.healthService.Stop(ctx)
	if err != nil {
		a.logger.Errorf("Failed stopping healthcheck, err: %s", err.Error())
	}

	return a.server.Shutdown(ctx)
}

func NewGrpcApp(opts ...GrpcAppOption) (*grpcApp, error) {
	p := GrpcAppParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.Config == nil {
		return nil, fmt.Errorf("invalid config")
	}

	config := &GrpcAppConfig{
		AppName:        fmt.Sprintf("%s-grpc", p.Config.AppName),
		AppVersion:     p.Config.AppVersion,
		AppHost:        p.Config.GRPCAppHost,
		AppPort:        p.Config.GRPCAppPort,
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

	fileManager := filesystem.NewFileManager()
	dirManager := filesystem.NewDirectoryManager()
	identifier := text.NewKsuid()
	validator := validation.NewGoValidator()
	locator := file.NewDailyRotate(file.NewDailyRotateParam{})

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

	encoder := encoding.NewBase64Encoder()
	hasher := hashing.NewBcryptHasher()

	basicAuth, err := auth.NewBasicAuth(auth.NewBasicAuthParam{
		AuthRepo: repo.GetAuthRepo(),
		Encoder:  encoder,
		Hasher:   hasher,
	})
	if err != nil {
		return nil, err
	}

	grpcLogOpt := []grpc_log.LogInterceptorOption{
		grpc_log.WithLogger(logger),
		grpc_log.IgnoredMethod([]string{
			"/health.v1.HealthService/CheckHealth",
		}),
		grpc_log.AllowedMetadata([]string{
			"X-Correlation-Id",
		}),
	}
	grpcBasicAuth := grpc_auth.WithAuth(BasicAuth(basicAuth))
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_log.UnaryServerInterceptor(grpcLogOpt...),
			grpc_auth.UnaryServerInterceptor(grpcBasicAuth),
		),
		grpc.ChainStreamInterceptor(
			grpc_log.StreamServerInterceptor(grpcLogOpt...),
			grpc_auth.StreamServerInterceptor(grpcBasicAuth),
		),
	)
	healthCheckHandler := NewHealthHandler(healthService)
	fileHandler := NewFileHandler(fileService, config)
	grpc_v1.RegisterHealthServiceServer(grpcServer, healthCheckHandler)
	grpc_v1.RegisterFileServiceServer(grpcServer, fileHandler)

	svr := p.Server
	if svr == nil {
		svr = &server{
			grpcServer: grpcServer,
			address:    config.GetAddress(),
		}
	}

	app := &grpcApp{
		server:        svr,
		logger:        logger,
		config:        config,
		repository:    repo,
		healthService: healthService,
	}
	return app, nil
}
