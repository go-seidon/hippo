package filesystem

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"time"
)

type FileManager interface {
	IsFileExists(ctx context.Context, p IsFileExistsParam) (bool, error)
	OpenFile(ctx context.Context, p OpenFileParam) (*OpenFileResult, error)
	SaveFile(ctx context.Context, p SaveFileParam) (*SaveFileResult, error)
	RemoveFile(ctx context.Context, p RemoveFileParam) (*RemoveFileResult, error)
}

type IsFileExistsParam struct {
	Path string
}

type OpenFileParam struct {
	Path string
}

type OpenFileResult struct {
	File *os.File
}

type SaveFileParam struct {
	Name       string
	Data       []byte
	Permission fs.FileMode
}

type SaveFileResult struct {
	SavedAt time.Time
}

type RemoveFileParam struct {
	Path string
}

type RemoveFileResult struct {
	RemovedAt time.Time
}

type fileManager struct {
}

func (fm *fileManager) IsFileExists(ctx context.Context, p IsFileExistsParam) (bool, error) {
	_, err := os.Stat(p.Path)
	if err == nil {
		return true, nil
	}

	notExists := errors.Is(err, os.ErrNotExist)
	if notExists {
		return false, nil
	}

	return false, err
}

func (fm *fileManager) OpenFile(ctx context.Context, p OpenFileParam) (*OpenFileResult, error) {
	file, err := os.Open(p.Path)
	if err == nil {
		res := &OpenFileResult{
			File: file,
		}
		return res, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrorFileNotFound
	}

	return nil, err
}

// @note: save file/overwrite if exists
func (fm *fileManager) SaveFile(ctx context.Context, p SaveFileParam) (*SaveFileResult, error) {
	err := os.WriteFile(p.Name, p.Data, p.Permission)
	if err != nil {
		return nil, err
	}

	currentTimestamp := time.Now()
	res := &SaveFileResult{
		SavedAt: currentTimestamp,
	}
	return res, nil
}

func (fm *fileManager) RemoveFile(ctx context.Context, p RemoveFileParam) (*RemoveFileResult, error) {
	err := os.Remove(p.Path)
	if err == nil {
		res := &RemoveFileResult{
			RemovedAt: time.Now(),
		}
		return res, nil
	}

	notExists := errors.Is(err, os.ErrNotExist)
	if notExists {
		return nil, ErrorFileNotFound
	}
	return nil, err
}

func NewFileManager() *fileManager {
	s := &fileManager{}
	return s
}
