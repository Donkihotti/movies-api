package handlers

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

type ActorHandler struct {
	Service *service.ActorService
}

func (h *ActorHandler) GetActors(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")

	actors, errStruct := h.Service.GetActors(name)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) GetActorByID(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.ErrorStruct{
			ErrType: errs.BadRequest,
			ErrMsg:  fmt.Errorf("invalid actor id: %v", idString),
		})
		return
	}

	actor, errStruct := h.Service.GetActorByID(id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) GetActorsByName(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")
	ctx := r.Context()

	actors, errStruct := h.Service.GetActorsByName(ctx, name)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) GetActorsByBirthdate(w http.ResponseWriter, r *http.Request) {

	date := r.PathValue("birthdate")
	ctx := r.Context()

	actors, errStruct := h.Service.GetActorsByBirthdate(ctx, date)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) PostActor(w http.ResponseWriter, r *http.Request) {

	var req models.Actor
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("error posting an actor"),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	actor, errStruct := h.Service.PostActor(r.Context(), req)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {

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
	force := r.URL.Query().Get("force") == "true"

	errStruct := h.Service.DeleteActor(r.Context(), id, force)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
}

func (h *ActorHandler) PatchActor(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			fmt.Errorf("invalid id: %v", idString),
		)
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	var req models.PatchActorReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		log.Println(err)
		errStruct := errs.NewErrorStruct(
			errs.BadRequest,
			errors.New("error updating actor"),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	patchedActor, errStruct := h.Service.PatchActor(r.Context(), req, id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(patchedActor)
}
