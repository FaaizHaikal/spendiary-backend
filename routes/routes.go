package routes

import (
	"github.com/FaaizHaikal/spendiary-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/refresh", controllers.Refresh)
}
