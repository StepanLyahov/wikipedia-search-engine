package zap

import "github.com/wikipedia-search-engine/search-api/internal/logger"

// Error writes an error structured log entry.
func (l *Logger) Error(message string, fields ...logger.Field) {
	l.logger.Error(message, toZapFields(fields)...)
}
