package app

import (
	"context"
	"net/http"

	"github.com/wikipedia-search-engine/search-api/config"
	searchhandler "github.com/wikipedia-search-engine/search-api/internal/handler/http"
	zaplogger "github.com/wikipedia-search-engine/search-api/internal/logger/zap"
	"github.com/wikipedia-search-engine/search-api/internal/repository/elasticsearch"
	"github.com/wikipedia-search-engine/search-api/internal/service/search"
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

	searchIndex := elasticsearch.NewSearchIndex(cfg.ElasticsearchURL, cfg.ElasticsearchIndex, cfg.RequestTimeout)
	service := search.New(searchIndex, applicationLogger)

	handler := searchhandler.NewHandler(service, applicationLogger, searchhandler.Config{
		DefaultSize: cfg.DefaultSize,
		MaxSize:     cfg.MaxSize,
	})
	router := searchhandler.NewRouter(handler)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: cfg.RequestTimeout,
	}

	return &App{server: server, logger: applicationLogger, shutdownTimeout: cfg.ShutdownTimeout}, nil
}
