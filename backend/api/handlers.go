package api

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)


type Handler struct {
	MovieService *internal.MovieService
	GenreService *internal.GenreService
}


func NewHandler(MovieService *internal.MovieService, GenreService *internal.GenreService) *Handler {
	return &Handler{
		MovieService: MovieService,
		GenreService: GenreService,
	}
}


func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
		
}


func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)	
	return
	}

	movies, err := h.MovieService.GetMovies()
	if err != nil {
	log.Println("Server error",err)
	http.Error(w, "Server error", http.StatusInternalServerError)
	return
	}
	
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movies)
}


func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
	http.Error(w, "Status not allowed", http.StatusMethodNotAllowed)
	return
	}

	var req models.CreateMovieReq
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
	}
	
	movie, err := h.MovieService.PostMovie(r.Context(), req)
	if err != nil { 
	http.Error(w, "Error posting movie", http.StatusBadRequest)
	return 
	}

	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}


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


func (h *Handler) GenreHandler(w http.ResponseWriter, r *http.Request) {

	genre, err := h.GenreService.GetGenres(r.Context())
	if err != nil {
	http.Error(w, "Internal Server error", http.StatusInternalServerError) 	
	}
	
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(genre)
}


func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost { 
	http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
	}
	
	var req models.Genre
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
	http.Error(w, "Bad request", http.StatusBadRequest)
	return
	}

	genre, err := h.GenreService.PostGenre(r.Context(), req)
	if err != nil {
	log.Println(err)
	http.Error(w, "Server error", http.StatusInternalServerError)
	return 
	}

	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(genre)

}


func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
	log.Println(err)
	http.Error(w, "Invalid id", http.StatusBadRequest)
	return
	}

	err = h.GenreService.DeleteGenre(r.Context(), id)
	if err != nil {
	http.Error(w, "Bad request", http.StatusBadRequest)
	return
	}

	w.WriteHeader(http.StatusNoContent)
}

