package services

import (
	"fmt"
	"log"
	"math/rand"
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
		_ = s.workerRepo.UpdateStatus(w.ID, models.WorkerBusy)

		go func(workerID, taskID, payload, modelID string) {
			defer func() {
				_ = s.workerRepo.UpdateStatus(workerID, models.WorkerIdle)
			}()

			result, err := s.executeTask(modelID, payload)
			if err != nil {
				_ = s.taskRepo.Fail(taskID, err.Error())
				log.Printf("[WorkerService] Worker %s failed task %s: %v", workerID, taskID, err)
				return
			}

			_ = s.taskRepo.Complete(taskID, result)
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
		log.Printf("[WorkerService] Ollama error, falling back to simulation: %v", err)
	}

	delay := time.Duration(2+rand.Intn(4)) * time.Second
	time.Sleep(delay)
	return fmt.Sprintf("[Симуляція] Відповідь на запит \"%s\" від моделі %s. Згенеровано за %v", payload, aiModel.Name, delay), nil
}

// shouldSimulateFailure is used only for testing and demo scenarios.
func shouldSimulateFailure(payload string) bool {
	return strings.Contains(strings.ToLower(payload), strings.ToLower("__SIMULATE_FAILURE__"))
}
