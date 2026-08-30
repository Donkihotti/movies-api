package api

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/movies", h.CreateMovie)
	mux.HandleFunc("GET /api/movies", h.MoviesHandler)
	mux.HandleFunc("GET /api/movies/{id}", h.MovieHandler)
	mux.HandleFunc("PATCH /api/movies/{id}", h.PatchMovie)
	mux.HandleFunc("DELETE /api/movies/{id}", h.DeleteMovie)
	mux.HandleFunc("GET /api/movies/{movieId}/actors", h.MovieActors)

	mux.HandleFunc("GET /api/genres", h.GenresHandler)
	mux.HandleFunc("GET /api/genres/{id}", h.GenreHandler)
	mux.HandleFunc("POST /api/genres", h.CreateGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", h.DeleteGenre)
	mux.HandleFunc("PATCH /api/genres/{id}", h.PatchGenre)

	mux.HandleFunc("POST /api/actors", h.PostActor)
	mux.HandleFunc("GET /api/actors", h.GetActorsHandler)
	mux.HandleFunc("GET /api/actors/name/{name}", h.GetActorsByNameHandler)
	mux.HandleFunc("GET /api/actors/birthdate/{birthdate}", h.GetActorsByBirthdateHandler)
	mux.HandleFunc("GET /api/actors/{id}", h.GetActorHandler)
	mux.HandleFunc("DELETE /api/actors/{id}", h.DeleteActorHandler)
	mux.HandleFunc("PATCH /api/actors/{id}", h.PatchActorHandler)

	return mux
}
