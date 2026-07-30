package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Load reads and validates the application configuration.
func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:    env("DATABASE_URL", "postgres://wiki:wiki@localhost:5432/wiki?sslmode=disable"),
		SeedURL:        env("CRAWLER_SEED_URL", "https://en.wikipedia.org/wiki/Elasticsearch"),
		MaxDepth:       envInt("CRAWLER_MAX_DEPTH", 2),
		MaxPages:       envInt("CRAWLER_MAX_PAGES", 100),
		RequestTimeout: envDuration("CRAWLER_REQUEST_TIMEOUT", 15*time.Second),
		UserAgent:      env("CRAWLER_USER_AGENT", "wikipedia-search-engine-crawler/1.0"),
	}
	if cfg.MaxDepth < 0 || cfg.MaxPages <= 0 || cfg.RequestTimeout <= 0 {
		return Config{}, fmt.Errorf("max depth must be non-negative; max pages and request timeout must be positive")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(env(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(env(key, fallback.String()))
	if err != nil {
		return fallback
	}
	return value
}
