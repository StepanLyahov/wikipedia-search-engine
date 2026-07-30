// Package crawler implements the page traversal use case.
package crawler

import (
	"github.com/wikipedia-search-engine/crawler/internal/logger"
	"github.com/wikipedia-search-engine/crawler/internal/repository"
)

// Config configures crawl limits.
type Config struct{ MaxDepth, MaxPages int }

// Service implements the crawler use case.
type Service struct {
	pages   repository.PageRepository
	fetcher Fetcher
	logger  logger.Logger
	cfg     Config
}
