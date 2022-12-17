package grpcapp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-seidon/hippo/api/grpcapp"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/provider/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func NewHealthHandler(healthClient healthcheck.HealthCheck) *healthHandler {
	return &healthHandler{
		healthClient: healthClient,
	}
}

type fileHandler struct {
	grpcapp.UnimplementedFileServiceServer
	fileService file.File
	config      *GrpcAppConfig
}

func (h *fileHandler) DeleteFileById(ctx context.Context, p *grpcapp.DeleteFileByIdParam) (*grpcapp.DeleteFileByIdResult, error) {
	deletion, err := h.fileService.DeleteFile(ctx, file.DeleteFileParam{
		FileId: p.FileId,
	})
	if err != nil {
		res := &grpcapp.DeleteFileByIdResult{
			Code:    err.Code,
			Message: err.Message,
		}
		return res, nil
	}

	res := &grpcapp.DeleteFileByIdResult{
		Code:    deletion.Success.Code,
		Message: deletion.Success.Message,
		Data: &grpcapp.DeleteFileByIdData{
			DeletedAt: deletion.DeletedAt.UnixMilli(),
		},
	}
	return res, nil
}

func (h *fileHandler) RetrieveFileById(p *grpcapp.RetrieveFileByIdParam, stream grpcapp.FileService_RetrieveFileByIdServer) error {
	retrieval, rerr := h.fileService.RetrieveFile(stream.Context(), file.RetrieveFileParam{
		FileId: p.FileId,
	})
	if rerr != nil {
		res := &grpcapp.RetrieveFileByIdResult{
			Code:    rerr.Code,
			Message: rerr.Message,
		}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		return nil
	}

	err := stream.SendHeader(metadata.New(map[string]string{
		"file_name":      retrieval.Name,
		"file_mimetype":  retrieval.MimeType,
		"file_extension": retrieval.Extension,
		"file_size":      fmt.Sprintf("%d", retrieval.Size),
	}))
	if err != nil {
		return err
	}

	err = stream.Send(&grpcapp.RetrieveFileByIdResult{
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
			err = stream.Send(&grpcapp.RetrieveFileByIdResult{
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
			err = stream.Send(&grpcapp.RetrieveFileByIdResult{
				Chunks: chunks,
			})
			if err != nil {
				return err
			}
			continue
		}

		if errors.Is(err, io.EOF) {
			err = stream.Send(&grpcapp.RetrieveFileByIdResult{
				Code:    retrieval.Success.Code,
				Message: retrieval.Success.Message,
			})
			if err != nil {
				return err
			}
			break
		}

		err = stream.Send(&grpcapp.RetrieveFileByIdResult{
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

func (h *fileHandler) UploadFile(stream grpcapp.FileService_UploadFileServer) error {
	fileSize := int64(0)
	fileInfo := grpcapp.UploadFileInfo{}
	fileReader := &bytes.Buffer{}

	for {
		err := stream.Context().Err()
		if err != nil {
			err = stream.SendAndClose(&grpcapp.UploadFileResult{
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
				fileInfo = grpcapp.UploadFileInfo{
					Name:      info.GetName(),
					Mimetype:  info.GetMimetype(),
					Extension: info.GetExtension(),
				}
			}

			chunks := param.GetChunks()
			if chunks != nil {
				fileSize += int64(len(chunks))

				if fileSize > h.config.UploadFormSize {
					err = stream.SendAndClose(&grpcapp.UploadFileResult{
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

		err = stream.SendAndClose(&grpcapp.UploadFileResult{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		})
		if err != nil {
			return err
		}
		return nil
	}

	upload, uerr := h.fileService.UploadFile(
		stream.Context(),
		file.WithFileInfo(
			fileInfo.Name,
			fileInfo.Mimetype,
			fileInfo.Extension,
			fileSize,
		),
		file.WithReader(fileReader),
	)
	if uerr != nil {
		res := &grpcapp.UploadFileResult{
			Code:    uerr.Code,
			Message: uerr.Message,
		}

		err := stream.SendAndClose(res)
		if err != nil {
			return err
		}
		return nil
	}

	err := stream.SendAndClose(&grpcapp.UploadFileResult{
		Code:    upload.Success.Code,
		Message: upload.Success.Message,
		Data: &grpcapp.UploadFileData{
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
