package file

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exists")
)
