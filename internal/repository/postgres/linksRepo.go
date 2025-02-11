package postgres

import (
	"database/sql"
	"fmt"

	"url-shortener/internal/domain"
	"url-shortener/internal/repository"

	"github.com/lib/pq"
)

type PostgresLinksRepo struct {
	db        *sql.DB
	tableName string
}

func NewPostgresLinksRepo(host string, port int, user, password, name, tableName string) (*PostgresLinksRepo, error) {
	db, err := connectToDB(host, port, user, password, name)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	if err := migrateSchema(db, tableName); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return &PostgresLinksRepo{db: db, tableName: tableName}, nil
}

func connectToDB(host string, port int, user, password, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Postgres: %w", err)
	}
	return db, nil
}

func migrateSchema(db *sql.DB, tableName string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			short_link CHARACTER(10) NOT NULL UNIQUE,
			original_url TEXT NOT NULL UNIQUE
		);
		CREATE INDEX IF NOT EXISTS idx_short_link ON %s (short_link);
		CREATE INDEX IF NOT EXISTS idx_original_url ON %s (original_url);
	`, tableName, tableName, tableName)

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("error executing migration: %w", err)
	}
	return nil
}

func (p *PostgresLinksRepo) Add(link domain.Link) (string, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (short_link, original_url) 
		VALUES ($1, $2)
		ON CONFLICT (original_url) DO NOTHING
		RETURNING short_link;
	`, p.tableName)

	var shortLink string
	err := p.db.QueryRow(query, link.ShortLink, link.OriginalURL).Scan(&shortLink)
	if err != nil {
		return p.handleAddError(err, link.OriginalURL)
	}

	return shortLink, nil
}

func (p *PostgresLinksRepo) handleAddError(err error, originalURL string) (string, error) {
	if err == sql.ErrNoRows {
		return p.retrieveShortLink(originalURL)
	}

	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return "", repository.ErrShortURLExists
	}

	return "", fmt.Errorf("error adding link to Postgres: %w", err)
}

func (p *PostgresLinksRepo) retrieveShortLink(originalURL string) (string, error) {
	query := fmt.Sprintf(`
		SELECT short_link FROM %s WHERE original_url = $1;
	`, p.tableName)

	var shortLink string
	if err := p.db.QueryRow(query, originalURL).Scan(&shortLink); err != nil {
		return "", fmt.Errorf("error retrieving short link: %w", err)
	}

	return shortLink, nil
}

func (p *PostgresLinksRepo) GetByShortLink(shortLink string) (*domain.Link, error) {
	query := fmt.Sprintf(`
		SELECT original_url FROM %s WHERE short_link = $1;
	`, p.tableName)

	var originalURL string
	err := p.db.QueryRow(query, shortLink).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error retrieving original URL: %w", err)
	}

	return &domain.Link{ShortLink: shortLink, OriginalURL: originalURL}, nil
}
