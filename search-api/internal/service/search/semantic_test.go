package search

import (
	"context"
	"errors"
	"testing"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/repository/mocks"
	servicemocks "github.com/wikipedia-search-engine/search-api/internal/service/search/mocks"
)

func TestSemanticReturnsHits(t *testing.T) {
	hits := []domain.Hit{
		{Title: "Search Engine", URL: "https://en.wikipedia.org/wiki/Search_engine", Score: 0.93},
	}
	queryVector := []float32{0.12, -0.45, 0.91}

	embedder := &servicemocks.Embedder{
		EmbedFunc: func(_ context.Context, text string) ([]float32, error) {
			if text != "how search engines work" {
				t.Fatalf("unexpected embed text: %q", text)
			}
			return queryVector, nil
		},
	}

	var gotVector []float32
	var gotK, gotNumCandidates int
	index := &mocks.SearchIndex{
		VectorSearchFunc: func(_ context.Context, vector []float32, k, numCandidates int) ([]domain.Hit, error) {
			gotVector, gotK, gotNumCandidates = vector, k, numCandidates
			return hits, nil
		},
	}

	service := New(index, embedder, noopLogger(), Config{NumCandidates: 100})

	got, err := service.Semantic(context.Background(), "how search engines work", 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(gotVector) != 3 || gotVector[0] != 0.12 {
		t.Fatalf("unexpected query vector passed to VectorSearch: %v", gotVector)
	}

	if gotK != 10 || gotNumCandidates != 100 {
		t.Fatalf("unexpected k=%d numCandidates=%d", gotK, gotNumCandidates)
	}

	if len(got) != 1 || got[0] != hits[0] {
		t.Fatalf("unexpected hits: %+v", got)
	}
}

func TestSemanticWrapsEmbedError(t *testing.T) {
	embedErr := errors.New("embedding-service unavailable")
	embedder := &servicemocks.Embedder{
		EmbedFunc: func(context.Context, string) ([]float32, error) { return nil, embedErr },
	}
	index := &mocks.SearchIndex{
		VectorSearchFunc: func(context.Context, []float32, int, int) ([]domain.Hit, error) {
			t.Fatal("VectorSearch should not be called when embedding fails")
			return nil, nil
		},
	}

	service := New(index, embedder, noopLogger(), Config{NumCandidates: 100})

	_, err := service.Semantic(context.Background(), "query", 10)
	if !errors.Is(err, embedErr) {
		t.Fatalf("expected wrapped embed error, got %v", err)
	}
}

func TestSemanticWrapsVectorSearchError(t *testing.T) {
	searchErr := errors.New("elasticsearch unavailable")
	embedder := &servicemocks.Embedder{
		EmbedFunc: func(context.Context, string) ([]float32, error) { return []float32{0.1}, nil },
	}
	index := &mocks.SearchIndex{
		VectorSearchFunc: func(context.Context, []float32, int, int) ([]domain.Hit, error) { return nil, searchErr },
	}

	service := New(index, embedder, noopLogger(), Config{NumCandidates: 100})

	_, err := service.Semantic(context.Background(), "query", 10)
	if !errors.Is(err, searchErr) {
		t.Fatalf("expected wrapped vector search error, got %v", err)
	}
}
