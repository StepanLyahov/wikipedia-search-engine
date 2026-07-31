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
		DatabaseURL:             env("DATABASE_URL", "postgres://wiki:wiki@localhost:5432/wiki?sslmode=disable"),
		ElasticsearchURL:        env("ELASTICSEARCH_URL", "http://localhost:9200"),
		ElasticsearchIndex:      env("ELASTICSEARCH_INDEX", "wiki_pages"),
		RequestTimeout:          envDuration("INDEXER_REQUEST_TIMEOUT", 15*time.Second),
		BatchSize:               envInt("INDEXER_BATCH_SIZE", 500),
		EmbeddingServiceAddr:    env("EMBEDDING_SERVICE_ADDR", "localhost:50051"),
		EmbeddingRequestTimeout: envDuration("INDEXER_EMBEDDING_TIMEOUT", 15*time.Second),
		EmbeddingDims:           envInt("INDEXER_EMBEDDING_DIMS", 384),
	}
	if cfg.BatchSize <= 0 || cfg.RequestTimeout <= 0 {
		return Config{}, fmt.Errorf("batch size and request timeout must be positive")
	}

	if cfg.EmbeddingRequestTimeout <= 0 || cfg.EmbeddingDims <= 0 {
		return Config{}, fmt.Errorf("embedding request timeout and dims must be positive")
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
