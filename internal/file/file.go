package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-seidon/hippo/internal/filesystem"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/identity"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/status"
	"github.com/go-seidon/provider/system"
	"github.com/go-seidon/provider/validation"
)

type File interface {
	UploadFile(ctx context.Context, opts ...UploadFileOption) (*UploadFileResult, *system.Error)
	RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, *system.Error)
	DeleteFile(ctx context.Context, p DeleteFileParam) (*DeleteFileResult, *system.Error)
}

type UploadFileOption = func(*UploadFileParam)

func WithReader(w io.Reader) UploadFileOption {
	return func(p *UploadFileParam) {
		p.fileReader = w
	}
}

func WithFileInfo(name, mimetype, extension string, size int64) UploadFileOption {
	return func(p *UploadFileParam) {
		p.fileName = name
		p.fileMimetype = mimetype
		p.fileExtension = extension
		p.fileSize = size
	}
}

type UploadFileParam struct {
	fileReader    io.Reader
	fileName      string `validate:"max=4096" label:"name"`
	fileMimetype  string `validate:"max=256" label:"mimetype"`
	fileExtension string `validate:"max=128" label:"extension"`
	fileSize      int64  `validate:"min=0" label:"size"`
}

type UploadFileResult struct {
	Success    system.Success
	UniqueId   string
	Name       string
	Path       string
	Mimetype   string
	Extension  string
	Size       int64
	UploadedAt time.Time
}

type RetrieveFileParam struct {
	FileId string `validate:"required,min=5,max=64" label:"file_id"`
}

type RetrieveFileResult struct {
	Success   system.Success
	Data      io.ReadCloser
	UniqueId  string
	Name      string
	Path      string
	MimeType  string
	Extension string
	Size      int64
}

type DeleteFileParam struct {
	FileId string `validate:"required,min=5,max=64" label:"file_id"`
}

type DeleteFileResult struct {
	Success   system.Success
	DeletedAt time.Time
}

type file struct {
	fileRepo    repository.File
	fileManager filesystem.FileManager
	dirManager  filesystem.DirectoryManager
	identifier  identity.Identifier
	log         logging.Logger
	locator     UploadLocation
	validator   validation.Validator
	config      *FileConfig
}

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
		}
		return nil, &system.Error{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		}
	}

	if retrieve.DeletedAt != nil {
		return nil, &system.Error{
			Code:    status.RESOURCE_NOTFOUND,
			Message: "file is deleted",
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
		MimeType:  retrieve.Mimetype,
		Extension: retrieve.Extension,
		Size:      retrieve.Size,
	}
	return res, nil
}

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

	deletion, err := s.fileRepo.DeleteFile(ctx, repository.DeleteFileParam{
		UniqueId: p.FileId,
		DeleteFn: NewDeleteFn(s.fileManager),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDeleted) {
			return nil, &system.Error{
				Code:    status.RESOURCE_NOTFOUND,
				Message: "file is deleted",
			}
		} else if errors.Is(err, repository.ErrNotFound) || errors.Is(err, ErrNotFound) {
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
		DeletedAt: deletion.DeletedAt,
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

type FileConfig struct {
	UploadDir string
}

type FileParam struct {
	FileRepo    repository.File
	FileManager filesystem.FileManager
	DirManager  filesystem.DirectoryManager
	Logger      logging.Logger
	Identifier  identity.Identifier
	Locator     UploadLocation
	Validator   validation.Validator
	Config      *FileConfig
}

func NewFile(p FileParam) *file {
	return &file{
		fileRepo:    p.FileRepo,
		fileManager: p.FileManager,
		dirManager:  p.DirManager,
		identifier:  p.Identifier,
		log:         p.Logger,
		locator:     p.Locator,
		config:      p.Config,
		validator:   p.Validator,
	}
}
