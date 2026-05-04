# Supabase & Railway Setup Checklist

## ✅ Pre-Deployment Checklist

### Supabase Storage Buckets
- [ ] Log in to Supabase console at supabase.com
- [ ] Select your project: Sevima-AI-Content-Creator
- [ ] Create bucket: **videos** (Public)
- [ ] Create bucket: **thumbnails** (Public)
- [ ] Create bucket: **assets** (Public)
- [ ] Configure CORS in Supabase Settings → Storage → CORS
- [ ] Test upload to one bucket manually
- [ ] Verify public URL works: `https://wkmvrwiesfpnfnaybpbx.supabase.co/storage/v1/object/public/videos/test.mp4`

### Environment Variables Review
- [ ] `SUPABASE_URL` = `https://wkmvrwiesfpnfnaybpbx.supabase.co` ✓
- [ ] `SUPABASE_ANON_KEY` exists ✓
- [ ] `SUPABASE_SERVICE_ROLE_KEY` exists ✓
- [ ] `STORAGE_BUCKET_VIDEOS` = `videos` ✓
- [ ] `STORAGE_BUCKET_THUMBNAILS` = `thumbnails` ✓
- [ ] `STORAGE_BUCKET_ASSETS` = `assets` ✓
- [ ] JWT_SECRET is NOT "dev-secret" for production ⚠️ *Change this*
- [ ] JWT_REFRESH_SECRET is NOT "dev-refresh-key" for production ⚠️ *Change this*
- [ ] DB credentials are correct ✓

### Local Testing
- [ ] Run backend locally: `go run cmd/main.go`
- [ ] Check logs - should connect to Supabase DB successfully
- [ ] Test API endpoint: `GET http://localhost:5000/api/health` (if exists)
- [ ] No errors in terminal

### Railway Setup
- [ ] Create account at railway.app with GitHub
- [ ] Create new project
- [ ] Connect your GitHub repository
- [ ] Add all environment variables
- [ ] Dockerfile detected correctly
- [ ] Port set to 8080

### Docker Build Test
```bash
# Optional: Test Docker build locally
docker build -t sevima-backend .
docker run -p 8080:8080 --env-file .env sevima-backend
```

### Post-Deployment
- [ ] Deployment shows "Success" in Railway dashboard
- [ ] Check Railway logs - look for "Connected to Supabase"
- [ ] Get Railway domain URL from dashboard
- [ ] Test backend health check: `GET https://your-railway-domain/api/health`
- [ ] Update frontend API URL to use Railway domain

### Integration Tests
- [ ] Frontend can authenticate with deployed backend
- [ ] Can create a project via API
- [ ] Can create a storyboard
- [ ] Can upload files to Supabase (test via API endpoint)
- [ ] Can retrieve files from Supabase

---

## 🚨 Critical Security Notes

**Before deploying to production:**

1. **Generate strong JWT secrets:**
   ```bash
   # On Windows (PowerShell)
   [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((Get-Random -Count 32 | ForEach-Object {[char](48..122 | Get-Random)})) -join '')
   ```

2. **In Railway:**
   - ✅ Set `APP_ENV=production`
   - ✅ Use strong, unique JWT secrets
   - ✅ Never expose `.env` file
   - ✅ Use Railway's built-in secrets management

3. **In Supabase:**
   - ✅ Verify bucket is public only if needed
   - ✅ Configure proper CORS
   - ✅ Review Row-Level Security (RLS) policies
   - ✅ Monitor usage/costs

4. **In Frontend:**
   - ✅ Only use SUPABASE_ANON_KEY (never service role)
   - ✅ Use HTTPS
   - ✅ Implement token refresh logic

---

## 📋 Next Steps After Deployment

1. **Monitor Performance**
   - Check Railway dashboard CPU/Memory usage
   - Review logs for errors
   - Set up alerts for failures

2. **Database Monitoring**
   - In Supabase, check database stats
   - Monitor connection pool
   - Set up slow query logs

3. **CI/CD Pipeline (Optional)**
   - Railway auto-deploys on push (enable in settings)
   - Add tests to GitHub Actions
   - Add build status checks

4. **Scaling (Future)**
   - If needed, upgrade Railway plan
   - Enable autoscaling
   - Consider load balancing

---

## 🆘 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Bucket upload fails | 1. Check CORS 2. Verify bucket is public 3. Check service role key |
| Railway build fails | Check Dockerfile and ensure go.mod is in root |
| Database connection timeout | Verify DB credentials and network access |
| File not accessible via URL | Ensure bucket is set to "Public" |
| 403 Unauthorized on upload | Use SERVICE_ROLE_KEY, not ANON_KEY |
| Frontend can't reach backend | 1. Check CORS headers 2. Verify Railway domain 3. Check firewall |

---

## Quick Links

- 🔗 Supabase Console: https://supabase.com
- 🔗 Railway Dashboard: https://railway.app
- 🔗 Go Documentation: https://golang.org/doc
- 🔗 Supabase Storage Docs: https://supabase.com/docs/guides/storage
- 🔗 Railway Docs: https://docs.railway.app

