package scraper

import (
	"crime-dashboard-backend/internal/database"
	"crime-dashboard-backend/internal/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// DistrictCoordinates maps district names to Lat/Long.
// Simplified list for demo purposes. Ideally complete this list.
var DistrictCoordinates = map[string][2]float64{
	"Andir":            {-6.9126, 107.5777},
	"Astana Anyar":     {-6.9389, 107.6009},
	"Antapani":         {-6.9122, 107.6596},
	"Arcamanik":        {-6.9250, 107.6750},
	"Babakan Ciparay":  {-6.9360, 107.5750},
	"Bandung Kidul":    {-6.9500, 107.6300},
	"Bandung Kulon":    {-6.9300, 107.5700},
	"Bandung Wetan":    {-6.9050, 107.6150},
	"Batununggal":      {-6.9150, 107.6350},
	"Bojongloa Kaler":  {-6.9400, 107.5900},
	"Bojongloa Kidul":  {-6.9500, 107.6000},
	"Buahbatu":         {-6.9600, 107.6400},
	"Cibeunying Kaler": {-6.8900, 107.6300},
	"Cibeunying Kidul": {-6.9000, 107.6400},
	"Cibiru":           {-6.9300, 107.7200},
	"Cicendo":          {-6.9050, 107.5900},
	"Cidadap":          {-6.8700, 107.6000},
	"Cinambo":          {-6.9400, 107.6900},
	"Coblong":          {-6.8900, 107.6150},
	"Gedebage":         {-6.9600, 107.6800},
	"Kiaracondong":     {-6.9200, 107.6450},
	"Lengkong":         {-6.9300, 107.6200},
	"Mandalajati":      {-6.9100, 107.6700},
	"Panyileukan":      {-6.9400, 107.7000},
	"Rancasari":        {-6.9500, 107.6700},
	"Regol":            {-6.9400, 107.6100},
	"Sukajadi":         {-6.8900, 107.5950},
	"Sukasari":         {-6.8700, 107.5850},
	"Sumur Bandung":    {-6.9150, 107.6100},
	"Ujung Berung":     {-6.9100, 107.7000},
}

func Run() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.detik.com", "detik.com"),
	)

	// Using a search URL for demonstration.
	// In production, might cycle through multiple news sites.
	targetURL := "https://www.detik.com/search/searchall?query=kriminal+bandung&siteid=2"

	c.OnHTML("article", func(e *colly.HTMLElement) {
		// This selector needs to be adjusted based on actual Detik HTML structure
		// Assuming typical article structure for demo.
		title := e.ChildText(".title")
		link := e.ChildAttr("a", "href")
		// dateStr := e.ChildText(".date") // Parsing this is tricky, depends on format

		if title == "" || link == "" {
			return
		}

		// Basic deduplication check
		var exists int64
		database.DB.Model(&models.CrimeReport{}).Where("source_url = ?", link).Count(&exists)
		if exists > 0 {
			return
		}

		fmt.Printf("Found article: %s\n", title)

		// District Parsing (Naive/Simple regex or substring search)
		var districtID uint
		var districtName string

		titleLower := strings.ToLower(title)

		for name, coords := range DistrictCoordinates {
			if strings.Contains(titleLower, strings.ToLower(name)) {
				districtName = name

				// Find or Create District
				var district models.District
				dbRes := database.DB.FirstOrCreate(&district, models.District{Name: name})
				if dbRes.RowsAffected > 0 {
					// New district, update coords
					district.Latitude = coords[0]
					district.Longitude = coords[1]
					database.DB.Save(&district)
				}
				districtID = district.ID
				break
			}
		}

		if districtID == 0 {
			// Default to 'Unknown' or skip if strictly district mapped
			// For now, let's skip extracting if no district found to keep data clean
			// Or log it.
			return
		}

		report := models.CrimeReport{
			Title:        title,
			SourceURL:    link,
			DistrictID:   districtID,
			Category:     "Uncategorized", // Could perform keyword analysis
			IncidentDate: time.Now(),      // Placeholder, needs parsing logic
		}

		// Try to parse date if possible, otherwise default to Now
		// Parsing dateStr...

		if err := database.DB.Create(&report).Error; err != nil {
			log.Printf("Error saving report: %v", err)
		} else {
			log.Printf("Saved report: %s (%s)", title, districtName)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with:", err)
	})

	c.Visit(targetURL)
}
