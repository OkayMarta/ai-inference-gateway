package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
)

// InferenceService — це ядро системи (білінг + оркестрація). Він об'єднує роботу одразу 4-х репозиторіїв
type InferenceService struct {
	userRepo  *repositories.UserRepository
	modelRepo *repositories.ModelRepository
	taskRepo  *repositories.TaskRepository
	txRepo    *repositories.TransactionRepository
}

func NewInferenceService(
	userRepo *repositories.UserRepository,
	modelRepo *repositories.ModelRepository,
	taskRepo *repositories.TaskRepository,
	txRepo *repositories.TransactionRepository,
) *InferenceService {
	return &InferenceService{
		userRepo:  userRepo,
		modelRepo: modelRepo,
		taskRepo:  taskRepo,
		txRepo:    txRepo,
	}
}

// generateID створює унікальний випадковий рядок для ID (замість автоінкременту БД)
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// SubmitPrompt — це головний бізнес-сценарій: перевірка балансу → списання → створення транзакції → створення задачі
func (s *InferenceService) SubmitPrompt(userID, modelID, payload string) (*models.PromptTask, error) {
	// 1. Шукаємо користувача. Якщо його немає — повертаємо помилку 404/422
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 2. Шукаємо модель. Нам потрібно знати її вартість (TokenCost)
	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return nil, fmt.Errorf("model not found: %w", err)
	}

	// 3. БІЗНЕС-ПРАВИЛО: Перевіряємо, чи вистачає у юзера грошей (токенів)
	if user.TokenBalance < model.TokenCost {
		return nil, fmt.Errorf("insufficient token balance: have %.2f, need %.2f",
			user.TokenBalance, model.TokenCost)
	}

	// 4. Списуємо токени: оновлюємо баланс користувача
	if err := s.userRepo.UpdateBalance(userID, user.TokenBalance-model.TokenCost); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// 5. Створюємо задачу (PromptTask). Статус обов'язково "Queued" (У черзі), бо ми ще не відправили її в Ollama. Це зробить фоновий Worker.
	task := &models.PromptTask{
		ID:      generateID(),
		UserID:  userID,
		ModelID: modelID,
		Payload: payload,
		Status:  models.StatusQueued,
	}
	s.taskRepo.Create(task)

	// 6. Створюємо запис про транзакцію (чек про оплату)
	tx := &models.Transaction{
		ID:     generateID(),
		UserID: userID,
		TaskID: task.ID,
		Amount: model.TokenCost,
	}
	s.txRepo.Create(tx)

	// 7. Повертаємо створену задачу (щоб контролер міг віддати її користувачу у форматі JSON)
	return task, nil
}

// GetTaskByID дозволяє перевірити статус конкретної задачі
func (s *InferenceService) GetTaskByID(id string) (*models.PromptTask, error) {
	return s.taskRepo.GetByID(id)
}

// GetTasksByUserID повертає історію задач конкретного користувача
func (s *InferenceService) GetTasksByUserID(userID string) []*models.PromptTask {
	return s.taskRepo.GetByUserID(userID)
}