package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"gitea.kood.tech/timdanielfiander/movies-api.git/service"
)

type GenreHandler struct {
	Service *service.GenreService
}

func (h *GenreHandler) GenresHandler(w http.ResponseWriter, r *http.Request) {

	genre, errStruct := h.Service.GetGenres(r.Context())
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) GenreHandler(w http.ResponseWriter, r *http.Request) {

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

	genre, errStruct := h.Service.GetGenre(r.Context(), id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {

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

	genre, errStruct := h.Service.PostGenre(r.Context(), req)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(genre)

}

func (h *GenreHandler) PatchGenre(w http.ResponseWriter, r *http.Request) {

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

	patchedGenre, errStruct := h.Service.PatchGenre(r.Context(), newGenre, id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(patchedGenre)
}

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {

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

	errStruct := h.Service.DeleteGenre(r.Context(), id, force)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
}
