package indexer

import (
	"context"
	"testing"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
	"github.com/wikipedia-search-engine/indexer/internal/logger"
	loggermocks "github.com/wikipedia-search-engine/indexer/internal/logger/mocks"
	"github.com/wikipedia-search-engine/indexer/internal/repository/mocks"
)

func TestRunIndexesPagesInBatches(t *testing.T) {
	pages := []domain.Page{
		{ID: 1, URL: "https://en.wikipedia.org/wiki/Elasticsearch", HTML: `<h1 id="firstHeading">Elasticsearch</h1><p>A search engine.</p>`},
		{ID: 2, URL: "https://en.wikipedia.org/wiki/Lucene", HTML: `<h1 id="firstHeading">Lucene</h1><p>A search library.</p>`},
	}

	pageRepository := &mocks.PageRepository{
		FetchAllFunc: func(context.Context) ([]domain.Page, error) { return pages, nil },
	}

	ensureIndexCalled := false
	var batches [][]domain.Document
	index := &mocks.DocumentIndex{
		EnsureIndexFunc: func(context.Context) error { ensureIndexCalled = true; return nil },
		BulkIndexFunc: func(_ context.Context, documents []domain.Document) error {
			batches = append(batches, documents)
			return nil
		},
	}

	log := &loggermocks.Logger{
		InfoFunc:  func(string, ...logger.Field) {},
		ErrorFunc: func(string, ...logger.Field) {},
	}

	service := New(pageRepository, index, log, Config{BatchSize: 1})
	if err := service.Run(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !ensureIndexCalled {
		t.Fatal("expected EnsureIndex to be called")
	}

	if len(batches) != 2 {
		t.Fatalf("got %d batches, want 2", len(batches))
	}

	// extractBody joins all visible text on the page, including the heading itself,
	// so the heading text is expected to appear at the start of the body too.
	if batches[0][0].Title != "Elasticsearch" || batches[0][0].Body != "Elasticsearch A search engine." {
		t.Fatalf("unexpected first document: %+v", batches[0][0])
	}

	if batches[1][0].Title != "Lucene" || batches[1][0].Body != "Lucene A search library." {
		t.Fatalf("unexpected second document: %+v", batches[1][0])
	}
}

func TestRunSkipsPagesWithoutContent(t *testing.T) {
	pages := []domain.Page{{ID: 1, URL: "https://en.wikipedia.org/wiki/Empty", HTML: `<html></html>`}}

	pageRepository := &mocks.PageRepository{
		FetchAllFunc: func(context.Context) ([]domain.Page, error) { return pages, nil },
	}

	bulkIndexCalled := false
	index := &mocks.DocumentIndex{
		EnsureIndexFunc: func(context.Context) error { return nil },
		BulkIndexFunc: func(context.Context, []domain.Document) error {
			bulkIndexCalled = true
			return nil
		},
	}

	log := &loggermocks.Logger{
		InfoFunc:  func(string, ...logger.Field) {},
		ErrorFunc: func(string, ...logger.Field) {},
	}

	service := New(pageRepository, index, log, Config{BatchSize: 10})
	if err := service.Run(context.Background()); err != nil {
		t.Fatal(err)
	}

	if bulkIndexCalled {
		t.Fatal("expected BulkIndex not to be called when no documents were built")
	}
}
