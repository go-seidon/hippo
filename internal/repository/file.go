package repository

import (
	"context"
	"time"
)

type (
	DeleteFn func(ctx context.Context, p DeleteFnParam) error
	CreateFn func(ctx context.Context, p CreateFnParam) error
)

type File interface {
	CreateFile(ctx context.Context, p CreateFileParam) (*CreateFileResult, error)
	RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, error)
	DeleteFile(ctx context.Context, p DeleteFileParam) (*DeleteFileResult, error)
}

type CreateFileParam struct {
	UniqueId  string
	Name      string
	Path      string
	Mimetype  string
	Extension string
	Size      int64
	CreatedAt time.Time
	CreateFn  CreateFn
}

type CreateFnParam struct {
	FilePath string
}

type CreateFileResult struct {
	UniqueId  string
	Name      string
	Path      string
	Mimetype  string
	Extension string
	Size      int64
	CreatedAt time.Time
}

type RetrieveFileParam struct {
	UniqueId string
}

type RetrieveFileResult struct {
	UniqueId  string
	Name      string
	Path      string
	Mimetype  string
	Extension string
	Size      int64
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DeleteFileParam struct {
	UniqueId  string
	DeletedAt time.Time
	DeleteFn  DeleteFn
}

type DeleteFnParam struct {
	FilePath string
}

type DeleteFileResult struct {
	DeletedAt time.Time
}
