package services

import (
	"errors"
	"strings"
)

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrInsufficientBalance       = errors.New("insufficient token balance")
	ErrInvalidUserUpdate         = errors.New("invalid user update")
	ErrInvalidLoginInput         = errors.New("Enter a valid email and password.")
	ErrAccountNotFound           = errors.New("No account found with this email.")
	ErrInvalidCredentials        = errors.New("Incorrect password.")
	ErrEmailAlreadyExists        = errors.New("An account with this email already exists.")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrForbidden                 = errors.New("forbidden")
	ErrInvalidRegisterInput      = errors.New("Please enter a username, valid email, and password.")
	ErrUsernameRequired          = errors.New("Username is required.")
	ErrInvalidEmail              = errors.New("Enter a valid email address.")
	ErrWeakPassword              = errors.New("Password must be at least 8 characters long and include a letter and a number.")
	ErrInvalidAuthInput          = errors.New("invalid auth input")
	ErrInvalidToken              = errors.New("invalid token")
	ErrInvalidBillingInput       = errors.New("invalid billing input")
	ErrInvalidPasswordResetToken = errors.New("invalid or expired password reset token")
	ErrInvalidPasswordResetInput = errors.New("Enter a valid reset token and password.")
)

func isRepoNotFoundError(err error, prefix string) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(strings.ToLower(err.Error()), strings.ToLower(prefix))
}
