package handlers

import (
	"encoding/json"
	"net/http"

	"billing-service/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusCreated, authResponse{User: user, Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, authResponse{User: user, Token: token})
}
