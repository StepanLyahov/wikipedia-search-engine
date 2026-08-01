package search

import (
	"context"
	"fmt"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/logger"
)

// Semantic finds documents whose meaning is closest to query via a kNN search over embeddings,
// even when the documents share no exact words with query.
func (s *Service) Semantic(ctx context.Context, query string, k int) ([]domain.Hit, error) {
	s.logger.Info("semantic search started",
		logger.Field{Key: "query", Value: query},
		logger.Field{Key: "k", Value: k},
	)

	vector, err := s.embedder.Embed(ctx, query)
	if err != nil {
		s.logger.Error("query embedding failed", logger.Field{Key: fieldError, Value: err})

		return nil, fmt.Errorf("embed query: %w", err)
	}

	hits, err := s.index.VectorSearch(ctx, vector, k, s.cfg.NumCandidates)
	if err != nil {
		s.logger.Error("semantic search failed", logger.Field{Key: fieldError, Value: err})

		return nil, fmt.Errorf("vector search: %w", err)
	}

	s.logger.Info("semantic search completed", logger.Field{Key: "hits", Value: len(hits)})

	return hits, nil
}
