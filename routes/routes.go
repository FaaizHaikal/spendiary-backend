package routes

import (
	"github.com/FaaizHaikal/spendiary-backend/controllers"
	"github.com/FaaizHaikal/spendiary-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/delete", controllers.DeleteUser)
	api.Post("/refresh", controllers.Refresh)

	protected := api.Group("/user", middleware.RequireAuth)
	protected.Get("/verify", controllers.VerifyAccessToken)

	expense := protected.Group("/expenses")

	expense.Get("/monthly", controllers.GetExpensesByMonth)
	expense.Get("/group", controllers.GetExpensesGroupByPeriod)
	expense.Get("/all", controllers.GetExpenses)
	expense.Post("/create", controllers.CreateExpense)
	expense.Post("/delete", controllers.DeleteExpense)
	expense.Post("/update", controllers.UpdateExpense)
}
