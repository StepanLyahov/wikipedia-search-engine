package http

import "github.com/wikipedia-search-engine/search-api/internal/logger"

// Config configures request handling defaults and limits.
type Config struct {
	DefaultSize int
	MaxSize     int
	DefaultK    int
	MaxK        int
}

// Handler serves the search-api HTTP endpoints.
type Handler struct {
	searcher Searcher
	logger   logger.Logger
	cfg      Config
}
