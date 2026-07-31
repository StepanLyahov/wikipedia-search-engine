// Package search implements the document search use case.
package search

import (
	"github.com/wikipedia-search-engine/search-api/internal/logger"
	"github.com/wikipedia-search-engine/search-api/internal/repository"
)

// Service implements the search use case.
type Service struct {
	index  repository.SearchIndex
	logger logger.Logger
}
