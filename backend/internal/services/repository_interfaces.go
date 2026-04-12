package services

import (
	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
)

type UserRepository interface {
	GetAll() []*models.User
	GetByID(id string) (*models.User, error)
	UpdateBalance(id string, balance float64) error
}

type ModelRepository interface {
	GetAll() []*models.AIModel
	GetByID(id string) (*models.AIModel, error)
}

type TaskRepository interface {
	GetAll() []*models.PromptTask
	GetByID(id string) (*models.PromptTask, error)
	GetByUserID(userID string) []*models.PromptTask
	Create(task *models.PromptTask)
	Complete(id, result string) error
	Fail(id, result string) error
	GetNextQueued(supportedModels []string) *models.PromptTask
}

type TransactionRepository interface {
	Create(tx *models.Transaction)
}

type WorkerRepository interface {
	GetAll() []*models.WorkerNode
	Create(w *models.WorkerNode)
	UpdateStatus(id string, status models.WorkerStatus) error
	GetIdle() []*models.WorkerNode
}

var (
	_ UserRepository        = (*repositories.UserRepository)(nil)
	_ ModelRepository       = (*repositories.ModelRepository)(nil)
	_ TaskRepository        = (*repositories.TaskRepository)(nil)
	_ TransactionRepository = (*repositories.TransactionRepository)(nil)
	_ WorkerRepository      = (*repositories.WorkerRepository)(nil)
)
