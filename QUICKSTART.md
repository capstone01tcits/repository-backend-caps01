# Quick Start: Supabase Storage & Railway Deployment

## 📍 What You Need to Do RIGHT NOW

### 1️⃣ Create Supabase Buckets (5 minutes)

Go to **supabase.com** → Your Project → **Storage**

```
Create 3 buckets with these exact names (make them PUBLIC):
✅ videos
✅ thumbnails  
✅ assets
```

**That's it!** Your storage service code is already ready to use them.

---

### 2️⃣ Deploy to Railway (10 minutes)

1. **Sign up:** railway.app (use GitHub)
2. **Create project** → Deploy from GitHub
3. **Select repo:** Sevima-BackEnd Ai Video Gen
4. **Add env vars** (copy from your `.env` file)
5. **Deploy!** ✅

Railway will automatically:
- Detect your Dockerfile
- Build your Go app
- Deploy to production
- Give you a public URL

---

## 🔗 Quick Reference

Your **current setup**:
- ✅ DB: Connected to Supabase pool (`aws-1-ap-southeast-2.pooler.supabase.com`)
- ✅ Storage code: Ready (just waiting for buckets)
- ✅ Docker: Configured
- ⏳ Buckets: Need to create
- ⏳ Railway: Need to deploy

---

## 📝 Environment Variables Ready

**In your `.env` file, you already have:**
```
SUPABASE_URL=https://wkmvrwiesfpnfnaybpbx.supabase.co
SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
STORAGE_BUCKET_VIDEOS=videos
STORAGE_BUCKET_THUMBNAILS=thumbnails
STORAGE_BUCKET_ASSETS=assets
```

Just make sure the bucket names match! ✅

---

## 🧪 Test After Setup

### Supabase Bucket Test
```bash
curl -X POST https://wkmvrwiesfpnfnaybpbx.supabase.co/storage/v1/object/videos/test.txt \
  -H "Authorization: Bearer YOUR_SERVICE_ROLE_KEY" \
  -d "test content"
```

### Railway Deployment Test
```bash
curl https://your-railway-domain.railway.app/
# Should get your API response
```

---

## 🚀 Deployment Flow

```
GitHub Commit → Push
        ↓
    Railway Detects Change
        ↓
    Builds Docker Image
        ↓
    Uploads to Supabase ✓
        ↓
    Serves to Frontend ✓
```

---

## 📚 Full Documentation

See complete guides in:
- 📄 **SETUP_SUPABASE_RAILWAY.md** - Detailed steps
- ✅ **DEPLOYMENT_CHECKLIST.md** - Full checklist
- 🚂 **railway.json** - Auto-configuration

---

## ❓ Most Common Questions

**Q: Do I need to modify the code?**
A: No! Your storage service is already implemented and ready to use.

**Q: Where do I create the buckets?**
A: Supabase console → Storage → New Bucket

**Q: How do I get the Railway URL?**
A: Railway dashboard → Your project → Domains

**Q: Will my frontend work with the deployed backend?**
A: Yes, just update the API URL in `src/lib/axios.ts`

**Q: Is my database connected?**
A: Yes! You're already using Supabase pool connection.

---

## ✨ Next Steps

1. ✅ Create the 3 buckets in Supabase (2 min)
2. ✅ Push code to GitHub
3. ✅ Create Railway project (5 min)
4. ✅ Add env variables to Railway (2 min)
5. ✅ Wait for deployment (3-5 min)
6. ✅ Test with frontend
7. 🎉 Done!

**Total time: ~15-20 minutes**

---

## 🆘 If Something Goes Wrong

### Buckets not working?
1. Make sure bucket is **Public** (not Private)
2. Check CORS settings in Supabase
3. Verify service role key in code

### Railway deployment failed?
1. Check logs in Railway dashboard
2. Verify Docker builds locally: `docker build -t test .`
3. Ensure all env vars are set

### Can't connect DB?
1. Verify credentials in Railway env vars
2. Check if pool is active in Supabase
3. Test locally: `go run cmd/main.go`

---

**Questions? Check:**
- SETUP_SUPABASE_RAILWAY.md (detailed guide)
- DEPLOYMENT_CHECKLIST.md (step-by-step)
- storage_service.go (implementation details)

Good luck! 🚀
