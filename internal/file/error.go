package file

import "errors"

var (
	ErrorNotFound = errors.New("not found")
	ErrorExists   = errors.New("already exists")
)
