package services

import (
	"os"

	"github.com/FaaizHaikal/spendiary-backend/database"
	"github.com/FaaizHaikal/spendiary-backend/models"
	"github.com/FaaizHaikal/spendiary-backend/utils"
	"github.com/golang-jwt/jwt/v4"
)

func RegisterUser(username, password string) error {
	hashedPassword := utils.HashPassword(password)

	user := models.User{
		Username: username,
		Password: hashedPassword,
	}

	return database.DB.Create(&user).Error
}

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error

	return &user, err
}

func GenerateTokens(userID uint) (string, string, error) {
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ParseRefreshToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	return userID, nil
}
