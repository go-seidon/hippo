package repository

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrDeleted      = errors.New("record deleted")
	ErrInvalidParam = errors.New("invalid param")
	ErrExists       = errors.New("resource already exists")
)
