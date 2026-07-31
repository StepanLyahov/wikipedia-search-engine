package elasticsearch

import (
	"net/http"
	"time"
)

// NewDocumentIndex creates an Elasticsearch document index adapter.
func NewDocumentIndex(baseURL, index string, timeout time.Duration, embeddingDims int) *DocumentIndex {
	return &DocumentIndex{
		client:        &http.Client{Timeout: timeout},
		baseURL:       baseURL,
		index:         index,
		embeddingDims: embeddingDims,
	}
}
