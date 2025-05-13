package services

import (
	"fmt"
	"time"

	"github.com/FaaizHaikal/spendiary-backend/database"
	"github.com/FaaizHaikal/spendiary-backend/models"
	"gorm.io/gorm"
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

func GetRecentExpenses(userID uint, count int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := database.DB.Where("user_id = ?", userID).Order("date DESC").Limit(count).Find(&expenses).Error
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

func GetExpensesGroupByPeriod(userID uint, period string) ([]models.ChartPoint, error) {
	var points []models.ChartPoint
	var rows *gorm.DB
	db := database.DB

	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 1)
		rows = db.Raw(`
			SELECT to_char(date, 'HH24:MI') as label, SUM(amount) as total
			FROM expenses
			WHERE user_id = ? AND date BETWEEN ? AND ?
			GROUp BY label ORDER BY label`, userID, startDate, endDate)

	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}

		daysSinceMonday := weekday - 1 // Monday is 1 in time.Weekday

		startDate := now.AddDate(0, 0, -daysSinceMonday).Truncate(24 * time.Hour)
		endDate := startDate.AddDate(0, 0, 6)
		rows = db.Raw(`
			SELECT to_char(date, 'Dy') as label, SUM(amount) as total
			FROM expenses
			WHERE user_id = ? AND date BETWEEN ? AND ?
			GROUP BY label ORDER BY MIN(date)`, userID, startDate, endDate)

	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)

		rows = db.Raw(`
			SELECT to_char(date, 'FMDDth') as label, SUM(amount) as total
			FROM expenses
			WHERE user_id = ? AND date BETWEEN ? AND ?
			GROUP BY label ORDER BY MIN(date)`, userID, startDate, endDate)

	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)
		rows = db.Raw(`
			SELECT to_char(date, 'Mon') as label, SUM(amount) as total
			FROM expenses
			WHERE user_id = ? AND date BETWEEN ? AND ?
			GROUP BY label ORDER BY MIN(date)`, userID, startDate, endDate)

	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	if err := rows.Scan(&points).Error; err != nil {
		return nil, err
	}

	return points, nil
}
