package indexer

import (
	"context"
	"fmt"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
	"github.com/wikipedia-search-engine/indexer/internal/logger"
)

const fieldURL = "url"

// Run extracts, cleans and indexes every stored page into Elasticsearch.
func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("indexing started")

	if err := s.index.EnsureIndex(ctx); err != nil {
		return fmt.Errorf("ensure index: %w", err)
	}

	pages, err := s.pages.FetchAll(ctx)
	if err != nil {
		return fmt.Errorf("fetch pages: %w", err)
	}

	s.logger.Info("pages loaded", logger.Field{Key: "count", Value: len(pages)})

	documents := s.buildDocuments(pages)

	embedded := s.embedDocuments(ctx, documents)
	s.logger.Info("documents embedded", logger.Field{Key: "count", Value: len(embedded)})

	indexed, err := s.indexInBatches(ctx, embedded)
	if err != nil {
		return err
	}

	s.logger.Info("indexing finished", logger.Field{Key: "indexed", Value: indexed})

	return nil
}

func (s *Service) buildDocuments(pages []domain.Page) []domain.Document {
	documents := make([]domain.Document, 0, len(pages))
	for _, page := range pages {
		title := extractTitle(page.HTML)
		body := extractBody(page.HTML)

		if title == "" && body == "" {
			s.logger.Info("page skipped: empty content", logger.Field{Key: fieldURL, Value: page.URL})

			continue
		}

		documents = append(documents, domain.Document{ID: page.ID, URL: page.URL, Title: title, Body: body})
	}

	return documents
}

func (s *Service) indexInBatches(ctx context.Context, documents []domain.Document) (int, error) {
	indexed := 0
	for start := 0; start < len(documents); start += s.cfg.BatchSize {
		end := start + s.cfg.BatchSize
		if end > len(documents) {
			end = len(documents)
		}

		batch := documents[start:end]
		if err := s.index.BulkIndex(ctx, batch); err != nil {
			return indexed, fmt.Errorf("bulk index batch %d-%d: %w", start, end, err)
		}

		indexed += len(batch)
		s.logger.Info("batch indexed",
			logger.Field{Key: "indexed", Value: indexed},
			logger.Field{Key: "total", Value: len(documents)},
		)
	}

	return indexed, nil
}
