package file

import (
	"context"
	"fmt"
	"io"

	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/repository"
)

func (s *file) UploadFile(ctx context.Context, opts ...UploadFileOption) (*UploadFileResult, error) {
	s.log.Debug("In function: UploadFile")
	defer s.log.Debug("Returning function: UploadFile")

	p := UploadFileParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if p.fileData == nil && p.fileReader == nil {
		return nil, fmt.Errorf("invalid file is not specified")
	}
	if p.fileDir == "" {
		return nil, fmt.Errorf("invalid upload directory is not specified")
	}

	exists, err := s.dirManager.IsDirectoryExists(ctx, filesystem.IsDirectoryExistsParam{
		Path: p.fileDir,
	})
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err := s.dirManager.CreateDir(ctx, filesystem.CreateDirParam{
			Path:       p.fileDir,
			Permission: 0644,
		})
		if err != nil {
			return nil, err
		}
	}

	data := p.fileData
	if p.fileReader != nil {
		byte, err := io.ReadAll(p.fileReader)
		if err != nil {
			return nil, err
		}
		data = byte
	}

	uniqueId, err := s.identifier.GenerateId()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", p.fileDir, uniqueId)
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
		CreateFn:  NewCreateFn(data, s.fileManager),
	})
	if err != nil {
		return nil, err
	}

	res := &UploadFileResult{
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

type UploadFileOption = func(*UploadFileParam)

func WithData(d []byte) UploadFileOption {
	return func(ufp *UploadFileParam) {
		ufp.fileData = d
		ufp.fileReader = nil
	}
}

func WithReader(w io.Reader) UploadFileOption {
	return func(ufp *UploadFileParam) {
		ufp.fileReader = w
		ufp.fileData = nil
	}
}

func WithDirectory(path string) UploadFileOption {
	return func(ufp *UploadFileParam) {
		ufp.fileDir = path
	}
}

func WithFileInfo(name, mimetype, extension string, size int64) UploadFileOption {
	return func(ufp *UploadFileParam) {
		ufp.fileName = name
		ufp.fileMimetype = mimetype
		ufp.fileExtension = extension
		ufp.fileSize = size
	}
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
			return ErrorExists
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
