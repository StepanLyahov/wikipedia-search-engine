package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

// Searcher is a test mock for http.Searcher.
type Searcher struct {
	SearchFunc   func(ctx context.Context, query string, from, size int) ([]domain.Hit, error)
	SemanticFunc func(ctx context.Context, query string, k int) ([]domain.Hit, error)
}

// Search delegates to SearchFunc.
func (m *Searcher) Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error) {
	return m.SearchFunc(ctx, query, from, size)
}

// Semantic delegates to SemanticFunc.
func (m *Searcher) Semantic(ctx context.Context, query string, k int) ([]domain.Hit, error) {
	return m.SemanticFunc(ctx, query, k)
}
