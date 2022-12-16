package grpclog

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...LogInterceptorOption) grpc.UnaryServerInterceptor {
	cfg := buildConfig(opts...)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := cfg.Clock.Now()

		var deadlineAt *time.Time
		dlTime, dlOccured := ctx.Deadline()
		if dlOccured {
			deadlineAt = &dlTime
		}

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
			Request:   req,
			Response:  res,
		})

		cfg.SendLog(ctx, SendLogParam{
			Logger:     cfg.Logger,
			LogInfo:    *logInfo,
			Error:      err,
			DeadlineAt: deadlineAt,
		})

		return res, err
	}
}

func StreamServerInterceptor(opts ...LogInterceptorOption) grpc.StreamServerInterceptor {
	cfg := buildConfig(opts...)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := cfg.Clock.Now()

		ctx := ss.Context()
		var deadlineAt *time.Time
		dlTime, dlOccured := ctx.Deadline()
		if dlOccured {
			deadlineAt = &dlTime
		}

		lss := NewLogServerStream(ss, cfg.Logger)
		err := handler(srv, lss)

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

		cfg.SendLog(ctx, SendLogParam{
			Logger:     cfg.Logger,
			LogInfo:    *logInfo,
			Error:      err,
			DeadlineAt: deadlineAt,
		})

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
		SendLog:      DefaultSendLog,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
