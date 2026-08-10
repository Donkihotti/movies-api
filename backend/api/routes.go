package api

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.HomeHandler)
	mux.HandleFunc("GET /movies", h.MoviesHandler)
	mux.HandleFunc("POST /movies", h.CreateMovie)

	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("GET /movies", h.MoviesHandler)
	mux.HandleFunc("GET /movies/{id}", h.MovieHandler)

	return mux
}

