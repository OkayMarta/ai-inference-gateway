package services

import (
	"context"
	"database/sql"
	"fmt"

	"billing-service/internal/models"
)

type BillingService struct {
	db       *sql.DB
	userRepo UserRepository
	txRepo   TransactionRepository
}

func NewBillingService(db *sql.DB, userRepo UserRepository, txRepo TransactionRepository) *BillingService {
	return &BillingService{db: db, userRepo: userRepo, txRepo: txRepo}
}

func (s *BillingService) Charge(userID, taskID string, amount float64) (*models.Transaction, error) {
	return s.applyBillingEvent("charge", userID, taskID, amount)
}

func (s *BillingService) Refund(userID, taskID string, amount float64) (*models.Transaction, error) {
	return s.applyBillingEvent("refund", userID, taskID, amount)
}

func (s *BillingService) applyBillingEvent(txType, userID, taskID string, amount float64) (*models.Transaction, error) {
	if userID == "" || taskID == "" || amount <= 0 {
		return nil, ErrInvalidBillingInput
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("begin billing transaction: %w", err)
	}

	if _, err := s.userRepo.GetByIDTx(tx, userID); err != nil {
		_ = tx.Rollback()
		if isRepoNotFoundError(err, "user not found:") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	switch txType {
	case "charge":
		if err := s.userRepo.DeductBalanceTx(tx, userID, amount); err != nil {
			_ = tx.Rollback()
			if err.Error() == ErrInsufficientBalance.Error() {
				return nil, ErrInsufficientBalance
			}
			return nil, fmt.Errorf("deduct balance: %w", err)
		}
	case "refund":
		if err := s.userRepo.AddBalanceTx(tx, userID, amount); err != nil {
			_ = tx.Rollback()
			if isRepoNotFoundError(err, "user not found:") {
				return nil, ErrUserNotFound
			}
			return nil, fmt.Errorf("add balance: %w", err)
		}
	default:
		_ = tx.Rollback()
		return nil, ErrInvalidBillingInput
	}

	transaction := &models.Transaction{
		ID:     generateID(),
		UserID: userID,
		TaskID: taskID,
		Amount: amount,
		Type:   txType,
	}
	if err := s.txRepo.CreateTx(tx, transaction); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit billing transaction: %w", err)
	}

	return transaction, nil
}
