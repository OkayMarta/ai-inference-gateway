package services

import (
	"strings"

	"billing-service/internal/models"
)

type UpdateUserInput struct {
	Username *string
}

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
	if input.Username == nil {
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

	if err := s.repo.UpdateProfile(user); err != nil {
		if isRepoNotFoundError(err, "user not found:") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
