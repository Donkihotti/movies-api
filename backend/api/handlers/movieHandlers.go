package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"errors"
	"strings"

	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/service"
)

type MovieHandler struct {
	Service *service.MovieService
}

func (h *MovieHandler) MoviesHandler(w http.ResponseWriter, r *http.Request) {

	if strings.HasSuffix(r.RequestURI, "?") {
    errStruct := errs.NewErrorStruct(
        errs.BadRequest,
        errors.New("empty query string"),
    )
    WriteErrorStatus(w, errStruct)
    return
	}

	q := r.URL.Query()

	for key := range q {
	switch key {
	case "genre", "actor", "releaseYear":
	default: 
    errStruct := errs.NewErrorStruct(
        errs.BadRequest,
        errors.New("bad query"),
	    )
		WriteErrorStatus(w, errStruct)
		return 
	}

	}

	filters := models.MovieFilters{}
	actorid := q.Get("actor")

	releaseYear := q.Get("releaseYear")
	if genreValues, exists := q["genre"]; exists {
		if len(genreValues) == 0 || genreValues[0] == "" {
		errStruct := errs.ErrorStruct{
		errs.BadRequest, 
		errors.New("No ID present in request"), 
		}
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)	
		return
		}
	filters.GenreID = &genreValues[0]
	}


	if actorid != "" {
		filters.ActorID = &actorid
	}
	if releaseYear != "" {
		filters.ReleaseYear = &releaseYear
	}

	movies, errStruct := h.Service.GetMovies(filters)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) MovieHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	movie, errStruct := h.Service.GetMovieByID(id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {

	var req models.Movie
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("something went wrong creating a movie"),
		)
		WriteErrorStatus(w, errStruct) 
		return
	}


	movie, errStruct := h.Service.PostMovie(r.Context(), req)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) PatchMovie(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	var req models.PatchMovieReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("something went wrong updating movie"),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	movie, errStruct := h.Service.PatchMovie(r.Context(), req, id)
	if errStruct.ErrType != nil {
		WriteErrorStatus(w, errStruct)
		return
	}

	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.ErrorStruct{
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		}) 
		return
	}


	if errStruct := h.Service.DeleteMovie(r.Context(), id); errStruct.ErrType != nil {
		WriteErrorStatus(w, errStruct)
		return
	}

	WriteStatus(w, r.Method)
}

func (h *MovieHandler) MovieActors(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("movieId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.ErrorStruct{
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		}) 
		return
	}

	actors, errStruct := h.Service.MovieActors(r.Context(), id)
	if errStruct.ErrType != nil {
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)

}
