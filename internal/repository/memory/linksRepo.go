package memory

import (
	"sync"

	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

type MemoryLinksRepo struct {
	aliasMap sync.Map // Short to long url mapping
	urlsMap  sync.Map // Long of short url mapping
}

func NewMemoryLinksRepo() *MemoryLinksRepo {
	return &MemoryLinksRepo{}
}

func (p *MemoryLinksRepo) Add(linkDTO domain.Link) (string, error) {
	loadedLink, isLoaded := p.urlsMap.LoadOrStore(linkDTO.OriginalURL, linkDTO.ShortLink)
	if !isLoaded {
		actualLink, _ := p.aliasMap.LoadOrStore(linkDTO.ShortLink, linkDTO.OriginalURL)
		v, _ := actualLink.(string)
		return v, nil
	}
	v, _ := loadedLink.(string)
	return v, nil
}

func (p *MemoryLinksRepo) GetByShortLink(shortLink string) (*domain.Link, error) {
	if originalURL, ok := p.aliasMap.Load(shortLink); ok {
		v, _ := originalURL.(string)
		return &domain.Link{ShortLink: shortLink, OriginalURL: v}, nil
	}
	return nil, repository.ErrShortURLNotFound
}
