package main

import (
	"flag"
	"fmt"
	"log"
	"url-shortener/internal/cache"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/generator"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/repository"
	"url-shortener/internal/repository/memory"
	"url-shortener/internal/repository/postgres"
	"url-shortener/internal/routers"
	"url-shortener/internal/services"
)

// @title Golang URL Shortener
// @version 1.0
// @description Rest URL shortener
// @termsOfService https://github.com/zaqbez39me/UrlShortener/
// @license.name MIT
// @license.url https://github.com/zaqbez39me/UrlShortener/LICENSE
// @BasePath	/api/v1
// @accept		json
// @produce	json
func main() {
	// Parse command-line arguments
	var storageType, cacheType string
	flag.StringVar(&storageType, "storage-type", "memory", "Type of storage (memory, postgres)")
	flag.StringVar(&cacheType, "cache-type", "redis", "Type of cache (redis, none)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	sLog := logger.Setup(cfg.App.Env)

	// Initialize dependencies
	linkService, err := initDependencies(cfg, storageType, cacheType)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Initialize and start the router
	r := routers.InitRouter(sLog, linkService)
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDependencies(cfg *config.Config, storageType, cacheType string) (*services.LinkService, error) {
	// Initialize link repository
	linkRepo, err := initLinkRepo(storageType, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("link repo initialization error: %w", err)
	}

	// Initialize cache
	cache, err := initCache(cacheType, cfg.Cache)
	if err != nil {
		return nil, fmt.Errorf("cache initialization error: %w", err)
	}

	// Initialize short link generator and service
	generator := generator.NewRandomGenerator(cfg.App.ShortLinkAlphabet)
	linkService, err := services.NewLinkService(
		linkRepo,
		cache,
		generator,
		cfg.App.ShortLinkAlphabet,
		cfg.App.ShortLinkLength,
		cfg.App.Domain,
	)
	if err != nil {
		return nil, fmt.Errorf("link service initialization error: %w", err)
	}

	return linkService, nil
}

func initLinkRepo(storageType string, dbCfg config.DatabaseConfig) (repository.LinksRepo, error) {
	switch storageType {
	case "memory":
		return memory.NewMemoryLinksRepo(), nil
	case "postgres":
		return postgres.NewPostgresLinksRepo(
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Name,
			"links", // Can be extracted as a configuration parameter
		)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

func initCache(cacheType string, cacheCfg config.CacheConfig) (cache.Cache, error) {
	switch cacheType {
	case "redis":
		return cache.NewRedisCache(
			cacheCfg.Host,
			cacheCfg.Port,
			cacheCfg.Password,
			cacheCfg.DB,
			cacheCfg.TTL,
		), nil
	case "none":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}
