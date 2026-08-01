package search

import (
	"context"
	"errors"
	"testing"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/logger"
	loggermocks "github.com/wikipedia-search-engine/search-api/internal/logger/mocks"
	"github.com/wikipedia-search-engine/search-api/internal/repository/mocks"
	servicemocks "github.com/wikipedia-search-engine/search-api/internal/service/search/mocks"
)

func noopLogger() *loggermocks.Logger {
	return &loggermocks.Logger{
		InfoFunc:  func(string, ...logger.Field) {},
		ErrorFunc: func(string, ...logger.Field) {},
	}
}

func unusedEmbedder(t *testing.T) *servicemocks.Embedder {
	return &servicemocks.Embedder{
		EmbedFunc: func(context.Context, string) ([]float32, error) {
			t.Fatal("embedder should not be called by keyword search")
			return nil, nil
		},
	}
}

func TestSearchReturnsHits(t *testing.T) {
	hits := []domain.Hit{
		{Title: "Distributed computing", URL: "https://en.wikipedia.org/wiki/Distributed_computing", Score: 13.42},
	}

	var gotQuery string
	var gotFrom, gotSize int
	index := &mocks.SearchIndex{
		SearchFunc: func(_ context.Context, query string, from, size int) ([]domain.Hit, error) {
			gotQuery, gotFrom, gotSize = query, from, size
			return hits, nil
		},
	}

	service := New(index, unusedEmbedder(t), noopLogger(), Config{NumCandidates: 100})

	got, err := service.Search(context.Background(), "distributed systems", 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	if gotQuery != "distributed systems" || gotFrom != 0 || gotSize != 10 {
		t.Fatalf("unexpected search parameters: query=%q from=%d size=%d", gotQuery, gotFrom, gotSize)
	}

	if len(got) != 1 || got[0] != hits[0] {
		t.Fatalf("unexpected hits: %+v", got)
	}
}

func TestSearchWrapsIndexError(t *testing.T) {
	indexErr := errors.New("elasticsearch unavailable")
	index := &mocks.SearchIndex{
		SearchFunc: func(context.Context, string, int, int) ([]domain.Hit, error) { return nil, indexErr },
	}

	service := New(index, unusedEmbedder(t), noopLogger(), Config{NumCandidates: 100})

	_, err := service.Search(context.Background(), "query", 0, 10)
	if !errors.Is(err, indexErr) {
		t.Fatalf("expected wrapped index error, got %v", err)
	}
}
