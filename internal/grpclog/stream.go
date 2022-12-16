package grpclog

import (
	"github.com/go-seidon/provider/logging"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

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
			"grpcStream": stream,
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
			"grpcStream": stream,
		})
		logger.Info("receive stream")
	}
	return err
}

func NewLogServerStream(ss grpc.ServerStream, logger logging.Logger) *logServerStream {
	return &logServerStream{ss, logger}
}
