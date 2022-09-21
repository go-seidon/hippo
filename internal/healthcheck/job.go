package healthcheck

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/InVisionApp/go-health/checkers"
	diskchk "github.com/InVisionApp/go-health/checkers/disk"
)

type NewHttpPingJobParam struct {
	Name     string
	Interval time.Duration
	Url      string
}

func NewHttpPingJob(p NewHttpPingJobParam) (*HealthJob, error) {
	if strings.TrimSpace(p.Name) == "" {
		return nil, fmt.Errorf("invalid name")
	}

	pingUrl, err := url.Parse(p.Url)
	if err != nil {
		return nil, err
	}

	internetConnection, err := checkers.NewHTTP(&checkers.HTTPConfig{
		URL: pingUrl,
	})
	if err != nil {
		return nil, err
	}

	job := &HealthJob{
		Name:     p.Name,
		Interval: p.Interval,
		Checker:  internetConnection,
	}
	return job, err
}

type NewDiskUsageJobParam struct {
	Name      string
	Interval  time.Duration
	Directory string
}

func NewDiskUsageJob(p NewDiskUsageJobParam) (*HealthJob, error) {
	if strings.TrimSpace(p.Name) == "" {
		return nil, fmt.Errorf("invalid name")
	}
	if strings.TrimSpace(p.Directory) == "" {
		return nil, fmt.Errorf("invalid directory")
	}

	appDiskChecker, err := diskchk.NewDiskUsage(&diskchk.DiskUsageConfig{
		Path:              p.Directory,
		WarningThreshold:  50,
		CriticalThreshold: 20,
	})
	if err != nil {
		return nil, err
	}

	job := &HealthJob{
		Name:     p.Name,
		Interval: p.Interval,
		Checker:  appDiskChecker,
	}
	return job, err
}

type NewRepoPingJobParam struct {
	Name       string
	Interval   time.Duration
	DataSource DataSource
}

func NewRepoPingJob(p NewRepoPingJobParam) (*HealthJob, error) {
	if strings.TrimSpace(p.Name) == "" {
		return nil, fmt.Errorf("invalid name")
	}
	if p.DataSource == nil {
		return nil, fmt.Errorf("invalid data source")
	}

	pingChecker := NewRepoPingChecker(p.DataSource)

	job := &HealthJob{
		Name:     p.Name,
		Interval: p.Interval,
		Checker:  pingChecker,
	}
	return job, nil
}
