package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai-inference-gateway/internal/services"
)

func TestMapErrorToStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "user not found", err: services.ErrUserNotFound, want: http.StatusNotFound},
		{name: "task not found", err: services.ErrTaskNotFound, want: http.StatusNotFound},
		{name: "model not found", err: services.ErrModelNotFound, want: http.StatusNotFound},
		{name: "insufficient balance", err: services.ErrInsufficientBalance, want: http.StatusUnprocessableEntity},
		{name: "task cannot be updated", err: services.ErrTaskCannotBeUpdated, want: http.StatusConflict},
		{name: "task cannot be deleted", err: services.ErrTaskCannotBeDeleted, want: http.StatusConflict},
		{name: "invalid pagination", err: services.ErrInvalidPagination, want: http.StatusBadRequest},
		{name: "user update not allowed", err: services.ErrUserUpdateNotAllowed, want: http.StatusForbidden},
		{name: "invalid user update", err: services.ErrInvalidUserUpdate, want: http.StatusBadRequest},
		{name: "invalid task update", err: services.ErrInvalidTaskUpdate, want: http.StatusBadRequest},
		{name: "invalid login input", err: services.ErrInvalidLoginInput, want: http.StatusBadRequest},
		{name: "account not found", err: services.ErrAccountNotFound, want: http.StatusNotFound},
		{name: "invalid credentials", err: services.ErrInvalidCredentials, want: http.StatusUnauthorized},
		{name: "unauthorized", err: services.ErrUnauthorized, want: http.StatusUnauthorized},
		{name: "forbidden", err: services.ErrForbidden, want: http.StatusForbidden},
		{name: "email already exists", err: services.ErrEmailAlreadyExists, want: http.StatusConflict},
		{name: "invalid register input", err: services.ErrInvalidRegisterInput, want: http.StatusBadRequest},
		{name: "username required", err: services.ErrUsernameRequired, want: http.StatusBadRequest},
		{name: "invalid email", err: services.ErrInvalidEmail, want: http.StatusBadRequest},
		{name: "weak password", err: services.ErrWeakPassword, want: http.StatusBadRequest},
		{name: "unknown error", err: errors.New("unexpected failure"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapErrorToStatus(tt.err); got != tt.want {
				t.Fatalf("mapErrorToStatus(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/task-123", nil)
	rec := httptest.NewRecorder()

	respondError(rec, req, http.StatusBadRequest, "bad request")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}

	var response ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response.Timestamp == "" {
		t.Fatal("timestamp is empty")
	}
	if _, err := time.Parse(time.RFC3339, response.Timestamp); err != nil {
		t.Fatalf("timestamp %q is not RFC3339: %v", response.Timestamp, err)
	}
	if response.Status != http.StatusBadRequest {
		t.Fatalf("response status = %d, want %d", response.Status, http.StatusBadRequest)
	}
	if response.Message != "bad request" {
		t.Fatalf("response message = %q, want %q", response.Message, "bad request")
	}
	if response.Path != req.URL.Path {
		t.Fatalf("response path = %q, want %q", response.Path, req.URL.Path)
	}
}
