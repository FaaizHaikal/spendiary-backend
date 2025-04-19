package controllers

import (
	"os"

	"github.com/FaaizHaikal/spendiary-backend/database"
	"github.com/FaaizHaikal/spendiary-backend/models"
	"github.com/FaaizHaikal/spendiary-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(ctx *fiber.Ctx) error {
	var req AuthRequest
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
	var req AuthRequest
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Refresh(ctx *fiber.Ctx) error {
	type TokenInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	var input TokenInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	token, err := jwt.Parse(input.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate access token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": newAccessToken})
}

func Profile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")

	return ctx.JSON(fiber.Map{"message": "Hello from profile!", "user_id": userID})
}
