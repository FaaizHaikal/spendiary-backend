package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:50;not null;unique"`
	Password string `gorm:"size:255;not null"`
	Expenses []Expense
	Savings  []Saving
}
