// Package config loads search-api application settings.
package config

import "time"

// Config contains application settings.
type Config struct {
	HTTPAddr                string
	ElasticsearchURL        string
	ElasticsearchIndex      string
	RequestTimeout          time.Duration
	ShutdownTimeout         time.Duration
	DefaultSize             int
	MaxSize                 int
	DefaultK                int
	MaxK                    int
	NumCandidates           int
	EmbeddingServiceAddr    string
	EmbeddingRequestTimeout time.Duration
}
