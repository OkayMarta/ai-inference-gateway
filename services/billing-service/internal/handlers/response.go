package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"billing-service/internal/services"
)

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
		Message:   message,
		Path:      r.URL.Path,
	}

	_ = json.NewEncoder(w).Encode(response)
}

func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrInsufficientBalance):
		return http.StatusUnprocessableEntity
	case errors.Is(err, services.ErrInvalidUserUpdate):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, services.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, services.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, services.ErrEmailAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, services.ErrInvalidRegisterInput):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrInvalidBillingInput):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrInvalidPasswordResetInput):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrInvalidPasswordResetToken):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func respondServiceError(w http.ResponseWriter, r *http.Request, err error) {
	status := mapErrorToStatus(err)
	message := err.Error()
	if status == http.StatusInternalServerError {
		message = "internal server error"
	}

	respondError(w, r, status, message)
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				respondError(w, r, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
