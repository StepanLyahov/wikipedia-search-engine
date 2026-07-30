package app

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/config"
	zaplogger "github.com/wikipedia-search-engine/crawler/internal/logger/zap"
	"github.com/wikipedia-search-engine/crawler/internal/repository/postgres"
	"github.com/wikipedia-search-engine/crawler/internal/service/crawler"
	"github.com/wikipedia-search-engine/crawler/internal/transport/wikipedia"
	"go.uber.org/zap"
)

// New builds the application and performs all dependency injection.
func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	rawLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	applicationLogger := zaplogger.New(rawLogger)

	pages, err := postgres.NewPageRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		applicationLogger.Sync()

		return nil, err
	}

	fetcher := wikipedia.NewClient(cfg.RequestTimeout, cfg.UserAgent)
	service := crawler.New(pages, fetcher, applicationLogger, crawler.Config{MaxDepth: cfg.MaxDepth, MaxPages: cfg.MaxPages})

	return &App{crawler: service, pages: pages, logger: applicationLogger, seedURL: cfg.SeedURL}, nil
}
