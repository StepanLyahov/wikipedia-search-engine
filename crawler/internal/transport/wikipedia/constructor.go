package wikipedia

import (
	"net/http"
	"time"
)

// NewClient creates a Wikipedia HTTP client.
func NewClient(timeout time.Duration, userAgent string) *Client {
	return &Client{client: &http.Client{Timeout: timeout}, userAgent: userAgent}
}
