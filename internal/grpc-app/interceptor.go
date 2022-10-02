package grpc_app

import (
	"context"

	"github.com/go-seidon/local/internal/auth"
	grpc_auth "github.com/go-seidon/local/internal/grpc-auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BasicAuth(basicAuth auth.BasicAuth) grpc_auth.CheckCredential {
	return func(ctx context.Context) error {
		token, err := grpc_auth.AuthFromMD(ctx, grpc_auth.BasicKey)
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
			return status.Errorf(codes.Unauthenticated, grpc_auth.ErrorInvalidCredential.Error())
		}
		return nil
	}
}
