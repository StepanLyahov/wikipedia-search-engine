package zap

import (
	"github.com/wikipedia-search-engine/search-api/internal/logger"
	gozap "go.uber.org/zap"
)

func toZapFields(fields []logger.Field) []gozap.Field {
	zapFields := make([]gozap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, gozap.Any(field.Key, field.Value))
	}

	return zapFields
}
