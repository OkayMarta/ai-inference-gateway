package services

import (
	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
)

// ModelService відповідає за роботу з ШІ-моделями (читання списку доступних)
type ModelService struct {
	repo *repositories.ModelRepository
}

// Конструктор сервісу
func NewModelService(repo *repositories.ModelRepository) *ModelService {
	return &ModelService{repo: repo}
}

// GetAll повертає всі доступні моделі
func (s *ModelService) GetAll() []*models.AIModel {
	return s.repo.GetAll()
}

// GetByID шукає конкретну модель (потрібно, щоб дізнатись її вартість)
func (s *ModelService) GetByID(id string) (*models.AIModel, error) {
	return s.repo.GetByID(id)
}