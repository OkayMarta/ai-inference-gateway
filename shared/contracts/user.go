package contracts

import "time"

type UserDTO struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	TokenBalance float64   `json:"tokenBalance"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AuthResponseDTO struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}
