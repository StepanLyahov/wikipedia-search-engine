package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/handler/http/mocks"
	"github.com/wikipedia-search-engine/search-api/internal/logger"
	loggermocks "github.com/wikipedia-search-engine/search-api/internal/logger/mocks"
)

func noopLogger() *loggermocks.Logger {
	return &loggermocks.Logger{
		InfoFunc:  func(string, ...logger.Field) {},
		ErrorFunc: func(string, ...logger.Field) {},
	}
}

func TestSearchReturnsHits(t *testing.T) {
	hits := []domain.Hit{
		{Title: "Distributed computing", URL: "https://en.wikipedia.org/wiki/Distributed_computing", Score: 13.42},
	}

	var gotFrom, gotSize int
	searcher := &mocks.Searcher{
		SearchFunc: func(_ context.Context, _ string, from, size int) ([]domain.Hit, error) {
			gotFrom, gotSize = from, size
			return hits, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultSize: 10, MaxSize: 100})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/search?q=distributed+systems", nil)
	rec := httptest.NewRecorder()

	handler.Search(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	if gotFrom != 0 || gotSize != 10 {
		t.Fatalf("unexpected defaults: from=%d size=%d", gotFrom, gotSize)
	}

	var body searchResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if len(body.Hits) != 1 || body.Hits[0] != hits[0] {
		t.Fatalf("unexpected response body: %+v", body)
	}
}

func TestSearchRequiresQuery(t *testing.T) {
	searcher := &mocks.Searcher{
		SearchFunc: func(context.Context, string, int, int) ([]domain.Hit, error) {
			t.Fatal("searcher should not be called without q")
			return nil, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultSize: 10, MaxSize: 100})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()

	handler.Search(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestSearchClampsSizeToMax(t *testing.T) {
	var gotSize int
	searcher := &mocks.Searcher{
		SearchFunc: func(_ context.Context, _ string, _, size int) ([]domain.Hit, error) {
			gotSize = size
			return nil, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultSize: 10, MaxSize: 50})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/search?q=test&size=1000", nil)
	rec := httptest.NewRecorder()

	handler.Search(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	if gotSize != 50 {
		t.Fatalf("got size %d, want clamped 50", gotSize)
	}
}

func TestSearchRejectsInvalidPagination(t *testing.T) {
	searcher := &mocks.Searcher{
		SearchFunc: func(context.Context, string, int, int) ([]domain.Hit, error) {
			t.Fatal("searcher should not be called with invalid pagination")
			return nil, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultSize: 10, MaxSize: 50})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/search?q=test&from=-1", nil)
	rec := httptest.NewRecorder()

	handler.Search(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
