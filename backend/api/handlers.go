package api

import (
	"net/http"
	"gitea.kood.tech/timdanielfiander/movies-api.git/internal"
)

type Handler struct {
	service *internal.Service
}

func NewHandler(service *internal.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	
}
