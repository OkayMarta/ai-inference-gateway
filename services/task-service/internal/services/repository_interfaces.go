package services

import (
	appdb "task-service/internal/db"
	"task-service/internal/models"
)

type TaskListFilter struct {
	UserID string
	Status string
	Limit  int
	Offset int
	Sort   string
}

type ModelRepository interface {
	GetAll() ([]*models.AIModel, error)
	GetByID(id string) (*models.AIModel, error)
	GetByIDTx(tx appdb.DBTX, id string) (*models.AIModel, error)
	Create(model *models.AIModel) error
	ReplaceAll(models []*models.AIModel) error
	Update(model *models.AIModel) error
	Delete(id string) error
}

type TaskRepository interface {
	GetByID(id string) (*models.PromptTask, error)
	GetByIDTx(tx appdb.DBTX, id string) (*models.PromptTask, error)
	List(filter TaskListFilter) ([]*models.PromptTask, error)
	Create(task *models.PromptTask) error
	CreateTx(tx appdb.DBTX, task *models.PromptTask) error
	Update(task *models.PromptTask) error
	UpdateTx(tx appdb.DBTX, task *models.PromptTask) error
	Delete(id string) error
	Complete(id, result string) error
	Fail(id, result string) error
	GetNextQueued(supportedModels []string) (*models.PromptTask, error)
}

type WorkerRepository interface {
	GetAll() ([]*models.WorkerNode, error)
	GetIdle() ([]*models.WorkerNode, error)
	GetByID(id string) (*models.WorkerNode, error)
	EnsureDefaultWorker(id string) error
	Create(worker *models.WorkerNode) error
	Update(worker *models.WorkerNode) error
	UpdateStatus(id string, status models.WorkerStatus) error
	ReplaceSupportedModelsForAllWorkers(modelIDs []string) error
	Delete(id string) error
}
