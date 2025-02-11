package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Host              string
	Port              int
	Domain            string
	ShortLinkLength   int
	ShortLinkAlphabet string
	Env               string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	TTL      int
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

// LoadConfig initializes and returns the full configuration.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		App: AppConfig{
			Host:              getEnv("APP_HOST", "localhost"),
			Port:              getEnvAsInt("APP_PORT", 8080),
			Domain:            getEnv("APP_DOMAIN", "example.com"),
			ShortLinkLength:   getEnvAsInt("APP_LINK_LENGTH", 10),
			ShortLinkAlphabet: getEnv("APP_LINK_ALPHABET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"),
			Env:               getEnv("APP_ENV", "prod"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			Name:     getEnv("POSTGRES_DATABASE", "url_shortener"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
		},
		Cache: CacheConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			TTL:      getEnvAsInt("REDIS_TTL", 604800),
		},
	}, nil
}

// getEnv returns the value of an environment variable or a fallback value if it's not set.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt returns the value of an environment variable as an integer or a fallback value if it's not set or invalid.
func getEnvAsInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
