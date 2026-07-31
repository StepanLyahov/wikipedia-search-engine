package http

import "github.com/wikipedia-search-engine/search-api/internal/logger"

// NewHandler creates a search-api HTTP handler.
func NewHandler(searcher Searcher, log logger.Logger, cfg Config) *Handler {
	return &Handler{searcher: searcher, logger: log, cfg: cfg}
}
