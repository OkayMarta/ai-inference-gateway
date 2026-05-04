package models

import "time"

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

type Transaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	TaskID    string    `json:"taskId"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
