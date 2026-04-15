package repositories

import (
	"database/sql"
	"errors"
	"fmt"

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
	model := &models.AIModel{}

	err := r.db.QueryRow(`
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
