package app

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/config"
	zaplogger "github.com/wikipedia-search-engine/indexer/internal/logger/zap"
	"github.com/wikipedia-search-engine/indexer/internal/repository/elasticsearch"
	"github.com/wikipedia-search-engine/indexer/internal/repository/postgres"
	"github.com/wikipedia-search-engine/indexer/internal/service/indexer"
	"github.com/wikipedia-search-engine/indexer/internal/transport/embedding"
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

	embedder, err := embedding.NewClient(cfg.EmbeddingServiceAddr, cfg.EmbeddingRequestTimeout)
	if err != nil {
		pages.Close()
		applicationLogger.Sync()

		return nil, err
	}

	documentIndex := elasticsearch.NewDocumentIndex(cfg.ElasticsearchURL, cfg.ElasticsearchIndex, cfg.RequestTimeout, cfg.EmbeddingDims)
	service := indexer.New(pages, documentIndex, embedder, applicationLogger, indexer.Config{
		BatchSize:     cfg.BatchSize,
		EmbeddingDims: cfg.EmbeddingDims,
	})

	return &App{indexer: service, pages: pages, embedder: embedder, logger: applicationLogger}, nil
}
