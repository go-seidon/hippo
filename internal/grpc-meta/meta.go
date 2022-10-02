package grpc_meta

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

type Metadata metadata.MD

func (m Metadata) Get(key string) string {
	k := strings.ToLower(key)
	v, ok := m[k]
	if !ok {
		return ""
	}
	return v[0]
}

// extracts an inbound metadata from the server-side context.
func ExtractIncoming(ctx context.Context) Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Metadata(metadata.Pairs())
	}
	return Metadata(md)
}

// extracts an outbound metadata from the client-side context.
func ExtractOutgoing(ctx context.Context) Metadata {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return Metadata(metadata.Pairs())
	}
	return Metadata(md)
}
