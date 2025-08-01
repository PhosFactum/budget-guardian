package interfaces

type BudgetService interface {
	CalculateDailyLimit(userID uint) (float64, float64, error)
	RecordTransaction(userID uint, amount float64) error
}
