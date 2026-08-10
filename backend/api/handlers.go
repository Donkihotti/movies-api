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
	service *internal.Service
}


func NewHandler(service *internal.Service) *Handler {
	return &Handler{
		service: service,
	}
}


func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
		
}


func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)	
	return
	}

	movies, err := h.service.GetMovies()
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
	
	movie, err := h.service.PostMovie(r.Context(), req)
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

	movie, err := h.service.GetMovieByID(id)
	if err != nil {
	log.Println("Server error", err)
	http.Error(w, "Server error", http.StatusInternalServerError)
	return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}



