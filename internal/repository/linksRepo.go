package repository

import (
	"url-shortener/internal/domain"
)

type LinksRepo interface {
	Add(domain.Link) (string, error)
	GetByShortLink(string) (*domain.Link, error)
}
