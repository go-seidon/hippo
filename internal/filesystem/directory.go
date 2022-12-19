package filesystem

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"time"
)

type DirectoryManager interface {
	IsDirectoryExists(ctx context.Context, p IsDirectoryExistsParam) (bool, error)
	CreateDir(ctx context.Context, p CreateDirParam) (*CreateDirResult, error)
}

type IsDirectoryExistsParam struct {
	Path string
}

type CreateDirParam struct {
	Path       string
	Permission fs.FileMode
}

type CreateDirResult struct {
	CreatedAt time.Time
}

type directoryManager struct {
}

func (dm *directoryManager) IsDirectoryExists(ctx context.Context, p IsDirectoryExistsParam) (bool, error) {
	_, err := os.Stat(p.Path)
	if err != nil {
		notExists := errors.Is(err, os.ErrNotExist)
		if notExists {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (dm *directoryManager) CreateDir(ctx context.Context, p CreateDirParam) (*CreateDirResult, error) {
	err := os.MkdirAll(p.Path, p.Permission)
	if err != nil {
		return nil, err
	}

	currentTs := time.Now()
	res := &CreateDirResult{
		CreatedAt: currentTs,
	}
	return res, nil
}

func NewDirectoryManager() *directoryManager {
	s := &directoryManager{}
	return s
}
