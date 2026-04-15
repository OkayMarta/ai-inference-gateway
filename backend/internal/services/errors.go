package services

import (
	"errors"
	"strings"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrModelNotFound       = errors.New("model not found")
	ErrTaskNotFound        = errors.New("task not found")
	ErrInsufficientBalance = errors.New("insufficient token balance")
	ErrTaskCannotBeUpdated = errors.New("task cannot be updated")
	ErrTaskCannotBeDeleted = errors.New("task cannot be deleted")
	ErrInvalidPagination   = errors.New("invalid pagination")
	ErrUserUpdateNotAllowed = errors.New("user update not allowed")
	ErrInvalidUserUpdate    = errors.New("invalid user update")
	ErrInvalidTaskUpdate    = errors.New("invalid task update")
)

// Репозиторії поки повертають текстові not found помилки, тому тимчасово
// нормалізуємо їх тут до доменних помилок сервісного шару.
func isRepoNotFoundError(err error, prefix string) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(strings.ToLower(err.Error()), strings.ToLower(prefix))
}
