package repositories

import (
	"database/sql"
	"fmt"

	"ai-inference-gateway/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	// Поточна доменна модель ще не містить поля типу транзакції.
	// На цьому етапі репозиторій зберігає єдиний наявний сценарій: списання за задачу.
	_, err := r.db.Exec(`
		INSERT INTO transactions (id, user_id, task_id, amount, type)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5)
	`, tx.ID, tx.UserID, tx.TaskID, tx.Amount, "charge")
	if err != nil {
		return fmt.Errorf("create transaction %s: %w", tx.ID, err)
	}

	return nil
}
