package grpcauth

import (
	"context"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...AuthInterceptorOption) grpc.UnaryServerInterceptor {
	cfg := buildConfig(opts...)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := cfg.CheckCredential(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(opts ...AuthInterceptorOption) grpc.StreamServerInterceptor {
	cfg := buildConfig(opts...)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := cfg.CheckCredential(ss.Context())
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func buildConfig(opts ...AuthInterceptorOption) *AuthInterceptorConfig {
	cfg := &AuthInterceptorConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
