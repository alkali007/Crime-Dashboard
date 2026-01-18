package main

import (
	"crime-dashboard-backend/internal/database"
	"crime-dashboard-backend/internal/models"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load("../../.env"); err != nil {
		// Try looking in current directory too, handling different run contexts
		_ = godotenv.Load()
		// We don't fatal here because on production (Render), there is no .env file,
		// it uses real env vars. So we just ignore the error or log it as info.
		log.Println("Info: No .env file found or error loading it. Using system env vars.")
	}

	database.Connect()

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bandung Crime Dashboard API")
	})

	app.Get("/incidents", func(c *fiber.Ctx) error {
		var incidents []models.CrimeReport
		// Preload District to get details
		if err := database.DB.Preload("District").Order("incident_date desc").Find(&incidents).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(incidents)
	})

	app.Get("/stats/district", func(c *fiber.Ctx) error {
		type DistrictStat struct {
			DistrictName string  `json:"district_name"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Count        int64   `json:"count"`
		}

		var stats []DistrictStat
		// Join districts and count reports
		// Correct SQL: SELECT d.name, d.latitude, d.longitude, count(r.id) as count FROM districts d LEFT JOIN crime_reports r ON r.district_id = d.id GROUP BY d.id
		if err := database.DB.Table("districts").
			Select("districts.name as district_name, districts.latitude, districts.longitude, count(crime_reports.id) as count").
			Joins("left join crime_reports on crime_reports.district_id = districts.id").
			Group("districts.id").
			Scan(&stats).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(stats)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
