package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"ai-inference-gateway/internal/services"
)

type TaskHandler struct {
	inference *services.InferenceService
}

func NewTaskHandler(s *services.InferenceService) *TaskHandler {
	return &TaskHandler{inference: s}
}

type submitRequest struct {
	UserID  string `json:"userId"`
	ModelID string `json:"modelId"`
	Payload string `json:"payload"`
}

func (h *TaskHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req submitRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == "" || req.ModelID == "" || req.Payload == "" {
		respondError(w, r, http.StatusBadRequest, "userId, modelId, and payload are required")
		return
	}

	task, err := h.inference.SubmitPrompt(req.UserID, req.ModelID, req.Payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound), errors.Is(err, services.ErrModelNotFound):
			respondError(w, r, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrInsufficientBalance):
			respondError(w, r, http.StatusUnprocessableEntity, err.Error())
		default:
			respondError(w, r, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.inference.GetTaskByID(id)
	if err != nil {
		respondError(w, r, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		tasks, err := h.inference.GetAllTasks()
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}

		respondJSON(w, http.StatusOK, tasks)
		return
	}

	tasks, err := h.inference.GetTasksByUserID(userID)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}
