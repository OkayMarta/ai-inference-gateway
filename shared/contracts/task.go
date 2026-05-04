package contracts

import "time"

type TaskDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ModelID   string    `json:"modelId"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	Result    string    `json:"result,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type SubmitTaskRequestDTO struct {
	ModelID string `json:"modelId"`
	Payload string `json:"payload"`
}
