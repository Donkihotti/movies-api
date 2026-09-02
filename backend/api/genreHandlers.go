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

// done
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

// done
func (h *Handler) GenreHandler(w http.ResponseWriter, r *http.Request) {

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

// done
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

// CHANGE
// not yet done
func (h *Handler) PatchGenre(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, errs.BadRequest)
		return
	}

	var newGenre models.Genre

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&newGenre); err != nil {
		log.Println(err)
		//WriteErrorStatus(w, errs.BadRequest)
		return
	}

	err = h.GenreService.PatchGenre(r.Context(), newGenre, id)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}

// CHANGE
func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, errs.BadRequest)
		return
	}

	err = h.GenreService.DeleteGenre(r.Context(), id)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}
