package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	appdb "ai-inference-gateway/internal/db"
	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/services"

	"github.com/lib/pq"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetByID(id string) (*models.PromptTask, error) {
	return r.getByID(r.db, id)
}

func (r *TaskRepository) GetByIDTx(tx appdb.DBTX, id string) (*models.PromptTask, error) {
	return r.getByID(tx, id)
}

func (r *TaskRepository) getByID(exec appdb.DBTX, id string) (*models.PromptTask, error) {
	task := &models.PromptTask{}
	var result sql.NullString

	err := exec.QueryRow(`
		SELECT id, user_id, model_id, payload, status, result, created_at
		FROM prompt_tasks
		WHERE id = $1
	`, id).Scan(
		&task.ID,
		&task.UserID,
		&task.ModelID,
		&task.Payload,
		&task.Status,
		&result,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found: %s", id)
		}
		return nil, fmt.Errorf("get task by id %s: %w", id, err)
	}

	if result.Valid {
		task.Result = result.String
	}

	return task, nil
}

func (r *TaskRepository) List(filter services.TaskListFilter) ([]*models.PromptTask, error) {
	orderBy, err := taskSortClause(filter.Sort)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, model_id, payload, status, result, created_at
		FROM prompt_tasks
	`
	var conditions []string
	var args []any
	argPos := 1

	if filter.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, filter.UserID)
		argPos++
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, filter.Status)
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY " + orderBy

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.PromptTask
	for rows.Next() {
		task, err := scanPromptTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Create(task *models.PromptTask) error {
	return r.create(r.db, task)
}

func (r *TaskRepository) CreateTx(tx appdb.DBTX, task *models.PromptTask) error {
	return r.create(tx, task)
}

func (r *TaskRepository) create(exec appdb.DBTX, task *models.PromptTask) error {
	if task.CreatedAt.IsZero() {
		if err := exec.QueryRow(`
			INSERT INTO prompt_tasks (id, user_id, model_id, payload, status, result)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''))
			RETURNING created_at
		`, task.ID, task.UserID, task.ModelID, task.Payload, task.Status, task.Result).Scan(&task.CreatedAt); err != nil {
			return fmt.Errorf("create task %s: %w", task.ID, err)
		}
		return nil
	}

	_, err := exec.Exec(`
		INSERT INTO prompt_tasks (id, user_id, model_id, payload, status, result, created_at)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7)
	`, task.ID, task.UserID, task.ModelID, task.Payload, task.Status, task.Result, task.CreatedAt)
	if err != nil {
		return fmt.Errorf("create task %s: %w", task.ID, err)
	}

	return nil
}

func (r *TaskRepository) Update(task *models.PromptTask) error {
	return r.update(r.db, task)
}

func (r *TaskRepository) UpdateTx(tx appdb.DBTX, task *models.PromptTask) error {
	return r.update(tx, task)
}

func (r *TaskRepository) update(exec appdb.DBTX, task *models.PromptTask) error {
	result, err := exec.Exec(`
		UPDATE prompt_tasks
		SET user_id = $2,
		    model_id = $3,
		    payload = $4,
		    status = $5,
		    result = NULLIF($6, '')
		WHERE id = $1
	`, task.ID, task.UserID, task.ModelID, task.Payload, task.Status, task.Result)
	if err != nil {
		return fmt.Errorf("update task %s: %w", task.ID, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("task not found: %s", task.ID)); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM prompt_tasks
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete task %s: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("task not found: %s", id)); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Complete(id, resultText string) error {
	return r.updateTaskResult(id, models.StatusCompleted, resultText)
}

func (r *TaskRepository) Fail(id, resultText string) error {
	return r.updateTaskResult(id, models.StatusFailed, resultText)
}

func (r *TaskRepository) GetNextQueued(supportedModels []string) (*models.PromptTask, error) {
	if len(supportedModels) == 0 {
		return nil, nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin task selection transaction: %w", err)
	}

	task := &models.PromptTask{}
	var result sql.NullString

	err = tx.QueryRow(`
		SELECT id, user_id, model_id, payload, status, result, created_at
		FROM prompt_tasks
		WHERE status = $1
		  AND model_id = ANY($2)
		ORDER BY created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`, models.StatusQueued, pq.Array(supportedModels)).Scan(
		&task.ID,
		&task.UserID,
		&task.ModelID,
		&task.Payload,
		&task.Status,
		&result,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = tx.Rollback()
			return nil, nil
		}
		_ = tx.Rollback()
		return nil, fmt.Errorf("select next queued task: %w", err)
	}

	if result.Valid {
		task.Result = result.String
	}

	if _, err := tx.Exec(`
		UPDATE prompt_tasks
		SET status = $2
		WHERE id = $1
	`, task.ID, models.StatusProcessing); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("mark task %s as processing: %w", task.ID, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit task selection transaction: %w", err)
	}

	task.Status = models.StatusProcessing
	return task, nil
}

func (r *TaskRepository) updateTaskResult(id string, status models.TaskStatus, resultText string) error {
	result, err := r.db.Exec(`
		UPDATE prompt_tasks
		SET status = $2,
		    result = NULLIF($3, '')
		WHERE id = $1
	`, id, status, resultText)
	if err != nil {
		return fmt.Errorf("update task %s state: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("task not found: %s", id)); err != nil {
		return err
	}

	return nil
}

func scanPromptTask(scanner interface {
	Scan(dest ...any) error
}) (*models.PromptTask, error) {
	task := &models.PromptTask{}
	var result sql.NullString
	var createdAt time.Time

	if err := scanner.Scan(
		&task.ID,
		&task.UserID,
		&task.ModelID,
		&task.Payload,
		&task.Status,
		&result,
		&createdAt,
	); err != nil {
		return nil, fmt.Errorf("scan task: %w", err)
	}

	task.CreatedAt = createdAt
	if result.Valid {
		task.Result = result.String
	}

	return task, nil
}

func taskSortClause(sortValue string) (string, error) {
	switch strings.ToLower(sortValue) {
	case "", "created_at_desc":
		return "created_at DESC", nil
	case "created_at_asc":
		return "created_at ASC", nil
	default:
		return "", services.ErrInvalidPagination
	}
}
