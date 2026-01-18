package database

import (
	"fmt"
	"log"
	"os"

	"crime-dashboard-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Fallback for local testing if env not set, though ideally it should be.
		// Asking user to provide it is better, but dev defaults help.
		// Assuming standard Postgres port.
		dsn = "host=localhost user=postgres password=postgres dbname=crime_dashboard port=5432 sslmode=disable"
		fmt.Println("No DB_DSN environment variable found. Using default local DSN.")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		// We don't panic here so the scraper or server can maybe retry or start up in a limited mode if needed,
		// but for this app panic might be appropriate. Let's return.
		return
	}

	log.Println("Database connection established")

	log.Println("Running migrations...")
	err = DB.AutoMigrate(&models.District{}, &models.CrimeReport{})
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated")
}
