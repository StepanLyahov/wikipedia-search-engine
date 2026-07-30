package crawler

import (
	"github.com/wikipedia-search-engine/crawler/internal/logger"
	"github.com/wikipedia-search-engine/crawler/internal/repository"
)

// New creates a crawler service with its required ports.
func New(pages repository.PageRepository, fetcher Fetcher, log logger.Logger, cfg Config) *Service {
	return &Service{pages: pages, fetcher: fetcher, logger: log, cfg: cfg}
}
