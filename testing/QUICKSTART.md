# Quick Setup Guide - Testing Suite

**Updated:** April 8, 2026

---

## Quick Start

### Step 1: Ensure Backend & AI Service Are Running

```powershell
# Terminal 1 - Backend
cd "d:\ATHA ITS\Capstone\Backend\Sevima-BackEnd Ai Video Gen"
go run cmd/main.go

# Terminal 2 - AI Service
cd "d:\ATHA ITS\Capstone\Backend\Sevima-BackEnd Ai Video Gen\ai-service"
python main.py
```

Verify both are running:
- Backend: http://localhost:3000/health (should return {"status":"ok"})
- AI: http://localhost:8000/health

---

### Step 2: Run Test Suite

```powershell
cd "d:\ATHA ITS\Capstone\Backend\Sevima-BackEnd Ai Video Gen"
powershell -ExecutionPolicy Bypass -File testing/complete_test.ps1
```

Expected output will show test results per endpoint.

---

### Step 3: Add Credits (If Needed)

If video generation test fails with "insufficient credits":

**Method 1: Database Direct Update (Fastest)**

```sql
-- Connect to PostgreSQL
psql -U postgres -d sevima_ai_video_gen

-- Add credits
UPDATE user_credits 
SET credits = 1000 
WHERE user_id = '2c4a7410-ef75-4f02-86ff-d54445528c98';
```

**Method 2: Run Credit Script**

```powershell
powershell -ExecutionPolicy Bypass -File testing/add_credits.ps1
```

---

### Step 4: Re-run Tests

```powershell
powershell -ExecutionPolicy Bypass -File testing/complete_test.ps1
```

Should see 18/18 PASS (100% success) ✓

---

## Pre-Flight Checklist

Before running tests:
- ✓ Backend running on port 3000
- ✓ AI Service running on port 8000
- ✓ PostgreSQL running on port 5432
- ✓ 1000+ credits (for full video generation testing)

---

## Expected Results

### Test Execution Output

```
= Test Starting =

[1] HEALTH & AUTHENTICATION
  [OK] GET /health
  [OK] GET /api/auth/me

[2] PROJECTS
  [OK] POST /api/projects
  [OK] GET /api/projects
  [OK] GET /api/projects/:id

...

= Test Summary =
Total Tests: 18
Passed: 18 ✓
Failed: 0 ✗
Pass Rate: 100%
```

### Report File Location

Auto-generated report: `reports/ENDPOINT_TEST_REPORT_YYYYMMDD_HHMMSS.txt`

---

## Troubleshooting

### Problem: "Backend not responding"
- Check backend is running: `go run cmd/main.go`
- Check firewall allows port 3000
- Check no other service using port 3000

### Problem: "AI service not responding"
- Check Python installed: `python --version`
- Check requirements installed: `pip install -r ai-service/requirements.txt`
- Check Python process running

### Problem: "Database connection error"
- Check PostgreSQL running: `docker-compose up -d`
- Check connection credentials in `config/config.go`

### Problem: "Insufficient credits error"
- Run: `testing/add_credits.ps1`
- Or use SQL UPDATE command above
- Re-run tests

---

Status: ✓ Ready for Testing
Last Updated: April 8, 2026
