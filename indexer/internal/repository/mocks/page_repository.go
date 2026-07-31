package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

// PageRepository is a test mock for repository.PageRepository.
type PageRepository struct {
	FetchAllFunc func(context.Context) ([]domain.Page, error)
}
