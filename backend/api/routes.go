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
	mux.HandleFunc("GET /actors", h.GetActorsHandler)
	mux.HandleFunc("GET /actors/{id}", h.GetActorHandler)
	mux.HandleFunc("DELETE /actors/{id}", h.DeleteActorHandler)
	mux.HandleFunc("PATCH /actors/{id}", h.PatchActorHandler)
	

	return mux
}

