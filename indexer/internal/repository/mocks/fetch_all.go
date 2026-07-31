package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

func (m *PageRepository) FetchAll(ctx context.Context) ([]domain.Page, error) {
	return m.FetchAllFunc(ctx)
}
