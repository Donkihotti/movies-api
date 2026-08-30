package api

import (
	"encoding/json"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
	"gitea.kood.tech/timdanielfiander/movies-api.git/errs"
)

func (h *Handler) GetActorsHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")	
	
	actors, err := h.ActorService.GetActors(name)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}
    w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actors)
}

func (h *Handler) GetActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("invalid id: ", id)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}

	actor, err := h.ActorService.GetActorByID(id)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
	json.NewEncoder(w).Encode(actor)
}

func (h *Handler) GetActorsByNameHandler(w http.ResponseWriter, r *http.Request) {

    name := r.PathValue("name")
    ctx := r.Context()

    actors, err := h.ActorService.GetActorsByName(ctx, name) 
    if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
        return
    }
        
    w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
    json.NewEncoder(w).Encode(actors)
}

//GET ACTORS BIRTHDATE
func (h *Handler) GetActorsByBirthdateHandler(w http.ResponseWriter, r *http.Request) {

    date := r.PathValue("birthdate")
    ctx := r.Context()

    actors, err := h.ActorService.GetActorsByBirthdate(ctx, date)
    if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
    }
    w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
    json.NewEncoder(w).Encode(actors)
}


//CREATE AN ACTOR
func (h *Handler) PostActor(w http.ResponseWriter, r *http.Request) {

    var req models.Actor
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&req); err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
        return
    }

    actor, err := h.ActorService.PostActor(r.Context(), req)
    if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
	WriteStatus(w, r.Method)
    json.NewEncoder(w).Encode(actor)
}

//DELETE AN ACTOR
func (h *Handler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
		return
	}
	err = h.ActorService.DeleteActor(r.Context(), id)
	if err != nil {
		log.Println(err)
		WriteErrorStatus(w, err)
		return
	}

	WriteStatus(w, r.Method)
}

//PATCH AN ACTOR
func (h *Handler) PatchActorHandler(w http.ResponseWriter, r *http.Request) {

    idString := r.PathValue("id")
    id, err := strconv.Atoi(idString)
    if err != nil {
		log.Println(err)
		WriteErrorStatus(w, errs.BadRequest)
        return
    }

    var actor models.PatchActorReq

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&actor); err != nil {
        log.Println(err)
		WriteErrorStatus(w, err)
        return
    }

    err = h.ActorService.PatchActor(r.Context(), actor, id)
    if err != nil {
        log.Println(err)
		WriteErrorStatus(w, err)
        return
    }
	WriteStatus(w, r.Method)
}

