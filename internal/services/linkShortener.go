package services

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"url-shortener/internal/cache"
	"url-shortener/internal/domain"
	"url-shortener/internal/lib/generator"
	"url-shortener/internal/repository"

	linkDomain "github.com/chmike/domain"
)

func isValidShortLink(shortLink string, linkSize int, alphabetSet map[rune]bool) bool {
	if len(shortLink) != linkSize {
		return false
	}
	for _, ch := range shortLink {
		if !alphabetSet[ch] {
			return false
		}
	}
	return true
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

type LinkService struct {
	repo        repository.LinksRepo
	cache       cache.Cache
	generator   generator.Generator
	linkSize    int
	alphabetSet map[rune]bool
	host        string
}

func NewLinkService(r repository.LinksRepo, c cache.Cache, g generator.Generator, linkAlphabet string, linkSize int, host string) (*LinkService, error) {
	if err := linkDomain.Check(host); err != nil {
		return nil, ErrInvalidHost
	}
	if linkSize <= 0 {
		return nil, ErrInvalidLinkSize
	}

	alphabetSet := make(map[rune]bool, len(linkAlphabet))
	for _, ch := range linkAlphabet {
		alphabetSet[ch] = true
	}

	return &LinkService{
		repo:        r,
		cache:       c,
		generator:   g,
		linkSize:    linkSize,
		host:        host,
		alphabetSet: alphabetSet,
	}, nil
}

func (s *LinkService) Save(originalURL string, retries int) (string, error) {
	if !isValidURL(originalURL) {
		return "", ErrInvalidURL
	}

	shortURL := url.URL{Scheme: "https", Host: s.host}
	logger := log.Default()

	for i := 0; i < retries; i++ {
		newLink := domain.Link{
			ShortLink:   s.generator.Generate(s.linkSize),
			OriginalURL: originalURL,
		}

		shortLink, err := s.repo.Add(newLink)
		if err == nil {
			logger.Printf("Successfully saved link %s with short link %s", originalURL, shortLink)
			return s.saveToCacheAndReturnURL(shortLink, originalURL, shortURL)
		}

		if errors.Is(err, repository.ErrShortURLExists) {
			logger.Printf("Short link collision occurred: %s", newLink.ShortLink)
			continue
		}

		logger.Printf("Failed to save short link: %v", err)
	}

	return "", ErrMaxRetriesExceeded
}

func (s *LinkService) GetOriginalURL(shortLink string) (string, error) {
	if !isValidShortLink(shortLink, s.linkSize, s.alphabetSet) {
		return "", ErrInvalidLink
	}

	logger := log.Default()
	logger.Printf("Fetching original URL for short link: %s", shortLink)

	// Check cache first
	if s.cache != nil {
		if originalURL, err := s.cache.Get(shortLink); err == nil {
			logger.Printf("Found in cache: %s", originalURL)
			return originalURL, nil
		}
		logger.Printf("Cache miss for short link: %s", shortLink)
	}

	// Check repository
	link, err := s.repo.GetByShortLink(shortLink)
	if err != nil {
		if errors.Is(err, repository.ErrShortURLNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("failed to get original URL from repository for '%s': %w", shortLink, err)
	}

	return link.OriginalURL, nil
}

func (s *LinkService) saveToCacheAndReturnURL(shortLink, originalURL string, shortURL url.URL) (string, error) {
	if s.cache != nil {
		if err := s.cache.Set(shortLink, originalURL); err != nil {
			return "", fmt.Errorf("failed to set short link in cache for '%s': %w", shortLink, err)
		}
	}
	shortURL.Path = shortLink
	return shortURL.String(), nil
}
