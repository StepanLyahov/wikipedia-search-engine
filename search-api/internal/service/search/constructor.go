package search

import (
	"github.com/wikipedia-search-engine/search-api/internal/logger"
	"github.com/wikipedia-search-engine/search-api/internal/repository"
)

// New creates a search service with its required ports.
func New(index repository.SearchIndex, embedder Embedder, log logger.Logger, cfg Config) *Service {
	return &Service{index: index, embedder: embedder, logger: log, cfg: cfg}
}
