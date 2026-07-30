package mocks

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
)

// PageRepository is a test mock for repository.PageRepository.
type PageRepository struct {
	ExistsFunc func(context.Context, string) (bool, error)
	SaveFunc   func(context.Context, domain.Page) error
}
