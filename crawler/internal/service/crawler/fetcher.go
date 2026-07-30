package crawler

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
)

// Fetcher is the outbound HTTP dependency needed by the crawler service.
type Fetcher interface {
	Fetch(ctx context.Context, url string) (domain.Page, error)
}
