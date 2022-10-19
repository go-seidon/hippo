package app

import (
	"time"

	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/logging"
)

func NewDefaultHealthCheck(logger logging.Logger, repo repository.Provider) (healthcheck.HealthCheck, error) {

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
		Interval:  10 * time.Minute,
		Directory: "/",
	})
	if err != nil {
		return nil, err
	}

	repoPingJob, err := healthcheck.NewRepoPingJob(healthcheck.NewRepoPingJobParam{
		Name:       "repository-connection",
		Interval:   15 * time.Minute,
		DataSource: repo,
	})
	if err != nil {
		return nil, err
	}

	healthService, err := healthcheck.NewGoHealthCheck(
		healthcheck.WithLogger(logger),
		healthcheck.AddJob(inetPingJob),
		healthcheck.AddJob(appDiskJob),
		healthcheck.AddJob(repoPingJob),
	)
	if err != nil {
		return nil, err
	}

	return healthService, nil
}
