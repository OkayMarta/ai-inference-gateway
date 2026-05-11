package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"

	"task-service/internal/models"
)

const allModelsCacheKey = "ai_models:all"

type ModelCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

// ModelService handles model-related business logic.
type ModelService struct {
	repo   ModelRepository
	ollama *OllamaClient
	cache  ModelCache
}

func NewModelService(repo ModelRepository, ollama *OllamaClient, cache ModelCache) *ModelService {
	return &ModelService{
		repo:   repo,
		ollama: ollama,
		cache:  cache,
	}
}

func (s *ModelService) GetAll() ([]*models.AIModel, error) {
	ctx := context.Background()

	if s.cache != nil {
		cachedModels, err := s.cache.Get(ctx, allModelsCacheKey)
		if err == nil {
			var items []*models.AIModel
			if err := json.Unmarshal([]byte(cachedModels), &items); err == nil {
				log.Println("models cache hit: loading from Redis")
				return items, nil
			}

			log.Printf("models cache decode error: %v", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Printf("models Redis cache error: %v", err)
		}
	}

	log.Println("models cache miss: loading from PostgreSQL")
	items, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		payload, err := json.Marshal(items)
		if err != nil {
			log.Printf("models cache encode error: %v", err)
			return items, nil
		}

		if err := s.cache.Set(ctx, allModelsCacheKey, string(payload)); err != nil {
			log.Printf("models cache write error: %v", err)
		} else {
			log.Println("models cached in Redis")
		}
	}

	return items, nil
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

	s.invalidateModelsCache(context.Background())

	return nil
}

func (s *ModelService) invalidateModelsCache(ctx context.Context) {
	if s.cache == nil {
		return
	}

	if err := s.cache.Delete(ctx, allModelsCacheKey); err != nil {
		log.Printf("models cache invalidation error: %v", err)
		return
	}

	log.Println("models cache invalidated after Ollama sync")
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
