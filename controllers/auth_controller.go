package controllers

import (
	"os"
	"strings"

	"github.com/FaaizHaikal/spendiary-backend/services"
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

	if err := services.RegisterUser(req.Username, req.Password); err != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username taken"})
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func Login(ctx *fiber.Ctx) error {
	var req AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	user, err := services.FindUserByUsername(req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username not found"})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect password"})
	}

	accessToken, refreshToken, err := services.GenerateTokens(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func DeleteUser(ctx *fiber.Ctx) error {
	var req AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	user, err := services.FindUserByUsername(req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username not found"})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect password"})
	}

	if err := services.DeleteUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete user"})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func Refresh(ctx *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	userID, err := services.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate access token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": newAccessToken})
}

func VerifyAccessToken(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
