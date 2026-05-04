package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	ModelID string `json:"modelId"`
	Payload string `json:"payload"`
}

type updateTaskRequest struct {
	Payload *string `json:"payload"`
}

func (h *TaskHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req submitRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	currentUser, ok := CurrentUser(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

	if req.ModelID == "" || req.Payload == "" {
		respondError(w, r, http.StatusBadRequest, "modelId and payload are required")
		return
	}

	task, err := h.inference.SubmitPrompt(currentUser.ID, req.ModelID, req.Payload)
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

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.inference.CancelTask(id)
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
	query := r.URL.Query()

	limit := 20
	if rawLimit := query.Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 {
			respondError(w, r, http.StatusBadRequest, services.ErrInvalidPagination.Error())
			return
		}
		limit = parsedLimit
	}

	offset := 0
	if rawOffset := query.Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil || parsedOffset < 0 {
			respondError(w, r, http.StatusBadRequest, services.ErrInvalidPagination.Error())
			return
		}
		offset = parsedOffset
	}

	sort := query.Get("sort")
	if sort == "" {
		sort = "created_at_desc"
	}

	filter := services.TaskListFilter{
		UserID: query.Get("userId"),
		Status: query.Get("status"),
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	tasks, err := h.inference.ListTasks(filter)
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
