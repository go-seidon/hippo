package grpc_log

import (
	"context"
	"time"

	"github.com/go-seidon/local/internal/logging"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...LogInterceptorOption) grpc.UnaryServerInterceptor {
	cfg := buildConfig(opts...)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := cfg.Clock.Now()
		dlTime, dlOccured := ctx.Deadline()

		res, err := handler(ctx, req)

		shouldLog := cfg.ShouldLog(ctx, ShouldLogParam{
			Method:       info.FullMethod,
			Error:        err,
			IgnoreMethod: cfg.IgnoreMethod,
		})
		if !shouldLog {
			return res, err
		}

		logInfo := cfg.CreateLog(ctx, CreateLogParam{
			Method:    info.FullMethod,
			Error:     err,
			StartTime: startTime,
			Metadata:  cfg.Metadata,
		})

		grpcRequest := map[string]interface{}{
			"requestService": logInfo.Service,
			"requestMethod":  logInfo.Method,
			"status":         logInfo.Status,
			"receivedAt":     logInfo.ReceivedAt.UTC().Format(time.RFC3339),
			"duration":       logInfo.Duration,
			"remoteAddr":     logInfo.RemoteAddress,
			"protocol":       logInfo.Protocol,
		}
		if dlOccured {
			grpcRequest["deadlineAt"] = dlTime.UTC().Format(time.RFC3339)
		}

		logger := cfg.Logger
		if err != nil {
			logger = logger.WithFields(map[string]interface{}{
				logging.FIELD_ERROR: err,
			})
		}
		logger = logger.WithFields(map[string]interface{}{
			"grpcRequest": grpcRequest,
		})
		logger.Logf(logInfo.Level, "request: %s@%s", logInfo.Service, logInfo.Method)

		return res, err
	}
}

func StreamServerInterceptor(opts ...LogInterceptorOption) grpc.StreamServerInterceptor {
	cfg := buildConfig(opts...)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := cfg.Clock.Now()

		ctx := ss.Context()
		dlTime, dlOccured := ctx.Deadline()

		err := handler(srv, ss)

		shouldLog := cfg.ShouldLog(ctx, ShouldLogParam{
			Method:       info.FullMethod,
			Error:        err,
			IgnoreMethod: cfg.IgnoreMethod,
		})
		if !shouldLog {
			return err
		}

		logInfo := cfg.CreateLog(ctx, CreateLogParam{
			Method:    info.FullMethod,
			Error:     err,
			StartTime: startTime,
			Metadata:  cfg.Metadata,
		})

		grpcRequest := map[string]interface{}{
			"requestService": logInfo.Service,
			"requestMethod":  logInfo.Method,
			"status":         logInfo.Status,
			"receivedAt":     logInfo.ReceivedAt.UTC().Format(time.RFC3339),
			"duration":       logInfo.Duration,
			"remoteAddr":     logInfo.RemoteAddress,
			"protocol":       logInfo.Protocol,
		}
		if dlOccured {
			grpcRequest["deadlineAt"] = dlTime.UTC().Format(time.RFC3339)
		}

		logger := cfg.Logger
		if err != nil {
			logger = logger.WithFields(map[string]interface{}{
				logging.FIELD_ERROR: err,
			})
		}
		logger = logger.WithFields(map[string]interface{}{
			"grpcRequest": grpcRequest,
		})
		logger.Logf(logInfo.Level, "request: %s@%s", logInfo.Service, logInfo.Method)

		return err
	}
}

func buildConfig(opts ...LogInterceptorOption) *LogInterceptorConfig {
	cfg := &LogInterceptorConfig{
		Clock:        DefaultClock,
		IgnoreMethod: map[string]bool{},
		Metadata:     map[string]string{},
		ShouldLog:    DefaultShouldLog,
		CreateLog:    DefaultCreateLog,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
