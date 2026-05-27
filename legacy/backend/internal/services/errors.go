package services

import (
	"errors"
	"strings"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrModelNotFound        = errors.New("model not found")
	ErrTaskNotFound         = errors.New("task not found")
	ErrInsufficientBalance  = errors.New("insufficient token balance")
	ErrTaskCannotBeUpdated  = errors.New("task cannot be updated")
	ErrTaskCannotBeDeleted  = errors.New("task cannot be deleted")
	ErrInvalidPagination    = errors.New("invalid pagination")
	ErrUserUpdateNotAllowed = errors.New("user update not allowed")
	ErrInvalidUserUpdate    = errors.New("invalid user update")
	ErrInvalidTaskUpdate    = errors.New("invalid task update")
	ErrInvalidLoginInput    = errors.New("Enter a valid email and password.")
	ErrAccountNotFound      = errors.New("No account found with this email.")
	ErrInvalidCredentials   = errors.New("Incorrect password.")
	ErrEmailAlreadyExists   = errors.New("An account with this email already exists.")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidRegisterInput = errors.New("Please enter a username, valid email, and password.")
	ErrUsernameRequired     = errors.New("Username is required.")
	ErrInvalidEmail         = errors.New("Enter a valid email address.")
	ErrWeakPassword         = errors.New("Password must be at least 8 characters long and include a letter and a number.")
	ErrInvalidAuthInput     = errors.New("invalid auth input")
	ErrInvalidToken         = errors.New("invalid token")
)

// Репозиторії поки повертають текстові not found помилки, тому тимчасово
// нормалізуємо їх тут до доменних помилок сервісного шару.
func isRepoNotFoundError(err error, prefix string) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(strings.ToLower(err.Error()), strings.ToLower(prefix))
}
