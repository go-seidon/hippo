package grpc_auth

import "context"

type AuthInterceptorConfig struct {
	CheckCredential CheckCredential
}

type AuthInterceptorOption = func(*AuthInterceptorConfig)

type CheckCredential = func(ctx context.Context) error

func WithAuth(cc CheckCredential) AuthInterceptorOption {
	return func(cfg *AuthInterceptorConfig) {
		cfg.CheckCredential = cc
	}
}
