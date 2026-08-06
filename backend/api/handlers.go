package api

import (
	"net/http"
	"encoding/json"
	"log"

	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
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
