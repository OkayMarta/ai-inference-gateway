package services

import (
	"fmt"
	"strings"

	"ai-inference-gateway/internal/models"
)

// ModelService handles model-related business logic.
type ModelService struct {
	repo   ModelRepository
	ollama *OllamaClient
}

func NewModelService(repo ModelRepository, ollama *OllamaClient) *ModelService {
	return &ModelService{
		repo:   repo,
		ollama: ollama,
	}
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

func (s *ModelService) SyncFromOllama() error {
	if s.ollama == nil {
		return fmt.Errorf("ollama client is not configured")
	}

	ollamaModels, err := s.ollama.ListModels()
	if err != nil {
		return err
	}

	syncedModels := make([]*models.AIModel, 0, len(ollamaModels))
	for _, ollamaModel := range ollamaModels {
		syncedModels = append(syncedModels, &models.AIModel{
			ID:          sanitizeModelID(ollamaModel.Name),
			Name:        ollamaModel.Name,
			Description: buildModelDescription(ollamaModel),
			TokenCost:   costByModelSize(ollamaModel.Size),
		})
	}

	if err := s.repo.ReplaceAll(syncedModels); err != nil {
		return fmt.Errorf("replace synced models: %w", err)
	}

	return nil
}

func sanitizeModelID(name string) string {
	replacer := strings.NewReplacer(":", "-", "/", "-", " ", "-")
	return strings.ToLower(replacer.Replace(name))
}

func buildModelDescription(model OllamaModelInfo) string {
	parts := []string{"Ollama model"}

	if model.Model != "" && model.Model != model.Name {
		parts = append(parts, model.Model)
	}

	if model.Size > 0 {
		parts = append(parts, formatModelSize(model.Size))
	}

	return strings.Join(parts, " · ")
}

func costByModelSize(bytes int64) float64 {
	gb := float64(bytes) / (1024 * 1024 * 1024)

	switch {
	case gb < 2:
		return 3
	case gb < 5:
		return 5
	case gb < 15:
		return 10
	default:
		return 15
	}
}

func formatModelSize(bytes int64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%.1f GB", gb)
	}

	mb := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.0f MB", mb)
}
