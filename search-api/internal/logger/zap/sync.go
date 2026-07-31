package zap

// Sync flushes buffered log entries.
func (l *Logger) Sync() { _ = l.logger.Sync() }
