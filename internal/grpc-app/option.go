package grpc_app

import (
	"fmt"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/repository"
)

type GrpcAppConfig struct {
	AppName    string
	AppVersion string
	AppHost    string
	AppPort    int
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
