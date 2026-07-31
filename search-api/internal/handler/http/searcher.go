// Package http contains the search-api HTTP handlers.
package http

import (
	"context"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

// Searcher is the application port required to serve search requests.
type Searcher interface {
	Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error)
}
