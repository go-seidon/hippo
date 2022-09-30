package grpc_log

import (
	"context"
	"path"
	"time"

	"github.com/go-seidon/local/internal/datetime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var DefaultClock = datetime.NewClock()

type ShouldLog = func(ctx context.Context, p ShouldLogParam) bool

type ShouldLogParam struct {
	Method       string
	Error        error
	IgnoreMethod map[string]bool
}

var DefaultShouldLog = func(ctx context.Context, p ShouldLogParam) bool {
	return !p.IgnoreMethod[p.Method]
}

type CreateLog = func(ctx context.Context, p CreateLogParam) *LogInfo

type CreateLogParam struct {
	Method    string
	Error     error
	StartTime time.Time
	Metadata  map[string]string
}

type LogInfo struct {
	Service       string
	Method        string
	Status        string
	Level         string
	ReceivedAt    time.Time
	Duration      int64
	RemoteAddress string
	Protocol      string
	Metadata      map[string]interface{}
}

var DefaultCreateLog = func(ctx context.Context, p CreateLogParam) *LogInfo {

	timeElapsed := time.Since(p.StartTime)
	duration := int64(timeElapsed) / int64(time.Millisecond)
	service := path.Dir(p.Method)[1:]
	method := path.Base(p.Method)

	code := status.Code(p.Error)
	status := code.String()
	receivedAt := p.StartTime

	remoteAddr := ""
	protocol := ""
	peer, ok := peer.FromContext(ctx)
	if ok {
		remoteAddr = peer.Addr.String()
		protocol = peer.Addr.Network()
	}

	meta := map[string]interface{}{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for mdKey, logKey := range p.Metadata {
			values := md.Get(mdKey)
			if len(values) > 0 {
				meta[logKey] = values[0]
			}
		}
	}

	level := "error"
	switch code {
	case
		codes.OK, codes.Canceled, codes.InvalidArgument,
		codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
		level = "info"
	case
		codes.Unknown, codes.DeadlineExceeded,
		codes.Unimplemented, codes.Internal, codes.DataLoss:
		level = "error"
	case
		codes.PermissionDenied, codes.ResourceExhausted,
		codes.FailedPrecondition, codes.Aborted,
		codes.OutOfRange, codes.Unavailable:
		level = "warn"
	}

	return &LogInfo{
		Service:       service,
		Method:        method,
		Status:        status,
		ReceivedAt:    receivedAt,
		Duration:      duration,
		RemoteAddress: remoteAddr,
		Protocol:      protocol,
		Metadata:      meta,
		Level:         level,
	}
}