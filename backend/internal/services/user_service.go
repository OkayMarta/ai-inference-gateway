package services

import "ai-inference-gateway/internal/models"

// UserService handles user-related business logic.
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetAll() ([]*models.User, error) {
	return s.repo.GetAll()
}
