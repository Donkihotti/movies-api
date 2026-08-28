package api

import (
	"encoding/json"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
)

//GET ALL MOVIES
func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {

	year := r.URL.Query().Get("releaseYear")

	movies, err := h.MovieService.GetMovies(year)
	if err != nil {
		log.Println("Server error", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movies)
}

//GET MOVIE BY ID
func (h *Handler) MovieHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("Invalid id: %v\n", idString)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	movie, err := h.MovieService.GetMovieByID(id)
	if err != nil {
		log.Println("Server error", err)
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {

    var req models.Movie
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }   

    movie, err := h.MovieService.PostMovie(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }   
    w.Header().Set("Content-type", "application/json")

    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(movie)
}

//PATCH MOVIE
func (h *Handler) PatchMovie(w http.ResponseWriter, r *http.Request) {

    idString := r.PathValue("id")
    id, err := strconv.Atoi(idString)
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }   

    var req models.PatchMovieReq

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    decoder.Decode(&req)

    if err := h.MovieService.PatchMovie(r.Context(), req, id); err != nil {
        log.Println(err)
        http.Error(w, "Internal Server error", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

//DELETE MOVIE
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString) 
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := h.MovieService.DeleteMovie(r.Context(), id); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

