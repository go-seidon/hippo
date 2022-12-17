package file

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/status"
	"github.com/go-seidon/provider/system"
)

func (s *file) UploadFile(ctx context.Context, opts ...UploadFileOption) (*UploadFileResult, *system.Error) {
	s.log.Debug("In function: UploadFile")
	defer s.log.Debug("Returning function: UploadFile")

	p := UploadFileParam{}
	for _, opt := range opts {
		opt(&p)
	}

	err := s.validator.Validate(p)
	if err != nil {
		return nil, &system.Error{
			Code:    status.INVALID_PARAM,
			Message: err.Error(),
		}
	}

	if p.fileReader == nil {
		return nil, &system.Error{
			Code:    status.INVALID_PARAM,
			Message: "file is not specified",
		}
	}

	uploadDir := fmt.Sprintf("%s/%s", s.config.UploadDir, s.locator.GetLocation())

	exists, err := s.dirManager.IsDirectoryExists(ctx, filesystem.IsDirectoryExistsParam{
		Path: uploadDir,
	})
	if err != nil {
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	if !exists {
		_, err := s.dirManager.CreateDir(ctx, filesystem.CreateDirParam{
			Path:       uploadDir,
			Permission: 0644,
		})
		if err != nil {
			return nil, &system.Error{
				Code:    status.ACTION_FAILED,
				Message: err.Error(),
			}
		}
	}

	data := bytes.NewBuffer([]byte{})
	_, err = io.Copy(data, p.fileReader)
	if err != nil {
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	uniqueId, err := s.identifier.GenerateId()
	if err != nil {
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	path := fmt.Sprintf("%s/%s", uploadDir, uniqueId)
	if p.fileExtension != "" {
		path = fmt.Sprintf("%s.%s", path, p.fileExtension)
	}

	cRes, err := s.fileRepo.CreateFile(ctx, repository.CreateFileParam{
		UniqueId:  uniqueId,
		Path:      path,
		Name:      p.fileName,
		Mimetype:  p.fileMimetype,
		Extension: p.fileExtension,
		Size:      p.fileSize,
		CreateFn:  NewCreateFn(data.Bytes(), s.fileManager),
	})
	if err != nil {
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	res := &UploadFileResult{
		Success: system.Success{
			Code:    status.ACTION_SUCCESS,
			Message: "success upload file",
		},
		UniqueId:   cRes.UniqueId,
		Name:       cRes.Name,
		Path:       cRes.Path,
		Mimetype:   cRes.Mimetype,
		Extension:  cRes.Extension,
		Size:       cRes.Size,
		UploadedAt: cRes.CreatedAt,
	}
	return res, nil
}

func NewCreateFn(data []byte, fileManager filesystem.FileManager) repository.CreateFn {
	return func(ctx context.Context, cp repository.CreateFnParam) error {
		exists, err := fileManager.IsFileExists(ctx, filesystem.IsFileExistsParam{
			Path: cp.FilePath,
		})
		if err != nil {
			return err
		}
		if exists {
			return ErrExists
		}

		_, err = fileManager.SaveFile(ctx, filesystem.SaveFileParam{
			Name:       cp.FilePath,
			Data:       data,
			Permission: 0644,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
