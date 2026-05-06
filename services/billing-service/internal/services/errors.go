package services

import (
	"errors"
	"strings"
)

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrInsufficientBalance       = errors.New("insufficient token balance")
	ErrInvalidUserUpdate         = errors.New("invalid user update")
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrForbidden                 = errors.New("forbidden")
	ErrInvalidRegisterInput      = errors.New("invalid registration input")
	ErrInvalidAuthInput          = errors.New("invalid auth input")
	ErrInvalidToken              = errors.New("invalid token")
	ErrInvalidBillingInput       = errors.New("invalid billing input")
	ErrInvalidPasswordResetToken = errors.New("invalid or expired password reset token")
	ErrInvalidPasswordResetInput = errors.New("invalid password reset input")
)

func isRepoNotFoundError(err error, prefix string) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(strings.ToLower(err.Error()), strings.ToLower(prefix))
}
