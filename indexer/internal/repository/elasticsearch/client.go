// Package elasticsearch contains the Elasticsearch document index adapter.
package elasticsearch

import "net/http"

// DocumentIndex adapts Elasticsearch HTTP APIs to the repository.DocumentIndex port.
type DocumentIndex struct {
	client        *http.Client
	baseURL       string
	index         string
	embeddingDims int
}
