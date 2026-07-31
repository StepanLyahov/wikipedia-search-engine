package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

// SearchIndex is a test mock for repository.SearchIndex.
type SearchIndex struct {
	SearchFunc func(ctx context.Context, query string, from, size int) ([]domain.Hit, error)
}

// Search delegates to SearchFunc.
func (m *SearchIndex) Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error) {
	return m.SearchFunc(ctx, query, from, size)
}
