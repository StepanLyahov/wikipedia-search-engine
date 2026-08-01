package http

import (
	"net/http"

	"github.com/wikipedia-search-engine/search-api/internal/logger"
)

// Semantic handles GET /semantic: it embeds the query and runs a kNN search over the
// embedding field, finding documents by meaning even without exact word matches.
func (h *Handler) Semantic(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "q is required")

		return
	}

	k, err := intQueryParam(r, "k", h.cfg.DefaultK)
	if err != nil || k <= 0 {
		writeError(w, http.StatusBadRequest, "k must be a positive integer")

		return
	}

	if k > h.cfg.MaxK {
		k = h.cfg.MaxK
	}

	hits, err := h.searcher.Semantic(r.Context(), query, k)
	if err != nil {
		h.logger.Error("semantic search request failed", logger.Field{Key: "error", Value: err})
		writeError(w, http.StatusInternalServerError, "semantic search failed")

		return
	}

	writeJSON(w, http.StatusOK, searchResponse{Hits: hits})
}
