// Package domain contains indexer business entities.
package domain

// Page is the raw crawled Wikipedia article as stored by the crawler service.
type Page struct {
	ID   int64
	URL  string
	HTML string
}
