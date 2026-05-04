package repositories

import (
	"database/sql"
	"fmt"

	appdb "billing-service/internal/db"
	"billing-service/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.create(r.db, tx)
}

func (r *TransactionRepository) CreateTx(exec appdb.DBTX, tx *models.Transaction) error {
	return r.create(exec, tx)
}

func (r *TransactionRepository) create(exec appdb.DBTX, tx *models.Transaction) error {
	err := exec.QueryRow(`
		INSERT INTO transactions (id, user_id, task_id, amount, type)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5)
		RETURNING created_at
	`, tx.ID, tx.UserID, tx.TaskID, tx.Amount, tx.Type).Scan(&tx.CreatedAt)
	if err != nil {
		return fmt.Errorf("create transaction %s: %w", tx.ID, err)
	}

	tx.CreatedAt = tx.CreatedAt.UTC()
	return nil
}
