// Package repository defines persistence ports.
package repository

import (
	"context"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

// SearchIndex is the persistence port required by the search use case.
type SearchIndex interface {
	Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error)
	VectorSearch(ctx context.Context, vector []float32, k, numCandidates int) ([]domain.Hit, error)
}
