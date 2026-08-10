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

	mux.HandleFunc("POST /actors", h.CreateActor)
	mux.HandleFunc("GET /actors", h.ActorsHandler)
	mux.HandleFunc("GET /actors/{id}", h.ActorHandler)
	

	return mux
}

