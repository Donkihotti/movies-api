package api

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.HomeHandler)
	mux.HandleFunc("GET /movies", h.MoviesHandler)
	mux.HandleFunc("POST /movie", h.CreateMovie)
	mux.HandleFunc("GET /movie", h.CreateMovie)
	mux.HandleFunc("PUT /movie", h.CreateMovie)

	return mux
}
