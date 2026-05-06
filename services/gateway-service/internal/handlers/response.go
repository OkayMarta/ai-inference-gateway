package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

func RespondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
		Message:   message,
		Path:      r.URL.Path,
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				RespondError(w, r, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
