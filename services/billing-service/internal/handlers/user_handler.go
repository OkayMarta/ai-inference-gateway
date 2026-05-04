package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"billing-service/internal/services"
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

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := CurrentUser(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

	user, err := h.service.GetByID(currentUser.ID)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.GetByID(id)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

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
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}
