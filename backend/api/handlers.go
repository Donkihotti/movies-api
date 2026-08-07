package api

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
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

func (h *Handler) MovieHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString) 

	//write the error os.Stderr if there was one
	if err != nil {
	log.Printf("Invalid id: %v\n", idString)
	//Bad request because the user requested an id that is not valid
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




