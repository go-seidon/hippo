package file

import (
	"context"
	"errors"

	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/status"
	"github.com/go-seidon/provider/system"
)

func (s *file) DeleteFile(ctx context.Context, p DeleteFileParam) (*DeleteFileResult, *system.Error) {
	s.log.Debug("In function: DeleteFile")
	defer s.log.Debug("Returning function: DeleteFile")

	err := s.validator.Validate(p)
	if err != nil {
		return nil, &system.Error{
			Code:    status.INVALID_PARAM,
			Message: err.Error(),
		}
	}

	deleteion, err := s.fileRepo.DeleteFile(ctx, repository.DeleteFileParam{
		UniqueId: p.FileId,
		DeleteFn: NewDeleteFn(s.fileManager),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDeleted) {
			return nil, &system.Error{
				Code:    status.RESOURCE_NOTFOUND,
				Message: "file is deleted",
			}
		} else if errors.Is(err, repository.ErrNotFound) {
			return nil, &system.Error{
				Code:    status.RESOURCE_NOTFOUND,
				Message: "file is not found",
			}
		} else if errors.Is(err, ErrNotFound) {
			return nil, &system.Error{
				Code:    status.RESOURCE_NOTFOUND,
				Message: "file is not found",
			}
		}
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	res := &DeleteFileResult{
		Success: system.Success{
			Code:    status.ACTION_SUCCESS,
			Message: "success delete file",
		},
		DeletedAt: deleteion.DeletedAt,
	}
	return res, nil
}

func NewDeleteFn(fileManager filesystem.FileManager) repository.DeleteFn {
	return func(ctx context.Context, r repository.DeleteFnParam) error {
		exists, err := fileManager.IsFileExists(ctx, filesystem.IsFileExistsParam{
			Path: r.FilePath,
		})
		if err != nil {
			return err
		}

		if !exists {
			return ErrNotFound
		}

		_, err = fileManager.RemoveFile(ctx, filesystem.RemoveFileParam{
			Path: r.FilePath,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
