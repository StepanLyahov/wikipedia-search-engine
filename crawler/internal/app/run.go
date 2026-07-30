package app

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/logger"
)

// Run starts the application use case.
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("crawler service starting", logger.Field{Key: "seed_url", Value: a.seedURL})

	err := a.crawler.Crawl(ctx, a.seedURL)
	if err != nil {
		a.logger.Error("crawler service failed", logger.Field{Key: "error", Value: err})

		return err
	}

	a.logger.Info("crawler service completed")

	return nil
}
