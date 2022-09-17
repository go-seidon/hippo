package healthcheck

import "github.com/InVisionApp/go-health"

type goHealthClient struct {
	h *health.Health
}

func (c *goHealthClient) AddChecks(cfgs []*HealthConfig) error {
	hcfgs := []*health.Config{}
	for _, cfg := range cfgs {
		hcfgs = append(hcfgs, &health.Config{
			Name:     cfg.Name,
			Checker:  cfg.Checker,
			Interval: cfg.Interval,
			Fatal:    cfg.Fatal,
			OnComplete: func(state *health.State) {
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
			},
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
