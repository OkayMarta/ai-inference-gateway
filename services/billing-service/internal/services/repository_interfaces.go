package services

import (
	appdb "billing-service/internal/db"
	"billing-service/internal/models"
)

type UserRepository interface {
	GetAll() ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByIDTx(tx appdb.DBTX, id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	DeductBalanceTx(tx appdb.DBTX, id string, amount float64) error
	AddBalanceTx(tx appdb.DBTX, id string, amount float64) error
}

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	CreateTx(exec appdb.DBTX, tx *models.Transaction) error
}
