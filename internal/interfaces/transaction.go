package interfaces

import "github.com/PhosFactum/budget-guardian/internal/domain"

type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	GetByUserID(userID uint) ([]domain.Transaction, error)
}
