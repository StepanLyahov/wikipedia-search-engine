package indexer

import (
	"context"
	"strings"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
	"github.com/wikipedia-search-engine/indexer/internal/logger"
)

// embedDocuments enriches documents with embeddings from the embedding-service, skipping any
// document whose embedding fails or does not match the configured vector size.
func (s *Service) embedDocuments(ctx context.Context, documents []domain.Document) []domain.Document {
	embedded := make([]domain.Document, 0, len(documents))
	for _, document := range documents {
		vector, err := s.embedder.Embed(ctx, embeddingText(document))
		if err != nil {
			s.logger.Error("embedding failed",
				logger.Field{Key: fieldURL, Value: document.URL},
				logger.Field{Key: "error", Value: err},
			)

			continue
		}

		if len(vector) != s.cfg.EmbeddingDims {
			s.logger.Error("embedding dimension mismatch",
				logger.Field{Key: fieldURL, Value: document.URL},
				logger.Field{Key: "got", Value: len(vector)},
				logger.Field{Key: "want", Value: s.cfg.EmbeddingDims},
			)

			continue
		}

		document.Embedding = vector
		embedded = append(embedded, document)
	}

	return embedded
}

func embeddingText(document domain.Document) string {
	return strings.TrimSpace(document.Title + " " + document.Body)
}
