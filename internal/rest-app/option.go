package rest_app

import (
	"fmt"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/repository"
)

type RestAppConfig struct {
	AppName        string
	AppVersion     string
	AppHost        string
	AppPort        int
	UploadDir      string
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

func WithConfig(cfg app.Config) RestAppOption {
	return func(rao *RestAppParam) {
		rao.Config = &cfg
	}
}

func WithLogger(logger logging.Logger) RestAppOption {
	return func(rao *RestAppParam) {
		rao.Logger = logger
	}
}

func WithServer(server Server) RestAppOption {
	return func(rao *RestAppParam) {
		rao.Server = server
	}
}

func WithService(healthService healthcheck.HealthCheck) RestAppOption {
	return func(rao *RestAppParam) {
		rao.HealthService = healthService
	}
}

func WithRepository(repo repository.Provider) RestAppOption {
	return func(rao *RestAppParam) {
		rao.Repository = repo
	}
}
