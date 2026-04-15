package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	appdb "ai-inference-gateway/internal/db"
	"ai-inference-gateway/internal/models"
)

type ModelRepository struct {
	db *sql.DB
}

func NewModelRepository(db *sql.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) GetAll() ([]*models.AIModel, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, token_cost
		FROM ai_models
		ORDER BY name, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list models: %w", err)
	}
	defer rows.Close()

	var items []*models.AIModel
	for rows.Next() {
		model := &models.AIModel{}
		if err := rows.Scan(&model.ID, &model.Name, &model.Description, &model.TokenCost); err != nil {
			return nil, fmt.Errorf("scan model: %w", err)
		}
		items = append(items, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate models: %w", err)
	}

	return items, nil
}

func (r *ModelRepository) GetByID(id string) (*models.AIModel, error) {
	return r.getByID(r.db, id)
}

func (r *ModelRepository) GetByIDTx(tx appdb.DBTX, id string) (*models.AIModel, error) {
	return r.getByID(tx, id)
}

func (r *ModelRepository) getByID(exec appdb.DBTX, id string) (*models.AIModel, error) {
	model := &models.AIModel{}

	err := exec.QueryRow(`
		SELECT id, name, description, token_cost
		FROM ai_models
		WHERE id = $1
	`, id).Scan(&model.ID, &model.Name, &model.Description, &model.TokenCost)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("model not found: %s", id)
		}
		return nil, fmt.Errorf("get model by id %s: %w", id, err)
	}

	return model, nil
}

func (r *ModelRepository) Create(model *models.AIModel) error {
	_, err := r.db.Exec(`
		INSERT INTO ai_models (id, name, description, token_cost)
		VALUES ($1, $2, $3, $4)
	`, model.ID, model.Name, model.Description, model.TokenCost)
	if err != nil {
		return fmt.Errorf("create model %s: %w", model.ID, err)
	}

	return nil
}

func (r *ModelRepository) ReplaceAll(models []*models.AIModel) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin replace models transaction: %w", err)
	}

	for _, model := range models {
		if _, err := tx.Exec(`
			INSERT INTO ai_models (id, name, description, token_cost)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE
			SET name = EXCLUDED.name,
			    description = EXCLUDED.description,
			    token_cost = EXCLUDED.token_cost
		`, model.ID, model.Name, model.Description, model.TokenCost); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("upsert synced model %s: %w", model.ID, err)
		}
	}

	if err := pruneObsoleteModels(tx, models); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit replace models transaction: %w", err)
	}

	return nil
}

func pruneObsoleteModels(tx *sql.Tx, models []*models.AIModel) error {
	// З БД видаляємо лише моделі, яких більше немає в актуальному Ollama sync
	// і які вже не використовуються історичними prompt_tasks. Це не ламає FK
	// і зберігає цілісність уже створених задач.
	if len(models) == 0 {
		if _, err := tx.Exec(`
			DELETE FROM ai_models AS m
			WHERE NOT EXISTS (
				SELECT 1
				FROM prompt_tasks AS t
				WHERE t.model_id = m.id
			)
		`); err != nil {
			return fmt.Errorf("prune obsolete models after empty sync: %w", err)
		}

		return nil
	}

	placeholders := make([]string, 0, len(models))
	args := make([]any, 0, len(models))
	for idx, model := range models {
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx+1))
		args = append(args, model.ID)
	}

	query := fmt.Sprintf(`
		DELETE FROM ai_models AS m
		WHERE m.id NOT IN (%s)
		  AND NOT EXISTS (
			SELECT 1
			FROM prompt_tasks AS t
			WHERE t.model_id = m.id
		  )
	`, strings.Join(placeholders, ", "))

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("prune obsolete models: %w", err)
	}

	return nil
}

func (r *ModelRepository) Update(model *models.AIModel) error {
	result, err := r.db.Exec(`
		UPDATE ai_models
		SET name = $2,
		    description = $3,
		    token_cost = $4
		WHERE id = $1
	`, model.ID, model.Name, model.Description, model.TokenCost)
	if err != nil {
		return fmt.Errorf("update model %s: %w", model.ID, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("model not found: %s", model.ID)); err != nil {
		return err
	}

	return nil
}

func (r *ModelRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM ai_models
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete model %s: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("model not found: %s", id)); err != nil {
		return err
	}

	return nil
}
