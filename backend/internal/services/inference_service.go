package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"ai-inference-gateway/internal/models"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrModelNotFound       = errors.New("model not found")
	ErrInsufficientBalance = errors.New("insufficient token balance")
)

// InferenceService coordinates billing and task orchestration.
type InferenceService struct {
	userRepo  UserRepository
	modelRepo ModelRepository
	taskRepo  TaskRepository
	txRepo    TransactionRepository
}

func NewInferenceService(
	userRepo UserRepository,
	modelRepo ModelRepository,
	taskRepo TaskRepository,
	txRepo TransactionRepository,
) *InferenceService {
	return &InferenceService{
		userRepo:  userRepo,
		modelRepo: modelRepo,
		taskRepo:  taskRepo,
		txRepo:    txRepo,
	}
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *InferenceService) SubmitPrompt(userID, modelID, payload string) (*models.PromptTask, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUserNotFound, userID)
	}

	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrModelNotFound, modelID)
	}

	if user.TokenBalance < model.TokenCost {
		return nil, fmt.Errorf("%w: have %.2f, need %.2f", ErrInsufficientBalance,
			user.TokenBalance, model.TokenCost)
	}

	if err := s.userRepo.UpdateBalance(userID, user.TokenBalance-model.TokenCost); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	task := &models.PromptTask{
		ID:      generateID(),
		UserID:  userID,
		ModelID: modelID,
		Payload: payload,
		Status:  models.StatusQueued,
	}
	s.taskRepo.Create(task)

	tx := &models.Transaction{
		ID:     generateID(),
		UserID: userID,
		TaskID: task.ID,
		Amount: model.TokenCost,
	}
	s.txRepo.Create(tx)

	return task, nil
}

func (s *InferenceService) GetTaskByID(id string) (*models.PromptTask, error) {
	return s.taskRepo.GetByID(id)
}

func (s *InferenceService) GetAllTasks() []*models.PromptTask {
	return s.taskRepo.GetAll()
}

func (s *InferenceService) GetTasksByUserID(userID string) []*models.PromptTask {
	return s.taskRepo.GetByUserID(userID)
}
