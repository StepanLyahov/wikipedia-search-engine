package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

func (m *DocumentIndex) BulkIndex(ctx context.Context, documents []domain.Document) error {
	return m.BulkIndexFunc(ctx, documents)
}
