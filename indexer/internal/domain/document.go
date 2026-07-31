package domain

// Document is the cleaned representation of a page as stored in Elasticsearch.
type Document struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
