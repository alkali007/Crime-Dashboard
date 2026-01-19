# 🚀 Automated Deployment & Architecture Guide

This guide explains how to deploy the **Waspada Bandung** system and how to avoid common "Mixed Content" security errors.

---

## 🧩 Part 1: How the System Works (Architecture)

1.  **Memory (Database - Supabase)**: PostgreSQL storage.
2.  **Brain (Backend API - Koyeb)**: Serves data via JSON.
3.  **Face (Frontend UI - Vercel)**: Next.js Dashboard.
4.  **Worker (Scraper - GitHub Actions)**: Updates data every 6 hours.

### 🐳 The Role of Docker (Optional)
The **Dockerfile** in the `backend` folder is an **optional** "Recipe".
- **With Docker**: You have total control over the environment (Go version, dependencies).
- **Without Docker**: Cloud providers (Koyeb/Render) can auto-detect Go code and build it themselves (Buildpacks).
*This project supports both methods.*

---

## 🛠️ Part 2: Deployment Steps

### Step 1: Push to GitHub
Ensure your code is in a public or private GitHub repository.

### Step 2: Backend API (Koyeb)
1.  **Join**: [koyeb.com](https://www.koyeb.com/).
2.  **Create Service**: Click **GitHub**, select `Crime-Dashboard`.
3.  **Settings**:
    -   **Work Directory**: `backend`.
    -   **Deployment Method**: You can choose **Docker** (uses our Dockerfile) or **Buildpack** (auto-detects Go). Both work.
4.  **Environment Variables**:
    -   `DB_DSN`: Your PostgreSQL connection string.
5.  **Result**: You get a URL (e.g., `https://xxxx.koyeb.app`).

### Step 3: Frontend (Vercel)
1.  **Import**: [vercel.com](https://vercel.com/) -> Import your repo.
2.  **Settings**:
    -   **Root Directory**: `frontend`.
3.  **Environment Variables**:
    -   `NEXT_PUBLIC_API_URL`: **IMPORTANT: MUST START WITH HTTPS**. 
        -   ✅ Correct: `https://xxxx.koyeb.app`
        -   ❌ Incorrect: `http://xxxx.koyeb.app` (This will cause a "Mixed Content" error in your browser).
4.  **Deploy**: Your website is live!

---

## ⚠️ Troubleshooting: The "Not Found" or Blank Map Error
If your website loads but the Map/Feed stays empty, check your browser's **Network tab**:
- **Mixed Content Error**: If your Vercel site is `https` but your API link is `http`, the browser will block the data for security. **Always use `https://`**.
- **Trailing Slash**: Ensure your API URL does **not** end with a `/`.
- **Redeploy**: If you update an Environment Variable in Vercel, you **must** "Redeploy" for the change to take effect.
