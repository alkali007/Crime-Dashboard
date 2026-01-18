# Deployment Guide: Waspada Bandung

This guide explains how to deploy the Crime Dashboard for free using **Supabase** (Database), **Render** (Backend), and **Vercel** (Frontend).

## 1. Database (Supabase)

1.  **Create Account/Login**: Go to [supabase.com](https://supabase.com/).
2.  **New Project**: Create a new project. Give it a name (e.g., `bandung-crime-db`) and a secure password.
3.  **Get Connection String**:
    -   Go to **Project Settings** -> **Database**.
    -   Under **Connection Parameters**, find the **URI** (Mode: Session).
    -   It should look like: `postgresql://postgres:[YOUR-PASSWORD]@db.xxxx.supabase.co:5432/postgres`
    -   **Save this string**. You will need it for the Backend.

## 2. Backend (Render)

1.  **Create Account/Login**: Go to [render.com](https://render.com/).
2.  **Connect GitHub**: detailed instructions are better if you push this code to your own GitHub repo first.
    *   *If you haven't pushed to GitHub yet, do checking "Create a new repository on the command line" on GitHub and push this code.*
3.  **New Web Service**:
    -   Click **New +** -> **Web Service**.
    -   Select your repository `Crime-Dashboard`.
4.  **Configure**:
    -   **Name**: `crime-dashboard-backend`
    -   **Runtime**: **Go**
    -   **Build Command**: `cd backend && go build -o server cmd/server/main.go`
    -   **Start Command**: `cd backend && ./server`
    -   **Region**: Singapore (likely closest to Bandung) or whatever is default free.
    -   **Instance Type**: **Free**.
5.  **Environment Variables**:
    -   Scroll down to **Environment Variables**.
    -   Add Key: `DB_DSN`
    -   Value: Paste your Supabase Connection String from Step 1.
6.  **Deploy**: Click **Create Web Service**.
    -   Wait for deployment to finish.
    -   **Copy the URL** (e.g., `https://crime-dashboard-backend.onrender.com`). You need this for the Frontend.

## 2. Alternative Backend (Koyeb) - If Render fails
**Koyeb** is excellent and often doesn't require a credit card for the "Free Forever" tier.

1.  **Login**: [koyeb.com](https://www.koyeb.com/).
2.  **Deploy**: Click **Create App** -> **GitHub**.
3.  **Select Repo**: Choose `Crime-Dashboard`.
4.  **Builder**: Choose **Go** (Buildpack).
5.  **Settings**:
    -   **Work Directory**: `backend`
    -   **Build Command**: `go build -o server cmd/server/main.go`
    -   **Run Command**: `./server`
    -   **Privileged**: Unchecked (Leave blank).
6.  **Environment Variables**:
    -   Add `DB_DSN` = (Your Supabase URL).
7.  **Deploy**. Copy the `xxxx.koyeb.app` URL.

### 2.1 Backend Scraper (GitHub Actions - Free)
We will use **GitHub Actions** to run the scraper automatically for free.

1.  **Go to your GitHub Repo**.
2.  Click **Settings** -> **Secrets and variables** -> **Actions**.
3.  Click **New repository secret**.
    -   Name: `DB_DSN`
    -   Value: (Your Supabase Connection String).
4.  **Done!**
    -   The scraper is already configured in `.github/workflows/scraper.yml`.
    -   It will run everyday every 6 hours.
    -   To test it immediately: Go to **Actions** tab -> **Run Scraper** -> **Run workflow**.

## 3. Frontend (Vercel)

1.  **Create Account/Login**: Go to [vercel.com](https://vercel.com/).
2.  **Import Project**:
    -   Click **Add New...** -> **Project**.
    -   Select your `Crime-Dashboard` repository.
3.  **Configure**:
    -   **Framework Preset**: Next.js.
    -   **Root Directory**: **IMPORTANT**. Click **Edit** and select `frontend`. Use **`frontend`** as the root, otherwise the build will fail.
4.  **Environment Variables**:
    -   Add Key: `NEXT_PUBLIC_API_URL`
    -   Value: The Render Backend URL from Step 2 (e.g., `https://crime-dashboard-backend.onrender.com`). **Note**: Ensure no trailing slash, or handle it carefully.
5.  **Deploy**: Click **Deploy**.
6.  **Done!** Your dashboard is live.

## Troubleshooting

-   **CORS Issues**: The backend allows all origins (`*`) by default, so it should work.
-   **Database**: If backend fails to start, check `DB_DSN` in Render logs.
-   **Scraper**: If no data appears, run the scraper locally once pointing to the remote DB to populate initial data, or wait for the Cron Job.
