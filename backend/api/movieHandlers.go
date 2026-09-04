package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"errors"

	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)

func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	actorid := q.Get("actor")
	genreid := q.Get("genre")
	releaseYear := q.Get("releaseYear")

	filters := models.MovieFilters{}
	if actorid != "" {
		filters.ActorID = &actorid
	}
	if genreid != "" {
		filters.GenreID = &genreid
	}
	if releaseYear != "" {
		filters.ReleaseYear = &releaseYear
	}

	movies, errStruct := h.MovieService.GetMovies(filters)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movies)
}

func (h *Handler) MovieHandler(w http.ResponseWriter, r *http.Request) {

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

	movie, errStruct := h.MovieService.GetMovieByID(id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {

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

	movie, errStruct := h.MovieService.PostMovie(r.Context(), req)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

func (h *Handler) PatchMovie(w http.ResponseWriter, r *http.Request) {

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

	movie, errStruct := h.MovieService.PatchMovie(r.Context(), req, id)
	if errStruct.ErrType != nil {
		WriteErrorStatus(w, errStruct)
		return
	}

	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {

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

	force := r.URL.Query().Get("force") == "true"

	if errStruct := h.MovieService.DeleteMovie(r.Context(), id, force); errStruct.ErrType != nil {
		WriteErrorStatus(w, errStruct)
		return
	}

	WriteStatus(w, r.Method)
}

func (h *Handler) MovieActors(w http.ResponseWriter, r *http.Request) {

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

	actors, errStruct := h.MovieService.MovieActors(r.Context(), id)
	if errStruct.ErrType != nil {
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)

}
