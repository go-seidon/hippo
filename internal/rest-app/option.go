package rest_app

import (
	"fmt"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/logging"
	"github.com/go-seidon/hippo/internal/repository"
)

type RestAppConfig struct {
	AppName        string
	AppVersion     string
	AppHost        string
	AppPort        int
	UploadFormSize int64
}

func (c *RestAppConfig) GetAppName() string {
	return c.AppName
}

func (c *RestAppConfig) GetAppVersion() string {
	return c.AppVersion
}

func (c *RestAppConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.AppHost, c.AppPort)
}

type RestAppParam struct {
	Config        *app.Config
	Logger        logging.Logger
	Server        Server
	Repository    repository.Provider
	HealthService healthcheck.HealthCheck
}

type RestAppOption func(*RestAppParam)

func WithConfig(cfg *app.Config) RestAppOption {
	return func(p *RestAppParam) {
		p.Config = cfg
	}
}

func WithLogger(logger logging.Logger) RestAppOption {
	return func(p *RestAppParam) {
		p.Logger = logger
	}
}

func WithServer(server Server) RestAppOption {
	return func(p *RestAppParam) {
		p.Server = server
	}
}

func WithService(healthService healthcheck.HealthCheck) RestAppOption {
	return func(p *RestAppParam) {
		p.HealthService = healthService
	}
}

func WithRepository(repo repository.Provider) RestAppOption {
	return func(p *RestAppParam) {
		p.Repository = repo
	}
}
