package file

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-seidon/local/internal/filesystem"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/repository"
	"github.com/go-seidon/local/internal/text"
	"github.com/go-seidon/local/internal/validation"
)

type File interface {
	UploadFile(ctx context.Context, opts ...UploadFileOption) (*UploadFileResult, error)
	RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, error)
	DeleteFile(ctx context.Context, p DeleteFileParam) (*DeleteFileResult, error)
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

func WithFileInfo(name, mimetype, extension string, size int64) UploadFileOption {
	return func(ufp *UploadFileParam) {
		ufp.fileName = name
		ufp.fileMimetype = mimetype
		ufp.fileExtension = extension
		ufp.fileSize = size
	}
}

type UploadFileParam struct {
	fileData      []byte
	fileReader    io.Reader
	fileName      string `validate:"max=4096"`
	fileMimetype  string `validate:"max=256"`
	fileExtension string `validate:"max=128"`
	fileSize      int64  `validate:"min=0"`
}

type UploadFileResult struct {
	UniqueId   string
	Name       string
	Path       string
	Mimetype   string
	Extension  string
	Size       int64
	UploadedAt time.Time
}

type RetrieveFileParam struct {
	FileId string `validate:"required,min=5,max=64"`
}

type RetrieveFileResult struct {
	Data      io.ReadCloser
	UniqueId  string
	Name      string
	Path      string
	MimeType  string
	Extension string
	DeletedAt *int64
}

type DeleteFileParam struct {
	FileId string `validate:"required,min=5,max=64"`
}

type DeleteFileResult struct {
	DeletedAt time.Time
}

type file struct {
	fileRepo    repository.FileRepository
	fileManager filesystem.FileManager
	dirManager  filesystem.DirectoryManager
	identifier  text.Identifier
	log         logging.Logger
	locator     UploadLocation
	validator   validation.Validator
	config      *FileConfig
}

type FileConfig struct {
	UploadDir string
}

type NewFileParam struct {
	FileRepo    repository.FileRepository
	FileManager filesystem.FileManager
	DirManager  filesystem.DirectoryManager
	Logger      logging.Logger
	Identifier  text.Identifier
	Locator     UploadLocation
	Validator   validation.Validator
	Config      *FileConfig
}

func NewFile(p NewFileParam) (*file, error) {
	if p.FileRepo == nil {
		return nil, fmt.Errorf("file repo is not specified")
	}
	if p.FileManager == nil {
		return nil, fmt.Errorf("file manager is not specified")
	}
	if p.DirManager == nil {
		return nil, fmt.Errorf("directory manager is not specified")
	}
	if p.Logger == nil {
		return nil, fmt.Errorf("logger is not specified")
	}
	if p.Identifier == nil {
		return nil, fmt.Errorf("identifier is not specified")
	}
	if p.Config == nil {
		return nil, fmt.Errorf("config is not specified")
	}
	if p.Locator == nil {
		return nil, fmt.Errorf("locator is not specified")
	}
	if p.Validator == nil {
		return nil, fmt.Errorf("validator is not specified")
	}

	s := &file{
		fileRepo:    p.FileRepo,
		fileManager: p.FileManager,
		dirManager:  p.DirManager,
		identifier:  p.Identifier,
		log:         p.Logger,
		locator:     p.Locator,
		config:      p.Config,
		validator:   p.Validator,
	}
	return s, nil
}
