package grpc_app

import (
	"context"
	"errors"

	grpc_v1 "github.com/go-seidon/local/generated/proto/api/grpc/v1"
	"github.com/go-seidon/local/internal/file"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/status"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
)

type healthHandler struct {
	grpc_v1.UnimplementedHealthServiceServer
	healthService healthcheck.HealthCheck
}

func (s *healthHandler) CheckHealth(ctx context.Context, p *grpc_v1.CheckHealthParam) (*grpc_v1.CheckHealthResult, error) {
	checkRes, err := s.healthService.Check()
	if err != nil {
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	details := map[string]*grpc_v1.CheckHealthDetail{}
	for _, item := range checkRes.Items {
		details[item.Name] = &grpc_v1.CheckHealthDetail{
			Name:      item.Name,
			Status:    item.Status,
			CheckedAt: item.CheckedAt.UnixMilli(),
			Error:     item.Error,
		}
	}

	res := &grpc_v1.CheckHealthResult{
		Code:    status.ACTION_SUCCESS,
		Message: "success check service health",
		Data: &grpc_v1.CheckHealthData{
			Status:  checkRes.Status,
			Details: details,
		},
	}
	return res, nil
}

func NewHealthHandler(healthService healthcheck.HealthCheck) *healthHandler {
	return &healthHandler{
		healthService: healthService,
	}
}

type fileHandler struct {
	grpc_v1.UnimplementedFileServiceServer
	fileService file.File
}

func (h *fileHandler) DeleteFile(ctx context.Context, p *grpc_v1.DeleteFileParam) (*grpc_v1.DeleteFileResult, error) {
	deletion, err := h.fileService.DeleteFile(ctx, file.DeleteFileParam{
		FileId: p.FileId,
	})
	if err == nil {
		res := &grpc_v1.DeleteFileResult{
			Code:    status.ACTION_SUCCESS,
			Message: "success delete file",
			Data: &grpc_v1.DeleteFileData{
				DeletedAt: deletion.DeletedAt.UnixMilli(),
			},
		}
		return res, nil
	}

	if errors.Is(err, file.ErrorNotFound) {
		res := &grpc_v1.DeleteFileResult{
			Code:    status.RESOURCE_NOTFOUND,
			Message: err.Error(),
		}
		return res, nil
	}

	res := &grpc_v1.DeleteFileResult{
		Code:    status.ACTION_FAILED,
		Message: err.Error(),
	}
	return res, nil
}

func NewFileHandler(fileService file.File) *fileHandler {
	return &fileHandler{fileService: fileService}
}
