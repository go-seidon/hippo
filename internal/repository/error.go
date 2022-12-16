package repository

import "errors"

var (
	ErrNotFound = errors.New("record not found")
	ErrDeleted  = errors.New("record deleted")
)
