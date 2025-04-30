package database

import (
	"log"
	"math/rand"
	"time"

	"github.com/FaaizHaikal/spendiary-backend/models"
)

func SeedExpense(userID uint) {
	for i := range 200 {
		expense := models.Expense{
			UserID:      userID,
			Amount:      float64(rand.Intn(20000)) / 1.0,
			Description: "Seeded",
			Date:        randomDateThisYear(),
		}

		if err := DB.Create(&expense).Error; err != nil {
			log.Printf("Failed to seed expense %d: %v", i, err)
		}
	}
}

func randomDateThisYear() time.Time {
	start := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	delta := end.Sub(start)
	randomSeconds := rand.Int63n(int64(delta.Seconds()))

	return start.Add(time.Duration(randomSeconds) * time.Second)
}
