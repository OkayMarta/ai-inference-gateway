package models

import "time"

// User представляє користувача платформи з балансом токенів.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	TokenBalance float64   `json:"tokenBalance"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// AIModel представляє ШІ-модель, доступну для обробки запитів.
type AIModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TokenCost   float64 `json:"tokenCost"`
}

// TaskStatus описує поточний стан виконання завдання.
type TaskStatus string

const (
	StatusQueued     TaskStatus = "Queued"
	StatusProcessing TaskStatus = "Processing"
	StatusCompleted  TaskStatus = "Completed"
	StatusFailed     TaskStatus = "Failed"
	StatusCancelled  TaskStatus = "Cancelled"
)

// PromptTask представляє запит користувача до ШІ-моделі.
type PromptTask struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	ModelID   string     `json:"modelId"`
	Payload   string     `json:"payload"`
	Status    TaskStatus `json:"status"`
	Result    string     `json:"result,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// WorkerStatus описує стан обчислювального вузла.
type WorkerStatus string

const (
	WorkerIdle WorkerStatus = "Idle"
	WorkerBusy WorkerStatus = "Busy"
)

// WorkerNode представляє фоновий процес, який бере завдання з черги.
type WorkerNode struct {
	ID              string       `json:"id"`
	SupportedModels []string     `json:"supportedModels"`
	Status          WorkerStatus `json:"status"`
}

// Transaction фіксує фінансову операцію по задачі.
type Transaction struct {
	ID     string  `json:"id"`
	UserID string  `json:"userId"`
	TaskID string  `json:"taskId"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type,omitempty"`
}
