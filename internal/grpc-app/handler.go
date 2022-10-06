package grpc_app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	grpc_v1 "github.com/go-seidon/local/generated/grpc-v1"
	"github.com/go-seidon/local/internal/file"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/status"
	"github.com/go-seidon/local/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpc_status "google.golang.org/grpc/status"
)

type healthHandler struct {
	grpc_v1.UnimplementedHealthServiceServer
	healthService healthcheck.HealthCheck
}

func (s *healthHandler) CheckHealth(ctx context.Context, p *grpc_v1.CheckHealthParam) (*grpc_v1.CheckHealthResult, error) {
	checkRes, err := s.healthService.Check(ctx)
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
	config      *GrpcAppConfig
}

func (h *fileHandler) DeleteFile(ctx context.Context, p *grpc_v1.DeleteFileParam) (*grpc_v1.DeleteFileResult, error) {
	deletion, err := h.fileService.DeleteFile(ctx, file.DeleteFileParam{
		FileId: p.FileId,
	})
	if err != nil {
		res := &grpc_v1.DeleteFileResult{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}

		if errors.Is(err, file.ErrorNotFound) {
			res.Code = status.RESOURCE_NOTFOUND
			res.Message = err.Error()
		}

		switch e := err.(type) {
		case *validation.ValidationError:
			res.Code = status.INVALID_PARAM
			res.Message = e.Error()
		}

		return res, nil
	}

	res := &grpc_v1.DeleteFileResult{
		Code:    status.ACTION_SUCCESS,
		Message: "success delete file",
		Data: &grpc_v1.DeleteFileData{
			DeletedAt: deletion.DeletedAt.UnixMilli(),
		},
	}
	return res, nil
}

func (h *fileHandler) RetrieveFile(p *grpc_v1.RetrieveFileParam, stream grpc_v1.FileService_RetrieveFileServer) error {
	retrieval, err := h.fileService.RetrieveFile(stream.Context(), file.RetrieveFileParam{
		FileId: p.FileId,
	})
	if err != nil {
		res := &grpc_v1.RetrieveFileResult{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
		if errors.Is(err, file.ErrorNotFound) {
			res.Code = status.RESOURCE_NOTFOUND
			res.Message = err.Error()
		}

		switch e := err.(type) {
		case *validation.ValidationError:
			res.Code = status.INVALID_PARAM
			res.Message = e.Error()
		}

		err = stream.Send(res)
		if err != nil {
			return err
		}
		return nil
	}

	err = stream.SendHeader(metadata.New(map[string]string{
		"file_name":      retrieval.Name,
		"file_mimetype":  retrieval.MimeType,
		"file_extension": retrieval.Extension,
		"file_size":      fmt.Sprintf("%d", retrieval.Size),
	}))
	if err != nil {
		return err
	}

	err = stream.Send(&grpc_v1.RetrieveFileResult{
		Code:    status.ACTION_PENDING,
		Message: "retrieving file",
	})
	if err != nil {
		return err
	}

	defer retrieval.Data.Close()

	const chunkSize = 102400 //100KB
	for {
		err = stream.Context().Err()
		if err != nil {
			err = stream.Send(&grpc_v1.RetrieveFileResult{
				Code:    status.ACTION_FAILED,
				Message: err.Error(),
			})
			if err != nil {
				return err
			}
			break
		}

		chunks := make([]byte, chunkSize)
		_, err := retrieval.Data.Read(chunks)
		if err == nil {
			err = stream.Send(&grpc_v1.RetrieveFileResult{
				Chunks: chunks,
			})
			if err != nil {
				return err
			}
			continue
		}

		if errors.Is(err, io.EOF) {
			err = stream.Send(&grpc_v1.RetrieveFileResult{
				Code:    status.ACTION_SUCCESS,
				Message: "success retrieve file",
			})
			if err != nil {
				return err
			}
			break
		}

		err = stream.Send(&grpc_v1.RetrieveFileResult{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		})
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func (h *fileHandler) UploadFile(stream grpc_v1.FileService_UploadFileServer) error {
	fileSize := int64(0)
	fileInfo := grpc_v1.UploadFileInfo{}
	fileReader := &bytes.Buffer{}

	for {
		err := stream.Context().Err()
		if err != nil {
			err = stream.SendAndClose(&grpc_v1.UploadFileResult{
				Code:    status.ACTION_FAILED,
				Message: err.Error(),
			})
			if err != nil {
				return err
			}
			return nil
		}

		param, err := stream.Recv()
		if err == nil {
			info := param.GetInfo()
			if info != nil {
				fileInfo = grpc_v1.UploadFileInfo{
					Name:      info.GetName(),
					Mimetype:  info.GetMimetype(),
					Extension: info.GetExtension(),
				}
			}

			chunks := param.GetChunks()
			if chunks != nil {
				fileSize += int64(len(chunks))

				if fileSize > h.config.UploadFormSize {
					err = stream.SendAndClose(&grpc_v1.UploadFileResult{
						Code:    status.ACTION_FAILED,
						Message: "file is too large",
					})
					if err != nil {
						return err
					}
					return nil
				}

				fileReader.Write(chunks)
			}
			continue
		}

		if errors.Is(err, io.EOF) {
			break
		}

		err = stream.SendAndClose(&grpc_v1.UploadFileResult{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		})
		if err != nil {
			return err
		}
		return nil
	}

	upload, err := h.fileService.UploadFile(
		stream.Context(),
		file.WithFileInfo(
			fileInfo.Name,
			fileInfo.Mimetype,
			fileInfo.Extension,
			fileSize,
		),
		file.WithReader(fileReader),
	)
	if err != nil {
		res := &grpc_v1.UploadFileResult{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}

		switch e := err.(type) {
		case *validation.ValidationError:
			res.Code = status.INVALID_PARAM
			res.Message = e.Error()
		}

		err = stream.SendAndClose(res)
		if err != nil {
			return err
		}
		return nil
	}

	err = stream.SendAndClose(&grpc_v1.UploadFileResult{
		Code:    status.ACTION_SUCCESS,
		Message: "success upload file",
		Data: &grpc_v1.UploadFileData{
			Id:         upload.UniqueId,
			Name:       upload.Name,
			Path:       upload.Path,
			Mimetype:   upload.Mimetype,
			Extension:  upload.Extension,
			Size:       upload.Size,
			UploadedAt: upload.UploadedAt.UnixMilli(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func NewFileHandler(fileService file.File, config *GrpcAppConfig) *fileHandler {
	return &fileHandler{
		fileService: fileService,
		config:      config,
	}
}
