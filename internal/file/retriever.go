package file

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/repository"
)

func (s *file) RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, error) {
	s.log.Debug("In function: RetrieveFile")
	defer s.log.Debug("Returning function: RetrieveFile")

	if p.FileId == "" {
		return nil, fmt.Errorf("invalid file id parameter")
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
	}

	return res, nil
}