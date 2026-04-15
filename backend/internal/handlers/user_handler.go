package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"ai-inference-gateway/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// GetAll обробляє запит GET /api/users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.service.GetAll())
}

// GetByID обробляє запит GET /api/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, r, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}
