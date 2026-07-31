package mocks

import "github.com/wikipedia-search-engine/search-api/internal/logger"

func (m *Logger) Error(message string, fields ...logger.Field) { m.ErrorFunc(message, fields...) }
