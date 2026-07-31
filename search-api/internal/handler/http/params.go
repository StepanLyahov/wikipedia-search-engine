package http

import (
	"net/http"
	"strconv"
)

func intQueryParam(r *http.Request, key string, fallback int) (int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback, nil
	}

	return strconv.Atoi(raw)
}
