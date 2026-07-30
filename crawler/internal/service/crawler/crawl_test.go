package crawler

import (
	"context"
	"testing"
	"time"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
	"github.com/wikipedia-search-engine/crawler/internal/logger"
	loggermocks "github.com/wikipedia-search-engine/crawler/internal/logger/mocks"
	repositorymocks "github.com/wikipedia-search-engine/crawler/internal/repository/mocks"
	"github.com/wikipedia-search-engine/crawler/internal/service/crawler/mocks"
)

func TestCrawlRespectsDepth(t *testing.T) {
	seed := "https://en.wikipedia.org/wiki/Elasticsearch"
	child := "https://en.wikipedia.org/wiki/Lucene"
	documents := map[string]domain.Page{
		seed:  {URL: seed, HTML: `<a href="/wiki/Lucene">Lucene</a>`, Status: 200, CreatedAt: time.Now()},
		child: {URL: child, Status: 200, CreatedAt: time.Now()},
	}
	saved := make(map[string]domain.Page)
	repository := &repositorymocks.PageRepository{
		ExistsFunc: func(_ context.Context, url string) (bool, error) { _, exists := saved[url]; return exists, nil },
		SaveFunc:   func(_ context.Context, page domain.Page) error { saved[page.URL] = page; return nil },
	}
	fetcher := &mocks.Fetcher{FetchFunc: func(_ context.Context, url string) (domain.Page, error) { return documents[url], nil }}
	log := &loggermocks.Logger{
		InfoFunc:  func(string, ...logger.Field) {},
		ErrorFunc: func(string, ...logger.Field) {},
	}

	service := New(repository, fetcher, log, Config{MaxDepth: 1, MaxPages: 10})
	if err := service.Crawl(context.Background(), seed); err != nil {
		t.Fatal(err)
	}
	if len(saved) != 2 {
		t.Fatalf("saved %d pages, want 2", len(saved))
	}
}
