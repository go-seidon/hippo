package grpcapp

import (
	"context"
	"fmt"

	"github.com/go-seidon/hippo/api/grpcapp"
	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/auth"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/grpcauth"
	"github.com/go-seidon/hippo/internal/grpchandler"
	"github.com/go-seidon/hippo/internal/grpclog"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/encoding/base64"
	"github.com/go-seidon/provider/hashing/bcrypt"
	"github.com/go-seidon/provider/health"
	"github.com/go-seidon/provider/identity/ksuid"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/validation/govalidator"
	"google.golang.org/grpc"
)

type grpcApp struct {
	server       Server
	config       *GrpcAppConfig
	logger       logging.Logger
	repository   repository.Repository
	healthClient health.HealthCheck
}

func (a *grpcApp) Run(ctx context.Context) error {
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
	err = a.server.ListenAndServe()
	if err != grpc.ErrServerStopped {
		return err
	}

	return nil
}

func (a *grpcApp) Stop(ctx context.Context) error {
	a.logger.Infof("Stopping %s on: %s", a.config.GetAppName(), a.config.GetAddress())

	err := a.healthClient.Stop(ctx)
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

	healthClient := p.HealthClient
	if healthClient == nil {
		healthClient, err = app.NewDefaultHealthCheck(logger, repo)
		if err != nil {
			return nil, err
		}
	}

	fileManager := filesystem.NewFileManager()
	dirManager := filesystem.NewDirectoryManager()
	ksuIdentifier := ksuid.NewIdentifier()
	govalidator := govalidator.NewValidator()
	locator := file.NewDailyRotate(file.DailyRotateParam{})

	fileClient := file.NewFile(file.FileParam{
		FileRepo:    repo.GetFile(),
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

	base64Encoder := base64.NewEncoder()
	bcryptHasher := bcrypt.NewHasher()

	basicClient := auth.NewBasicAuth(auth.NewBasicAuthParam{
		AuthRepo: repo.GetAuth(),
		Encoder:  base64Encoder,
		Hasher:   bcryptHasher,
	})

	grpcLogOpt := []grpclog.LogInterceptorOption{
		grpclog.WithLogger(logger),
		grpclog.IgnoredMethod([]string{
			"/health.v1.HealthService/CheckHealth",
		}),
		grpclog.AllowedMetadata([]string{
			"X-Correlation-Id",
		}),
	}
	grpcBasicAuth := grpcauth.WithAuth(BasicAuth(basicClient))
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclog.UnaryServerInterceptor(grpcLogOpt...),
			grpcauth.UnaryServerInterceptor(grpcBasicAuth),
		),
		grpc.ChainStreamInterceptor(
			grpclog.StreamServerInterceptor(grpcLogOpt...),
			grpcauth.StreamServerInterceptor(grpcBasicAuth),
		),
	)
	healthCheck := healthcheck.NewHealthCheck(healthcheck.HealthCheckParam{
		HealthClient: healthClient,
	})
	healthCheckHandler := grpchandler.NewHealth(grpchandler.HealthParam{
		HealthClient: healthCheck,
	})
	fileHandler := grpchandler.NewFile(grpchandler.FileParam{
		FileClient: fileClient,
		Config: &grpchandler.FileConfig{
			UploadFormSize: config.UploadFormSize,
		},
	})
	grpcapp.RegisterHealthServiceServer(grpcServer, healthCheckHandler)
	grpcapp.RegisterFileServiceServer(grpcServer, fileHandler)

	svr := p.Server
	if svr == nil {
		svr = &server{
			grpcServer: grpcServer,
			address:    config.GetAddress(),
		}
	}

	app := &grpcApp{
		server:       svr,
		logger:       logger,
		config:       config,
		repository:   repo,
		healthClient: healthClient,
	}
	return app, nil
}
