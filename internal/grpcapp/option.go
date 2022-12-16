package grpcapp

import (
	"fmt"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/logging"
)

type GrpcAppConfig struct {
	AppName        string
	AppVersion     string
	AppHost        string
	AppPort        int
	UploadFormSize int64
}

func (c *GrpcAppConfig) GetAppName() string {
	return c.AppName
}

func (c *GrpcAppConfig) GetAppVersion() string {
	return c.AppVersion
}

func (c *GrpcAppConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.AppHost, c.AppPort)
}

type GrpcAppParam struct {
	Config        *app.Config
	Logger        logging.Logger
	Server        Server
	Repository    repository.Provider
	HealthService healthcheck.HealthCheck
}

type GrpcAppOption = func(*GrpcAppParam)

func WithConfig(cfg *app.Config) GrpcAppOption {
	return func(p *GrpcAppParam) {
		p.Config = cfg
	}
}

func WithLogger(logger logging.Logger) GrpcAppOption {
	return func(p *GrpcAppParam) {
		p.Logger = logger
	}
}

func WithService(healthService healthcheck.HealthCheck) GrpcAppOption {
	return func(p *GrpcAppParam) {
		p.HealthService = healthService
	}
}

func WithServer(server Server) GrpcAppOption {
	return func(p *GrpcAppParam) {
		p.Server = server
	}
}

func WithRepository(repo repository.Provider) GrpcAppOption {
	return func(p *GrpcAppParam) {
		p.Repository = repo
	}
}