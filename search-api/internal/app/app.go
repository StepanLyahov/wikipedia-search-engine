// Package app assembles and runs the search-api application.
package app

import (
	"net/http"
	"time"

	zaplogger "github.com/wikipedia-search-engine/search-api/internal/logger/zap"
)

// App owns the assembled dependencies and application lifecycle.
type App struct {
	server          *http.Server
	logger          *zaplogger.Logger
	shutdownTimeout time.Duration
}
