package elasticsearch

import (
	"net/http"
	"time"
)

// NewSearchIndex creates an Elasticsearch search index adapter.
func NewSearchIndex(baseURL, index string, timeout time.Duration) *SearchIndex {
	return &SearchIndex{client: &http.Client{Timeout: timeout}, baseURL: baseURL, index: index}
}
