package healthcheck

import (
	"fmt"

	"github.com/InVisionApp/go-health"
)

type HealthClient interface {
	AddChecks(cfgs []*health.Config) error
	Start() error
	Stop() error
	State() (map[string]health.State, bool, error)
}

type goHealthCheck struct {
	client    HealthClient
	jobs      []*HealthJob
	runStatus bool
}

func (s *goHealthCheck) Start() error {
	if s.runStatus {
		return nil
	}

	cfgs := []*health.Config{}
	for _, job := range s.jobs {
		cfgs = append(cfgs, &health.Config{
			Name:     job.Name,
			Checker:  job.Checker,
			Interval: job.Interval,
		})
	}

	err := s.client.AddChecks(cfgs)
	if err != nil {
		return err
	}

	err = s.client.Start()
	if err != nil {
		return err
	}

	s.runStatus = true
	return nil
}

func (s *goHealthCheck) Stop() error {
	return s.client.Stop()
}

func (s *goHealthCheck) Check() (*CheckResult, error) {
	states, isFailed, err := s.client.State()
	if err != nil {
		return nil, err
	}

	res := &CheckResult{
		Status: STATUS_FAILED,
		Items:  make(map[string]CheckResultItem),
	}
	if isFailed {
		return res, nil
	}

	totalFailed := 0
	for key, state := range states {

		status := STATUS_OK
		if state.Status == "failed" {
			status = STATUS_FAILED
			totalFailed++
		}

		res.Items[key] = CheckResultItem{
			Name:      state.Name,
			Status:    status,
			Error:     state.Err,
			CheckedAt: state.CheckTime.UTC(),
		}
	}

	if totalFailed == 0 {
		res.Status = STATUS_OK
	} else if totalFailed != len(states) {
		res.Status = STATUS_WARNING
	}

	return res, nil
}

func NewGoHealthCheck(opts ...Option) (*goHealthCheck, error) {
	p := HealthCheckParam{
		Jobs: []*HealthJob{},
	}
	for _, opt := range opts {
		opt(&p)
	}
	if len(p.Jobs) == 0 {
		return nil, fmt.Errorf("invalid jobs specified")
	}
	if p.Logger == nil {
		return nil, fmt.Errorf("invalid logger specified")
	}

	client := p.Client
	if client == nil {
		h := health.New()
		hlog, err := NewGoHealthLog(p.Logger)
		if err != nil {
			return nil, err
		}
		h.Logger = hlog
		client = h
	}

	s := &goHealthCheck{
		client: client,
		jobs:   p.Jobs,
	}
	return s, nil
}
