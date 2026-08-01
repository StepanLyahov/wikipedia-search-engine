package http

import "net/http"

// NewRouter builds the search-api HTTP route table.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search", h.Search)
	mux.HandleFunc("GET /semantic", h.Semantic)

	return mux
}
