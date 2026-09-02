package api

import (
	"encoding/json"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"errors"
)

//done
func (h *Handler) GetActorsHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")

	actors, errStruct := h.ActorService.GetActors(name) 
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct) 
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}

//done
func (h *Handler) GetActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString) 
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.ErrorStruct{
			ErrType: errs.BadRequest,
			ErrMsg: fmt.Errorf("invalid actor id: %v", idString),
		})
		return
	}

	actor, errStruct := h.ActorService.GetActorByID(id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actor)
}

//done
func (h *Handler) GetActorsByNameHandler(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")
	ctx := r.Context()

	actors, errStruct := h.ActorService.GetActorsByName(ctx, name)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}


//NEEDS TO CHANGE
func (h *Handler) GetActorsByBirthdateHandler(w http.ResponseWriter, r *http.Request) {

	date := r.PathValue("birthdate")
	ctx := r.Context()

	actors, err := h.ActorService.GetActorsByBirthdate(ctx, date)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}



//done
func (h *Handler) PostActor(w http.ResponseWriter, r *http.Request) {

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

	actor, errStruct := h.ActorService.PostActor(r.Context(), req)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actor)
}




//NEEDS TO CHANGE
func (h *Handler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, errs.BadRequest)
		return
	}
	err = h.ActorService.DeleteActor(r.Context(), id)
	if err != nil {
		log.Println(err)
		//WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}




//CHANGE!!
func (h *Handler) PatchActorHandler(w http.ResponseWriter, r *http.Request) {

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
			errs.ServerError,
			errors.New("error updating actor"),
		)
		WriteErrorStatus(w, errStruct)
		return
	}

	errStruct := h.ActorService.PatchActor(r.Context(), req, id)
	if errStruct.ErrType != nil {
		log.Println(errStruct.Error())
		WriteErrorStatus(w, errStruct)
		return
	}
	WriteStatus(w, r.Method)
}
