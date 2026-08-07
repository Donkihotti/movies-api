package api

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("/movies", h.MoviesHandler)
	//Ehdotus laittaa noihin handlereihin toi automaattinen tarkistus, mp?
	mux.HandleFunc("GET /movies/{id}", h.MovieHandler)

	return mux
}

