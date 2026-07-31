// Package domain contains the search-api business entities.
package domain

// Hit is a single search result ranked by Elasticsearch relevance score.
type Hit struct {
	Title string  `json:"title"`
	URL   string  `json:"url"`
	Score float64 `json:"score"`
}
