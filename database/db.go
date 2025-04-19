package database

import (
	"log"
	"os"

	"github.com/FaaizHaikal/spendiary-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	err = db.AutoMigrate(&models.User{}, &models.Expense{}, &models.Saving{})

	if err != nil {
		log.Fatal("Failed to migrate database!")
	}

	DB = db
}
