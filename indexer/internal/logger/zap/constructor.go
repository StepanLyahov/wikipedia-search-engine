package zap

import "go.uber.org/zap"

// New creates a Zap logger adapter.
func New(zapLogger *zap.Logger) *Logger { return &Logger{logger: zapLogger} }
