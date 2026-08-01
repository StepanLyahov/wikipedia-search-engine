package app

import (
	"context"
	"net/http"

	"github.com/wikipedia-search-engine/search-api/config"
	searchhandler "github.com/wikipedia-search-engine/search-api/internal/handler/http"
	zaplogger "github.com/wikipedia-search-engine/search-api/internal/logger/zap"
	"github.com/wikipedia-search-engine/search-api/internal/repository/elasticsearch"
	"github.com/wikipedia-search-engine/search-api/internal/service/search"
	"github.com/wikipedia-search-engine/search-api/internal/transport/embedding"
	"go.uber.org/zap"
)

// New builds the application and performs all dependency injection.
func New(_ context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	rawLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	applicationLogger := zaplogger.New(rawLogger)

	embedder, err := embedding.NewClient(cfg.EmbeddingServiceAddr, cfg.EmbeddingRequestTimeout)
	if err != nil {
		applicationLogger.Sync()

		return nil, err
	}

	searchIndex := elasticsearch.NewSearchIndex(cfg.ElasticsearchURL, cfg.ElasticsearchIndex, cfg.RequestTimeout)
	service := search.New(searchIndex, embedder, applicationLogger, search.Config{NumCandidates: cfg.NumCandidates})

	handler := searchhandler.NewHandler(service, applicationLogger, searchhandler.Config{
		DefaultSize: cfg.DefaultSize,
		MaxSize:     cfg.MaxSize,
		DefaultK:    cfg.DefaultK,
		MaxK:        cfg.MaxK,
	})
	router := searchhandler.NewRouter(handler)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: cfg.RequestTimeout,
	}

	return &App{server: server, embedder: embedder, logger: applicationLogger, shutdownTimeout: cfg.ShutdownTimeout}, nil
}
