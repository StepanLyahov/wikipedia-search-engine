// Package repository defines persistence ports.
package repository

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

// PageRepository is the persistence port required to read crawled pages.
type PageRepository interface {
	FetchAll(ctx context.Context) ([]domain.Page, error)
}
