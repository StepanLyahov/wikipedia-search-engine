package app

import "github.com/wikipedia-search-engine/search-api/internal/logger"

// Close releases application resources.
func (a *App) Close() {
	a.logger.Info("search-api service stopping")

	if err := a.embedder.Close(); err != nil {
		a.logger.Error("embedding client close failed", logger.Field{Key: fieldError, Value: err})
	}

	a.logger.Sync()
}
