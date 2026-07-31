// Package app assembles and runs the indexer application.
package app

import (
	zaplogger "github.com/wikipedia-search-engine/indexer/internal/logger/zap"
	"github.com/wikipedia-search-engine/indexer/internal/repository/postgres"
	"github.com/wikipedia-search-engine/indexer/internal/service/indexer"
)

// App owns the assembled dependencies and application lifecycle.
type App struct {
	indexer *indexer.Service
	pages   *postgres.PageRepository
	logger  *zaplogger.Logger
}
