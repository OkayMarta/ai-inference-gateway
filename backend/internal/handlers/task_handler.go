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

// Структура DTO (Data Transfer Object) для прийняття тіла запиту
type submitRequest struct {
	UserID  string `json:"userId"`
	ModelID string `json:"modelId"`
	Payload string `json:"payload"`
}

// Submit обробляє POST /api/tasks
func (h *TaskHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req submitRequest
	
	// 1. Декодуємо JSON з тіла запиту
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// 2. Валідація: перевіряємо, чи всі обов'язкові поля передані
	if req.UserID == "" || req.ModelID == "" || req.Payload == "" {
		respondError(w, http.StatusBadRequest, "userId, modelId, and payload are required")
		return
	}

	// 3. Викликаємо бізнес-логіку (Сервіс)
	task, err := h.inference.SubmitPrompt(req.UserID, req.ModelID, req.Payload)
	if err != nil {
		// Якщо не вистачило грошей або юзер не знайдений — повертаємо 422 Unprocessable Entity
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	
	// 4. Повертаємо створену задачу зі статусом 201 Created
	respondJSON(w, http.StatusCreated, task)
}

// GetByID обробляє GET /api/tasks/{id}
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.inference.GetTaskByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, task)
}

// GetByUserID обробляє GET /api/tasks?userId={id}
func (h *TaskHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	// r.URL.Query().Get дістає параметри після знаку питання в URL
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "userId query parameter is required")
		return
	}
	tasks := h.inference.GetTasksByUserID(userID)
	respondJSON(w, http.StatusOK, tasks)
}