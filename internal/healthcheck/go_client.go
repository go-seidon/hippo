package healthcheck

import (
	"fmt"

	"github.com/InVisionApp/go-health"
)

type goHealthClient struct {
	h *health.Health
}

func (c *goHealthClient) AddChecks(cfgs []*HealthConfig) error {
	if len(cfgs) == 0 {
		return fmt.Errorf("configs are invalid")
	}

	hcfgs := []*health.Config{}
	for _, cfg := range cfgs {
		onComplete := func(state *health.State) {
			if cfg.OnComplete != nil {
				cfg.OnComplete(&HealthState{
					Name:               state.Name,
					Status:             state.Status,
					Err:                state.Err,
					Fatal:              state.Fatal,
					Details:            state.Details,
					CheckTime:          state.CheckTime,
					ContiguousFailures: state.ContiguousFailures,
					TimeOfFirstFailure: state.TimeOfFirstFailure,
				})
			}
		}
		hcfgs = append(hcfgs, &health.Config{
			Name:       cfg.Name,
			Checker:    cfg.Checker,
			Interval:   cfg.Interval,
			Fatal:      cfg.Fatal,
			OnComplete: onComplete,
		})
	}
	return c.h.AddChecks(hcfgs)
}

func (c *goHealthClient) Start() error {
	return c.h.Start()
}

func (c *goHealthClient) Stop() error {
	return c.h.Stop()
}

func (c *goHealthClient) State() (map[string]HealthState, bool, error) {
	states, success, err := c.h.State()
	hs := map[string]HealthState{}
	for _, state := range states {
		hs[state.Name] = HealthState{
			Name:               state.Name,
			Status:             state.Status,
			Err:                state.Err,
			Fatal:              state.Fatal,
			Details:            state.Details,
			CheckTime:          state.CheckTime,
			ContiguousFailures: state.ContiguousFailures,
			TimeOfFirstFailure: state.TimeOfFirstFailure,
		}
	}
	return hs, success, err
}

func NewGohealthClient(h *health.Health) (*goHealthClient, error) {
	if h == nil {
		return nil, fmt.Errorf("invalid health client")
	}
	return &goHealthClient{h}, nil
}
