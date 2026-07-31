// Package indexer implements the page indexing use case.
package indexer

import (
	"github.com/wikipedia-search-engine/indexer/internal/logger"
	"github.com/wikipedia-search-engine/indexer/internal/repository"
)

// Config configures indexing behaviour.
type Config struct{ BatchSize int }

// Service implements the indexing use case.
type Service struct {
	pages  repository.PageRepository
	index  repository.DocumentIndex
	logger logger.Logger
	cfg    Config
}
