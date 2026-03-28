package models

import "time"

// User представляє користувача платформи з балансом токенів
type User struct {
	ID           string  `json:"id"`
	Username     string  `json:"username"`
	TokenBalance float64 `json:"tokenBalance"`
}

// AIModel представляє ШІ-модель, доступну для обробки запитів
type AIModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TokenCost   float64 `json:"tokenCost"` // Вартість одного виклику в токенах
}

// TaskStatus описує поточний стан виконання завдання
type TaskStatus string

const (
	StatusQueued     TaskStatus = "Queued"     // У черзі
	StatusProcessing TaskStatus = "Processing" // Обробляється воркером
	StatusCompleted  TaskStatus = "Completed"  // Успішно завершено
	StatusFailed     TaskStatus = "Failed"     // Помилка генерації
)

// PromptTask представляє запит користувача до ШІ-моделі (доменна сутність)
type PromptTask struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	ModelID   string     `json:"modelId"`
	Payload   string     `json:"payload"`           // Сам текст запиту (промпт)
	Status    TaskStatus `json:"status"`
	Result    string     `json:"result,omitempty"`  // Відповідь моделі (не показується, якщо пуста)
	CreatedAt time.Time  `json:"createdAt"`
}

// WorkerStatus описує стан обчислювального вузла (воркера)
type WorkerStatus string

const (
	WorkerIdle WorkerStatus = "Idle" // Вільний
	WorkerBusy WorkerStatus = "Busy" // Зайнятий обробкою
)

// WorkerNode представляє фоновий процес, який бере завдання з черги та виконує їх
type WorkerNode struct {
	ID              string       `json:"id"`
	SupportedModels []string     `json:"supportedModels"` // Які моделі цей воркер може обробити
	Status          WorkerStatus `json:"status"`
}

// Transaction фіксує факт списання токенів у користувача за конкретне завдання
type Transaction struct {
	ID     string  `json:"id"`
	UserID string  `json:"userId"`
	TaskID string  `json:"taskId"`
	Amount float64 `json:"amount"` // Скільки токенів було списано
}