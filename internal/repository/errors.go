package repository

import "errors"

var (
	ErrInternal         = errors.New("internal error")
	ErrShortURLNotFound = errors.New("not found")
	ErrShortURLExists   = errors.New("already exists")
)
