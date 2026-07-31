package search

import (
	"github.com/wikipedia-search-engine/search-api/internal/logger"
	"github.com/wikipedia-search-engine/search-api/internal/repository"
)

// New creates a search service with its required ports.
func New(index repository.SearchIndex, log logger.Logger) *Service {
	return &Service{index: index, logger: log}
}
