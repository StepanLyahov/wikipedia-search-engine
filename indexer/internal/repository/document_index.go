package repository

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

// DocumentIndex is the persistence port required to index documents for search.
type DocumentIndex interface {
	EnsureIndex(ctx context.Context) error
	BulkIndex(ctx context.Context, documents []domain.Document) error
}
