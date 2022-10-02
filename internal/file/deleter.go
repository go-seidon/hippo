package file

import (
	"context"
	"errors"

	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/repository"
)

func (s *file) DeleteFile(ctx context.Context, p DeleteFileParam) (*DeleteFileResult, error) {
	s.log.Debug("In function: DeleteFile")
	defer s.log.Debug("Returning function: DeleteFile")

	err := s.validator.Validate(p)
	if err != nil {
		return nil, err
	}

	delRes, err := s.fileRepo.DeleteFile(ctx, repository.DeleteFileParam{
		UniqueId: p.FileId,
		DeleteFn: NewDeleteFn(s.fileManager),
	})

	if err == nil {
		res := &DeleteFileResult{
			DeletedAt: delRes.DeletedAt,
		}
		return res, nil
	}

	if errors.Is(err, repository.ErrorRecordNotFound) {
		return nil, ErrorNotFound
	}
	return nil, err
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
			return ErrorNotFound
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
