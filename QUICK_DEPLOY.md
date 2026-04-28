# Quick Setup Guide: Supabase + Railway

Panduan cepat dalam 5 langkah untuk deploy ke production.

---

## Quick Steps (30 menit)

### Step 1: Create Supabase Database (5 min)

```
1. Buka https://supabase.com → Sign Up
2. Create Project: "sevima-video-gen"
   - Region: Singapore
   - Password: Catat! (gunakan password generator)
3. Settings > Database > Copy Connection String (URI format)
```

**Simpan format ini:**
```
DB_HOST = db.XXXXXX.supabase.co
DB_PORT = 5432
DB_USER = postgres
DB_PASSWORD = [password_kamu]
DB_NAME = postgres
```

---

### Step 2: GitHub Push (2 min)

```bash
cd d:\ATHA ITS\Capstone\Backend\Sevima-BackEnd Ai Video Gen

# Verify changes
git status

# Commit
git add .
git commit -m "Configure for production deployment"

# Push
git push origin main
```

---

### Step 3: Create Railway Project (3 min)

```
1. Buka https://railway.app → Sign Up (use GitHub)
2. Create New Project → Deploy from GitHub repo
3. Select: Sevima-BackEnd Ai Video Gen
4. Railway auto-detects Dockerfile ✓
```

---

### Step 4: Configure Environment Variables (10 min)

Di Railway Dashboard → Your Project → **Variables**

Copy & Paste ini, update dengan nilai Supabase + secrets:

```
APP_PORT=8080
APP_ENV=production

# Dari Supabase
DB_HOST=db.XXXXXXXXXX.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=XXXXXXXXXXXXX
DB_NAME=postgres

# Generate strong secrets (copy random strings)
JWT_SECRET=9f7a8c3e1d2b5f4c6a9e2d1c8f4a3b5c7d9e1f2a3b4c5d6e7f8a9b0c1d2e3
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=3b7c9e1a5f2d8c4b0a6e9f3d1c7a4e2b5f8a0c3e6b9d2f4a7c0e3b6d9
JWT_REFRESH_EXPIRE_HOURS=168

# AI Service URL (sesuaikan deployment)
AI_SERVICE_URL=http://localhost:8000

# API Keys (optional)
LTX_API_KEY=your-key
RUNWAY_API_KEY=your-key
```

**Generate Strong Secret di PowerShell:**
```powershell
[Convert]::ToBase64String([guid]::NewGuid().ToString() + [guid]::NewGuid().ToString())
```

---

### Step 5: Deploy! (5 min)

```bash
# Railway auto-deploys ketika ada push ke main
# Tunggu di Railway Dashboard → Deployments

# Atau manual trigger:
npm install -g @railway/cli
railway login
railway trigger
```

**Monitor:**
- Railway Dashboard → **Logs** (real-time logs)
- Tunggu status: "✓ Running" (hijau)

---

## ✓ Verification

Setelah deployment selesai, test API:

```bash
# Get Railway URL dari Dashboard
# Format: https://[project-name].up.railway.app

# Test health check
curl https://[project-name].up.railway.app/health

# Should return:
# {"status":"healthy"}
```

---

## Update Frontend

Setelah backend live, update frontend `.env.local`:

```
NEXT_PUBLIC_API_URL=https://[your-railway-project].up.railway.app
```

Maka frontend akan connect ke backend production.

---

## Environment Variables Reference

| Variable | Nilai | Sumber |
|----------|-------|--------|
| `APP_PORT` | 8080 | Fixed untuk Railway |
| `APP_ENV` | production | Fixed |
| `DB_HOST` | db.XXXX.supabase.co | Supabase Settings |
| `DB_PASSWORD` | [catat saat create] | Supabase |
| `JWT_SECRET` | [random string] | Generate random |
| `JWT_REFRESH_SECRET` | [random string] | Generate random |
| `AI_SERVICE_URL` | Sesuaikan deployment | Your AI service |

---

## Common Mistakes

**DON'T:**
- Copy `.env` ke repo (secret exposed!)
- Use default JWT secrets di production
- Forget catat password Supabase

**DO:**
- Add environment variables di Railway Dashboard
- Use strong random secrets (min 32 chars)
- Test health check setelah deploy

---

## Troubleshooting

###  "Connection refused"
```
→ CheckDB credentials di Railway Variables
→ Verify Supabase project active
```

### "SSL certificate error"
```
→ Supabase requires SSL
→ Credentials in Railway should auto-handle
```
### "Port 5000: address already in use"
```
→ Code updated untuk Railway PORT env
→ Rebuild Docker: railway trigger
```

### Deployment stuck/failed
```
→ Check Logs di Railway Dashboard
→ Recent logs akan show error
→ Fix & git push main again
```

---

## 📱 Final URLs

Setelah deploy:

```
Backend API:  https://[project].up.railway.app
Database:     PostgreSQL di Supabase
Frontend:     Akan connect ke backend API
```

---

## Next Steps

1. Database live (Supabase)
2. API live (Railway)
3. Deploy Frontend ke Vercel/Railway
4. Setup monitoring & alerts
5. Configure custom domain

---

Butuh bantuan? Buka `DEPLOYMENT_SETUP.md` untuk detail lengkap!
