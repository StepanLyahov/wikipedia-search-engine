package indexer

import (
	"github.com/wikipedia-search-engine/indexer/internal/logger"
	"github.com/wikipedia-search-engine/indexer/internal/repository"
)

// New creates an indexer service with its required ports.
func New(pages repository.PageRepository, index repository.DocumentIndex, embedder Embedder, log logger.Logger, cfg Config) *Service {
	return &Service{pages: pages, index: index, embedder: embedder, logger: log, cfg: cfg}
}
