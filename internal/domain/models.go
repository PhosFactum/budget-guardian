package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ChatID    int64 `gorm:"uniqueIndex"`
	Income    float64
	FixedCost float64
	PayDay1   int // First day of moneyearning (5th, 20th and etc.)
	PayDay2   int // Second day of moneyearning
}

type Transaction struct {
	gorm.Model
	UserID  uint
	Amount  float64
	Comment string
}
