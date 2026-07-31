package http

import (
	"net/http"

	"github.com/wikipedia-search-engine/search-api/internal/logger"
)

// Search handles GET /search: it runs a title/body BM25 query and returns paginated hits.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "q is required")

		return
	}

	from, err := intQueryParam(r, "from", 0)
	if err != nil || from < 0 {
		writeError(w, http.StatusBadRequest, "from must be a non-negative integer")

		return
	}

	size, err := intQueryParam(r, "size", h.cfg.DefaultSize)
	if err != nil || size <= 0 {
		writeError(w, http.StatusBadRequest, "size must be a positive integer")

		return
	}

	if size > h.cfg.MaxSize {
		size = h.cfg.MaxSize
	}

	hits, err := h.searcher.Search(r.Context(), query, from, size)
	if err != nil {
		h.logger.Error("search request failed", logger.Field{Key: "error", Value: err})
		writeError(w, http.StatusInternalServerError, "search failed")

		return
	}

	writeJSON(w, http.StatusOK, searchResponse{Hits: hits})
}
