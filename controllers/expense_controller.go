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

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Expense created successfully"})
}

func GetExpenses(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	expenses, err := services.GetAllExpenses(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}

	return ctx.JSON(expenses)
}

func UpdateExpense(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
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

	return ctx.JSON(existing)
}

func DeleteExpense(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	id, _ := strconv.Atoi(ctx.Params("id"))

	err := services.DeleteExpense(uint(id), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete expense"})
	}

	return ctx.JSON(fiber.Map{"message": "Expense deleted"})
}

func GetMonthlyExpenses(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	now := time.Now()
	expenses, err := services.GetExpensesByMonth(userID, now.Year(), now.Month())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch expenses")
	}

	return ctx.JSON(expenses)
}
