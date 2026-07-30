package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
)

// Fetcher is a test mock for crawler.Fetcher.
type Fetcher struct {
	FetchFunc func(context.Context, string) (domain.Page, error)
}

func (m *Fetcher) Fetch(ctx context.Context, url string) (domain.Page, error) {
	return m.FetchFunc(ctx, url)
}
