package file

import (
	"context"
	"errors"

	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/status"
	"github.com/go-seidon/provider/system"
)

func (s *file) RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, *system.Error) {
	s.log.Debug("In function: RetrieveFile")
	defer s.log.Debug("Returning function: RetrieveFile")

	err := s.validator.Validate(p)
	if err != nil {
		return nil, &system.Error{
			Code:    status.INVALID_PARAM,
			Message: err.Error(),
		}
	}

	retrieve, err := s.fileRepo.RetrieveFile(ctx, repository.RetrieveFileParam{
		UniqueId: p.FileId,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &system.Error{
				Code:    status.RESOURCE_NOTFOUND,
				Message: "file is not found",
			}
		} else if errors.Is(err, repository.ErrDeleted) {
			return nil, &system.Error{
				Code:    status.RESOURCE_NOTFOUND,
				Message: "file is deleted",
			}
		}
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	open, err := s.fileManager.OpenFile(ctx, filesystem.OpenFileParam{
		Path: retrieve.Path,
	})
	if err != nil {
		if errors.Is(err, filesystem.ErrorFileNotFound) {
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

	res := &RetrieveFileResult{
		Success: system.Success{
			Code:    status.ACTION_SUCCESS,
			Message: "success retrieve file",
		},
		Data:      open.File,
		UniqueId:  retrieve.UniqueId,
		Name:      retrieve.Name,
		Path:      retrieve.Path,
		MimeType:  retrieve.MimeType,
		Extension: retrieve.Extension,
		Size:      retrieve.Size,
	}
	return res, nil
}
