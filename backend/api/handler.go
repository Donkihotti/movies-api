package api

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
	"net/http"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
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

//Write status and print error
func WriteErrorStatus(w http.ResponseWriter, err error) {
	switch err {
	case errs.NotFound: 
		w.WriteHeader(http.StatusNotFound)
	case errs.BadRequest: 
		w.WriteHeader(http.StatusBadRequest)
	case errs.ServerError: 
		w.WriteHeader(http.StatusInternalServerError)
	case errs.Conflict: 
		w.WriteHeader(http.StatusConflict)
	default: w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func WriteStatus(w http.ResponseWriter, method string) {
	switch method {
 	case http.MethodGet: w.WriteHeader(http.StatusOK)
	case http.MethodPost: w.WriteHeader(http.StatusCreated)
	case http.MethodDelete, http.MethodPatch: w.WriteHeader(http.StatusNoContent)
	default: w.WriteHeader(http.StatusOK)
	}
	return
} 
