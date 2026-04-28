# Deployment Setup: Supabase + Railway

Panduan lengkap untuk deploy backend ke Railway dengan Supabase sebagai database.

---

## Part 1: Supabase Setup (Database)

### Step 1: Create Supabase Project

1. Buka https://supabase.com
2. Sign up / Login
3. Create New Project:
   - Project Name: `sevima-video-gen`
   - Database Password: Catat password ini! (Generate strong password)
   - Region: `Singapore` (Asia Tenggara terdekat)
   - Pricing Plan: Free tier OK untuk development
4. Tunggu project selesai dibuat (~1-2 menit)

### Step 2: Get Database Connection URL

1. Di Supabase Dashboard, buka **Settings > Database**
2. Cari section **Connection String** → **URI**
3. Pilih format **Postgres** atau **Psql**
4. Copy connection string:
   ```
   postgresql://postgres:[PASSWORD]@[HOST]:[PORT]/postgres
   ```

**Format untuk .env:**
```
DB_HOST=db.[PROJECT_ID].supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=[PASSWORD_YANG_DISIMPAN]
DB_NAME=postgres
```

### Step 3: Enable SSL Connection (Production)

Supabase memerlukan SSL untuk production. Update connection string:
```
postgresql://postgres:[PASSWORD]@db.[PROJECT_ID].supabase.co:5432/postgres?sslmode=require
```

Untuk GORM, tambahkan di config.go:
```go
dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
    Cfg.DBHost, Cfg.DBPort, Cfg.DBUser, Cfg.DBPassword, Cfg.DBName)
```

---

## Part 2: Railway Setup (Deployment)

### Step 1: Create Railway Account

1. Buka https://railway.app
2. Sign up dengan GitHub (recommended)
3. Create New Project

### Step 2: Connect GitHub Repository

1. Di Railway Dashboard → Create New Project
2. Pilih **"Deploy from GitHub repo"**
3. Authorize Railway dengan GitHub account
4. Select repository: `Sevima-BackEnd Ai Video Gen`
5. Railway akan auto-detect Dockerfile

### Step 3: Configure Environment Variables

1. Di Railway Project → **Variables**
2. Tambahkan environment variables:

```env
# App Config
APP_PORT=8080
APP_ENV=production

# Database (dari Supabase)
DB_HOST=db.[PROJECT_ID].supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=[PASSWORD_DARI_SUPABASE]
DB_NAME=postgres

# JWT
JWT_SECRET=[GENERATE_RANDOM_STRONG_SECRET]
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=[GENERATE_RANDOM_STRONG_SECRET]
JWT_REFRESH_EXPIRE_HOURS=168

# AI Service
AI_SERVICE_URL=https://[AI_SERVICE_URL] # sesuaikan dengan deployment AI service

# API Keys
LTX_API_KEY=[YOUR_API_KEY]
RUNWAY_API_KEY=[YOUR_API_KEY]
```

**Tips Generate Strong Secret:**
```bash
# Di PowerShell
[Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes([guid]::NewGuid().ToString())) 
```

### Step 4: Configure Domain (Optional)

1. Railway → **Settings > Custom Domain**
2. Pilih domain atau gunakan Railway domain default
3. DNS akan di-configure automatically

### Step 5: Deploy

Railway akan auto-deploy ketika ada push ke GitHub:

```bash
# Push ke main branch
git push origin main

# Railway akan:
# 1. Detect perubahan
# 2. Build Docker image
# 3. Run database migrations (AutoMigrate)
# 4. Start server
```

Monitoring:
- Railway Dashboard → **Logs** untuk melihat real-time logs
- Cek status deployment di **Deployments** tab

---

## Part 3: Update .env untuk Production

Copy `.env.example` ke `.env` dan update:

```env
APP_PORT=8080
APP_ENV=production

# Supabase Database
DB_HOST=db.[YOUR_PROJECT_ID].supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=[PASSWORD]
DB_NAME=postgres

# JWT Secrets (Strong random strings)
JWT_SECRET=put-very-long-random-secret-here-min-32-chars
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=put-another-long-random-secret-here-min-32-chars
JWT_REFRESH_EXPIRE_HOURS=168

# AI Service
AI_SERVICE_URL=https://[your-ai-service-domain]

# API Keys
LTX_API_KEY=xxx
RUNWAY_API_KEY=xxx
```

---

## Part 4: Local Testing dengan Supabase

### Option A: Local Supabase (Docker)

```bash
# Install Supabase CLI
npm install -g supabase

# Start local Supabase
supabase start

# Cek connection string di terminal output
```

### Option B: Use Remote Supabase (Recommended for Testing)

Update `.env.local` dengan Supabase credentials:

```bash
# Copy .env.example ke .env.production
cp .env.example .env.production

# Edit dengan Supabase details
# Lalu update config.go untuk baca dari file ini
```

---

## Part 5: GitHub Actions untuk Auto-Deploy (Optional)

Tambahkan di `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Railway

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Deploy to Railway
        run: |
          npm install -g @railway/cli
          railway up
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

Get `RAILWAY_TOKEN`:
1. Railway Dashboard → **Account > API Tokens**
2. Create token → copy ke GitHub Secrets
3. Add ke repo: **Settings > Secrets > New secret**

---

## Part 6: Monitoring & Debugging

### Railway Logs
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# View logs
railway logs

# View status
railway status
```

### Database Connection Issues

**Error: "i/o timeout"**
- Pastikan SSL enabled: `sslmode=require`
- Check Supabase firewall settings

**Error: "connection refused"**
- Verifikasi DB_HOST, DB_PORT, DB_PASSWORD di Railway environment
- Pastikan database variables ter-update

### Check Database Migration

```bash
# Connect ke Supabase dari local:
psql "postgresql://postgres:[PASSWORD]@db.[ID].supabase.co:5432/postgres?sslmode=require"

# List tables
\dt

# Check schema
\d users
```

---

## Part 7: Common Issues & Solutions

### 1. SSL Certificate Error
**Solution:** Tambahkan `?sslmode=require` ke connection string

### 2. Port Conflict
Railway auto-assign PORT. Jangan hardcode port 5000:
```go
port := os.Getenv("PORT")  // Railway set PORT environment
if port == "" {
    port = "5000"  // fallback
}
```

### 3. Database Timeout
- Increase connection timeout di config.go
- Check Supabase max connections (Free tier: 100)

### 4. Cold Start Issue
- Railway hibernates inactive apps
- Upgrade untuk production apps

---

## Part 8: Frontend Configuration

Update frontend API base URL ke Railway deployment:

```typescript
// src/lib/axios.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 
  'https://[your-railway-app].up.railway.app';
```

Add environment variable di frontend `.env.local`:
```
NEXT_PUBLIC_API_URL=https://[your-railway-app].up.railway.app
```

---

## Checklist

- [ ] Supabase project created
- [ ] Database connection string obtained
- [ ] Railway account created
- [ ] GitHub repository connected
- [ ] Environment variables configured di Railway
- [ ] First deploy successful
- [ ] Database tables created (auto migration)
- [ ] API endpoints responding
- [ ] JWT tokens working
- [ ] Frontend API base URL updated

---

## Quick Deploy Command

```bash
# 1. Ensure all changes committed
git add .
git commit -m "Production configuration"

# 2. Push to main (Railway auto-deploys)
git push origin main

# 3. Monitor deployment
npm install -g @railway/cli
railway login
railway logs
```

---

## Useful Links

- Supabase Docs: https://supabase.com/docs
- Railway Docs: https://docs.railway.app
- Railway Logs: https://railway.app/dashboard
- Supabase Connection: https://supabase.com/docs/guides/connecting-to-postgres

---

**Questions?** Check Railway docs atau Supabase support dashboard.
