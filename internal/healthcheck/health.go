package healthcheck

import (
	"context"
	"time"

	"github.com/go-seidon/provider/logging"
)

const (
	STATUS_OK      = "OK"
	STATUS_WARNING = "WARNING"
	STATUS_FAILED  = "FAILED"
)

type HealthCheck interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Check(ctx context.Context) (*CheckResult, error)
}

type CheckResult struct {
	Status string
	Items  map[string]CheckResultItem
}

type CheckResultItem struct {
	Name      string
	Status    string
	Error     string
	Fatal     bool
	CheckedAt time.Time
}

type HealthCheckParam struct {
	Jobs   []*HealthJob
	Logger logging.Logger
	Client HealthClient
}

type HealthJob struct {
	Name     string
	Checker  Checker
	Interval time.Duration
}

type Checker interface {
	Status() (interface{}, error)
}

type Option func(*HealthCheckParam)

func WithLogger(logger logging.Logger) Option {
	return func(p *HealthCheckParam) {
		p.Logger = logger
	}
}

func AddJob(job *HealthJob) Option {
	return func(p *HealthCheckParam) {
		p.Jobs = append(p.Jobs, job)
	}
}

func WithJobs(jobs []*HealthJob) Option {
	return func(p *HealthCheckParam) {
		p.Jobs = jobs
	}
}

func WithClient(client HealthClient) Option {
	return func(p *HealthCheckParam) {
		p.Client = client
	}
}
