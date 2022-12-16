package grpclog

import (
	"context"
	"time"

	"github.com/go-seidon/provider/datetime"
	"github.com/go-seidon/provider/logging"
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

	// send log using the specified logger
	SendLog SendLog
}

type ShouldLog = func(ctx context.Context, p ShouldLogParam) bool

type ShouldLogParam struct {
	Method       string
	Error        error
	IgnoreMethod map[string]bool
}

type CreateLog = func(ctx context.Context, p CreateLogParam) *LogInfo

type CreateLogParam struct {
	Method    string
	Error     error
	StartTime time.Time
	Metadata  map[string]string
	Request   interface{}
	Response  interface{}
}

type LogInfo struct {
	Service       string
	Method        string
	Status        string
	Level         string
	ReceivedAt    time.Time
	Duration      int64
	RemoteAddress string
	Protocol      string
	Metadata      map[string]interface{}
	Request       Message
	Response      Message
}

type SendLog = func(ctx context.Context, p SendLogParam) error

type SendLogParam struct {
	Logger     logging.Logger
	LogInfo    LogInfo
	Error      error
	DeadlineAt *time.Time
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

func IgnoredMethod(ims []string) LogInterceptorOption {
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

func AllowedMetadata(mds []string) LogInterceptorOption {
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
