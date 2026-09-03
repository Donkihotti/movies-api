package api

import (
	"encoding/json"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/service"
	"net/http"
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

func WriteErrorStatus(w http.ResponseWriter, es errs.ErrorStruct) {

	switch es.ErrType {
	case errs.NotFound:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(es.ErrMsg.Error())
	case errs.BadRequest:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(es.ErrMsg.Error())
	case errs.ServerError:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(es.ErrMsg.Error())
	case errs.Conflict:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(es.ErrMsg.Error())
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func WriteStatus(w http.ResponseWriter, method string) {
	switch method {
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete, http.MethodPatch:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusOK) 
	}
	return
}
