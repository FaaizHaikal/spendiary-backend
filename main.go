package main

import (
	"github.com/FaaizHaikal/spendiary-backend/config"
	"github.com/FaaizHaikal/spendiary-backend/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	app := fiber.New(fiber.Config{})
	routes.Initialize(app)

	// database.Connect()

	port := config.GetEnv("port")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
