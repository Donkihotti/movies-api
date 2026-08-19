package api

import (
	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)

type Handler struct {
	service *internal.Service
	actorService *internal.ActorService
}


func NewHandler(service *internal.Service, actorService *internal.ActorService) *Handler {
	return &Handler{
		service: service,
		actorService: actorService,
	}
}


func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
		
}


func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)	
	return
	}

	movies, err := h.service.GetMovies()
	if err != nil {
	log.Println("Server error",err)
	http.Error(w, "Server error", http.StatusInternalServerError)
	return
	}
	
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movies)
}


func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
	http.Error(w, "Status not allowed", http.StatusMethodNotAllowed)
	return
	}

	var req models.CreateMovieReq
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
	}
	
	movie, err := h.service.PostMovie(r.Context(), req)
	if err != nil { 
	http.Error(w, "Error posting movie", http.StatusBadRequest)
	return 
	}
	//korjausehdotus claudelta, application-json --> application/json
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}

func (h *Handler) MovieHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString) 
	if err != nil {
	log.Printf("Invalid id: %v\n", idString)
	http.Error(w, "Invalid id", http.StatusBadRequest)
	return 
	}

	movie, err := h.service.GetMovieByID(id)
	if err != nil {
	log.Println("Server error", err)
	http.Error(w, "Server error", http.StatusInternalServerError)
	return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}


//BELOW GOES ALL THE ACTOR RELATED HANDLERS


func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {

	var req models.CreateActorReq
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
				return
		}

	actor, err := h.actorService.PostActor(r.Context(), req)
		if err != nil { 
			http.Error(w, "Error posting actor", http.StatusBadRequest)
				return 
		}
	//korjausehdotus claudelta, application-json --> application/json
	w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(actor)
}


func (h *Handler) GetActorsHandler(w http.ResponseWriter, r *http.Request) {

	actors, err := h.actorService.GetActors()
		if err != nil {
			fmt.Println("Server error",err)
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

		  actor, err := h.actorService.GetActorByID(id)
			  if err != nil {
				  log.Println("Server error", err)
					  http.Error(w, "Server error", http.StatusInternalServerError)
					  return
			  }

		  w.WriteHeader(http.StatusAccepted)
		  json.NewEncoder(w).Encode(actor)
}




func(h *Handler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)	
	if err != nil {
	http.Error(w, "invalid id", http.StatusBadRequest)
	log.Println(err)
	return
	}
	err = h.actorService.DeleteActor(r.Context(), id)
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
		log.Println(err)
		http.Error(w, "invalid id", http.StatusBadRequest)
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

	err = h.actorService.PatchActor(r.Context(), actor, id)
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














