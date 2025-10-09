package handler

import (
	"api/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	return nil
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {

}

func GetTasksByType(w http.ResponseWriter, r *http.Request)
