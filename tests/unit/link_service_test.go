package services_test

import (
	"errors"
	"testing"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
	"url-shortener/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo simulates the repository behavior
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Add(link domain.Link) (string, error) {
	args := m.Called(link)
	return args.String(0), args.Error(1)
}

func (m *MockRepo) GetByShortLink(shortLink string) (*domain.Link, error) {
	args := m.Called(shortLink)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Link), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockCache simulates cache behavior
type MockCache struct {
	mock.Mock
}

func (c *MockCache) Get(key string) (string, error) {
	args := c.Called(key)
	return args.String(0), args.Error(1)
}

func (c *MockCache) Set(key, value string) error {
	args := c.Called(key, value)
	return args.Error(0)
}

// MockGenerator simulates the link generator behavior
type MockGenerator struct {
	alphabet string
}

func (g *MockGenerator) Generate(size int) string {
	return g.alphabet[:size]
}

func TestSave_Success(t *testing.T) {
	repo := new(MockRepo)
	cache := new(MockCache)
	generator := &MockGenerator{alphabet: "abcdefghijklmnopqrstuvwxyz"}
	linkService, _ := services.NewLinkService(repo, cache, generator, "abcdefghijklmnopqrstuvwxyz", 10, "example.com")

	originalURL := "https://example.com"
	shortLink := "abcdefghij"

	repo.On("Add", mock.Anything).Return(shortLink, nil).Once()
	cache.On("Set", shortLink, originalURL).Return(nil)

	result, err := linkService.Save(originalURL, 3)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/abcdefghij", result)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestSave_InvalidURL(t *testing.T) {
	repo := new(MockRepo)
	cache := new(MockCache)
	generator := &MockGenerator{alphabet: "abcdefghijklmnopqrstuvwxyz"}
	linkService, _ := services.NewLinkService(repo, cache, generator, "abcdefghijklmnopqrstuvwxyz", 10, "example.com")

	invalidURL := "invalid-url"
	result, err := linkService.Save(invalidURL, 3)
	assert.ErrorIs(t, err, services.ErrInvalidURL)
	assert.Empty(t, result)
}

func TestGetOriginalURL_CacheHit(t *testing.T) {
	repo := new(MockRepo)
	cache := new(MockCache)
	generator := &MockGenerator{alphabet: "abcdefghijklmnopqrstuvwxyz"}
	linkService, _ := services.NewLinkService(repo, cache, generator, "abcdefghijklmnopqrstuvwxyz", 10, "example.com")

	shortLink := "abcdefghij"
	originalURL := "https://example.com"

	cache.On("Get", shortLink).Return("https://example.com", nil).Once()

	result, err := linkService.GetOriginalURL(shortLink)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)
	cache.AssertExpectations(t)
}

func TestGetOriginalURL_CacheMiss_RepoHit(t *testing.T) {
	repo := new(MockRepo)
	cache := new(MockCache)
	generator := &MockGenerator{alphabet: "abcdefghijklmnopqrstuvwxyz"}
	linkService, _ := services.NewLinkService(repo, cache, generator, "abcdefghijklmnopqrstuvwxyz", 10, "example.com")

	shortLink := "abcdefghij"
	originalURL := "https://example.com"

	cache.On("Get", shortLink).Return("", errors.New("cache miss"))
	repo.On("GetByShortLink", shortLink).Return(&domain.Link{ShortLink: shortLink, OriginalURL: originalURL}, nil).Once()

	result, err := linkService.GetOriginalURL(shortLink)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	repo := new(MockRepo)
	cache := new(MockCache)
	generator := &MockGenerator{alphabet: "abcdefghijklmnopqrstuvwxyz"}
	linkService, _ := services.NewLinkService(repo, cache, generator, "abcdefghijklmnopqrstuvwxyz", 10, "example.com")

	shortLink := "nonexisten"
	cache.On("Get", shortLink).Return("", errors.New("cache miss"))
	repo.On("GetByShortLink", shortLink).Return(nil, repository.ErrShortURLNotFound).Once()

	result, err := linkService.GetOriginalURL(shortLink)
	assert.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
