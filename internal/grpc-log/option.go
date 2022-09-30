package grpc_log

import (
	"github.com/go-seidon/local/internal/datetime"
	"github.com/go-seidon/local/internal/logging"
)

type LogInterceptorConfig struct {
	// required logger
	Logger logging.Logger

	// optional clock
	Clock datetime.Clock

	// key = method
	// value = set `true` to ignore the method being logged
	IgnoreMethod map[string]bool

	// key = metadata key
	// value = log key
	Metadata map[string]string

	// return true to specify the request should be logged
	ShouldLog ShouldLog

	// return log info based on request information
	CreateLog CreateLog
}

type LogInterceptorOption = func(*LogInterceptorConfig)

func WithLogger(logger logging.Logger) LogInterceptorOption {
	return func(cfg *LogInterceptorConfig) {
		cfg.Logger = logger
	}
}

func WithClock(clock datetime.Clock) LogInterceptorOption {
	return func(cfg *LogInterceptorConfig) {
		cfg.Clock = clock
	}
}

func WithIgnoredMethod(ims []string) LogInterceptorOption {
	return func(cfg *LogInterceptorConfig) {
		if len(ims) > 0 {
			im := map[string]bool{}
			for _, method := range ims {
				im[method] = true
			}
			cfg.IgnoreMethod = im
		}
	}
}

func WithAllowedMetadata(mds []string) LogInterceptorOption {
	return func(cfg *LogInterceptorConfig) {
		if len(mds) > 0 {
			md := map[string]string{}
			for _, meta := range mds {
				md[meta] = meta
			}
			cfg.Metadata = md
		}
	}
}
