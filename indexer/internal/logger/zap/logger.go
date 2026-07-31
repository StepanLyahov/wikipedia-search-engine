// Package zap adapts go.uber.org/zap to the application logger port.
package zap

import "go.uber.org/zap"

// Logger is a Zap-backed implementation of logger.Logger.
type Logger struct{ logger *zap.Logger }
