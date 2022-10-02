package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type ServerStream interface {
	SetHeader(metadata.MD) error
	SendHeader(metadata.MD) error
	SetTrailer(metadata.MD)
	Context() context.Context
	SendMsg(m interface{}) error
	RecvMsg(m interface{}) error
}
