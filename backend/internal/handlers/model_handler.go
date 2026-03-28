package handlers

import (
	"net/http"

	"ai-inference-gateway/internal/services"
)

type ModelHandler struct {
	service *services.ModelService
}

func NewModelHandler(s *services.ModelService) *ModelHandler {
	return &ModelHandler{service: s}
}

// GetAll обробляє запит GET /api/models
func (h *ModelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.service.GetAll())
}