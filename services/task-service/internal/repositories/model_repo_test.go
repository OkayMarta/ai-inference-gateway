package repositories

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestModelRepositoryGetAllReturnsOnlyActiveModels(t *testing.T) {
	db, mock := newModelRepoDB(t)
	repo := NewModelRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "token_cost"}).
		AddRow("active-model", "Active Model", "active", 5.0)

	mock.ExpectQuery(`SELECT id, name, description, token_cost\s+FROM ai_models\s+WHERE is_active = TRUE\s+ORDER BY name, id`).
		WillReturnRows(rows)

	models, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(models) != 1 || models[0].ID != "active-model" {
		t.Fatalf("expected only active model row, got %#v", models)
	}
	assertModelRepoExpectations(t, mock)
}

func TestModelRepositoryGetByIDReturnsActiveModel(t *testing.T) {
	db, mock := newModelRepoDB(t)
	repo := NewModelRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "token_cost"}).
		AddRow("active-model", "Active Model", "active", 5.0)

	mock.ExpectQuery(`SELECT id, name, description, token_cost\s+FROM ai_models\s+WHERE id = \$1\s+AND is_active = TRUE`).
		WithArgs("active-model").
		WillReturnRows(rows)

	model, err := repo.GetByID("active-model")
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if model.ID != "active-model" {
		t.Fatalf("expected active model, got %#v", model)
	}
	assertModelRepoExpectations(t, mock)
}

func TestModelRepositoryGetByIDDoesNotReturnInactiveModel(t *testing.T) {
	db, mock := newModelRepoDB(t)
	repo := NewModelRepository(db)

	mock.ExpectQuery(`SELECT id, name, description, token_cost\s+FROM ai_models\s+WHERE id = \$1\s+AND is_active = TRUE`).
		WithArgs("inactive-model").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID("inactive-model")
	if err == nil || !strings.Contains(err.Error(), "model not found: inactive-model") {
		t.Fatalf("expected model not found for inactive model, got %v", err)
	}
	assertModelRepoExpectations(t, mock)
}

func newModelRepoDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db, mock
}

func assertModelRepoExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
