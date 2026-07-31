package mocks

import "github.com/wikipedia-search-engine/indexer/internal/logger"

func (m *Logger) Error(message string, fields ...logger.Field) { m.ErrorFunc(message, fields...) }
