package controllers

import (
	"github.com/FaaizHaikal/spendiary-backend/database"
	"github.com/FaaizHaikal/spendiary-backend/models"
	"github.com/FaaizHaikal/spendiary-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	hashedPassword := utils.HashPassword(req.Password)

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username taken"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(ctx *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username not found"})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect password"})
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return ctx.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
