package file

import (
	"context"
	"errors"

	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/repository"
)

func (s *file) RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, error) {
	s.log.Debug("In function: RetrieveFile")
	defer s.log.Debug("Returning function: RetrieveFile")

	err := s.validator.Validate(p)
	if err != nil {
		return nil, err
	}

	file, err := s.fileRepo.RetrieveFile(ctx, repository.RetrieveFileParam{
		UniqueId: p.FileId,
	})
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return nil, ErrorNotFound
		}
		return nil, err
	}

	oRes, err := s.fileManager.OpenFile(ctx, filesystem.OpenFileParam{
		Path: file.Path,
	})
	if err != nil {
		if errors.Is(err, filesystem.ErrorFileNotFound) {
			return nil, ErrorNotFound
		}
		return nil, err
	}

	res := &RetrieveFileResult{
		Data:      oRes.File,
		UniqueId:  file.UniqueId,
		Name:      file.Name,
		Path:      file.Path,
		MimeType:  file.MimeType,
		Extension: file.Extension,
		Size:      file.Size,
	}
	return res, nil
}
