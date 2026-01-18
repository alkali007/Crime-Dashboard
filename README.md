# Waspada Bandung: Crime Dashboard

## Overview
A real-time dashboard visualization of crime incidents in Bandung, Indonesia.
Built with Golang (Fiber) backend and Next.js frontend.

## Prerequisites
- Go 1.20+
- Node.js 18+ & npm
- PostgreSQL (or just run backend which defaults to dry-run/local if configured)

## Setup

### Backend
1. Navigate to `backend/`.
2. Run `go mod tidy` to download dependencies.
3. Start the server (API):
   ```bash
   go run cmd/server/main.go
   ```
   Server runs on http://localhost:8080.

4. Run the scraper:
   ```bash
   go run cmd/scraper/main.go
   ```
   (Requires DB connection to work fully, configurable via `DB_DSN`).

### Frontend
1. Navigate to `frontend/`.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run development server:
   ```bash
   npm run dev
   ```
   Open http://localhost:3000.

## Configuration
Set `DB_DSN` environment variable for Postgres connection string.
Example: `host=localhost user=postgres password=secret dbname=crime_dashboard port=5432 sslmode=disable`
