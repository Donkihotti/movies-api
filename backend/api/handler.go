package api

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
)

type Handler struct {
	MovieService *internal.MovieService
	GenreService *internal.GenreService
	ActorService *internal.ActorService
}

func NewHandler(MovieService *internal.MovieService, GenreService *internal.GenreService, ActorService *internal.ActorService) *Handler {
	return &Handler{
		MovieService: MovieService,
		ActorService: ActorService,
		GenreService: GenreService,
	}
}
