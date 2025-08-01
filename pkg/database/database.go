package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PhosFactum/budget-guardian/pkg/config"
)

// InitDB: initializes PostgreSQL connection and returns *gorm.DB and error
func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Get("DB_HOST", "localhost"),
		config.Get("DB_USER", "postgres"),
		config.Get("DB_PASSWORD", "postgres"),
		config.Get("DB_NAME", "budget"),
		config.Get("DB_PORT", "5432"),
	)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return db, nil
}