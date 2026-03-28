package services

import (
	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
)

// UserService відповідає за бізнес-логіку роботи з користувачами. Поки що він просто передає виклики в репозиторій, але якщо нам знадобиться додати логіку "перевірити, чи не заблокований юзер", то її треба писати тут, а не в контролері
type UserService struct {
	repo *repositories.UserRepository
}

// Конструктор сервісу
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetByID отримує користувача за його ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

// GetAll повертає список всіх користувачів (для відображення на фронтенді)
func (s *UserService) GetAll() []*models.User {
	return s.repo.GetAll()
}