package service

import (
	"time"

	"github.com/PhosFactum/budget-guardian/internal/domain"
	"github.com/PhosFactum/budget-guardian/internal/interfaces"
)

type budgetService struct {
	userRepo interfaces.UserRepository
	txRepo   interfaces.TransactionRepository
}

func NewBudgetService(
	userRepo interfaces.UserRepository,
	txRepo interfaces.TransactionRepository,
) interfaces.BudgetService {
	return &budgetService{
		userRepo: userRepo,
		txRepo:   txRepo,
	}
}

func (s *budgetService) CalculateDailyLimit(userID uint) (float64, float64, error) {
	user, err := s.userRepo.FindByChatID(int64(userID))
	if err != nil {
		return 0, 0, err
	}

	days := calculateDaysUntilNextPayment(user.PayDay1, user.PayDay2)
	baseLimit := (user.Income - user.FixedCost) / float64(days+1)

	transactions, err := s.txRepo.GetByUserID(userID)
	if err != nil {
		return 0, 0, err
	}

	debt := calculateDebt(transactions, baseLimit)
	return baseLimit - (debt / float64(days)), debt, nil
}

func (s *budgetService) RecordTransaction(userID uint, amount float64) error {
	tx := &domain.Transaction{
		UserID: userID,
		Amount: amount,
	}
	return s.txRepo.Create(tx)
}

func calculateDaysUntilNextPayment(payDay1, payDay2 int) int {
	now := time.Now()
	currentDay := now.Day()

	if currentDay < payDay1 {
		return payDay1 - currentDay
	} else if currentDay < payDay2 {
		return payDay2 - currentDay
	}

	nextMonth := now.AddDate(0, 1, 0)
	daysInMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 0, 0, 0, 0, time.UTC).Day()
	return daysInMonth - currentDay + payDay1
}

func calculateDebt(transactions []domain.Transaction, dailyLimit float64) float64 {
	var debt float64
	for _, tx := range transactions {
		debt += tx.Amount - dailyLimit
	}
	return debt
}
