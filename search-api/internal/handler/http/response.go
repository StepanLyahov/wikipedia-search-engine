package http

import (
	"encoding/json"
	"net/http"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

type searchResponse struct {
	Hits []domain.Hit `json:"hits"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
