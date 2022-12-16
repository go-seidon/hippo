package grpcapp

import (
	"context"

	"github.com/go-seidon/hippo/internal/auth"
	"github.com/go-seidon/hippo/internal/grpcauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BasicAuth(basicAuth auth.BasicAuth) grpcauth.CheckCredential {
	return func(ctx context.Context) error {
		token, err := grpcauth.AuthFromMD(ctx, grpcauth.BasicKey)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, err.Error())
		}

		res, err := basicAuth.CheckCredential(ctx, auth.CheckCredentialParam{
			AuthToken: token,
		})
		if err != nil {
			return status.Errorf(codes.Unknown, err.Error())
		}

		if !res.IsValid() {
			return status.Errorf(codes.Unauthenticated, grpcauth.ErrorInvalidCredential.Error())
		}
		return nil
	}
}
