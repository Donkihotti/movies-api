package api

import (
	"encoding/json"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
)

// GET ALL MOVIES
func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	actorid := q.Get("actor")
	genreid := q.Get("genre")
	releaseYear := q.Get("releaseYear")
	duration := q.Get("duration")

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
	if duration != "" {
		filters.Duration = &duration
	}

	movies, err := h.MovieService.GetMovies(filters)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movies)
}

// GET MOVIE BY ID
func (h *Handler) MovieHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("Invalid id: %v\n", idString)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	movie, err := h.MovieService.GetMovieByID(id)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

// POST MOVIE
func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {

	var req models.Movie
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	movie, err := h.MovieService.PostMovie(r.Context(), req)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(movie)
}

// PATCH MOVIE
func (h *Handler) PatchMovie(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	var req models.PatchMovieReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&req)

	if err := h.MovieService.PatchMovie(r.Context(), req, id); err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}
	WriteStatus(w, r.Method)
}

// DELETE MOVIE
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	force := r.URL.Query().Get("force") == "true"

	if err := h.MovieService.DeleteMovie(r.Context(), id, force); err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}

// GET ALL ACTORS WITHIN A MOVIE
func (h *Handler) MovieActors(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("movieId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("Invalid id: %v\n", idString)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	actors, err := h.MovieService.MovieActors(r.Context(), id)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)

}
