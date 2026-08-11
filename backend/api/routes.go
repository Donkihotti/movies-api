package api

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("POST /movies", h.CreateMovie)
	mux.HandleFunc("GET /movies", h.MoviesHandler)
	mux.HandleFunc("GET /movies/{id}", h.MovieHandler)

	mux.HandleFunc("GET /genres", h.GenreHandler)
	mux.HandleFunc("POST /genres", h.CreateGenre)

	return mux
}

