package repository

import (
	"gorm.io/gorm"

	"github.com/PhosFactum/budget-guardian/internal/domain"
	"github.com/PhosFactum/budget-guardian/internal/interfaces"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

// Realizing user interface for Create method
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// Realizing user interface for FindByChatID method
func (r *userRepository) FindByChatID(chatID int64) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("chat_id = ?", chatID).First(&user).Error
	return &user, err
}
