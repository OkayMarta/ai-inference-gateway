package repositories

import (
	"sync"

	"ai-inference-gateway/internal/models"
)

// TransactionRepository зберігає історію транзакцій (списання токенів за задачі). Це основа для "білінгу"
type TransactionRepository struct {
	mu           sync.RWMutex
	transactions map[string]*models.Transaction
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{transactions: make(map[string]*models.Transaction)}
}

// Create просто зберігає запис про нову транзакцію
func (r *TransactionRepository) Create(tx *models.Transaction) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions[tx.ID] = tx
}

// GetByUserID дозволяє подивитись історію витрат конкретного користувача
func (r *TransactionRepository) GetByUserID(userID string) []*models.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var out []*models.Transaction
	for _, tx := range r.transactions {
		if tx.UserID == userID {
			cp := *tx
			out = append(out, &cp)
		}
	}
	return out
}