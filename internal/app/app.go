package app

import (
	"fmt"
	"log"

	"github.com/PhosFactum/budget-guardian/internal/domain"
	"github.com/PhosFactum/budget-guardian/internal/handler"
	"github.com/PhosFactum/budget-guardian/internal/rabbitmq"
	"github.com/PhosFactum/budget-guardian/internal/repository"
	"github.com/PhosFactum/budget-guardian/internal/service"
	"github.com/PhosFactum/budget-guardian/internal/telegram"
	"github.com/PhosFactum/budget-guardian/pkg/config"
	"github.com/PhosFactum/budget-guardian/pkg/database"
	"gorm.io/gorm"
)

func Run() error {
	// Loading configuration
	config.LoadConfig()

	// DB initializing
	db, err := database.InitDB()
	if err != nil {
		return err
	}

	// Automigrations
	if err := runMigrations(db); err != nil {
		return err
	}

	// Repositories init
	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// Service init
	budgetService := service.NewBudgetService(userRepo, txRepo)

	// После инициализации сервиса
	telegramBot, err := telegram.NewBot(
		config.Get("TELEGRAM_BOT_TOKEN", ""),
		budgetService,
	)
	if err != nil {
		return fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Println("Starting Telegram bot...")
	go telegramBot.Start()

	// RabbitMQ init
	rmq, err := rabbitmq.New(config.Get("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"))
	if err != nil {
		return err
	}
	defer rmq.Close()

	// Handlers init
	handler := handler.NewHandler(budgetService, rmq)
	router := handler.InitRoutes()

	// Launching server
	log.Println("Starting 'Budget Guardian' server on :8080")
	return router.Run(":8080")
}

// runMigrations: function to run automigration
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Transaction{},
	)
}
