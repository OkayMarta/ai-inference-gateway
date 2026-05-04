package services

import (
	appdb "ai-inference-gateway/internal/db"
	"ai-inference-gateway/internal/models"
)

// TaskListFilter описує базові параметри вибірки задач.
// Структура навмисно проста: її зручно мапити як на in-memory фільтрацію, так і на майбутні SQL WHERE / ORDER BY / LIMIT / OFFSET запити.
type TaskListFilter struct {
	UserID string
	Status string
	Limit  int
	Offset int
	Sort   string
}

// UserRepository описує контракт доступу до користувачів.
// Для Lab 3 інтерфейс розширено до CRUD-операцій, оскільки PostgreSQL-backed реалізації повинні вміти не лише читати, а й зберігати та оновлювати дані. Методи повертають error, бо робота з постійним сховищем може завершитися помилкою.
type UserRepository interface {
	GetAll() ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByIDTx(tx appdb.DBTX, id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
	UpdateBalance(id string, balance float64) error
	UpdateBalanceTx(tx appdb.DBTX, id string, balance float64) error
	DeductBalanceTx(tx appdb.DBTX, id string, amount float64) error
}

// ModelRepository описує доступ до AI-моделей як до повноцінних persistent entities.
// Це готує сервісний шар до переходу з in-memory зберігання на PostgreSQL.
type ModelRepository interface {
	GetAll() ([]*models.AIModel, error)
	GetByID(id string) (*models.AIModel, error)
	GetByIDTx(tx appdb.DBTX, id string) (*models.AIModel, error)
	Create(model *models.AIModel) error
	ReplaceAll(models []*models.AIModel) error
	Update(model *models.AIModel) error
	Delete(id string) error
}

// TaskRepository уніфікує операції над задачами так, щоб вони підходили для БД, фільтрації, пагінації, сортування та подальших транзакційних сценаріїв.
// List(filter) є єдиною точкою для отримання колекцій задач замість кількох вузькоспеціалізованих методів на кшталт GetAll/GetByUserID.
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

// TransactionRepository описує мінімальний контракт для збереження білінгових подій.
// Створення транзакції теж може падати, тому повертається error.
type TransactionRepository interface {
	Create(tx *models.Transaction) error
	CreateTx(exec appdb.DBTX, tx *models.Transaction) error
}

// WorkerRepository описує доступ до воркерів у DB-friendly стилі.
// Окремі методи читання та оновлення потрібні для майбутнього worker orchestration поверх PostgreSQL, де зміни статусів і вибір вільних воркерів можуть завершуватися помилками.
type WorkerRepository interface {
	GetAll() ([]*models.WorkerNode, error)
	GetIdle() ([]*models.WorkerNode, error)
	GetByID(id string) (*models.WorkerNode, error)
	Create(worker *models.WorkerNode) error
	Update(worker *models.WorkerNode) error
	UpdateStatus(id string, status models.WorkerStatus) error
	ReplaceSupportedModelsForAllWorkers(modelIDs []string) error
	Delete(id string) error
}

// Compile-time assertions тимчасово прибрано.
// Після розширення інтерфейсів під PostgreSQL поточні in-memory реалізації
// ще не відповідають новим контрактам, а фальшиво зберігати ці перевірки не можна.
