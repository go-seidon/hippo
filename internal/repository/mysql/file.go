package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/typeconv"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type file struct {
	gormClient *gorm.DB
}

func (r *file) CreateFile(ctx context.Context, p repository.CreateFileParam) (*repository.CreateFileResult, error) {
	tx := r.gormClient.
		WithContext(ctx).
		Clauses(dbresolver.Write).
		Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	currentFile := &File{}
	checkRes := tx.
		Select("id").
		First(currentFile, "id = ?", p.UniqueId)
	if !errors.Is(checkRes.Error, gorm.ErrRecordNotFound) {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		if checkRes.Error == nil {
			return nil, repository.ErrExists
		}
		return nil, checkRes.Error
	}

	createParam := &File{
		Id:        p.UniqueId,
		Name:      p.Name,
		Path:      p.Path,
		Mimetype:  p.Mimetype,
		Extension: p.Extension,
		Size:      p.Size,
		CreatedAt: p.CreatedAt.UnixMilli(),
		UpdatedAt: p.CreatedAt.UnixMilli(),
	}
	createRes := tx.Create(createParam)
	if createRes.Error != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, createRes.Error
	}

	file := &File{}
	findRes := tx.
		Select("id, name, path, mimetype, extension, size, created_at").
		First(file, "id = ?", p.UniqueId)
	if findRes.Error != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, findRes.Error
	}

	err := p.CreateFn(ctx, repository.CreateFnParam{
		FilePath: p.Path,
	})
	if err != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, err
	}

	txRes := tx.Commit()
	if txRes.Error != nil {
		return nil, txRes.Error
	}

	res := &repository.CreateFileResult{
		UniqueId:  file.Id,
		Path:      file.Path,
		Name:      file.Name,
		Mimetype:  file.Mimetype,
		Extension: file.Extension,
		Size:      file.Size,
		CreatedAt: time.UnixMilli(file.CreatedAt).UTC(),
	}
	return res, nil
}

func (r *file) RetrieveFile(ctx context.Context, p repository.RetrieveFileParam) (*repository.RetrieveFileResult, error) {
	query := r.gormClient.
		WithContext(ctx).
		Clauses(dbresolver.Read)

	file := &File{}
	findRes := query.
		Select("id, name, path, mimetype, extension, size, created_at, deleted_at").
		First(file, "id = ?", p.UniqueId)
	if findRes.Error != nil {
		if errors.Is(findRes.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, findRes.Error
	}

	var deletedAt *time.Time
	if file.DeletedAt.Valid {
		deletedAt = typeconv.Time(time.UnixMilli(file.DeletedAt.Int64).UTC())
	}

	res := &repository.RetrieveFileResult{
		UniqueId:  file.Id,
		Path:      file.Path,
		Name:      file.Name,
		Mimetype:  file.Mimetype,
		Extension: file.Extension,
		Size:      file.Size,
		CreatedAt: time.UnixMilli(file.CreatedAt).UTC(),
		DeletedAt: deletedAt,
	}
	return res, nil
}

func (r *file) DeleteFile(ctx context.Context, p repository.DeleteFileParam) (*repository.DeleteFileResult, error) {
	tx := r.gormClient.
		WithContext(ctx).
		Clauses(dbresolver.Write).
		Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	currentFile := &File{}
	findRes := tx.
		Select("id, deleted_at").
		First(currentFile, "id = ?", p.UniqueId)
	if findRes.Error != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		if errors.Is(findRes.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, findRes.Error
	}

	if currentFile.DeletedAt.Valid {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, repository.ErrDeleted
	}

	updateRes := tx.
		Model(&File{}).
		Where("id = ?", p.UniqueId).
		Updates(map[string]interface{}{
			"updated_at": p.DeletedAt.UnixMilli(),
			"deleted_at": p.DeletedAt.UnixMilli(),
		})
	if updateRes.Error != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, updateRes.Error
	}

	file := &File{}
	checkRes := tx.
		Select(`id, path, deleted_at`).
		First(file, "id = ?", p.UniqueId)
	if checkRes.Error != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, checkRes.Error
	}

	err := p.DeleteFn(ctx, repository.DeleteFnParam{
		FilePath: file.Path,
	})
	if err != nil {
		txRes := tx.Rollback()
		if txRes.Error != nil {
			return nil, txRes.Error
		}
		return nil, err
	}

	txRes := tx.Commit()
	if txRes.Error != nil {
		return nil, txRes.Error
	}

	res := &repository.DeleteFileResult{
		DeletedAt: time.UnixMilli(file.DeletedAt.Int64).UTC(),
	}
	return res, nil
}

type FileParam struct {
	GormClient *gorm.DB
}

func NewFile(p FileParam) *file {
	return &file{
		gormClient: p.GormClient,
	}
}

type File struct {
	Id        string        `gorm:"column:id;primaryKey"`
	Path      string        `gorm:"column:path"`
	Name      string        `gorm:"column:name"`
	Mimetype  string        `gorm:"column:mimetype"`
	Extension string        `gorm:"column:extension"`
	Size      int64         `gorm:"column:size"`
	CreatedAt int64         `gorm:"column:created_at"`
	UpdatedAt int64         `gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt sql.NullInt64 `gorm:"column:deleted_at;<-:update"`
}

func (File) TableName() string {
	return "file"
}
