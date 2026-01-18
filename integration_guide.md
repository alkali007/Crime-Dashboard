# Integration Guide: How It Works

This guide explains how the Frontend (Next.js) talks to the Backend (Go) and how you can verify it.

## The Flow

1.  **Frontend Request** (`frontend/src/app/page.js`)
    *   When the page loads, the `useEffect` function runs.
    *   It looks for the API URL: `const apiUrl = process.env.NEXT_PUBLIC_API_URL ...`
    *   It sends a request: `fetch(apiUrl + '/incidents')` -> e.g., `http://localhost:8080/incidents`

2.  **Backend Response** (`backend/cmd/server/main.go`)
    *   The Go server is listening on port `8080`.
    *   It receives the `GET /incidents` request.
    *   It asks the **Database** for the list of reports.
    *   It converts that list to **JSON** and sends it back.

3.  **Frontend Update**
    *   The frontend receives the JSON data.
    *   It updates the React state: `setIncidents(data)`.
    *   The Map component sees the new data and draws the pins.

## Integration Checklist

To ensure they are integrated correctly:

### 1. Check Backend
Open a terminal and run:
```bash
cd backend
go run cmd/server/main.go
```
*   **Success**: You see "Database connection established".
*   **Test**: Open your browser to `http://localhost:8080/incidents`. You should see `[]` (empty list) or some JSON data.

### 2. Check Frontend
Open a **new** terminal (keep backend running) and run:
```bash
cd frontend
npm run dev
```
*   **Success**: App starts at `http://localhost:3000`.
*   **Test**: Open `http://localhost:3000`. Open Developer Tools (F12) -> Network Tab. Refresh the page.
    *   Look for a request named `incidents`.
    *   Status should be `200 OK`.
    *   If it is `200`, **Integration is Working**.

## Common Issues
*   **CORS Error**: Check if the backend has `app.Use(cors.New())` (I have already enabled this).
*   **Connection Refused**: Backend is not running.
*   **Empty Map**: Backend is running but Database is empty. You need to run the **Scraper** (`go run cmd/scraper/main.go`) to get data.
