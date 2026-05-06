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

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type authResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

type messageResponse struct {
	Message string `json:"message"`
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

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.RequestPasswordReset(req.Email); err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, messageResponse{
		Message: "If an account with this email exists, a password reset link has been sent.",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, messageResponse{
		Message: "Password has been reset successfully.",
	})
}
