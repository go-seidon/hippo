package healthcheck

import (
	"time"
)

type HealthClient interface {
	AddChecks(cfgs []*HealthConfig) error
	Start() error
	Stop() error
	State() (map[string]HealthState, bool, error)
}

type HealthConfig struct {
	Name       string
	Checker    Checker
	Interval   time.Duration
	Fatal      bool
	OnComplete func(state *HealthState)
}

type HealthState struct {
	Name               string
	Status             string
	Err                string
	Fatal              bool
	Details            interface{}
	CheckTime          time.Time
	ContiguousFailures int64
	TimeOfFirstFailure time.Time
}
