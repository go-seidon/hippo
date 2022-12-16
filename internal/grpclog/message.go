package grpclog

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Message interface {
	MarshalJSON() ([]byte, error)
}

type messageMarshaller struct {
	proto.Message
}

func (m *messageMarshaller) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(m.Message)
}

func NewMessage(m proto.Message) *messageMarshaller {
	return &messageMarshaller{m}
}
