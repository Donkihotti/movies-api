package api

import (
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/service"
	"net/http"
	"fmt"
)

type Handler struct {
	MovieService *service.MovieService
	GenreService *service.GenreService
	ActorService *service.ActorService
}

func NewHandler(MovieService *service.MovieService, GenreService *service.GenreService, ActorService *service.ActorService) *Handler {
	return &Handler{
		MovieService: MovieService,
		ActorService: ActorService,
		GenreService: GenreService,
	}
}

// Write status and print error
func WriteErrorStatus(w http.ResponseWriter, es ErrorStruct) {
	switch es.ErrType {
	case errs.NotFound:
		w.WriteHeader(http.StatusNotFound)
	case errs.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case errs.ServerError:
		w.WriteHeader(http.StatusInternalServerError)
	case errs.Conflict:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func WriteStatus(w http.ResponseWriter, method string) {
	switch method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete, http.MethodPatch:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusOK)
	}
	return
}
