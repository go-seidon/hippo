package grpc_log

import (
	"context"

	"github.com/go-seidon/local/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type ServerStream interface {
	SetHeader(metadata.MD) error
	SendHeader(metadata.MD) error
	SetTrailer(metadata.MD)
	Context() context.Context
	SendMsg(m interface{}) error
	RecvMsg(m interface{}) error
}

type logServerStream struct {
	grpc.ServerStream
	logger logging.Logger
}

func (ss *logServerStream) SendMsg(m interface{}) error {
	err := ss.ServerStream.SendMsg(m)
	if err == nil {
		var stream *messageMarshaller
		msg, ok := m.(proto.Message)
		if ok {
			stream = NewMessage(msg)
		}

		logger := ss.logger.WithFields(map[string]interface{}{
			"stream": stream,
		})
		logger.Info("send stream")
	}
	return err
}

func (ss *logServerStream) RecvMsg(m interface{}) error {
	err := ss.ServerStream.RecvMsg(m)
	if err == nil {
		var stream *messageMarshaller
		msg, ok := m.(proto.Message)
		if ok {
			stream = NewMessage(msg)
		}

		logger := ss.logger.WithFields(map[string]interface{}{
			"stream": stream,
		})
		logger.Info("receive stream")
	}
	return err
}

func NewLogServerStream(ss grpc.ServerStream, logger logging.Logger) *logServerStream {
	return &logServerStream{ss, logger}
}
