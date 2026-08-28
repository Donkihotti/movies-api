package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GenresHandler(w http.ResponseWriter, r *http.Request) {

	genre, err := h.GenreService.GetGenres(r.Context())
	if err != nil {
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(genre)
}


func (h *Handler) GenreHandler(w http.ResponseWriter, r *http.Request) {
	
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)	
	if err != nil {
	log.Println(err)
	http.Error(w, "Bad Request", http.StatusBadRequest)
	return 
	}

	genre, err := h.GenreService.GetGenre(r.Context(), id)
	if err != nil {
	log.Println(err)
	http.Error(w, "Bad Request", http.StatusBadRequest)
	return
	}
	
	w.Header().Set("Content-type", "application/json")
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

func (h *Handler) PutGenre(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var newGenre models.Genre

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&newGenre); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.GenreService.PutGenre(r.Context(), newGenre, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
