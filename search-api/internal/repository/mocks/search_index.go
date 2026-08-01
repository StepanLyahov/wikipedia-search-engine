package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

// SearchIndex is a test mock for repository.SearchIndex.
type SearchIndex struct {
	SearchFunc       func(ctx context.Context, query string, from, size int) ([]domain.Hit, error)
	VectorSearchFunc func(ctx context.Context, vector []float32, k, numCandidates int) ([]domain.Hit, error)
}

// Search delegates to SearchFunc.
func (m *SearchIndex) Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error) {
	return m.SearchFunc(ctx, query, from, size)
}

// VectorSearch delegates to VectorSearchFunc.
func (m *SearchIndex) VectorSearch(ctx context.Context, vector []float32, k, numCandidates int) ([]domain.Hit, error) {
	return m.VectorSearchFunc(ctx, vector, k, numCandidates)
}
