package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"ai-inference-gateway/internal/models"
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
		if isRepoNotFoundError(err, "user not found:") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		if isRepoNotFoundError(err, "model not found:") {
			return nil, ErrModelNotFound
		}
		return nil, err
	}

	if user.TokenBalance < model.TokenCost {
		return nil, ErrInsufficientBalance
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
	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	tx := &models.Transaction{
		ID:     generateID(),
		UserID: userID,
		TaskID: task.ID,
		Amount: model.TokenCost,
	}
	if err := s.txRepo.Create(tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return task, nil
}

func (s *InferenceService) GetTaskByID(id string) (*models.PromptTask, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		if isRepoNotFoundError(err, "task not found:") {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return task, nil
}

func (s *InferenceService) ListTasks(filter TaskListFilter) ([]*models.PromptTask, error) {
	if filter.Limit < 0 || filter.Offset < 0 {
		return nil, ErrInvalidPagination
	}

	return s.taskRepo.List(filter)
}

func (s *InferenceService) GetAllTasks() ([]*models.PromptTask, error) {
	return s.ListTasks(TaskListFilter{})
}

func (s *InferenceService) GetTasksByUserID(userID string) ([]*models.PromptTask, error) {
	return s.ListTasks(TaskListFilter{UserID: userID})
}
