package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/wikipedia-search-engine/search-api/internal/logger"
)

const fieldError = "error"

// Run starts the HTTP server and blocks until ctx is cancelled or the server fails.
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("search-api service starting", logger.Field{Key: "addr", Value: a.server.Addr})

	serveErr := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err

			return
		}
		serveErr <- nil
	}()

	select {
	case <-ctx.Done():
		return a.shutdown()
	case err := <-serveErr:
		if err != nil {
			a.logger.Error("search-api service failed", logger.Field{Key: fieldError, Value: err})

			return err
		}

		a.logger.Info("search-api service completed")

		return nil
	}
}

func (a *App) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("search-api service shutdown failed", logger.Field{Key: fieldError, Value: err})

		return err
	}

	a.logger.Info("search-api service completed")

	return nil
}
