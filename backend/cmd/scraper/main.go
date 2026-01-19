package main

import (
	"crime-dashboard-backend/internal/database"
	"crime-dashboard-backend/internal/scraper"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Use Overload instead of Load because the user has system env vars set to 'no'
	// which blocks Load() from working. Overload() will prioritize the .env file.
	_ = godotenv.Overload()
	_ = godotenv.Overload("../../.env")

	log.Println("Starting Scraper Job...")
	database.Connect()
	scraper.Run()
	log.Println("Scraper Job Completed.")
}
