package search

import (
	"context"
	"fmt"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/logger"
)

// Search looks up documents matching query across title and body, ranked by BM25 score.
func (s *Service) Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error) {
	s.logger.Info("search started",
		logger.Field{Key: "query", Value: query},
		logger.Field{Key: "from", Value: from},
		logger.Field{Key: "size", Value: size},
	)

	hits, err := s.index.Search(ctx, query, from, size)
	if err != nil {
		s.logger.Error("search failed", logger.Field{Key: "error", Value: err})

		return nil, fmt.Errorf("search: %w", err)
	}

	s.logger.Info("search completed", logger.Field{Key: "hits", Value: len(hits)})

	return hits, nil
}
