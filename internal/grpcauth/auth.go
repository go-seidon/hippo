package grpcauth

import (
	"context"
	"fmt"
	"strings"

	grpc_meta "github.com/go-seidon/hippo/internal/grpc-meta"
)

const (
	AuthKey   = "authorization"
	BasicKey  = "basic"
	BearerKey = "bearer"
)

func AuthFromMD(ctx context.Context, scheme string) (string, error) {
	val := grpc_meta.ExtractIncoming(ctx).Get(AuthKey)
	if val == "" {
		return "", fmt.Errorf("request unauthenticated with %s", scheme)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", fmt.Errorf("bad authorization string")
	}
	if !strings.EqualFold(splits[0], scheme) {
		return "", fmt.Errorf("invalid scheme of %s", scheme)
	}
	return splits[1], nil
}
