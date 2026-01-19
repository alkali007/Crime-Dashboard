# 🚀 Automated Deployment & Architecture Guide

This guide explains how to deploy the **Waspada Bandung** system and, more importantly, **how it all works together**.

---

## 🧩 Part 1: How the System Works (Architecture)

To make this app global, we split it into four specialized parts:

1.  **Memory (Database - Supabase)**: Where we store districts and crime reports.
2.  **Brain (Backend API - Koyeb)**: A Go program that talks to the database and gives data to the website.
3.  **Face (Frontend UI - Vercel)**: The beautiful Map and Feed that users see.
4.  **Worker (Scraper - GitHub Actions)**: A background script that finds new news every 6 hours.

### 🐳 What is a Dockerfile? (The "Recipe")
Think of the **Dockerfile** in the `backend` folder as a **Recipe**. 
Instead of you manually installing Go, setting paths, and clicking "Run" on a remote server, the Dockerfile tells the cloud (Koyeb):
- *"Start with a clean Linux computer."*
- *"Install Go 1.24."*
- *"Copy my code into this folder."*
- *"Compile the code and start the server on port 8080."*

This ensures that if the app works on your computer, it **will** work on the server exactly the same way.

---

## 🛠️ Part 2: Deployment Steps

### Step 1: Push to GitHub
Before starting, ensure your local code is uploaded to a private or public **GitHub repository**. Both Vercel and Koyeb will "watch" this repo for updates.

### Step 2: Backend API (Koyeb)
Koyeb is your "Brain" hosting.
1.  **Join**: [koyeb.com](https://www.koyeb.com/).
2.  **Create Service**: Click **GitHub**, select `Crime-Dashboard`.
3.  **Settings**:
    -   **Work Directory**: `backend` (This tells Koyeb where the `Dockerfile` is).
    -   **Instance Type**: `Nano` (Free).
4.  **Environment Variables**:
    -   `DB_DSN`: Your PostgreSQL connection string.
5.  **Result**: You get a URL (e.g., `https://xxxx.koyeb.app`). **Keep this URL.**

### Step 3: Frontend (Vercel)
Vercel is your "Face" hosting.
1.  **Import**: [vercel.com](https://vercel.com/) -> Import your repo.
2.  **Settings**:
    -   **Root Directory**: `frontend`.
3.  **Environment Variables**:
    -   `NEXT_PUBLIC_API_URL`: Paste the Koyeb URL from Step 2.
4.  **Deploy**: Your website is now live!

### Step 4: Automated Scraper (GitHub Actions)
The "Worker" that runs for free.
1.  **Repo Settings**: Go to GitHub -> `Settings` -> `Secrets` -> `Actions`.
2.  **Add Secrets**:
    -   `DB_DSN`: (Your Database URL).
    -   `GOOGLE_API_KEY`: (Your Gemini API Key).
3.  **Check**: Go to the `Actions` tab in GitHub to see the scraper running!

---

## 🔄 Summary: The "Push-to-Live" Flow
Every time you change code locally and run `git push`:
1.  **GitHub** tells Vercel & Koyeb: *"Hey, there is new code!"*
2.  **Vercel** rebuilds the Map/UI.
3.  **Koyeb** uses the `Dockerfile` to rebuild the API.
4.  **Within 2 minutes**, your live website is updated automatically.
