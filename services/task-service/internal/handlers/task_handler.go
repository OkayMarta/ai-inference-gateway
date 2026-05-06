package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"task-service/internal/models"
	"task-service/internal/services"
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
	userID, _, ok := gatewayIdentity(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

	var req submitRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ModelID == "" || req.Payload == "" {
		respondError(w, r, http.StatusBadRequest, "modelId and payload are required")
		return
	}

	task, err := h.inference.SubmitPrompt(userID, req.ModelID, req.Payload)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := gatewayIdentity(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

	task, err := h.inference.GetTaskByID(chi.URLParam(r, "id"))
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	if !canAccessTask(userID, role, task) {
		respondError(w, r, http.StatusForbidden, services.ErrForbidden.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := gatewayIdentity(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

	id := chi.URLParam(r, "id")
	task, err := h.inference.GetTaskByID(id)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}
	if !canAccessTask(userID, role, task) {
		respondError(w, r, http.StatusForbidden, services.ErrForbidden.Error())
		return
	}

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

	updatedTask, err := h.inference.UpdateTaskPayload(id, *req.Payload)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := gatewayIdentity(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

	id := chi.URLParam(r, "id")
	task, err := h.inference.GetTaskByID(id)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}
	if !canAccessTask(userID, role, task) {
		respondError(w, r, http.StatusForbidden, services.ErrForbidden.Error())
		return
	}

	cancelledTask, err := h.inference.CancelTask(id)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, cancelledTask)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, role, ok := gatewayIdentity(r)
	if !ok {
		respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
		return
	}

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

	filterUserID := userID
	if role == "admin" && query.Get("userId") != "" {
		filterUserID = query.Get("userId")
	}

	sort := query.Get("sort")
	if sort == "" {
		sort = "created_at_desc"
	}

	tasks, err := h.inference.ListTasks(services.TaskListFilter{
		UserID: filterUserID,
		Status: query.Get("status"),
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	})
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}

func gatewayIdentity(r *http.Request) (string, string, bool) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return "", "", false
	}
	return userID, r.Header.Get("X-User-Role"), true
}

func canAccessTask(userID, role string, task *models.PromptTask) bool {
	return role == "admin" || task.UserID == userID
}
