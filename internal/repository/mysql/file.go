package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/datetime"
	db_mysql "github.com/go-seidon/provider/mysql"
)

type fileRepository struct {
	mClient db_mysql.Client
	rClient db_mysql.Client
	clock   datetime.Clock
}

func (r *fileRepository) DeleteFile(ctx context.Context, p repository.DeleteFileParam) (*repository.DeleteFileResult, error) {
	currentTimestamp := r.clock.Now()

	tx, err := r.mClient.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return nil, err
	}

	file, err := r.findFile(ctx, findFileParam{
		UniqueId:      p.UniqueId,
		DbTransaction: tx,
		ShouldLock:    true,
	})
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}

	if file.DeletedAt != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}

		return nil, repository.ErrDeleted
	}

	deleteQuery := `
		UPDATE file 
		SET deleted_at = ?
		WHERE id = ?
	`
	qRes, err := tx.Exec(
		deleteQuery,
		currentTimestamp.UnixMilli(),
		file.UniqueId,
	)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}

	// error is ommited since mysql driver is able to returning totalAffected
	totalAffected, _ := qRes.RowsAffected()
	if totalAffected != 1 {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, fmt.Errorf("record is not updated")
	}

	err = p.DeleteFn(ctx, repository.DeleteFnParam{
		FilePath: file.Path,
	})
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}

	txErr := tx.Commit()
	if txErr != nil {
		return nil, txErr
	}

	res := &repository.DeleteFileResult{
		DeletedAt: currentTimestamp,
	}
	return res, nil
}

func (r *fileRepository) RetrieveFile(ctx context.Context, p repository.RetrieveFileParam) (*repository.RetrieveFileResult, error) {
	file, err := r.findFile(ctx, findFileParam{
		UniqueId: p.UniqueId,
	})
	if err != nil {
		return nil, err
	}

	if file.DeletedAt != nil {
		return nil, repository.ErrDeleted
	}

	res := &repository.RetrieveFileResult{
		UniqueId:  file.UniqueId,
		Name:      file.Name,
		Path:      file.Path,
		MimeType:  file.MimeType,
		Extension: file.Extension,
		Size:      file.Size,
	}
	return res, nil
}

func (r *fileRepository) CreateFile(ctx context.Context, p repository.CreateFileParam) (*repository.CreateFileResult, error) {
	currentTimestamp := r.clock.Now()

	tx, err := r.mClient.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	insertQuery := `
		INSERT INTO file (
			id, name, path, 
			mimetype, extension, size, 
			created_at, updated_at
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(
		insertQuery,
		p.UniqueId,
		p.Name,
		p.Path,
		p.Mimetype,
		p.Extension,
		p.Size,
		currentTimestamp.UnixMilli(),
		currentTimestamp.UnixMilli(),
	)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}

	err = p.CreateFn(ctx, repository.CreateFnParam{
		FilePath: p.Path,
	})
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, txErr
		}
		return nil, err
	}

	txErr := tx.Commit()
	if txErr != nil {
		return nil, txErr
	}
	res := &repository.CreateFileResult{
		UniqueId:  p.UniqueId,
		Name:      p.Name,
		Path:      p.Path,
		Mimetype:  p.Mimetype,
		Extension: p.Extension,
		Size:      p.Size,
		CreatedAt: currentTimestamp,
	}
	return res, nil
}

// @note: using replica client by default
// when transaction occured switch to master client (through `DbTransaction`)
func (r *fileRepository) findFile(ctx context.Context, p findFileParam) (*findFileResult, error) {
	var q db_mysql.Queryable
	q = r.rClient

	if p.DbTransaction != nil {
		q = p.DbTransaction
	}

	sqlQuery := `
		SELECT 
			id, name, path,
			mimetype, extension, size,
			created_at, updated_at, deleted_at
		FROM file
		WHERE id = ?
	`
	if p.ShouldLock {
		sqlQuery += ` FOR UPDATE `
	}

	var res findFileResult
	row := q.QueryRow(sqlQuery, p.UniqueId)
	err := row.Scan(
		&res.UniqueId,
		&res.Name,
		&res.Path,
		&res.MimeType,
		&res.Extension,
		&res.Size,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)
	if err == nil {
		return &res, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return nil, err
}

type findFileParam struct {
	UniqueId      string
	ShouldLock    bool
	DbTransaction *sql.Tx
}

type findFileResult struct {
	UniqueId  string
	Name      string
	Path      string
	MimeType  string
	Extension string
	Size      int64
	CreatedAt int64
	UpdatedAt int64
	DeletedAt *int64
}

func NewFile(opts ...RepoOption) *fileRepository {
	p := RepositoryParam{}
	for _, opt := range opts {
		opt(&p)
	}

	clock := p.clock
	if clock == nil {
		clock = datetime.NewClock()
	}

	return &fileRepository{
		mClient: p.mClient,
		rClient: p.rClient,
		clock:   clock,
	}
}
