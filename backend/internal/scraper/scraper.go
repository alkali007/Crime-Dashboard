package scraper

import (
	"context"
	"crime-dashboard-backend/internal/database"
	"crime-dashboard-backend/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

var model = "gemini-2.5-pro"

// JSON Schema for response parsing
type ScrapedIncident struct {
	ID           int8   `json:"id"`
	DistrictName string `json:"district_name"`
	ArticleTitle string `json:"article_title"`
	Description  string `json:"description"`
	IncidentDate string `json:"incident_date"` // GenAI returns ISO string usually
	SourceURL    string `json:"source_url"`
	Category     string `json:"category"`
}

// InitialDistrictData from user input
var InitialDistrictData = []models.District{
	{Name: "Andir", Latitude: -6.9152, Longitude: 107.5857},
	{Name: "Astana Anyar", Latitude: -6.9366, Longitude: 107.6014},
	{Name: "Antapani", Latitude: -6.9142, Longitude: 107.6617},
	{Name: "Arcamanik", Latitude: -6.9148, Longitude: 107.6835},
	{Name: "Babakan Ciparay", Latitude: -6.9453, Longitude: 107.5794},
	{Name: "Bandung Kidul", Latitude: -6.9535, Longitude: 107.6322},
	{Name: "Bandung Kulon", Latitude: -6.9246, Longitude: 107.5649},
	{Name: "Batununggal", Latitude: -6.9272, Longitude: 107.6366},
	{Name: "Bojongloa Kaler", Latitude: -6.9248, Longitude: 107.5925},
	{Name: "Bojongloa Kidul", Latitude: -6.9535, Longitude: 107.5982},
	{Name: "Buah Batu", Latitude: -6.9479, Longitude: 107.6534},
	{Name: "Cibeunying Kidul", Latitude: -6.9038, Longitude: 107.6416},
	{Name: "Cibeunying Kaler", Latitude: -6.8927, Longitude: 107.6253},
	{Name: "Cibiru", Latitude: -6.9221, Longitude: 107.7126},
	{Name: "Cicendo", Latitude: -6.9103, Longitude: 107.5941},
	{Name: "Cidadap", Latitude: -6.8687, Longitude: 107.6056},
	{Name: "Cinambo", Latitude: -6.9254, Longitude: 107.6917},
	{Name: "Coblong", Latitude: -6.8837, Longitude: 107.6146},
	{Name: "Gedebage", Latitude: -6.9587, Longitude: 107.6946},
	{Name: "Kiaracondong", Latitude: -6.9246, Longitude: 107.6491},
	{Name: "Lengkong", Latitude: -6.9304, Longitude: 107.6183},
	{Name: "Mandalajati", Latitude: -6.9056, Longitude: 107.6749},
	{Name: "Panyileukan", Latitude: -6.9363, Longitude: 107.7058},
	{Name: "Rancasari", Latitude: -6.9532, Longitude: 107.6698},
	{Name: "Regol", Latitude: -6.9358, Longitude: 107.6105},
	{Name: "Sumur Bandung", Latitude: -6.9184, Longitude: 107.6110},
	{Name: "Ujung Berung", Latitude: -6.9137, Longitude: 107.7018},
}

func Run() {
	// Seed Districts
	seedDistricts()

	ctx := context.Background()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Println("Error: GOOGLE_API_KEY is not set in environment.")
	} else {
		// Log length and prefix for debugging
		log.Printf("DEBUG: Found GOOGLE_API_KEY (Length: %d, Prefix: %s...)", len(apiKey), apiKey[:2])
	}

	// Create client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Error creating GenAI client: %v", err)
		return
	}

	// Enable Google Search tool for source grounding
	var tools []*genai.Tool
	tools = append(tools, &genai.Tool{
		GoogleSearch: &genai.GoogleSearch{},
	})

	userPrompt := `
You are acting as a Data Scraping Specialist with experience in news data extraction and validation.

Objective:
Extract criminal case information occurring within districts of Bandung City from reputable news websites.

Context:
- Each district in Bandung City must have a minimum of 2 distinct criminal cases.
- Only criminal cases published within the last 12 months from the current date are allowed.
- To prevent model overload, randomly select exactly 15 districts per execution.
- Only the selected districts must be processed and returned.
- Cases must be sourced exclusively from the provided news websites.
- Only information explicitly stated in the article may be used.

Input:

Districts (Bandung City):
Andir, Astana Anyar, Antapani, Arcamanik, Babakan Ciparay, Bandung Kidul, Bandung Kulon,
Batununggal, Bojongloa Kaler, Bojongloa Kidul, Buah Batu, Cibeunying Kidul,
Cibeunying Kaler, Cibiru, Cicendo, Cidadap, Cinambo, Coblong, Gedebage,
Kiaracondong, Lengkong, Mandalajati, Panyileukan, Rancasari, Regol,
Sumur Bandung, Ujung Berung

News Websites:
- news.detik.com
- cnnindonesia.com
- en.antaranews.com
- bbc.com
- kompas.com
- sindonews.com
- tempo.com
- liputan6.com
- tribunnews.com

Crime Categories:
Assault, Battery, Rape, Murder, Burglary, Arson, Theft, Vandalism

Constraints:
- Randomly select exactly 15 districts from the list above per execution.
- Minimum of 2 cases per district within the last 12 months (hard requirement).
- If a selected district does not meet the minimum requirement,
  return records for that district with all fields set to NaN.
- Districts not selected must not appear in the output.
- Use English language only.
- Do not infer, assume, or fabricate any information.
- Summaries must be strictly derived from article content.
- Do not include cases older than 12 months.
- Do not include duplicate cases across districts.
- Output must strictly follow the defined schema and structure.

Output Format:
Return a JSON array of objects with the following schema. IMPORTANT: Output ONLY the raw JSON array, do not include any other text or explanation.

{
  "id": int8,
  "district_name": text,
  "article_title": text,
  "description": text,
  "incident_date": timestamptz,
  "source_url": text,
  "category": text
}

Evaluation Criteria:
- Output must be machine-readable and executable with minimal modification.
- Data accuracy must be verifiable against the source URL.
- Temporal filtering must be strictly enforced.
- Schema consistency is mandatory.
`

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				genai.NewPartFromText(userPrompt),
			},
		},
	}

	config := &genai.GenerateContentConfig{
		Tools: tools,
		// ResponseMIMEType: "application/json", // CRITICAL: Disabled because it conflicts with Tools (Google Search)
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				genai.NewPartFromText(
					"You are an expert data scraping and extraction specialist for criminal news from Indonesian media. " +
						"You must extract only verifiable information explicitly stated in articles, strictly limited to the last three months. " +
						"If coverage is insufficient, you must return NaN values rather than fabricating data. " +
						"All outputs must be source-traceable and strictly schema-compliant. " +
						"Your response MUST be a valid JSON array only.",
				),
			},
		},
	}

	log.Println("Calling Gemini API...")

	// We use GenerateContent (Unstreamed) to define easy JSON parsing
	resp, err := client.Models.GenerateContent(ctx, model, contents, config)
	if err != nil {
		log.Printf("Error generating content: %v", err)
		return
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("No content generated.")
		return
	}

	// Extract JSON Text
	jsonText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		jsonText += part.Text
	}

	// Clean Markdown code blocks if present
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")

	// Parse JSON
	var scrapedData []ScrapedIncident
	if err := json.Unmarshal([]byte(jsonText), &scrapedData); err != nil {
		log.Printf("Error unmarshalling JSON: %v. Raw text: %s", err, jsonText)
		return
	}

	log.Printf("Parsed %d items from Gemini. Saving to Database...", len(scrapedData))

	// Debug: Show first few items raw
	debugCount := 3
	if len(scrapedData) < debugCount {
		debugCount = len(scrapedData)
	}
	for i := 0; i < debugCount; i++ {
		log.Printf("DEBUG: Sample Item %d: Title='%s', District='%s', Source='%s'", i+1, scrapedData[i].ArticleTitle, scrapedData[i].DistrictName, scrapedData[i].SourceURL)
	}

	// Save to Database
	for _, item := range scrapedData {
		// Strict Validation: Skip NaN or Empty values
		if item.DistrictName == "" || item.DistrictName == "NaN" ||
			item.ArticleTitle == "" || item.ArticleTitle == "NaN" ||
			item.SourceURL == "" || item.SourceURL == "NaN" ||
			item.Description == "NaN" {
			continue
		}

		// Synchronize with existing Districts from Seeding
		var district models.District
		if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(item.DistrictName)).First(&district).Error; err != nil {
			log.Printf("District not found in DB: %s (Skipping report)", item.DistrictName)
			continue
		}

		// Deduplication: Retrieve all source_urls for this district to double check (though global check is usually enough)
		// For robustness, we do a global check on SourceURL since it should be unique.
		var exists int64
		database.DB.Model(&models.CrimeReport{}).Where("source_url = ?", item.SourceURL).Count(&exists)
		if exists > 0 {
			log.Printf("Duplicate found (skipping): %s", item.SourceURL)
			continue
		}

		// Parse Incident Date
		t, err := time.Parse(time.RFC3339, item.IncidentDate)
		if err != nil {
			// Try a few other formats if Gemini is inconsistent, or default to now
			// But user requested "The column IncidentDate should be using from Gemini's result"
			// warning: Gemini might return various string formats if not forced strictly by schema
			log.Printf("Warning: Failed to parse date %s for %s. Using Now.", item.IncidentDate, item.ArticleTitle)
			t = time.Now()
		}

		report := models.CrimeReport{
			Title:        item.ArticleTitle,
			Description:  item.Description,
			SourceURL:    item.SourceURL,
			DistrictID:   district.ID,
			Category:     item.Category,
			IncidentDate: t,
			CreatedAt:    time.Now(),
		}

		if err := database.DB.Create(&report).Error; err != nil {
			log.Printf("Error creating report: %v", err)
		} else {
			fmt.Printf("Saved: %s (%s)\n", item.ArticleTitle, item.DistrictName)
		}
	}
}

func seedDistricts() {
	log.Println("Seeding Districts...")
	for _, d := range InitialDistrictData {
		var district models.District
		// Check if exists by Name
		if err := database.DB.Where("name = ?", d.Name).First(&district).Error; err != nil {
			// Not found, create it
			if err := database.DB.Create(&d).Error; err != nil {
				log.Printf("Failed to seed district %s: %v", d.Name, err)
			}
		} else {
			// Found, do nothing as per user request to avoid unnecessary updates
			log.Printf("District %s already exists, skipping update.", d.Name)
		}
	}
	log.Println("District Seeding Completed.")
}
