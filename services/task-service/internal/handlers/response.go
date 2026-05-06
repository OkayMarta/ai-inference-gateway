package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"task-service/internal/services"
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
	case errors.Is(err, services.ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrModelNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrInsufficientBalance):
		return http.StatusUnprocessableEntity
	case errors.Is(err, services.ErrTaskCannotBeUpdated):
		return http.StatusConflict
	case errors.Is(err, services.ErrTaskCannotBeDeleted):
		return http.StatusConflict
	case errors.Is(err, services.ErrInvalidPagination):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrUserUpdateNotAllowed):
		return http.StatusForbidden
	case errors.Is(err, services.ErrInvalidUserUpdate):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrInvalidTaskUpdate):
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
	case errors.Is(err, services.ErrBillingUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, services.ErrBillingChargeFailed):
		return http.StatusBadGateway
	case errors.Is(err, services.ErrBillingRefundFailed):
		return http.StatusBadGateway
	case errors.Is(err, services.ErrTaskCompensationFailed):
		return http.StatusBadGateway
	case errors.Is(err, services.ErrTaskCreationFailed):
		return http.StatusInternalServerError
	case errors.Is(err, services.ErrTaskCancellationFailed):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func respondServiceError(w http.ResponseWriter, r *http.Request, err error) {
	status := mapErrorToStatus(err)
	message := err.Error()
	if status == http.StatusInternalServerError &&
		!errors.Is(err, services.ErrTaskCreationFailed) &&
		!errors.Is(err, services.ErrTaskCancellationFailed) {
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
