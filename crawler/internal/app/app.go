// Package app assembles and runs the crawler application.
package app

import (
	zaplogger "github.com/wikipedia-search-engine/crawler/internal/logger/zap"
	"github.com/wikipedia-search-engine/crawler/internal/repository/postgres"
	"github.com/wikipedia-search-engine/crawler/internal/service/crawler"
)

// App owns the assembled dependencies and application lifecycle.
type App struct {
	crawler *crawler.Service
	pages   *postgres.PageRepository
	logger  *zaplogger.Logger
	seedURL string
}
