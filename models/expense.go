package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Date        time.Time `gorm:"not null"`
	Description string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`

	UserID uint `gorm:"not null"`
}
