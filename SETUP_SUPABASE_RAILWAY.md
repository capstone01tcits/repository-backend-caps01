# Supabase Bucket & Railway Deployment Setup Guide

## Part 1: Supabase Storage Buckets Setup

### Step 1: Access Supabase Console
1. Go to [supabase.com](https://supabase.com)
2. Sign in with your account
3. Select your project: `Sevima-AI-Content-Creator`

### Step 2: Create Storage Buckets
Navigate to **Storage** section in left sidebar, then create these 3 buckets:

#### Bucket 1: Videos
- **Name**: `videos`
- **Access**: Public (to allow video playback)
- **File Size Limit**: Set to 5GB (adjust as needed)

Steps:
1. Click "New Bucket"
2. Name: `videos`
3. Check "Public bucket"
4. Click "Create bucket"

#### Bucket 2: Thumbnails
- **Name**: `thumbnails`
- **Access**: Public
- **File Size Limit**: 500MB

Steps:
1. Click "New Bucket"
2. Name: `thumbnails`
3. Check "Public bucket"
4. Click "Create bucket"

#### Bucket 3: Assets
- **Name**: `assets`
- **Access**: Public
- **File Size Limit**: 1GB

Steps:
1. Click "New Bucket"
2. Name: `assets`
3. Check "Public bucket"
4. Click "Create bucket"

### Step 3: Set Up CORS (Cross-Origin Resource Sharing)
For frontend to access these buckets:

1. In Supabase console, go to **Settings** → **Storage** → **CORS**
2. Add your frontend URLs:
   ```json
   {
     "allowedHeaders": ["*"],
     "allowedMethods": ["GET", "HEAD", "PUT", "POST", "DELETE"],
     "allowedOrigins": [
       "http://localhost:3000",
       "https://your-frontend-domain.vercel.app"
     ],
     "exposedHeaders": ["*"],
     "maxAgeSeconds": 86400
   }
   ```

### Step 4: Test Storage Connection
Create a simple test file:
1. Go to **Storage** → **videos** bucket
2. Click "Upload" and upload a test file
3. Verify you can access it via public URL

---

## Part 2: Railway Deployment Setup

### Step 1: Create Railway Account
1. Go to [railway.app](https://railway.app)
2. Sign up with GitHub (easiest method)
3. Grant permissions

### Step 2: Create New Project
1. Click **"New Project"** or **"+"**
2. Select **"Deploy from GitHub"**
3. Connect your GitHub account
4. Select your backend repository: `Sevima-BackEnd Ai Video Gen`

### Step 3: Configure Environment Variables
In Railway dashboard for your project:

1. Click **"Variables"** tab
2. Add all variables from your `.env` file:

```
APP_PORT=8080
APP_ENV=production

# Database (use your Supabase pool connection)
DB_HOST=aws-1-ap-southeast-2.pooler.supabase.com
DB_PORT=5432
DB_USER=postgres.wkmvrwiesfpnfnaybpbx
DB_PASSWORD=<your-password>
DB_NAME=postgres

# Supabase Storage
SUPABASE_URL=https://wkmvrwiesfpnfnaybpbx.supabase.co
SUPABASE_ANON_KEY=<your-anon-key>
SUPABASE_SERVICE_ROLE_KEY=<your-service-role-key>
STORAGE_BUCKET_VIDEOS=videos
STORAGE_BUCKET_THUMBNAILS=thumbnails
STORAGE_BUCKET_ASSETS=assets

# JWT (use strong keys in production!)
JWT_SECRET=<generate-strong-secret>
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=<generate-strong-refresh-secret>
JWT_REFRESH_EXPIRE_HOURS=168

# AI Service (if using external AI service)
AI_SERVICE_URL=<your-ai-service-url>
```

### Step 4: Configure Build Settings
1. In Railway, go to **Settings** tab
2. Set **"Dockerfile Path"**: `Dockerfile` (should auto-detect)
3. Set **"Start Command"**: Leave empty (uses CMD from Dockerfile)
4. Set **"Port"**: `8080`

### Step 5: Connect to Domain (Optional)
1. In Railway dashboard, go to **Settings** → **Domains**
2. Click **"Generate Domain"** for a Railway subdomain (free)
3. Or connect custom domain:
   - Add your domain (e.g., `api.yourdomain.com`)
   - Railway will show DNS records to update

### Step 6: Deploy
1. Commit your code to GitHub
2. Push to your repo
3. Railway auto-deploys on push (if you've enabled it)
4. Monitor deployment in **"Deployments"** tab

### Step 7: View Logs
1. Click on deployment in **"Deployments"** tab
2. View real-time logs
3. Check for any errors

---

## Part 3: Connect Frontend to Deployed Backend

### Update Frontend API URL
In your frontend project, update the API configuration:

**File**: `src/lib/axios.ts`

```typescript
const baseURL = process.env.NEXT_PUBLIC_API_URL || 
  (process.env.NODE_ENV === 'production' 
    ? 'https://your-railway-domain.railway.app'
    : 'http://localhost:5000');

const instance = axios.create({
  baseURL,
  // ... rest of config
});
```

---

## Part 4: Verify Everything Works

### Test Supabase Storage
1. After deployment, test file upload endpoint:
   ```bash
   curl -X POST http://localhost:8080/api/upload \
     -F "file=@test.mp4" \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

### Test Database Connection
1. In Railway logs, look for database connection messages
2. Should see: "Connected to Supabase database"

### Test API Endpoints
1. Try authentication endpoints
2. Create a project
3. Generate a storyboard
4. Try video generation

---

## Troubleshooting

### Supabase Buckets Issues
- **Buckets not appearing**: Clear cache and refresh
- **Upload fails**: Check CORS settings and bucket permissions
- **Files not accessible**: Ensure bucket is set to "Public"

### Railway Deployment Issues
- **Build fails**: Check logs for Go compilation errors
- **Port already in use**: Railway should handle this, but check PORT env var
- **Database connection fails**: Verify DB credentials in environment variables
- **Storage upload fails**: Check SUPABASE_SERVICE_ROLE_KEY is correct

### Database Connection Issues
```sql
-- Test connection from Railway
psql -h aws-1-ap-southeast-2.pooler.supabase.com -U postgres.wkmvrwiesfpnfnaybpbx -d postgres
```

---

## Security Checklist

- [ ] Use strong JWT secrets (not "dev-secret" in production)
- [ ] Store all secrets in Railway environment variables
- [ ] Never commit `.env` with real credentials
- [ ] Set `APP_ENV=production` for production deployment
- [ ] Use HTTPS for all API calls
- [ ] Enable CORS only for your frontend domain
- [ ] Regularly rotate JWT secrets
- [ ] Use `SUPABASE_SERVICE_ROLE_KEY` only for backend uploads
- [ ] Use `SUPABASE_ANON_KEY` only for frontend public access

---

## Useful Commands

### Check Railway Deployment Status
```bash
# View logs
railway logs

# See environment variables
railway env
```

### Local Testing Before Deploy
```bash
# Build Docker image locally
docker build -t sevima-backend .

# Run locally
docker run -p 8080:8080 --env-file .env sevima-backend
```

### Test Supabase Connection
```bash
# From your backend
go run cmd/main.go
# Should connect without errors
```
