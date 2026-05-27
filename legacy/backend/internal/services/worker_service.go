package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"ai-inference-gateway/internal/models"
)

type WorkerService struct {
	workerRepo WorkerRepository
	taskRepo   TaskRepository
	modelRepo  ModelRepository
	ollama     *OllamaClient
}

func NewWorkerService(
	workerRepo WorkerRepository,
	taskRepo TaskRepository,
	modelRepo ModelRepository,
	ollama *OllamaClient,
) *WorkerService {
	return &WorkerService{
		workerRepo: workerRepo,
		taskRepo:   taskRepo,
		modelRepo:  modelRepo,
		ollama:     ollama,
	}
}

func (s *WorkerService) Start() {
	go func() {
		log.Println("[WorkerService] Background processing started")
		for {
			s.processNext()
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

// RefreshSupportedModels синхронізує persisted mapping між воркерами та моделями.
// Для Lab 3 всі воркери підтримують усі доступні моделі, але цей зв'язок
// має зберігатися явно в PostgreSQL, а не лише припускатися в пам'яті.
func (s *WorkerService) RefreshSupportedModels() error {
	modelsList, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("load models for worker mapping refresh: %w", err)
	}

	modelIDs := make([]string, 0, len(modelsList))
	for _, model := range modelsList {
		modelIDs = append(modelIDs, model.ID)
	}

	if err := s.workerRepo.ReplaceSupportedModelsForAllWorkers(modelIDs); err != nil {
		return fmt.Errorf("replace worker model mappings: %w", err)
	}

	return nil
}

func (s *WorkerService) processNext() {
	workers, err := s.workerRepo.GetIdle()
	if err != nil {
		log.Printf("[WorkerService] Failed to load idle workers: %v", err)
		return
	}

	for _, w := range workers {
		task, err := s.taskRepo.GetNextQueued(w.SupportedModels)
		if err != nil {
			log.Printf("[WorkerService] Failed to fetch next queued task for worker %s: %v", w.ID, err)
			continue
		}
		if task == nil {
			continue
		}

		log.Printf("[WorkerService] Worker %s -> task %s (model %s)", w.ID, task.ID, task.ModelID)
		if err := s.workerRepo.UpdateStatus(w.ID, models.WorkerBusy); err != nil {
			log.Printf("[WorkerService] Failed to mark worker %s as busy: %v", w.ID, err)
			continue
		}

		go func(workerID, taskID, payload, modelID string) {
			defer func() {
				if err := s.workerRepo.UpdateStatus(workerID, models.WorkerIdle); err != nil {
					log.Printf("[WorkerService] Failed to mark worker %s as idle: %v", workerID, err)
				}
			}()

			result, err := s.executeTask(modelID, payload)
			if err != nil {
				if failErr := s.taskRepo.Fail(taskID, err.Error()); failErr != nil {
					log.Printf("[WorkerService] Worker %s failed task %s and could not persist failure: %v", workerID, taskID, failErr)
				}
				log.Printf("[WorkerService] Worker %s failed task %s: %v", workerID, taskID, err)
				return
			}

			if completeErr := s.taskRepo.Complete(taskID, result); completeErr != nil {
				log.Printf("[WorkerService] Worker %s completed task %s but could not persist completion: %v", workerID, taskID, completeErr)
				return
			}
			log.Printf("[WorkerService] Worker %s completed task %s", workerID, taskID)
		}(w.ID, task.ID, task.Payload, task.ModelID)
	}
}

func (s *WorkerService) executeTask(modelID, payload string) (string, error) {
	aiModel, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return "", fmt.Errorf("model %s not found", modelID)
	}

	if shouldSimulateFailure(payload) {
		return "", fmt.Errorf("task execution failed: simulated worker failure")
	}

	if s.ollama != nil {
		log.Printf("[WorkerService] Calling Ollama: model=%s", aiModel.Name)
		response, err := s.ollama.Generate(aiModel.Name, payload)
		if err == nil {
			return response, nil
		}
		log.Printf("[WorkerService] Ollama error: %v", err)
		return "", fmt.Errorf("task execution failed: ollama unavailable: %w", err)
	}

	return "", fmt.Errorf("task execution failed: ollama client is not configured")
}

// shouldSimulateFailure is used only for testing and demo scenarios.
func shouldSimulateFailure(payload string) bool {
	return strings.Contains(strings.ToLower(payload), strings.ToLower("__SIMULATE_FAILURE__"))
}
