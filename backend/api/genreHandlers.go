package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GenresHandler(w http.ResponseWriter, r *http.Request) {

	genre, errStruct := h.GenreService.GetGenres(r.Context())
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)
}

func (h *Handler) GenreHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	genre, errStruct := h.GenreService.GetGenre(r.Context(), id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
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
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("invalid input in posting genre"),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	genre, errStruct := h.GenreService.PostGenre(r.Context(), req)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
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
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	var newGenre models.Genre

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&newGenre); err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("something went wrong creating a genre"),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	patchedGenre, errStruct := h.GenreService.PatchGenre(r.Context(), newGenre, id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(patchedGenre)
}

func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {

	force := r.URL.Query().Get("force") == "true" 

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

	errStruct := h.GenreService.DeleteGenre(r.Context(), id, force)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
}
