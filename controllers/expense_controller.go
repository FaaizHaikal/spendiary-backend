package controllers

import (
	"strconv"
	"time"

	"github.com/FaaizHaikal/spendiary-backend/models"
	"github.com/FaaizHaikal/spendiary-backend/services"
	"github.com/gofiber/fiber/v2"
)

func CreateExpense(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	var expense models.Expense
	if err := ctx.BodyParser(&expense); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	expense.UserID = userID
	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}

	if err := services.CreateExpense(&expense); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create expense"})
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func GetExpenses(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	expenses, err := services.GetAllExpenses(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}

	return ctx.Status(fiber.StatusOK).JSON(expenses)
}

func UpdateExpense(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	id, _ := strconv.Atoi(ctx.Params("id"))

	var updated models.Expense
	if err := ctx.BodyParser(&updated); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	existing, err := services.GetExpenseByID(uint(id), userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expense not found"})
	}

	existing.Description = updated.Description
	existing.Amount = updated.Amount
	existing.Date = updated.Date

	err = services.UpdateExpense(existing)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update expense"})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func DeleteExpense(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	id, _ := strconv.Atoi(ctx.Params("id"))

	err := services.DeleteExpense(uint(id), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete expense"})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func GetExpensesByMonth(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	monthParam := ctx.Query("month")
	yearParam := ctx.Query("year")

	month, err := strconv.Atoi(monthParam)
	if err != nil || month < 1 || month > 12 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid month"})
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil || year < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid year"})
	}

	expenses, err := services.GetExpensesByMonth(userID, year, time.Month(month))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}

	return ctx.Status(fiber.StatusOK).JSON(expenses)
}

func GetExpensesGroupByPeriod(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	period := ctx.Query("period", "month")

	validPeriods := map[string]bool{
		"day":   true,
		"week":  true,
		"month": true,
		"year":  true,
	}

	if !validPeriods[period] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid period"})
	}

	points, err := services.GetExpensesGroupByPeriod(userID, period)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}

	return ctx.Status(fiber.StatusOK).JSON(points)
}
