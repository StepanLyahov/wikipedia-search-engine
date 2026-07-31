package app

import (
	"context"

	"github.com/wikipedia-search-engine/indexer/internal/logger"
)

// Run starts the application use case.
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("indexer service starting")

	err := a.indexer.Run(ctx)
	if err != nil {
		a.logger.Error("indexer service failed", logger.Field{Key: "error", Value: err})

		return err
	}

	a.logger.Info("indexer service completed")

	return nil
}
