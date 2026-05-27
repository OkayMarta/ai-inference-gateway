package handlers

import (
	"encoding/json"
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

type updateUserRequest struct {
	Username     *string  `json:"username"`
	TokenBalance *float64 `json:"tokenBalance"`
}

// GetAll обробляє запит GET /api/users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		respondError(w, r, mapErrorToStatus(err), "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, users)
}

// GetByID обробляє запит GET /api/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.GetByID(id)
	if err != nil {
		status := mapErrorToStatus(err)
		message := err.Error()
		if status == http.StatusInternalServerError {
			message = "internal server error"
		}
		respondError(w, r, status, message)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Update обробляє лабораторний PUT /api/users/{id} для CRUD-демонстрації.
// Оновлення tokenBalance тут розглядається як адміністративна/testing операція.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedUser, err := h.service.Update(id, services.UpdateUserInput{
		Username:     req.Username,
		TokenBalance: req.TokenBalance,
	})
	if err != nil {
		status := mapErrorToStatus(err)
		message := err.Error()
		if status == http.StatusInternalServerError {
			message = "internal server error"
		}
		respondError(w, r, status, message)
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}
