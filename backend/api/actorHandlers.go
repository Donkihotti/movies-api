package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetActorsHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")	
	
	actors, err := h.ActorService.GetActors(name)
	if err != nil {
		fmt.Println("Server error", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(actors)
}

func (h *Handler) GetActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("Invalid id: %v\n", idString)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	actor, err := h.ActorService.GetActorByID(id)
	if err != nil {
		log.Println("Server error", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(actor)
}

//GET ACTORS NAME 
func (h *Handler) GetActorsByNameHandler(w http.ResponseWriter, r *http.Request) {

    name := r.PathValue("name")
    ctx := r.Context()

    actors, err := h.ActorService.GetActorsByName(ctx, name) 
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            http.Error(w, "actors not found by this name", http.StatusNotFound)
            log.Println(err)
            return
        }
        log.Println(err)
        http.Error(w, "error occured when trying to find actors by name", http.StatusInternalServerError)
        return
    }
        
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(actors)
}

//GET ACTORS BIRTHDATE
func (h *Handler) GetActorsByBirthdateHandler(w http.ResponseWriter, r *http.Request) {

    date := r.PathValue("birthdate")
    ctx := r.Context()

    actors, err := h.ActorService.GetActorsByBirthdate(ctx, date)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            http.Error(w, "no actors have this birthdate", http.StatusBadRequest)
            log.Println(err)
            return
        }
        log.Println(err)
        http.Error(w, "error", http.StatusInternalServerError)
    }

    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(actors)
}


//CREATE AN ACTOR
func (h *Handler) PostActor(w http.ResponseWriter, r *http.Request) {

    var req models.Actor
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    actor, err := h.ActorService.PostActor(r.Context(), req)
    if err != nil {
        http.Error(w, "Error posting actor", http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(actor)
}




//DELETE AN ACTOR
func (h *Handler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = h.ActorService.DeleteActor(r.Context(), id)
	if err != nil {
		http.Error(w, "error deleting actor", http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//PATCH AN ACTOR
func (h *Handler) PatchActorHandler(w http.ResponseWriter, r *http.Request) {

    idString := r.PathValue("id")
    id, err := strconv.Atoi(idString)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var actor models.PatchActorReq

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&actor); err != nil {
        log.Println(err)
        http.Error(w, "error", http.StatusBadRequest)
        return
    }

    err = h.ActorService.PatchActor(r.Context(), actor, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            http.Error(w, "actor not found", http.StatusNotFound)
            log.Println(err)
            return
        }
        log.Println(err)
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

