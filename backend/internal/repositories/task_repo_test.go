package repositories

import (
	"errors"
	"testing"

	"ai-inference-gateway/internal/services"
)

func TestTaskSortClause(t *testing.T) {
	tests := []struct {
		name    string
		sort    string
		want    string
		wantErr error
	}{
		{name: "empty", sort: "", want: "created_at DESC"},
		{name: "created at desc", sort: "created_at_desc", want: "created_at DESC"},
		{name: "created at asc", sort: "created_at_asc", want: "created_at ASC"},
		{name: "mixed case desc", sort: "Created_At_Desc", want: "created_at DESC"},
		{name: "uppercase asc", sort: "CREATED_AT_ASC", want: "created_at ASC"},
		{name: "invalid", sort: "status_desc", wantErr: services.ErrInvalidPagination},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := taskSortClause(tt.sort)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("taskSortClause(%q) = %q, want %q", tt.sort, got, tt.want)
			}
		})
	}
}
