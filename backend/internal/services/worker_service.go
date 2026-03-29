package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
)

// WorkerService — сервіс, який працює у фоні. Він шукає задачі зі статусом Queued, бере вільного воркера і відправляє запит до Ollama
type WorkerService struct {
	workerRepo *repositories.WorkerRepository
	taskRepo   *repositories.TaskRepository
	modelRepo  *repositories.ModelRepository
	ollama     *OllamaClient
}

func NewWorkerService(
	workerRepo *repositories.WorkerRepository,
	taskRepo *repositories.TaskRepository,
	modelRepo *repositories.ModelRepository,
	ollama *OllamaClient,
) *WorkerService {
	return &WorkerService{
		workerRepo: workerRepo,
		taskRepo:   taskRepo,
		modelRepo:  modelRepo,
		ollama:     ollama,
	}
}

// Start запускає нескінченний цикл у фоновому потоці (горутині). Якби ми не написали "go func()", цей цикл заблокував би весь сервер і ми не змогли б приймати HTTP-запити.
func (s *WorkerService) Start() {
	go func() {
		log.Println("[WorkerService] Background processing started")
		for {
			s.processNext()
			// Робимо паузу 0.5 секунди, щоб не "спалити" процесор нескінченним пошуком
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

// processNext перевіряє, чи є вільні воркери і чи є для них робота
func (s *WorkerService) processNext() {
	workers := s.workerRepo.GetIdle() // Шукаємо вільних (Idle)
	for _, w := range workers {
		// Шукаємо найстарішу задачу в черзі, яку підтримує цей воркер
		task := s.taskRepo.GetNextQueued(w.SupportedModels)
		if task == nil {
			continue // Немає задач — йдемо далі
		}

		log.Printf("[WorkerService] Worker %s → task %s (model %s)", w.ID, task.ID, task.ModelID)
		
		// Ставимо воркеру статус "Зайнятий"
		_ = s.workerRepo.UpdateStatus(w.ID, models.WorkerBusy)

		// ЗАПУСКАЄМО ЩЕ ОДНУ ГОРУТИНУ!
		// Генерація тексту може тривати хвилину. Ми не хочемо, щоб через це інші воркери чекали. Тому сама генерація теж йде в окремий фоновий потік
		go func(workerID, taskID, payload, modelID string) {
			result := s.executeTask(modelID, payload) // Йдемо в Ollama або симулюємо
			_ = s.taskRepo.Complete(taskID, result)   // Зберігаємо результат і статус Completed
			_ = s.workerRepo.UpdateStatus(workerID, models.WorkerIdle) // Звільняємо воркера
			log.Printf("[WorkerService] Worker %s completed task %s", workerID, taskID)
		}(w.ID, task.ID, task.Payload, task.ModelID)
	}
}

// executeTask викликає реальну Ollama або повертає симульовану відповідь (якщо Ollama немає)
func (s *WorkerService) executeTask(modelID, payload string) string {
	aiModel, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return fmt.Sprintf("Error: model %s not found", modelID)
	}

	// Якщо Ollama підключена — робимо реальний запит
	if s.ollama != nil {
		log.Printf("[WorkerService] Calling Ollama: model=%s", aiModel.Name)
		response, err := s.ollama.Generate(aiModel.Name, payload)
		if err == nil {
			return response
		}
		log.Printf("[WorkerService] Ollama error, falling back to simulation: %v", err)
	}

	// Fallback: якщо Ollama не встановлена на ПК, ми просто імітуємо затримку від 2 до 5 секунд і віддаємо тестовий текст
	delay := time.Duration(2+rand.Intn(4)) * time.Second
	time.Sleep(delay)
	return fmt.Sprintf("[Симуляція] Відповідь на запит \"%s\" від моделі %s. Згенеровано за %v", payload, aiModel.Name, delay)
}