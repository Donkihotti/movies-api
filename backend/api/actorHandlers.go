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

	actors, err := h.ActorService.GetActors()
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

func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {

	var req models.CreateActorReq
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

func (h *Handler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = h.ActorService.DeleteActor(r.Context(), id)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "error deleting actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var actor models.Actor

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
