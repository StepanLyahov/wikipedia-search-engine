// Package logger defines the logging port used by the application layers.
package logger

// Field is a structured logging attribute.
type Field struct {
	Key   string
	Value any
}

// Logger is the logging dependency used by services without coupling to an implementation.
type Logger interface {
	Info(message string, fields ...Field)
	Error(message string, fields ...Field)
}
