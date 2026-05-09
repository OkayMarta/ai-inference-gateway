package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"task-service/internal/clients"
	"task-service/internal/models"
)

type BillingGateway interface {
	Charge(userID, taskID string, amount float64) error
	Refund(userID, taskID string, amount float64) error
	GetUser(userID string) (*clients.UserDTO, error)
}

type InferenceService struct {
	db        *sql.DB
	modelRepo ModelRepository
	taskRepo  TaskRepository
	billing   BillingGateway
}

func NewInferenceService(db *sql.DB, modelRepo ModelRepository, taskRepo TaskRepository, billing BillingGateway) *InferenceService {
	return &InferenceService{
		db:        db,
		modelRepo: modelRepo,
		taskRepo:  taskRepo,
		billing:   billing,
	}
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *InferenceService) SubmitPrompt(userID, modelID, payload string) (*models.PromptTask, error) {
	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		if isRepoNotFoundError(err, "model not found:") {
			return nil, ErrModelNotFound
		}
		return nil, err
	}

	if _, err := s.billing.GetUser(userID); err != nil {
		return nil, mapBillingError(err, ErrBillingUnavailable, ErrUserNotFound)
	}

	task := &models.PromptTask{
		ID:      generateID(),
		UserID:  userID,
		ModelID: modelID,
		Payload: strings.TrimSpace(payload),
		Status:  models.StatusQueued,
	}

	if err := s.billing.Charge(userID, task.ID, model.TokenCost); err != nil {
		return nil, mapBillingError(err, ErrBillingUnavailable, ErrBillingChargeFailed)
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		refundErr := s.billing.Refund(userID, task.ID, model.TokenCost)
		return nil, compensationError("begin submit prompt transaction", err, refundErr)
	}

	if err := s.taskRepo.CreateTx(tx, task); err != nil {
		_ = tx.Rollback()
		refundErr := s.billing.Refund(userID, task.ID, model.TokenCost)
		return nil, compensationError("failed to create task", err, refundErr)
	}

	if err := tx.Commit(); err != nil {
		refundErr := s.billing.Refund(userID, task.ID, model.TokenCost)
		return nil, compensationError("commit submit prompt transaction", err, refundErr)
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
	if filter.Limit <= 0 || filter.Offset < 0 {
		return nil, ErrInvalidPagination
	}

	switch filter.Sort {
	case "", "created_at_desc":
		filter.Sort = "created_at_desc"
	case "created_at_asc":
	default:
		return nil, ErrInvalidPagination
	}

	return s.taskRepo.List(filter)
}

func (s *InferenceService) UpdateTaskPayload(id string, payload string) (*models.PromptTask, error) {
	task, err := s.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	if task.Status != models.StatusQueued {
		return nil, ErrTaskCannotBeUpdated
	}

	trimmedPayload := strings.TrimSpace(payload)
	if trimmedPayload == "" {
		return nil, ErrInvalidTaskUpdate
	}

	task.Payload = trimmedPayload
	if err := s.taskRepo.Update(task); err != nil {
		if isRepoNotFoundError(err, "task not found:") {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return task, nil
}

func (s *InferenceService) CancelTask(id, userID, role string) (*models.PromptTask, error) {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: begin cancellation transaction: %v", ErrTaskCancellationFailed, err)
	}

	task, err := s.taskRepo.GetByIDForUpdateTx(tx, id)
	if err != nil {
		_ = tx.Rollback()
		if isRepoNotFoundError(err, "task not found:") {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if role != "admin" && task.UserID != userID {
		_ = tx.Rollback()
		return nil, ErrForbidden
	}

	if task.Status != models.StatusQueued {
		_ = tx.Rollback()
		return nil, ErrTaskCannotBeDeleted
	}

	model, err := s.modelRepo.GetByIDTx(tx, task.ModelID)
	if err != nil {
		_ = tx.Rollback()
		if isRepoNotFoundError(err, "model not found:") {
			return nil, ErrModelNotFound
		}
		return nil, err
	}

	if err := s.billing.Refund(task.UserID, task.ID, model.TokenCost); err != nil {
		_ = tx.Rollback()
		return nil, mapBillingError(err, ErrBillingUnavailable, ErrBillingRefundFailed)
	}

	task.Status = models.StatusCancelled
	task.Result = "Task was cancelled"
	if err := s.taskRepo.UpdateTx(tx, task); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %v", ErrTaskCancellationFailed, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%w: commit cancellation transaction: %v", ErrTaskCancellationFailed, err)
	}

	return task, nil
}

func mapBillingError(err error, unavailableError, downstreamError error) error {
	var downstream *clients.DownstreamError
	if errors.As(err, &downstream) {
		switch downstream.StatusCode {
		case 404:
			return ErrUserNotFound
		case 422:
			return ErrInsufficientBalance
		default:
			return fmt.Errorf("%w: %s", downstreamError, err.Error())
		}
	}

	if strings.Contains(strings.ToLower(err.Error()), "billing service unavailable") {
		return fmt.Errorf("%w: %v", unavailableError, err)
	}

	return fmt.Errorf("%w: %v", downstreamError, err)
}

func compensationError(action string, cause error, refundErr error) error {
	if refundErr != nil {
		return fmt.Errorf("%w: %s: %v; refund failed: %v", ErrTaskCompensationFailed, action, cause, refundErr)
	}

	return fmt.Errorf("%w: %s: %v", ErrTaskCreationFailed, action, cause)
}
