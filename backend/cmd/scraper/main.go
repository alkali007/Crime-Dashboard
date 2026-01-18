package main

import (
	"crime-dashboard-backend/internal/database"
	"crime-dashboard-backend/internal/scraper"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	log.Println("Starting Scraper Job...")
	database.Connect()
	scraper.Run()
	log.Println("Scraper Job Completed.")
}
