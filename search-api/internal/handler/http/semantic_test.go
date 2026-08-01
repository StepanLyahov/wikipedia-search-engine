package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/handler/http/mocks"
)

func TestSemanticReturnsHits(t *testing.T) {
	hits := []domain.Hit{
		{Title: "Search Engine", URL: "https://en.wikipedia.org/wiki/Search_engine", Score: 0.93},
	}

	var gotQuery string
	var gotK int
	searcher := &mocks.Searcher{
		SemanticFunc: func(_ context.Context, query string, k int) ([]domain.Hit, error) {
			gotQuery, gotK = query, k
			return hits, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultK: 10, MaxK: 100})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/semantic?q=how+search+engines+work", nil)
	rec := httptest.NewRecorder()

	handler.Semantic(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	if gotQuery != "how search engines work" || gotK != 10 {
		t.Fatalf("unexpected defaults: query=%q k=%d", gotQuery, gotK)
	}

	var body searchResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if len(body.Hits) != 1 || body.Hits[0] != hits[0] {
		t.Fatalf("unexpected response body: %+v", body)
	}
}

func TestSemanticRequiresQuery(t *testing.T) {
	searcher := &mocks.Searcher{
		SemanticFunc: func(context.Context, string, int) ([]domain.Hit, error) {
			t.Fatal("searcher should not be called without q")
			return nil, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultK: 10, MaxK: 100})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/semantic", nil)
	rec := httptest.NewRecorder()

	handler.Semantic(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestSemanticClampsKToMax(t *testing.T) {
	var gotK int
	searcher := &mocks.Searcher{
		SemanticFunc: func(_ context.Context, _ string, k int) ([]domain.Hit, error) {
			gotK = k
			return nil, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultK: 10, MaxK: 50})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/semantic?q=test&k=1000", nil)
	rec := httptest.NewRecorder()

	handler.Semantic(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	if gotK != 50 {
		t.Fatalf("got k %d, want clamped 50", gotK)
	}
}

func TestSemanticRejectsInvalidK(t *testing.T) {
	searcher := &mocks.Searcher{
		SemanticFunc: func(context.Context, string, int) ([]domain.Hit, error) {
			t.Fatal("searcher should not be called with invalid k")
			return nil, nil
		},
	}

	handler := NewHandler(searcher, noopLogger(), Config{DefaultK: 10, MaxK: 50})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/semantic?q=test&k=-1", nil)
	rec := httptest.NewRecorder()

	handler.Semantic(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
