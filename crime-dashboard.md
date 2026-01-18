To save this as a Markdown file, you can copy the code block below and save it as `PRD_Bandung_Crime_Dashboard.md`.

```markdown
# Product Requirements Document: Bandung Crime Dashboard (Waspada Bandung)

## 1. Project Objective
The goal is to build an automated crime reporting dashboard that visualizes the latest crime incidents across various districts in Bandung City. The system will act as a transparency tool, providing citizens with data-driven insights into local safety by scraping public news and official reports.

---

## 2. Technical Stack
| Layer | Technology |
| :--- | :--- |
| **Backend** | Golang (Gin/Fiber framework) |
| **Frontend** | JavaScript (React.js or Next.js) |
| **Database** | PostgreSQL (via Supabase Free Tier) |
| **Scraper/ETL** | Golang (Colly/Go-query) |
| **CI/CD** | GitHub Actions |
| **Hosting (BE)** | Render or Fly.io (Free Tier) |
| **Hosting (FE)** | Vercel or Netlify (Free Tier) |

---

## 3. System Architecture & Workflow
1.  **Scraper (ETL):** A specialized Golang script triggered by GitHub Actions on a Cron schedule. It fetches data from news portals and social feeds.
2.  **Database:** Structured storage for incidents, locations, and timestamps.
3.  **API (Backend):** Serves the processed data via RESTful endpoints.
4.  **Dashboard (Frontend):** Consumes the API to display an interactive map and statistical charts.

---

## 4. Functional Requirements

### 4.1 Data Scraper & ETL Pipeline
*   **Source Targeting:** Scrape news sites (e.g., Detik Jabar, PRFM, Tribun Jabar) using keywords like "Kriminal," "Begal," "Curanmor," and "Bandung."
*   **Parsing Logic:** Extract Title, Date, District Name (Kecamatan), and Source URL.
*   **Geocoding:** Convert Bandung district names into Latitude/Longitude coordinates using a static JSON lookup table to save API costs.
*   **Scheduler:** Must run automatically every 6 to 12 hours via GitHub Actions.

### 4.2 Backend API
*   **`GET /incidents`**: Retrieve a list of all crime reports with filtering by district and date range.
*   **`GET /stats/district`**: Provide aggregate counts of crimes per district for the heatmap.
*   **Health Checks**: Ensure the service stays active within free-tier sleep constraints.

### 4.3 Frontend Dashboard
*   **Interactive Map:** Use Leaflet.js (OpenStreetMap) to show incident pins or a district-based heatmap.
*   **District Sidebar:** List of the 30 districts in Bandung; clicking one filters the map and news feed.
*   **Crime Feed:** A chronological list of cards showing recent report summaries and links to original news sources.
*   **Responsive Design:** Fully functional on mobile and desktop browsers.

---

## 5. Data Schema (PostgreSQL)

```sql
CREATE TABLE districts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE,
    latitude DECIMAL,
    longitude DECIMAL
);

CREATE TABLE crime_reports (
    id SERIAL PRIMARY KEY,
    district_id INT REFERENCES districts(id),
    title TEXT,
    description TEXT,
    incident_date TIMESTAMP,
    source_url TEXT UNIQUE,
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 6. CI/CD & Deployment Strategy

### 6.1 GitHub Actions Workflows
1.  **Continuous Integration:** Run `go test` and `npm run lint` on every Pull Request.
2.  **ETL Runner (The Scraper):**
    *   **Trigger:** `schedule: - cron: '0 */6 * * *'`
    *   **Action:** Runs the Golang scraper, connects to the Supabase DB via secrets, and updates records.
3.  **Continuous Deployment:** 
    *   Push to `main` triggers an automatic build on Render (Backend) and Vercel (Frontend).

### 6.2 Free Resource Optimization
*   **Database:** Utilize Supabase's free tier (500MB storage).
*   **Compute:** Use Render "Blueprints" or Vercel "Deploy Hooks" to manage zero-cost scaling.
*   **Cold Starts:** Implement a lightweight ping in the scraper to wake up the Render web service if it has gone to sleep.

---

## 7. Non-Functional Requirements
*   **Performance:** Dashboard initial load time under 3 seconds.
*   **Reliability:** Scraper must include error handling for site structure changes (fail-safes).
*   **Security:** Database credentials and API keys must be managed strictly via GitHub Secrets and Environment Variables.
*   **Data Integrity:** Prevent duplicate entries by using `source_url` as a unique constraint in the database.
```