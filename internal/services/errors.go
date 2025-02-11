package services

import "errors"

var (
	ErrNotFound           = errors.New("link not found")
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
	ErrInvalidHost        = errors.New("invalid host domain")
	ErrInvalidLinkSize    = errors.New("invalid link size")
	ErrInvalidLink        = errors.New("invalid link")
	ErrInvalidURL         = errors.New("invalid url")
)
