package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
)

func (m *PageRepository) Save(ctx context.Context, page domain.Page) error {
	return m.SaveFunc(ctx, page)
}
