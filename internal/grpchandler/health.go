package grpchandler

import (
	"context"

	"github.com/go-seidon/hippo/api/grpcapp"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
)

type healthHandler struct {
	grpcapp.UnimplementedHealthServiceServer
	healthClient healthcheck.HealthCheck
}

func (s *healthHandler) CheckHealth(ctx context.Context, p *grpcapp.CheckHealthParam) (*grpcapp.CheckHealthResult, error) {
	checkRes, err := s.healthClient.Check(ctx)
	if err != nil {
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	details := map[string]*grpcapp.CheckHealthDetail{}
	for _, item := range checkRes.Items {
		details[item.Name] = &grpcapp.CheckHealthDetail{
			Name:      item.Name,
			Status:    item.Status,
			CheckedAt: item.CheckedAt.UnixMilli(),
			Error:     item.Error,
		}
	}

	res := &grpcapp.CheckHealthResult{
		Code:    checkRes.Success.Code,
		Message: checkRes.Success.Message,
		Data: &grpcapp.CheckHealthData{
			Status:  checkRes.Status,
			Details: details,
		},
	}
	return res, nil
}

type HealthParam struct {
	HealthClient healthcheck.HealthCheck
}

func NewHealth(p HealthParam) *healthHandler {
	return &healthHandler{
		healthClient: p.HealthClient,
	}
}
