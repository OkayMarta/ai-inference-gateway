package contracts

type ModelDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TokenCost   float64 `json:"tokenCost"`
	IsActive    bool    `json:"isActive"`
}
