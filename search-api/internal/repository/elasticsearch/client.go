// Package elasticsearch contains the Elasticsearch search index adapter.
package elasticsearch

import "net/http"

// SearchIndex adapts the Elasticsearch HTTP Search API to the repository.SearchIndex port.
type SearchIndex struct {
	client  *http.Client
	baseURL string
	index   string
}
