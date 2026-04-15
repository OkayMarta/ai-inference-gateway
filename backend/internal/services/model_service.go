package services

import "ai-inference-gateway/internal/models"

// ModelService handles model-related business logic.
type ModelService struct {
	repo ModelRepository
}

func NewModelService(repo ModelRepository) *ModelService {
	return &ModelService{repo: repo}
}

func (s *ModelService) GetAll() ([]*models.AIModel, error) {
	return s.repo.GetAll()
}

func (s *ModelService) GetByID(id string) (*models.AIModel, error) {
	model, err := s.repo.GetByID(id)
	if err != nil {
		if isRepoNotFoundError(err, "model not found:") {
			return nil, ErrModelNotFound
		}
		return nil, err
	}

	return model, nil
}
