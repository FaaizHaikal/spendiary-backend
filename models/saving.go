package models

import (
	"time"

	"gorm.io/gorm"
)

type Saving struct {
	gorm.Model
	Date   time.Time `gorm:"not null"`
	Amount float64   `gorm:"not null"`

	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`
}
