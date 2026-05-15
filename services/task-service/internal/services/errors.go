package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrModelNotFound          = errors.New("model not found")
	ErrTaskNotFound           = errors.New("task not found")
	ErrInsufficientBalance    = errors.New("insufficient token balance")
	ErrTaskCannotBeUpdated    = errors.New("task cannot be updated")
	ErrTaskCannotBeDeleted    = errors.New("task cannot be deleted")
	ErrInvalidPagination      = errors.New("invalid pagination")
	ErrUserUpdateNotAllowed   = errors.New("user update not allowed")
	ErrInvalidUserUpdate      = errors.New("invalid user update")
	ErrInvalidTaskUpdate      = errors.New("invalid task update")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrInvalidRegisterInput   = errors.New("invalid registration input")
	ErrInvalidAuthInput       = errors.New("invalid auth input")
	ErrInvalidToken           = errors.New("invalid token")
	ErrBillingUnavailable     = errors.New("billing service unavailable")
	ErrBillingLookupFailed    = errors.New("billing lookup failed")
	ErrBillingChargeFailed    = errors.New("billing charge failed")
	ErrBillingRefundFailed    = errors.New("billing refund failed")
	ErrTaskCreationFailed     = errors.New("task creation failed after billing charge; refund compensation succeeded")
	ErrTaskCompensationFailed = errors.New("task creation failed after billing charge; refund compensation failed")
	ErrTaskCancellationFailed = errors.New("billing refund succeeded but task cancellation failed")
)

type DownstreamError struct {
	StatusCode int
	Message    string
}

func (e *DownstreamError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if statusText := http.StatusText(e.StatusCode); statusText != "" {
		return statusText
	}
	return fmt.Sprintf("downstream service returned status %d", e.StatusCode)
}

// Репозиторії поки повертають текстові not found помилки, тому тимчасово
// нормалізуємо їх тут до доменних помилок сервісного шару.
func isRepoNotFoundError(err error, prefix string) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(strings.ToLower(err.Error()), strings.ToLower(prefix))
}
