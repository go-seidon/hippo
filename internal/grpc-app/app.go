package grpc_app

import (
	"context"
	"fmt"
	"time"

	grpc_v1 "github.com/go-seidon/local/generated/proto/api/grpc/v1"
	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"google.golang.org/grpc"
)

type grpcApp struct {
	server        Server
	config        *GrpcAppConfig
	logger        logging.Logger
	healthService healthcheck.HealthCheck
}

func (a *grpcApp) Run(ctx context.Context) error {
	a.logger.Infof("Running %s:%s", a.config.GetAppName(), a.config.GetAppVersion())

	err := a.healthService.Start()
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

	err := a.healthService.Stop()
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
		return nil, fmt.Errorf("invalid grpc app config")
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

	healthService := p.HealthService
	if healthService == nil {
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

		healthService, err = healthcheck.NewGoHealthCheck(
			healthcheck.WithLogger(logger),
			healthcheck.AddJob(inetPingJob),
			healthcheck.AddJob(appDiskJob),
		)
		if err != nil {
			return nil, err
		}
	}

	grpcServer := grpc.NewServer()
	healthCheckHandler := NewHealthHandler(healthService)
	grpc_v1.RegisterHealthServiceServer(grpcServer, healthCheckHandler)

	config := &GrpcAppConfig{
		AppName:    p.Config.AppName,
		AppVersion: p.Config.AppVersion,
		AppHost:    p.Config.RPCAppHost,
		AppPort:    p.Config.RPCAppPort,
	}

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
		healthService: healthService,
	}
	return app, nil
}
