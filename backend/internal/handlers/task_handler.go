package handlers

import (
	"encoding/json"
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

type updateTaskRequest struct {
	Payload *string `json:"payload"`
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
		status := mapErrorToStatus(err)
		message := err.Error()
		if status == http.StatusInternalServerError {
			message = "internal server error"
		}
		respondError(w, r, status, message)
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.inference.GetTaskByID(id)
	if err != nil {
		status := mapErrorToStatus(err)
		message := err.Error()
		if status == http.StatusInternalServerError {
			message = "internal server error"
		}
		respondError(w, r, status, message)
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Payload == nil {
		respondError(w, r, http.StatusBadRequest, services.ErrInvalidTaskUpdate.Error())
		return
	}

	task, err := h.inference.UpdateTaskPayload(id, *req.Payload)
	if err != nil {
		status := mapErrorToStatus(err)
		message := err.Error()
		if status == http.StatusInternalServerError {
			message = "internal server error"
		}
		respondError(w, r, status, message)
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		tasks, err := h.inference.GetAllTasks()
		if err != nil {
			status := mapErrorToStatus(err)
			message := err.Error()
			if status == http.StatusInternalServerError {
				message = "internal server error"
			}
			respondError(w, r, status, message)
			return
		}

		respondJSON(w, http.StatusOK, tasks)
		return
	}

	tasks, err := h.inference.GetTasksByUserID(userID)
	if err != nil {
		status := mapErrorToStatus(err)
		message := err.Error()
		if status == http.StatusInternalServerError {
			message = "internal server error"
		}
		respondError(w, r, status, message)
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}
