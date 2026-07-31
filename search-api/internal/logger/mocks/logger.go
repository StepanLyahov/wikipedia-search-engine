package mocks

import "github.com/wikipedia-search-engine/search-api/internal/logger"

// Logger is a test mock for logger.Logger.
type Logger struct {
	InfoFunc  func(string, ...logger.Field)
	ErrorFunc func(string, ...logger.Field)
}
