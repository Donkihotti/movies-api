package api

import (
	"encoding/json"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GenresHandler(w http.ResponseWriter, r *http.Request) {

	genre, err := h.GenreService.GetGenres(r.Context())
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)
}

func (h *Handler) GenreHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	genre, err := h.GenreService.GetGenre(r.Context(), id)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)
}

func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {

	var req models.Genre
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	genre, err := h.GenreService.PostGenre(r.Context(), req)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)

}

func (h *Handler) PatchGenre(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	var newGenre models.Genre

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&newGenre); err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	err = h.GenreService.PatchGenre(r.Context(), newGenre, id)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}

func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {

	force := r.URL.Query().Get("force") == "true" 

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	err = h.GenreService.DeleteGenre(r.Context(), id, force)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}
