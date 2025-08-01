package interfaces

import "github.com/PhosFactum/budget-guardian/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByChatID(chatid int64) (*domain.User, error)
}
