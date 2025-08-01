package repository

import (
	"github.com/PhosFactum/budget-guardian/internal/domain"
	"github.com/PhosFactum/budget-guardian/internal/interfaces"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) interfaces.TransactionRepository {
	return &transactionRepository{db: db}
}

// Realizing tx interface for Create method
func (r *transactionRepository) Create(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

// Realizing tx interface for GetByUserID method
func (r *transactionRepository) GetByUserID(userID uint) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&txs).Error
	return txs, err
}
