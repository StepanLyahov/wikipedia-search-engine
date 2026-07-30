package crawler

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/wikipedia-search-engine/crawler/internal/logger"
)

const fieldURL = "url"

// Crawl traverses Wikipedia pages breadth-first and saves unseen pages.
func (s *Service) Crawl(ctx context.Context, seed string) error {
	seedURL, err := wikipediaURL(seed)
	if err != nil {
		return fmt.Errorf("seed URL: %w", err)
	}

	s.logger.Info("crawl started",
		logger.Field{Key: "seed_url", Value: seedURL},
		logger.Field{Key: "max_depth", Value: s.cfg.MaxDepth},
		logger.Field{Key: "max_pages", Value: s.cfg.MaxPages},
	)

	queue := []queuedURL{{url: seedURL}}
	seen := map[string]struct{}{seedURL: {}}
	processed := 0

	for len(queue) > 0 && processed < s.cfg.MaxPages {
		item := queue[0]
		queue = queue[1:]

		s.logger.Info("page discovered",
			logger.Field{Key: fieldURL, Value: item.url},
			logger.Field{Key: "depth", Value: item.depth},
		)

		exists, err := s.pages.Exists(ctx, item.url)
		if err != nil {
			return fmt.Errorf("check page %s: %w", item.url, err)
		}

		if exists {
			s.logger.Info("page already stored", logger.Field{Key: fieldURL, Value: item.url})

			continue
		}

		page, err := s.fetcher.Fetch(ctx, item.url)
		if err != nil {
			s.logger.Error("page fetch failed",
				logger.Field{Key: fieldURL, Value: item.url},
				logger.Field{Key: "error", Value: err},
			)

			continue
		}

		s.logger.Info("page fetched",
			logger.Field{Key: fieldURL, Value: page.URL},
			logger.Field{Key: "title", Value: page.Title},
			logger.Field{Key: "status", Value: page.Status},
		)

		if err := s.pages.Save(ctx, page); err != nil {
			return fmt.Errorf("save page %s: %w", item.url, err)
		}

		processed++
		s.logger.Info("page saved",
			logger.Field{Key: fieldURL, Value: page.URL},
			logger.Field{Key: "title", Value: page.Title},
			logger.Field{Key: "saved_pages", Value: processed},
		)

		if item.depth >= s.cfg.MaxDepth || page.Status < 200 || page.Status >= 300 {
			continue
		}

		for _, link := range wikiLinks(page.HTML) {
			next := "https://en.wikipedia.org" + link
			if _, known := seen[next]; known {
				continue
			}

			seen[next] = struct{}{}
			queue = append(queue, queuedURL{url: next, depth: item.depth + 1})
		}
	}

	s.logger.Info("crawl finished", logger.Field{Key: "saved_pages", Value: processed})

	return nil
}

type queuedURL struct {
	url   string
	depth int
}

func wikipediaURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	if u.Scheme != "https" || u.Host != "en.wikipedia.org" || !strings.HasPrefix(u.Path, "/wiki/") {
		return "", fmt.Errorf("must be an en.wikipedia.org /wiki/ URL")
	}

	u.RawQuery, u.Fragment = "", ""

	return u.String(), nil
}
