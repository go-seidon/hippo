package app

import (
	"time"

	"github.com/go-seidon/local/internal/healthcheck"
)

func NewDefaultHealthCheck(opts ...healthcheck.Option) (healthcheck.HealthCheck, error) {
	p := healthcheck.HealthCheckParam{
		Jobs: []*healthcheck.HealthJob{},
	}
	for _, opt := range opts {
		opt(&p)
	}

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

	healthService, err := healthcheck.NewGoHealthCheck(
		healthcheck.WithLogger(p.Logger),
		healthcheck.AddJob(inetPingJob),
		healthcheck.AddJob(appDiskJob),
	)
	if err != nil {
		return nil, err
	}

	return healthService, nil
}
