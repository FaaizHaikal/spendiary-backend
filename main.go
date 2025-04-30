package main

import (
	"os"
	"strings"

	"github.com/FaaizHaikal/spendiary-backend/config"
	"github.com/FaaizHaikal/spendiary-backend/database"
	"github.com/FaaizHaikal/spendiary-backend/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	app := fiber.New(fiber.Config{})
	routes.Initialize(app)

	args := os.Args[1:] // Skip the first argument (program name)

	seed := false
	for _, arg := range args {
		if strings.ToLower(arg) == "seed=true" {
			seed = true
			break
		}
	}

	database.Connect()
	if seed {
		database.SeedExpense(1) // userID 1 is testuser
	}

	port := os.Getenv("port")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
