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
		HTTPAddr:                env("HTTP_ADDR", ":8080"),
		ElasticsearchURL:        env("ELASTICSEARCH_URL", "http://localhost:9200"),
		ElasticsearchIndex:      env("ELASTICSEARCH_INDEX", "wiki_pages"),
		RequestTimeout:          envDuration("SEARCH_API_REQUEST_TIMEOUT", 15*time.Second),
		ShutdownTimeout:         envDuration("SEARCH_API_SHUTDOWN_TIMEOUT", 5*time.Second),
		DefaultSize:             envInt("SEARCH_API_DEFAULT_SIZE", 10),
		MaxSize:                 envInt("SEARCH_API_MAX_SIZE", 100),
		DefaultK:                envInt("SEARCH_API_DEFAULT_K", 10),
		MaxK:                    envInt("SEARCH_API_MAX_K", 100),
		NumCandidates:           envInt("SEARCH_API_NUM_CANDIDATES", 100),
		EmbeddingServiceAddr:    env("EMBEDDING_SERVICE_ADDR", "localhost:50051"),
		EmbeddingRequestTimeout: envDuration("SEARCH_API_EMBEDDING_TIMEOUT", 15*time.Second),
	}
	if cfg.RequestTimeout <= 0 || cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("request timeout and shutdown timeout must be positive")
	}

	if cfg.DefaultSize <= 0 || cfg.MaxSize <= 0 || cfg.DefaultSize > cfg.MaxSize {
		return Config{}, fmt.Errorf("default size must be positive and not exceed max size")
	}

	if cfg.DefaultK <= 0 || cfg.MaxK <= 0 || cfg.DefaultK > cfg.MaxK {
		return Config{}, fmt.Errorf("default k must be positive and not exceed max k")
	}

	if cfg.NumCandidates <= 0 || cfg.EmbeddingRequestTimeout <= 0 {
		return Config{}, fmt.Errorf("num candidates and embedding request timeout must be positive")
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
