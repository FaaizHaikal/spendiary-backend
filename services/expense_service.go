package services

import (
	"time"

	"github.com/FaaizHaikal/spendiary-backend/database"
	"github.com/FaaizHaikal/spendiary-backend/models"
)

func CreateExpense(expense *models.Expense) error {
	return database.DB.Create(expense).Error
}

func GetExpenseByID(id uint, userID uint) (*models.Expense, error) {
	var expense models.Expense
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&expense).Error
	return &expense, err
}

func UpdateExpense(expense *models.Expense) error {
	return database.DB.Save(expense).Error
}

func DeleteExpense(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Expense{}).Error
}

func GetAllExpenses(userID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	err := database.DB.Where("user_id = ?", userID).Order("date DESC").Find(&expenses).Error
	return expenses, err
}

func GetExpensesByMonth(userID uint, year int, month time.Month) ([]models.Expense, error) {
	var expenses []models.Expense

	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // 1st of next month

	err := database.DB.Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Order("date ASC").Find(&expenses).Error

	return expenses, err
}
