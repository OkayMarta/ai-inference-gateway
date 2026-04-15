package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"ai-inference-gateway/internal/models"
)

type WorkerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) GetAll() ([]*models.WorkerNode, error) {
	rows, err := r.db.Query(`
		SELECT id, status
		FROM worker_nodes
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list workers: %w", err)
	}
	defer rows.Close()

	var workers []*models.WorkerNode
	for rows.Next() {
		worker := &models.WorkerNode{}
		if err := rows.Scan(&worker.ID, &worker.Status); err != nil {
			return nil, fmt.Errorf("scan worker: %w", err)
		}

		supportedModels, err := r.getSupportedModels(worker.ID)
		if err != nil {
			return nil, err
		}
		worker.SupportedModels = supportedModels
		workers = append(workers, worker)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workers: %w", err)
	}

	return workers, nil
}

func (r *WorkerRepository) GetIdle() ([]*models.WorkerNode, error) {
	rows, err := r.db.Query(`
		SELECT id, status
		FROM worker_nodes
		WHERE status = $1
		ORDER BY id
	`, models.WorkerIdle)
	if err != nil {
		return nil, fmt.Errorf("list idle workers: %w", err)
	}
	defer rows.Close()

	var workers []*models.WorkerNode
	for rows.Next() {
		worker := &models.WorkerNode{}
		if err := rows.Scan(&worker.ID, &worker.Status); err != nil {
			return nil, fmt.Errorf("scan idle worker: %w", err)
		}

		supportedModels, err := r.getSupportedModels(worker.ID)
		if err != nil {
			return nil, err
		}
		worker.SupportedModels = supportedModels
		workers = append(workers, worker)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate idle workers: %w", err)
	}

	return workers, nil
}

func (r *WorkerRepository) GetByID(id string) (*models.WorkerNode, error) {
	worker := &models.WorkerNode{}

	err := r.db.QueryRow(`
		SELECT id, status
		FROM worker_nodes
		WHERE id = $1
	`, id).Scan(&worker.ID, &worker.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("worker not found: %s", id)
		}
		return nil, fmt.Errorf("get worker by id %s: %w", id, err)
	}

	supportedModels, err := r.getSupportedModels(id)
	if err != nil {
		return nil, err
	}
	worker.SupportedModels = supportedModels

	return worker, nil
}

func (r *WorkerRepository) Create(worker *models.WorkerNode) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin worker create transaction: %w", err)
	}

	if _, err := tx.Exec(`
		INSERT INTO worker_nodes (id, status)
		VALUES ($1, $2)
	`, worker.ID, worker.Status); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("create worker %s: %w", worker.ID, err)
	}

	if err := replaceWorkerSupportedModels(tx, worker.ID, worker.SupportedModels); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit worker create transaction: %w", err)
	}

	return nil
}

func (r *WorkerRepository) Update(worker *models.WorkerNode) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin worker update transaction: %w", err)
	}

	result, err := tx.Exec(`
		UPDATE worker_nodes
		SET status = $2
		WHERE id = $1
	`, worker.ID, worker.Status)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update worker %s: %w", worker.ID, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("worker not found: %s", worker.ID)); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := replaceWorkerSupportedModels(tx, worker.ID, worker.SupportedModels); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit worker update transaction: %w", err)
	}

	return nil
}

func (r *WorkerRepository) UpdateStatus(id string, status models.WorkerStatus) error {
	result, err := r.db.Exec(`
		UPDATE worker_nodes
		SET status = $2
		WHERE id = $1
	`, id, status)
	if err != nil {
		return fmt.Errorf("update worker status %s: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("worker not found: %s", id)); err != nil {
		return err
	}

	return nil
}

func (r *WorkerRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM worker_nodes
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete worker %s: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("worker not found: %s", id)); err != nil {
		return err
	}

	return nil
}

func (r *WorkerRepository) getSupportedModels(workerID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT model_id
		FROM worker_supported_models
		WHERE worker_id = $1
		ORDER BY model_id
	`, workerID)
	if err != nil {
		return nil, fmt.Errorf("list supported models for worker %s: %w", workerID, err)
	}
	defer rows.Close()

	var modelIDs []string
	for rows.Next() {
		var modelID string
		if err := rows.Scan(&modelID); err != nil {
			return nil, fmt.Errorf("scan supported model for worker %s: %w", workerID, err)
		}
		modelIDs = append(modelIDs, modelID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate supported models for worker %s: %w", workerID, err)
	}

	return modelIDs, nil
}

func replaceWorkerSupportedModels(tx *sql.Tx, workerID string, modelIDs []string) error {
	if _, err := tx.Exec(`
		DELETE FROM worker_supported_models
		WHERE worker_id = $1
	`, workerID); err != nil {
		return fmt.Errorf("clear supported models for worker %s: %w", workerID, err)
	}

	for _, modelID := range modelIDs {
		if _, err := tx.Exec(`
			INSERT INTO worker_supported_models (worker_id, model_id)
			VALUES ($1, $2)
		`, workerID, modelID); err != nil {
			return fmt.Errorf("add supported model %s for worker %s: %w", modelID, workerID, err)
		}
	}

	return nil
}
