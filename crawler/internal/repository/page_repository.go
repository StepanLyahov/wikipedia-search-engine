// Package repository defines persistence ports.
package repository

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
)

// PageRepository is the persistence port required by the crawling use case.
type PageRepository interface {
	Exists(ctx context.Context, url string) (bool, error)
	Save(ctx context.Context, page domain.Page) error
}
