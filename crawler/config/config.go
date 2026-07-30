// Package config loads crawler application settings.
package config

import "time"

// Config contains application settings.
type Config struct {
	DatabaseURL    string
	SeedURL        string
	MaxDepth       int
	MaxPages       int
	RequestTimeout time.Duration
	UserAgent      string
}
