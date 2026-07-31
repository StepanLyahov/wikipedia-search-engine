package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

// DocumentIndex is a test mock for repository.DocumentIndex.
type DocumentIndex struct {
	EnsureIndexFunc func(context.Context) error
	BulkIndexFunc   func(context.Context, []domain.Document) error
}
