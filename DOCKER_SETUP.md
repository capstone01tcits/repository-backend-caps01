# Docker Setup Guide

Panduan lengkap setup Docker untuk local development dan production deployment.

---

## ⚡ Quick Start (Local Development)

### 1. Install Docker

**Windows:**
```
Download: https://www.docker.com/products/docker-desktop
- Buka Docker Desktop installer
- Follow setup wizard
- Restart Windows
- Verifikasi: docker --version
```

**Check Installation:**
```bash
docker --version
docker-compose --version
```

---

### 2. Start Local Development Environment

```bash
cd d:\ATHA ITS\Capstone\Backend\Sevima-BackEnd Ai Video Gen

# Start all services (Go backend + PostgreSQL + AI service)
docker-compose -f docker-compose.dev.yml up -d

# Check status
docker-compose -f docker-compose.dev.yml ps
```

**Services yang jalan:**
- ✓ Go Backend: http://localhost:5000
- ✓ PostgreSQL: localhost:5432 (postgres/postgres123)
- ✓ AI Service: http://localhost:8000

---

### 3. View Logs

```bash
# All services
docker-compose -f docker-compose.dev.yml logs -f

# Specific service
docker-compose -f docker-compose.dev.yml logs -f app
docker-compose -f docker-compose.dev.yml logs -f postgres
docker-compose -f docker-compose.dev.yml logs -f ai-service
```

---

### 4. Stop Everything

```bash
# Stop all containers
docker-compose -f docker-compose.dev.yml down

# Stop and remove volumes (reset database)
docker-compose -f docker-compose.dev.yml down -v
```

---

## 🏗️ Docker File Architecture

```
Dockerfile               # Go backend build (multi-stage)
Dockerfile.ai          # Python AI service
docker-compose.dev.yml # Local development environment
docker-compose.yml     # Production environment (Railway)
.dockerignore          # Files to exclude from Docker build
```

### Dockerfile Stages

```dockerfile
Stage 1: Builder
├── golang:1.24-alpine base image
├── Download dependencies
├── Build Go binary (CGO disabled for Alpine)
└── Output: main executable

Stage 2: Runtime
├── alpine:3.21 base image
├── Install ca-certificates (for HTTPS)
├── Copy main executable
└── Expose PORT (Railway or APP_PORT)
```

**Why multi-stage?**
- ✓ Smaller image size (~15MB vs ~500MB)
- ✓ Faster deployment
- ✓ No Go compiler in final image (security)

---

## 📁 File Structure Reference

### .dockerignore
```
.git
.gitignore
node_modules
.env
.env.local
.venv
venv
__pycache__
.pytest_cache
.coverage
reports/
README.md
QUICK_DEPLOY.md
DEPLOYMENT_SETUP.md
DOCKER_SETUP.md
```

**Purpose:** Exclude unnecessary files dari Docker build context
- Faster builds
- Smaller image
- No secrets exposed

---

## 🔧 Docker Compose Files

### Development (docker-compose.dev.yml)

```yaml
services:
  app:
    build: .
    ports:
      - "5000:5000"         # Expose port 5000 locally
    environment:
      - APP_PORT=5000       # Use APP_PORT for dev
      - APP_ENV=development
      - DB_HOST=postgres    # Service name
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres123
      - DB_NAME=go_auth
      - JWT_SECRET=dev-secret-change-in-production
      - JWT_REFRESH_SECRET=dev-refresh-secret
      - AI_SERVICE_URL=http://ai-service:8000
    depends_on:
      - postgres
      - ai-service
    volumes:
      - .:/app              # Mount code for hot reload (optional)
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres123
      POSTGRES_DB: go_auth
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  ai-service:
    build:
      context: .
      dockerfile: Dockerfile.ai
    ports:
      - "8000:8000"
    environment:
      - LTX_API_KEY=dummy_key_for_dev
      - RUNWAY_API_KEY=dummy_key_for_dev
    restart: unless-stopped

volumes:
  postgres_data:
```

### Production (docker-compose.yml)

**Used by Railway** - Railway reads this for multi-container apps (optional)

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      # Railway sets PORT, app reads from it
      - APP_ENV=production
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - JWT_SECRET=${JWT_SECRET}
      - JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
      - AI_SERVICE_URL=${AI_SERVICE_URL}
    depends_on:
      - postgres
    restart: always

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

---

## 🚀 Local Development Workflow

### Setup First Time

```bash
# 1. Clone repo
git clone https://github.com/[YOUR_ORG]/Sevima-BackEnd Ai Video Gen.git
cd "Sevima-BackEnd Ai Video Gen"

# 2. Create .env.local (development config)
cp .env.example .env.local

# 3. Edit .env.local untuk development
# DB_HOST=postgres (service name)
# DB_PASSWORD=postgres123
# JWT_SECRET=dev-secret-change-later
# etc.

# 4. Start Docker environment
docker-compose -f docker-compose.dev.yml up -d

# 5. Wait for services (~30 seconds)
docker-compose -f docker-compose.dev.yml logs app | grep "Server running"

# 6. Test API
curl http://localhost:5000/health
# Should return: {"status":"healthy"}
```

### Daily Development

```bash
# Start services
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f app

# Stop services
docker-compose -f docker-compose.dev.yml down

# Reset database
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

### Debug Running Container

```bash
# Execute command in container
docker-compose -f docker-compose.dev.yml exec app bash

# Inside container:
ls -la          # List files
cat .env        # View environment
go version      # Check Go version

# Exit
exit
```

---

## 🐳 Common Docker Commands

```bash
# Build images
docker-compose build                    # Build all services
docker-compose build --no-cache        # Force rebuild
docker build -t sevima-backend .       # Build single image

# Run containers
docker-compose up -d                   # Start in background
docker-compose up                      # Start with logs
docker-compose ps                      # List running containers

# Logs
docker-compose logs -f                 # Follow logs (all)
docker-compose logs app                # Specific service
docker logs -f [CONTAINER_ID]          # By container ID

# Stop/Remove
docker-compose stop                    # Stop containers (preserve data)
docker-compose down                    # Stop & remove containers
docker-compose down -v                 # Remove data volumes too

# Enter container
docker-compose exec app bash           # Execute command
docker-compose exec postgres psql -U postgres  # Connect to DB

# Image management
docker images                          # List images
docker rmi [IMAGE_ID]                 # Remove image
docker system prune                    # Clean up unused resources
```

---

## 🔍 Troubleshooting

### Port Already in Use

```
Error: bind: address already in use
→ Change port in docker-compose.yml: "5001:5000"
→ Or kill process using port:
  netstat -ano | findstr :5000  (Windows)
  taskkill /PID [PID] /F
```

### Container Exits Immediately

```bash
# Check logs
docker-compose logs app

# Common issues:
# 1. Database connection failed → Check DB_HOST, DB_PASSWORD
# 2. Port already in use → Use different port
# 3. Missing .env file → Create from .env.example
# 4. Build errors → Check docker build output
```

### Postgres Connection Refused

```
Error: connection refused
→ Ensure postgres service is running:
  docker-compose ps
→ Wait for healthcheck (~10 sec):
  docker-compose logs postgres
→ Check credentials in docker-compose file match .env
```

### Docker Desktop Not Running

```
Error: Cannot connect to Docker daemon
→ Start Docker Desktop
→ On Windows: Application menu → Docker Desktop
→ Wait for icon in system tray
```

### Volumes Not Mounting

```
Error: Volume not found
→ Check volume name matches in docker-compose
→ Recreate volumes:
  docker volume prune
  docker-compose down -v
  docker-compose up -d
```

---

## 📊 Docker Compose Flow

```
docker-compose up -d
    ↓
1. Read docker-compose.dev.yml
    ↓
2. Build images (if not exist)
    app:      golang:1.24 → build → alpine runtime
    postgres: postgres:15-alpine (pull from registry)
    ai-service: python:3.10 → build
    ↓
3. Create network (sevima-network)
    ↓
4. Create volumes (postgres_data)
    ↓
5. Start containers in order:
    postgres (has healthcheck)
    app (waits for postgres healthy)
    ai-service
    ↓
6. Container running:
    logs show in "docker-compose logs -f"
    access via localhost:5000, :5432, :8000
```

---

## 🎯 Production Deployment (Railway)

Railway reads `Dockerfile` and `Dockerfile.ai`:

```
1. Railway detects repo change (git push)
2. Railway builds image:
   - Uses Dockerfile (Go backend)
   - Runs: docker build -t sevima-backend .
3. Railway sets environment variables (from Dashboard)
4. Railway starts container:
   - docker run -e PORT=8080 -e DB_HOST=... sevima-backend
5. Container starts on PORT 8080
6. App available at: https://[project].up.railway.app
```

**Key Differences from Local:**
- `PORT` env var (Railway sets it)
- Supabase database (not local postgres)
- No volumes (ephemeral storage)
- Auto-restart on crash
- Health checks via HTTP endpoint

---

## ✅ Verification Checklist

- [ ] Docker Desktop installed
- [ ] `docker --version` works
- [ ] `docker-compose --version` works
- [ ] Repository cloned
- [ ] `.env.local` created and configured
- [ ] `docker-compose -f docker-compose.dev.yml up -d` succeeds
- [ ] `curl http://localhost:5000/health` returns status
- [ ] PostgreSQL connected: `docker-compose exec postgres psql -U postgres`
- [ ] View logs: `docker-compose logs -f app`
- [ ] Stop/start works: `docker-compose down && docker-compose up -d`

---

## 📚 Next Steps

1. **Development:**
   - Start with `docker-compose.dev.yml`
   - Use `docker-compose logs -f` for debugging
   - Access database on localhost:5432

2. **Production:**
   - Push to GitHub main branch
   - Railway auto-builds from Dockerfile
   - Monitor via Railway dashboard

3. **Team Development:**
   - Share `.env.example` (not `.env`)
   - Each dev creates own `.env.local`
   - Docker ensures consistent environment

---

**More help?** Check:
- Docker Docs: https://docs.docker.com/
- Docker Compose: https://docs.docker.com/compose/
- Railway Docker: https://docs.railway.app/deploy/dockerfiles
