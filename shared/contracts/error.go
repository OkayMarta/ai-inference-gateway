package contracts

type ErrorResponseDTO struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}
