package services

import "testing"

func TestExportedErrorsAreSet(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrUserNotFound", err: ErrUserNotFound},
		{name: "ErrModelNotFound", err: ErrModelNotFound},
		{name: "ErrTaskNotFound", err: ErrTaskNotFound},
		{name: "ErrInsufficientBalance", err: ErrInsufficientBalance},
		{name: "ErrTaskCannotBeUpdated", err: ErrTaskCannotBeUpdated},
		{name: "ErrTaskCannotBeDeleted", err: ErrTaskCannotBeDeleted},
		{name: "ErrInvalidPagination", err: ErrInvalidPagination},
		{name: "ErrUserUpdateNotAllowed", err: ErrUserUpdateNotAllowed},
		{name: "ErrInvalidUserUpdate", err: ErrInvalidUserUpdate},
		{name: "ErrInvalidTaskUpdate", err: ErrInvalidTaskUpdate},
		{name: "ErrInvalidLoginInput", err: ErrInvalidLoginInput},
		{name: "ErrAccountNotFound", err: ErrAccountNotFound},
		{name: "ErrInvalidCredentials", err: ErrInvalidCredentials},
		{name: "ErrEmailAlreadyExists", err: ErrEmailAlreadyExists},
		{name: "ErrUnauthorized", err: ErrUnauthorized},
		{name: "ErrForbidden", err: ErrForbidden},
		{name: "ErrInvalidRegisterInput", err: ErrInvalidRegisterInput},
		{name: "ErrUsernameRequired", err: ErrUsernameRequired},
		{name: "ErrInvalidEmail", err: ErrInvalidEmail},
		{name: "ErrWeakPassword", err: ErrWeakPassword},
		{name: "ErrInvalidAuthInput", err: ErrInvalidAuthInput},
		{name: "ErrInvalidToken", err: ErrInvalidToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("error is nil")
			}
			if tt.err.Error() == "" {
				t.Fatal("error message is empty")
			}
		})
	}
}
