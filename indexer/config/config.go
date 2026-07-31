// Package config loads indexer application settings.
package config

import "time"

// Config contains application settings.
type Config struct {
	DatabaseURL             string
	ElasticsearchURL        string
	ElasticsearchIndex      string
	RequestTimeout          time.Duration
	BatchSize               int
	EmbeddingServiceAddr    string
	EmbeddingRequestTimeout time.Duration
	EmbeddingDims           int
}
