package repositories

import (
	"errors"
	"testing"
)

type fakeResult struct {
	rowsAffected int64
	err          error
}

func (r fakeResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r fakeResult) RowsAffected() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	return r.rowsAffected, nil
}

func TestEnsureRowsAffected(t *testing.T) {
	rowsAffectedErr := errors.New("driver rows affected error")

	tests := []struct {
		name        string
		result      fakeResult
		notFoundMsg string
		wantErr     bool
		wantMessage string
	}{
		{
			name:        "one affected row",
			result:      fakeResult{rowsAffected: 1},
			notFoundMsg: "task not found",
		},
		{
			name:        "zero affected rows",
			result:      fakeResult{rowsAffected: 0},
			notFoundMsg: "task not found",
			wantErr:     true,
			wantMessage: "task not found",
		},
		{
			name:        "rows affected error",
			result:      fakeResult{err: rowsAffectedErr},
			notFoundMsg: "task not found",
			wantErr:     true,
			wantMessage: "read affected rows: driver rows affected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ensureRowsAffected(tt.result, tt.notFoundMsg)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantMessage {
					t.Fatalf("error = %q, want %q", err.Error(), tt.wantMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
