package app

import "github.com/wikipedia-search-engine/indexer/internal/logger"

// Close releases application resources.
func (a *App) Close() {
	a.logger.Info("indexer service stopping")

	a.pages.Close()

	if err := a.embedder.Close(); err != nil {
		a.logger.Error("embedding client close failed", logger.Field{Key: "error", Value: err})
	}

	a.logger.Sync()
}
