package api

import (
	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"errors"
	"database/sql"
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
	"gitea.kood.tech/timdanielfiander/movies-api.git/models"
)


type Handler struct {
	MovieService *internal.MovieService
	GenreService *internal.GenreService
	actorService *internal.ActorService
}


func NewHandler(MovieService *internal.MovieService, GenreService *internal.GenreService, actorService *internal.ActorService) *Handler {
	return &Handler{
		MovieService: MovieService,
		GenreService: GenreService,
		actorService: ActorService,
	}
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
		
}


func (h *Handler) MoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)	
	return
	}

	movies, err := h.MovieService.GetMovies()
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
	
	movie, err := h.MovieService.PostMovie(r.Context(), req)
	if err != nil { 
	http.Error(w, "Error posting movie", http.StatusBadRequest)
	return 
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}

func (h *Handler) PatchMovie(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString) 
	if err != nil {
	log.Println(err)
	http.Error(w, "Invalid id", http.StatusBadRequest)
	return
	}
	
	var req models.Movie

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&req)

	if err := h.MovieService.PatchMovie(r.Context(), req, id); err != nil {
	log.Println(err) 
	http.Error(w, "Internal Server error", http.StatusInternalServerError)
	return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MovieHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString) 
	if err != nil {
	log.Printf("Invalid id: %v\n", idString)
	http.Error(w, "Invalid id", http.StatusBadRequest)
	return 
	}

	movie, err := h.MovieService.GetMovieByID(id)
	if err != nil {
	log.Println("Server error", err)
	http.Error(w, "invalid id", http.StatusNotFound)
	return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}


func (h *Handler) GenreHandler(w http.ResponseWriter, r *http.Request) {

	genre, err := h.GenreService.GetGenres(r.Context())
	if err != nil {
	http.Error(w, "Internal Server error", http.StatusInternalServerError) 	
	}
	
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(genre)
}


func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost { 
	http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
	}
	
	var req models.Genre
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
	http.Error(w, "Bad request", http.StatusBadRequest)
	return
	}

	genre, err := h.GenreService.PostGenre(r.Context(), req)
	if err != nil {
	log.Println(err)
	http.Error(w, "Server error", http.StatusInternalServerError)
	return 
	}

	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(genre)

}


func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
	log.Println(err)
	http.Error(w, "Invalid id", http.StatusBadRequest)
	return
	}

	err = h.GenreService.DeleteGenre(r.Context(), id)
	if err != nil {
	http.Error(w, "Bad request", http.StatusBadRequest)
	return
	}

	w.WriteHeader(http.StatusNoContent)
}


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
	http.Error(w, "Invalid id", http.StatusBadRequest)
	return
	}

	err = h.GenreService.DeleteGenre(r.Context(), id)
	if err != nil {
	http.Error(w, "Bad request", http.StatusBadRequest)
	return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *Handler) PutGenre(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
	http.Error(w, "Invalid ID", http.StatusBadRequest)
	return
	}

	var newGenre models.Genre
	
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&newGenre); err != nil {
	log.Println(err)
	http.Error(w, "Bad Request", http.StatusBadRequest)
	return
	}

	err = h.GenreService.PutGenre(r.Context(), newGenre, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		http.Error(w, "Genre not found", http.StatusNotFound)
		return
		}
		
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchActorHandler(w http.ResponseWriter, r *http.Request) { 

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
	log.Error(err)
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

