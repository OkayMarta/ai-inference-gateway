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
	users, err := h.service.GetAll()
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondJSON(w, http.StatusOK, users)
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
