// Package wikipedia contains the Wikipedia HTTP client adapter.
package wikipedia

import "net/http"

// Client adapts Wikipedia HTTP responses to the crawler Fetcher port.
type Client struct {
	client    *http.Client
	userAgent string
}
