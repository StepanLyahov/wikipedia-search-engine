package mocks

import "github.com/wikipedia-search-engine/indexer/internal/logger"

func (m *Logger) Info(message string, fields ...logger.Field) { m.InfoFunc(message, fields...) }
