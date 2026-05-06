package models

import "time"

type AIModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TokenCost   float64 `json:"tokenCost"`
}

type TaskStatus string

const (
	StatusQueued     TaskStatus = "Queued"
	StatusProcessing TaskStatus = "Processing"
	StatusCompleted  TaskStatus = "Completed"
	StatusFailed     TaskStatus = "Failed"
	StatusCancelled  TaskStatus = "Cancelled"
)

type PromptTask struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	ModelID   string     `json:"modelId"`
	Payload   string     `json:"payload"`
	Status    TaskStatus `json:"status"`
	Result    string     `json:"result,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type WorkerStatus string

const (
	WorkerIdle WorkerStatus = "Idle"
	WorkerBusy WorkerStatus = "Busy"
)

type WorkerNode struct {
	ID              string       `json:"id"`
	SupportedModels []string     `json:"supportedModels"`
	Status          WorkerStatus `json:"status"`
}
