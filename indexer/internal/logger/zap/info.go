package zap

import "github.com/wikipedia-search-engine/indexer/internal/logger"

// Info writes an informational structured log entry.
func (l *Logger) Info(message string, fields ...logger.Field) {
	l.logger.Info(message, toZapFields(fields)...)
}
