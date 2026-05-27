package services

import (
	"strings"

	"ai-inference-gateway/internal/models"
)

type UpdateUserInput struct {
	Username     *string
	TokenBalance *float64
}

// UserService handles user-related business logic.
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if isRepoNotFoundError(err, "user not found:") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAll() ([]*models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) Update(id string, input UpdateUserInput) (*models.User, error) {
	// Для Lab 3 редагування tokenBalance лишаємо як адміністративну/testing операцію.
	// Звичайні зміни балансу й далі мають відбуватись через бізнес-процеси платформи.
	if input.Username == nil && input.TokenBalance == nil {
		return nil, ErrInvalidUserUpdate
	}

	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != nil {
		trimmedUsername := strings.TrimSpace(*input.Username)
		if trimmedUsername == "" {
			return nil, ErrInvalidUserUpdate
		}
		user.Username = trimmedUsername
	}

	if input.TokenBalance != nil {
		if *input.TokenBalance < 0 {
			return nil, ErrInvalidUserUpdate
		}
		user.TokenBalance = *input.TokenBalance
	}

	if err := s.repo.Update(user); err != nil {
		if isRepoNotFoundError(err, "user not found:") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
