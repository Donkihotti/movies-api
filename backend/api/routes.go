package api

import (
	"net/http"
)

func NewRouter(movieHandler *MovieHandler, genreHandler *GenreHandler, actorHandler *ActorHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/movies", movieHandler.CreateMovie)
	mux.HandleFunc("GET /api/movies", movieHandler.MoviesHandler)
	mux.HandleFunc("GET /api/movies/{id}", movieHandler.MovieHandler)
	mux.HandleFunc("PATCH /api/movies/{id}", movieHandler.PatchMovie)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.DeleteMovie)
	mux.HandleFunc("GET /api/movies/{movieId}/actors", movieHandler.MovieActors)

	mux.HandleFunc("GET /api/genres", genreHandler.GenresHandler)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GenreHandler)
	mux.HandleFunc("POST /api/genres", genreHandler.CreateGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.DeleteGenre)
	mux.HandleFunc("PATCH /api/genres/{id}", genreHandler.PatchGenre)

	mux.HandleFunc("POST /api/actors", actorHandler.PostActor)
	mux.HandleFunc("GET /api/actors", actorHandler.GetActorsHandler)
	mux.HandleFunc("GET /api/actors/name/{name}", actorHandler.GetActorsByNameHandler)
	mux.HandleFunc("GET /api/actors/birtactorHandlerdate/{birthdate}", actorHandler.GetActorsByBirthdateHandler)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetActorHandler)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.DeleteActorHandler)
	mux.HandleFunc("PATCH /api/actors/{id}", actorHandler.PatchActorHandler)

	return mux
}
