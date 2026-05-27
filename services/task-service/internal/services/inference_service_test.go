package services

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"task-service/internal/clients"
	appdb "task-service/internal/db"
	"task-service/internal/models"
)

type fakeModelRepo struct {
	model      *models.AIModel
	getByIDErr error
}

func (r *fakeModelRepo) GetAll() ([]*models.AIModel, error) { return nil, nil }
func (r *fakeModelRepo) GetByID(id string) (*models.AIModel, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	if r.model == nil || r.model.ID != id {
		return nil, fmt.Errorf("model not found: %s", id)
	}
	return cloneModel(r.model), nil
}
func (r *fakeModelRepo) GetByIDTx(tx appdb.DBTX, id string) (*models.AIModel, error) {
	return r.GetByID(id)
}
func (r *fakeModelRepo) Create(model *models.AIModel) error        { return nil }
func (r *fakeModelRepo) ReplaceAll(models []*models.AIModel) error { return nil }
func (r *fakeModelRepo) Update(model *models.AIModel) error        { return nil }
func (r *fakeModelRepo) Delete(id string) error                    { return nil }

type fakeTaskRepo struct {
	task      *models.PromptTask
	updated   *models.PromptTask
	updateErr error
}

func (r *fakeTaskRepo) GetByID(id string) (*models.PromptTask, error) {
	if r.task == nil || r.task.ID != id {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return cloneTask(r.task), nil
}
func (r *fakeTaskRepo) GetByIDTx(tx appdb.DBTX, id string) (*models.PromptTask, error) {
	return r.GetByID(id)
}
func (r *fakeTaskRepo) GetByIDForUpdateTx(tx appdb.DBTX, id string) (*models.PromptTask, error) {
	return r.GetByID(id)
}
func (r *fakeTaskRepo) List(filter TaskListFilter) ([]*models.PromptTask, error) { return nil, nil }
func (r *fakeTaskRepo) Create(task *models.PromptTask) error                     { return nil }
func (r *fakeTaskRepo) CreateTx(tx appdb.DBTX, task *models.PromptTask) error    { return nil }
func (r *fakeTaskRepo) Update(task *models.PromptTask) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	r.updated = cloneTask(task)
	return nil
}
func (r *fakeTaskRepo) UpdateTx(tx appdb.DBTX, task *models.PromptTask) error {
	return r.Update(task)
}
func (r *fakeTaskRepo) Delete(id string) error           { return nil }
func (r *fakeTaskRepo) Complete(id, result string) error { return nil }
func (r *fakeTaskRepo) Fail(id, result string) error     { return nil }
func (r *fakeTaskRepo) GetNextQueued(supportedModels []string) (*models.PromptTask, error) {
	return nil, nil
}

type fakeBillingGateway struct {
	chargeCalls  int
	getUserCalls int
	refundCalls  int
	refundUserID string
	refundTaskID string
	refundAmount float64
}

func (g *fakeBillingGateway) Charge(userID, taskID string, amount float64) error {
	g.chargeCalls++
	return nil
}
func (g *fakeBillingGateway) Refund(userID, taskID string, amount float64) error {
	g.refundCalls++
	g.refundUserID = userID
	g.refundTaskID = taskID
	g.refundAmount = amount
	return nil
}
func (g *fakeBillingGateway) GetUser(userID string) (*clients.UserDTO, error) {
	g.getUserCalls++
	return &clients.UserDTO{ID: userID}, nil
}

func TestSubmitPromptFailsForInactiveModel(t *testing.T) {
	modelRepo := &fakeModelRepo{getByIDErr: fmt.Errorf("model not found: inactive-model")}
	taskRepo := &fakeTaskRepo{}
	billing := &fakeBillingGateway{}
	svc := NewInferenceService(nil, modelRepo, taskRepo, billing)

	_, err := svc.SubmitPrompt("user-1", "inactive-model", "prompt")
	if !errors.Is(err, ErrModelNotFound) {
		t.Fatalf("expected ErrModelNotFound for inactive model, got %v", err)
	}
	if billing.getUserCalls != 0 || billing.chargeCalls != 0 {
		t.Fatalf("expected no billing calls for inactive model, got getUser=%d charge=%d", billing.getUserCalls, billing.chargeCalls)
	}
}
func TestUpdateTaskPayloadOwnerCanUpdateQueuedTask(t *testing.T) {
	taskRepo := &fakeTaskRepo{task: testTask(models.StatusQueued)}
	svc := NewInferenceService(nil, &fakeModelRepo{}, taskRepo, &fakeBillingGateway{})

	updated, err := svc.UpdateTaskPayload("task-1", "user-1", " updated payload ")
	if err != nil {
		t.Fatalf("UpdateTaskPayload returned error: %v", err)
	}

	if updated.Payload != "updated payload" {
		t.Fatalf("expected trimmed payload, got %q", updated.Payload)
	}
	if taskRepo.updated == nil || taskRepo.updated.Payload != "updated payload" {
		t.Fatalf("expected repository update with new payload")
	}
}

func TestCancelTaskOwnerCanCancelQueuedTask(t *testing.T) {
	db, mock := newTxDB(t)
	mock.ExpectBegin()
	mock.ExpectCommit()

	taskRepo := &fakeTaskRepo{task: testTask(models.StatusQueued)}
	modelRepo := &fakeModelRepo{model: &models.AIModel{ID: "model-1", TokenCost: 2.5}}
	billing := &fakeBillingGateway{}
	svc := NewInferenceService(db, modelRepo, taskRepo, billing)

	cancelled, err := svc.CancelTask("task-1", "user-1")
	if err != nil {
		t.Fatalf("CancelTask returned error: %v", err)
	}

	if cancelled.Status != models.StatusCancelled {
		t.Fatalf("expected cancelled status, got %s", cancelled.Status)
	}
	if taskRepo.updated == nil || taskRepo.updated.Status != models.StatusCancelled {
		t.Fatalf("expected repository update with cancelled status")
	}
	if billing.refundCalls != 1 || billing.refundUserID != "user-1" || billing.refundTaskID != "task-1" || billing.refundAmount != 2.5 {
		t.Fatalf("expected one refund for owner task, got calls=%d user=%q task=%q amount=%v", billing.refundCalls, billing.refundUserID, billing.refundTaskID, billing.refundAmount)
	}
	assertSQLExpectations(t, mock)
}

func TestUpdateTaskPayloadRejectsDifferentUser(t *testing.T) {
	taskRepo := &fakeTaskRepo{task: testTask(models.StatusQueued)}
	svc := NewInferenceService(nil, &fakeModelRepo{}, taskRepo, &fakeBillingGateway{})

	_, err := svc.UpdateTaskPayload("task-1", "user-2", "updated payload")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if taskRepo.updated != nil {
		t.Fatalf("expected no repository update for non-owner")
	}
}

func TestCancelTaskRejectsDifferentUser(t *testing.T) {
	db, mock := newTxDB(t)
	mock.ExpectBegin()
	mock.ExpectRollback()

	taskRepo := &fakeTaskRepo{task: testTask(models.StatusQueued)}
	billing := &fakeBillingGateway{}
	svc := NewInferenceService(db, &fakeModelRepo{model: &models.AIModel{ID: "model-1", TokenCost: 2.5}}, taskRepo, billing)

	_, err := svc.CancelTask("task-1", "user-2")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if taskRepo.updated != nil {
		t.Fatalf("expected no repository update for non-owner")
	}
	if billing.refundCalls != 0 {
		t.Fatalf("expected no refund for non-owner, got %d", billing.refundCalls)
	}
	assertSQLExpectations(t, mock)
}

func TestOwnerCannotUpdateOrCancelNonQueuedTask(t *testing.T) {
	t.Run("completed task cannot be updated", func(t *testing.T) {
		taskRepo := &fakeTaskRepo{task: testTask(models.StatusCompleted)}
		svc := NewInferenceService(nil, &fakeModelRepo{}, taskRepo, &fakeBillingGateway{})

		_, err := svc.UpdateTaskPayload("task-1", "user-1", "updated payload")
		if !errors.Is(err, ErrTaskCannotBeUpdated) {
			t.Fatalf("expected ErrTaskCannotBeUpdated, got %v", err)
		}
		if taskRepo.updated != nil {
			t.Fatalf("expected no repository update for completed task")
		}
	})

	t.Run("processing task cannot be cancelled", func(t *testing.T) {
		db, mock := newTxDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()

		taskRepo := &fakeTaskRepo{task: testTask(models.StatusProcessing)}
		billing := &fakeBillingGateway{}
		svc := NewInferenceService(db, &fakeModelRepo{model: &models.AIModel{ID: "model-1", TokenCost: 2.5}}, taskRepo, billing)

		_, err := svc.CancelTask("task-1", "user-1")
		if !errors.Is(err, ErrTaskCannotBeDeleted) {
			t.Fatalf("expected ErrTaskCannotBeDeleted, got %v", err)
		}
		if taskRepo.updated != nil {
			t.Fatalf("expected no repository update for processing task")
		}
		if billing.refundCalls != 0 {
			t.Fatalf("expected no refund for processing task, got %d", billing.refundCalls)
		}
		assertSQLExpectations(t, mock)
	})
}

func newTxDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db, mock
}

func assertSQLExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func testTask(status models.TaskStatus) *models.PromptTask {
	return &models.PromptTask{
		ID:        "task-1",
		UserID:    "user-1",
		ModelID:   "model-1",
		Payload:   "original payload",
		Status:    status,
		CreatedAt: time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC),
	}
}

func cloneTask(task *models.PromptTask) *models.PromptTask {
	if task == nil {
		return nil
	}
	clone := *task
	return &clone
}

func cloneModel(model *models.AIModel) *models.AIModel {
	if model == nil {
		return nil
	}
	clone := *model
	return &clone
}
