# AI Video Gen - Backend Service

Go REST API backend untuk AI Video Content Creator platform — Capstone project dengan SEVIMA. Meliputi authentication, project management, storyboard generation, video generation, credit system, dan admin management.

Status: Production Ready (May 2026) - Google Veo 3.1 Lite: High-coherence video generation with atomic project initialization and automated scene polling.

## Documentation

**Comprehensive API & Workflow References:**

- **[API Documentation](docs/API_DOCUMENTATION.md)** — Complete API reference dengan 25 active endpoints dan usage examples
- **[API Collection](docs/API_collection.json)** — Import ke Postman/Bruno untuk testing dengan pre-configured variables dan workflows
- **[Complete Workflow](docs/WORKFLOW_COMPLETE_FLOW.md)** — End-to-end flow dari register hingga video download dengan HTTP requests dan timeline

**Deployment & Production Setup:**

- **[Docker Setup Guide](DOCKER_SETUP.md)** 🐳 — Local development dengan Docker, troubleshooting, commands reference
- **[Quick Deploy Guide](QUICK_DEPLOY.md)** ⚡ — 5 langkah (30 min) setup Supabase + Railway untuk production
- **[Full Deployment Guide](DEPLOYMENT_SETUP.md)** — Detailed guide dengan configuration, monitoring, troubleshooting, dan CI/CD

## Tech Stack

**Go Backend (Port 5000) - Module: `Sevima-AI-Content-Creator`**
- Language: Go 1.21
- Framework: Fiber v2 (lightweight, fast HTTP framework)
- Database: PostgreSQL + GORM (Object-Relational Mapping with AutoMigrate)
- Authentication: JWT (Access + Refresh Token with bcrypt password hashing)
- Role-based Access Control: Middleware enforcing user/admin roles
- Middleware: CORS, Logger, Panic Recovery, JWT Validation
- Dependencies: 
  - `github.com/gofiber/fiber/v2` - Web framework
  - `github.com/golang-jwt/jwt/v5` - JWT token handling
  - `gorm.io/gorm` - ORM layer
  - `gorm.io/driver/postgres` - PostgreSQL driver
  - `golang.org/x/crypto` - Password hashing & cryptographic functions
  - `github.com/google/uuid` - UUID generation
  - `github.com/joho/godotenv` - Environment variable loading

**Python AI Service (Port 8000)**
- Language: Python 3.8+
- Framework: FastAPI + Uvicorn
- Purpose: AI processing microservice for video generation (Wavespeed wrapper)
- Integrated with: Go backend via HTTP calls in background worker

**Database**
- PostgreSQL (auto-created tables via GORM AutoMigrate on startup)
- 10 tables managed: users, projects, business_briefs, creative_briefs, storyboards, storyboard_sections, videos, generation_jobs, video_variants, scene_generations

**Testing & Documentation**
- PowerShell Test Scripts: `/testing/` folder with complete_test.ps1 and add_credits.ps1
- Postman Collection: Importable API collection in `/docs/API_COLLECTION.json`
- Test Reports: Execution reports in `/reports/` folder

## Struktur Folder

```
Sevima-BackEnd Ai Video Gen/
├── cmd/
│   └── main.go                              # Go server entry point
├── config/
│   └── config.go                            # Config loader & DB connection
├── internal/
│   ├── ai/                                  # AI provider integrations
│   │   ├── provider.go                      # Provider interface & models
│   │   └── veo3_provider.go                 # Veo 3 production provider
│   ├── handler/
│   │   ├── auth_handler.go                 # Auth endpoints (register, login, refresh, delete)
│   │   ├── brief_handler.go                # Unified project initialization (creates all briefs)
│   │   ├── credit_handler.go               # Credit management (balance, admin add)
│   │   ├── project_handler.go              # Project management (initialize, list, get)
│   │   ├── storyboard_handler.go           # Storyboard generation
│   │   └── video_handler.go                # Video generation & management (8 core methods)
│   ├── middleware/
│   │   └── auth.go                         # JWT validation & role-based access control
│   ├── migration/
│   │   └── video_generation.go             # Database migration helpers
│   ├── model/
│   │   ├── brief.go                        # BusinessBrief & CreativeBrief models
│   │   ├── generation_job.go               # GenerationJob, VideoVariant, SceneGeneration models
│   │   ├── project.go                      # Project model
│   │   ├── storyboard.go                   # Storyboard & Scene models
│   │   ├── user.go                         # User struct & request types
│   │   └── video.go                        # Video model
│   ├── queue/
│   │   └── job_queue.go                    # Asynchronous job queue system
│   ├── repository/
│   │   ├── brief_repository.go             # Brief database queries (business & creative)
│   │   ├── generation_job_repository.go    # Generation job tracking queries
│   │   ├── project_repository.go           # Project database queries
│   │   ├── scene_generation_repository.go  # Scene generation tracking queries
│   │   ├── storyboard_repository.go        # Storyboard & scene database queries
│   │   ├── user_repository.go              # User database queries
│   │   ├── video_repository.go             # Video database queries
│   │   └── video_variant_repository.go     # Video variant database queries
│   └── service/
│       ├── auth_service.go                 # Authentication logic
│       ├── brief_service.go                # Unified project initialization logic
│       ├── credit_service.go               # Credit management
│       ├── project_service.go              # Project management
│       ├── storyboard_service.go           # Storyboard generation
│       └── video_generation_service.go     # Video generation & regeneration
├── pkg/
│   └── utils/
│       ├── jwt.go                          # JWT token generation & validation
│       └── response.go                     # Standardized JSON response format
├── ai-service/
│   ├── main.py                             # Python AI service entry point
│   └── requirements.txt                    # Python dependencies (FastAPI, etc)
├── docs/
│   ├── API_COLLECTION.json                 # Postman v2.1 importable collection
│   ├── API_DOCUMENTATION.md                # Comprehensive API documentation
│   └── WORKFLOW_COMPLETE_FLOW.md           # End-to-end workflow documentation
├── testing/
│   ├── add_credits.ps1                     # PowerShell script for adding credits
│   ├── complete_test.ps1                   # Complete test suite script
│   ├── INDEX.md                            # Testing documentation index
│   └── QUICKSTART.md                       # Quick start guide for testing
├── reports/
│   ├── COMPLETE_TESTING_SUITE.md           # Comprehensive testing documentation
│   ├── CONSOLIDATED_TESTING_REFERENCE.md  # Testing reference guide
│   ├── ENDPOINT_TEST_REPORT_*.txt          # Endpoint test execution reports
│   └── INDEX.md                            # Reports index & navigation
├── .env                                    # Environment variables (copy from .env.example)
├── .env.example                            # Environment variables template
├── .gitignore                              # Git ignore rules
├── docker-compose.yml                      # Docker Compose configuration
├── Dockerfile                              # Docker image build configuration
├── go.mod                                  # Go module definition (Sevima-AI-Content-Creator)
├── go.sum                                  # Go dependencies lock file
├── README.md                               # This file
└── Dockerfile                              # Docker build configuration
```

## Database Initialization

Database tables are auto-created on startup via GORM AutoMigrate in `cmd/main.go`. 10 active tables are managed:

| Table | Purpose | Key Fields |
|-------|---------|-----------|
| `users` | User account & credentials | id, email (unique), password (bcrypt), role (user/admin), credits (default 10) |
| `projects` | Container for one video project | id, user_id (FK), name, theme, status |
| `business_briefs` | Business context from institution | id, project_id (FK), institute_name, school_level, target_audience, key_message, logo_path, environment_path |
| `creative_briefs` | Creative direction & copywriting | id, business_brief_id (FK), video_type, duration, style, tone, script, copywriting, hashtags |
| `storyboards` | Blueprint containing 3 mandatory scenes | id, project_id (FK), title, description |
| `storyboard_sections` | Individual section (Hook/Value/CTA) | id, storyboard_id (FK), section_type, content, duration |
| `generation_jobs` | Video generation queue task tracking | id, project_id (FK), storyboard_id (FK), user_id (FK), status, provider, model |
| `video_variants` | Single video output variant | id, status, video_url, thumbnail_url, file_size, credits_used, provider, model |
| `scene_generations` | Progress per scene | id, variant_id (FK), scene_number, status, video_url |
| `videos` | Final video record | id, project_id (FK), storyboard_id (FK), video_url |

### Database Workflow Usage

Each table is used at specific workflow stages:

```
Register
  └─> INSERT users (user, password hash, 10 initial credits)

Login
  └─> SELECT users (validate email/password)

Initialize Project (FE Wizard)
  └─> INSERT projects (container for campaign)
  └─> INSERT business_briefs (institution context)
  └─> INSERT creative_briefs (creative direction)
  └─> INSERT storyboards (auto-generated blueprint)
  └─> INSERT scenes (3 mandatory scenes: Hook, Value, CTA)

Generate Video (Async)
  └─> INSERT generation_jobs (status=queued, create job)
  └─> INSERT video_variants (1 variant: cinematic/Google Veo 3.1 Lite)
  └─> UPDATE users (deduct 1 credit)
  └─> Enqueue to background worker channel

Background Worker Processing
  └─> UPDATE generation_jobs (status=generating_assets)
  ├─> Create Veo 3 JSON Payload
  ├─> For each scene (Hook, Value, CTA):
  │   ├─> INSERT scene_generations (status=queued)
  │   ├─> Call AI Service (Google Veo 3.1 Lite)
  │   ├─> UPDATE scene_generations (status=completed, video_url)
  │   └─> UPDATE generation_jobs (progress++)
  ├─> UPDATE generation_jobs (status=stitching_video)
  ├─> Merge all scene videos (FFmpeg)
  ├─> INSERT videos (final MP4)
  ├─> UPDATE video_variants (status=completed, video_url)
  └─> UPDATE generation_jobs (status=completed)

Check Video Status
  └─> SELECT video_variants (get status & progress)
  └─> SELECT scene_generations (get per-scene status)

Download Video
  └─> SELECT videos (get final video_url)
```

No manual migration needed — just start the server and tables will be created/updated automatically via `AutoMigrate()` in main.go.

## Environment Configuration

Create `.env` file:
```bash
cp .env.example .env
```

Configure in `.env`:
```
APP_PORT=5000
APP_ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=go_auth
JWT_SECRET=your_jwt_secret
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=your_refresh_secret
JWT_REFRESH_EXPIRE_HOURS=168
AI_SERVICE_URL=http://localhost:8000
```

## Setup & Run

### Option 1: Docker (Recommended)

```bash
docker-compose up -d
```

This runs both PostgreSQL and Go service. For Python AI service:
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

### Option 2: Manual

**Go Backend:**
```bash
go mod tidy
go run cmd/main.go
```

Server berjalan di `http://localhost:5000`

**Python AI Service:**
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

Service berjalan di `http://localhost:8000`

## API Endpoints (27 Active)

### Health Check (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check Go backend status |

### Authentication (6 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | Public | Register new user (default role=user, credits=10) |
| POST | `/api/auth/login` | Public | Login and get access/refresh tokens |
| POST | `/api/auth/refresh` | Public | Refresh access token using refresh token |
| GET | `/api/auth/me` | Protected | Get current user profile (includes role & credits) |
| POST | `/api/auth/change-password` | Protected | Change user password |
| DELETE | `/api/auth/account` | Protected | Delete user account (soft delete) |

### Projects (5 endpoints - Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/initialize` | Initialize project, briefs, and auto-generate storyboard (Linear Flow) |
| GET | `/api/projects` | List current user's projects |
| GET | `/api/projects/:id` | Get project by ID |
| DELETE | `/api/projects/:id` | Soft delete project |
| POST | `/api/projects/:id/restore` | Restore soft-deleted project |

### Storyboard (6 endpoints - Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/storyboard/create` | Create manual storyboard (3 sections: hook/value/cta) |
| GET | `/api/storyboard/:project_id` | Get the storyboard for a project |
| GET | `/api/storyboard/detail/:storyboard_id` | Get single storyboard with sections |
| PUT | `/api/storyboard/:storyboard_id` | Update storyboard and sections |
| DELETE | `/api/storyboard/:storyboard_id` | Soft delete storyboard |
| POST | `/api/storyboard/:storyboard_id/restore` | Restore soft-deleted storyboard |
| GET | `/api/storyboard/:storyboard_id/sections` | Get sections for a storyboard |

### Videos (7 endpoints - Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/videos/generate` | Generate video from storyboard (creates job) |
| GET | `/api/videos/:id` | Get video details with scenes |
| GET | `/api/videos` | List user's videos |
| GET | `/api/videos/preview/:id` | Stream/Preview video directly via Supabase CDN |
| GET | `/api/videos/download/:id` | Get video download URL and metadata |
| POST | `/api/videos/:variantId/regenerate` | Regenerate video (using new prompt) |
| POST | `/api/videos/scene/:sceneId/regenerate` | Regenerate individual scene |

### Credits & Admin (2 endpoints - Protected)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/credits` | Protected | Get current user's credit balance |
| POST | `/api/admin/credits` | Admin | Add credits to a user |

### AI Gateway (2 endpoints)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/ai/health` | Check Python AI service connectivity |
| ANY | `/api/ai/*` | Proxy request to Python AI service (Protected) |

All AI gateway endpoints inject user context headers (`X-User-ID`, `X-User-Email`) automatically.

## Video Generation System

The backend includes a **scene-based video generation system** that creates a high-quality video variation from a single storyboard:

### Video Generation Overview

**Key Features:**
- **Single High-Quality Video**: Automatically generates one high-quality video per storyboard using the Veo 3 model.
- **Scene-Based Architecture**: Video composed of 3 mandatory scenes (Hook, Value, CTA).
- **Veo 3 & Wavespeed**: Integration with specialized `Veo3Payload` for cinematic continuity and automated FFmpeg stitching.
- **Queue System**: Background workers process videos asynchronously via `GenerationJob`.
- **Job Tracking**: Full job lifecycle: queued → generating_assets → stitching_video → completed/failed.
- **Regeneration**: Regenerate the full video or individual scenes with new prompts.
- **Retry Logic**: Automatic retry mechanism with configurable max retries per job.

### Models & Components

**Core Models** (in `internal/model/`):
- `GenerationJob` - Video generation task (tracks status, provider, credits, retries)
- `VideoVariant` - Individual video variant
- `SceneGeneration` - Individual scene generation tracking (progress at scene level)
- `Video` - Legacy model (can be deprecated in favor of VideoVariant)

**Repositories** (in `internal/repository/`):
- `GenerationJobRepository` - CRUD & query operations for jobs
- `VideoVariantRepository` - CRUD & query operations for variants
- `SceneGenerationRepository` - CRUD operations for scene tracking

**Service** (in `internal/service/`):
- `VideoGenerationService` - Business logic for generation, regeneration, status polling

**Queue System** (in `internal/queue/`):
- `JobQueue` - Manages async job processing and worker coordination

**Providers** (in `internal/ai/`):
- `VideoProvider` interface - Defines provider contract
- `Veo3Provider` - Production implementation hitting Python AI Service
- `ProviderFactory` - Factory for creating appropriate provider instances

### Video Generation Endpoints & Flow

Complete endpoints for video generation:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/videos/generate` | Generate video from storyboard (deducts 1 credit) |
| GET | `/api/videos/:id` | Check generation status and get video details |
| GET | `/api/videos` | List all user videos |
| GET | `/api/videos/preview/:id` | Preview video directly via Supabase CDN (Great for `<video>` tags) |
| GET | `/api/videos/download/:id` | Download video file |
| POST | `/api/videos/:id/regenerate` | Regenerate entire video |
| POST | `/api/videos/scene/:sceneId/regenerate` | Regenerate individual scene |

**Video CDN Integration (Supabase Storage)**
Semua video yang digenerate oleh AI Service akan secara otomatis di-download oleh Backend Go dan di-upload ke Supabase Storage Bucket. Frontend tidak perlu berurusan dengan local file path atau API AI secara langsung. FE cukup memanggil endpoint `GET /api/videos/preview/:id` yang akan otomatis merespon dengan `302 Redirect` ke **Public CDN URL** dari Supabase, sehingga video bisa langsung diputar di HTML `<video src="...">`.

### Quick Video Generation CLI Example

```bash
# 1. Register & Login
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"pass123"}'

# Save token from response: export TOKEN="eyJ..."

# 2. Create Project
curl -X POST http://localhost:5000/api/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Video Project"}' # Save project_id

# 3. Create Business Brief
curl -X POST http://localhost:5000/api/briefs/business \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "project_name": "My Video",
    "company_name": "Acme Corp",
    "industry": "Tech",
    "target_audience": "Developers"
  }' # Save business_brief_id

# 6. Generate Video
# Note: Storyboard ID is now returned directly from Project Initialization
curl -X POST http://localhost:5000/api/videos/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "storyboard_id": "{storyboard_id}"
  }' # Get generation_job_id

# 7. Poll Generation Status
curl http://localhost:5000/api/videos/{generation_job_id} \
  -H "Authorization: Bearer $TOKEN"

# 8. Download Video
curl http://localhost:5000/api/videos/download/{generation_job_id} \
  -H "Authorization: Bearer $TOKEN"
```

### Job Lifecycle

```
POST /api/videos/generate
    |
    v
1. API validates credits, storyboard, project
2. VideoGenerationService creates GenerationJob (status: queued)
3. Video deducts 1 credit from user
4. Creates 1 VideoVariant record
5. Enqueues job in JobQueue

    |
    v
Background Workers start processing
    |
    v
1. Fetch Storyboard Sections
2. Build Veo3Payload with all 3 scenes
3. Call AI Service /api/veo3/generate (Wavespeed/Veo 3)
4. AI Service processes scenes and joins them with FFmpeg
5. VideoVariant marked as completed
6. GenerationJob marked as completed
        
    |
    v
GET /api/videos/{jobId} returns status: completed
```

### Credit Calculation

**Standard Costs (single generation):**
- Google Veo 3.1 Lite (Wavespeed): 1 credit per full video generation

**Example (1 variant of 18-second video):**
- Total deducted: 1 credit (fixed per generation call)

## Quick Copy-Paste Commands

### 1. Health Check (No Auth)
```bash
curl http://localhost:5000/health
```

### 2. Register User
```bash
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure123"
  }'
```

### 3. Login (save tokens)
```bash
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secure123"
  }'
```

### 4. Get Profile (Protected)
```bash
curl http://localhost:5000/api/auth/me \
  -H "Authorization: Bearer {access_token}"
```

### 5. Create Project
```bash
curl -X POST http://localhost:5000/api/projects \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Brand Video Campaign Q3",
    "description": "A promotional video campaign for Q3 2026"
  }'
```

### 6. List Projects
```bash
curl http://localhost:5000/api/projects \
  -H "Authorization: Bearer {access_token}"
```

### 7. Generate Video (costs 1 credit)
```bash
curl -X POST http://localhost:5000/api/videos/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "storyboard_id": "{storyboard_id}"
  }'
```

### 8. View & Download Video
```bash
curl http://localhost:5000/api/videos/{video_id} \
  -H "Authorization: Bearer {access_token}"

curl http://localhost:5000/api/videos/download/{video_id} \
  -H "Authorization: Bearer {access_token}"
```

### 9. Check Credits
```bash
curl http://localhost:5000/api/credits \
  -H "Authorization: Bearer {access_token}"
```

### 10. Admin: Add Credits to User
```bash
curl -X POST http://localhost:5000/api/admin/credits \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "{target_user_id}", "amount": 10}'
```

## Usage Examples with Responses

### Register User

```bash
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10
    }
  }
}
```

### Login

```bash
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secure123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10
    }
  }
}
```

### Create Project

```bash
curl -X POST http://localhost:5000/api/projects \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Brand Video Q3", "description": "Q3 campaign"}'
```

Response:
```json
{
  "success": true,
  "message": "Project created successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "name": "Brand Video Q3",
    "description": "Q3 campaign",
    "status": "draft",
    "created_at": "2026-03-10T10:00:00Z"
  }
}
```



### Generate Video

```bash
curl -X POST http://localhost:5000/api/videos/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "storyboard_id": "{storyboard_id}",
    "title": "My Brand Video"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Video generated successfully",
  "data": {
    "id": "uuid",
    "project_id": "uuid",
    "storyboard_id": "uuid",
    "title": "My Brand Video",
    "status": "completed",
    "video_url": "https://storage.example.com/videos/uuid.mp4",
    "thumbnail_url": "https://storage.example.com/thumbnails/uuid.jpg",
    "duration": 120,
    "format": "mp4",
    "resolution": "1080p",
    "credits_used": 1
  }
}
```

## Response Format

All responses follow a standardized format:

```json
{
  "success": boolean,
  "message": "Descriptive message",
  "data": {} or []
}
```

Error Response:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error"
}
```

## Database Schema

All tables are auto-created via GORM AutoMigrate. See `cmd/main.go` for the 10 active models.
For detailed field definitions, see `internal/model/*.go` source files.

| Table | Source Model | Purpose |
|-------|--------------|---------|
| users | `model.User` | Accounts, auth, credits |
| projects | `model.Project` | Video project containers |
| business_briefs | `model.BusinessBrief` | Institution & business context |
| creative_briefs | `model.CreativeBrief` | Creative direction & copywriting |
| storyboards | `model.Storyboard` | Storyboard blueprints |
| storyboard_sections | `model.StoryboardSection` | 3-part sections (hook/value/cta) |
| videos | `model.Video` | Generated video records |
| generation_jobs | `model.GenerationJob` | Async job queue tracking |
| video_variants | `model.VideoVariant` | video variations per job |
| scene_generations | `model.SceneGeneration` | Per-scene generation tracking |

## Project Architecture

```
Client/Frontend (Web/Mobile)
    |
    v
Go Backend (Port 5000) - Sevima-AI-Content-Creator
    |
    +-- Middleware Layer
    |   ├── JWT Validation (Protected routes)
    |   ├── Role-Based Access Control (user/admin)
    |   ├── CORS, Logger, Panic Recovery
    |
    +-- Handler Layer (7 handlers)
    |   ├── Auth Handler        --> Auth Service --> User Repository --> PostgreSQL
    |   ├── Project Handler     --> Project Service --> Project Repository --> PostgreSQL
    |   ├── Brief Handler       --> Brief Service --> Brief Repository --> PostgreSQL
    |   ├── Storyboard Handler  --> Storyboard Service --> Storyboard Repository --> PostgreSQL
    |   ├── Video Handler       --> Video Gen Service --> Job Queue --> Generation Workers
    |   ├── Credit Handler      --> Credit Service --> User Repository --> PostgreSQL
    |
    +-- Service Layer (Business Logic)
    |   ├── Auth Service
    |   ├── Project Service
    |   ├── Brief Service
    |   ├── Storyboard Service
    |   ├── Video Generation Service
    |   ├── Storage Service
    |   └── Credit Service
    |
    +-- Repository Layer (Data Access)
    |   ├── User Repository
    |   ├── Project Repository
    |   ├── Brief Repository (Business + Creative)
    |   ├── Storyboard Repository (Storyboards + Sections)
    |   ├── Generation Job Repository
    |   ├── Video Variant Repository
    |   └── Scene Generation Repository
    |
    +-- Queue System (Async Processing)
    |   └── Job Queue --> Background Workers --> Video Generation Providers
    |
    +-- AI Provider Layer
    |   └── Veo 3 Provider (Wavespeed API)
    |
    v
PostgreSQL Database (10 tables via AutoMigrate)
    |
    v
Python AI Service (Port 8000) - FastAPI
    |
    └── Video Generation (via provider integrations)
```

**Key Design Patterns:**
- **Repository Pattern**: Abstraction layer for all database operations
- **Service Pattern**: Business logic separated from handlers
- **Middleware Pattern**: Centralized auth & access control
- **Proxy Pattern**: AI Gateway transparently routes to Python service
- **Queue Pattern**: Async job processing for video generation
- **Provider Pattern**: Pluggable AI provider implementations

## Video Generation Workflow

```
1. User registers/logs in (gets 10 credits by default)
        │
        ▼
2. Initialize Project (FE Wizard - creates project + briefs atomically)
   POST /api/projects/initialize
        │
        ▼
3. Generate Storyboard Templates OR Create Manual Storyboard
   POST /api/storyboard/templates/generate  (4 template options)
   POST /api/storyboard/create              (manual 3-section)
        │
        ▼
4. Select Storyboard & Generate Video (deducts credits)
   POST /api/storyboard/select
   POST /api/videos/generate
        │
        ▼
5. Poll Status & Download Video
   GET /api/videos/:id
   GET /api/videos/download/:id
```

## Authorization & Middleware

Protected endpoints require Authorization header:

```
Authorization: Bearer {access_token}
```

Token is obtained from `/api/auth/login` response. Access tokens expire after `JWT_EXPIRE_HOURS` (default 24 hours). Use refresh token to get new access token via `/api/auth/refresh`.

All protected routes are enforced through `middleware.Protected()` middleware that validates JWT token and attaches user context.

### Roles & Permissions

| Role | Default Credits | Capabilities |
|------|----------------|--------------|
| `user` | 10 | All standard endpoints (projects, storyboards, videos, credits) |
| `admin` | 10 | All user capabilities + admin routes (`POST /api/admin/credits`) |

Role enforcement via `middleware.RequireRole("admin")` middleware on admin routes.

## Credit System

- New users start with **10 credits** by default (configured in `model.User` with `default:10`)
- Each video generation costs **1 credit**
- Check balance: `GET /api/credits` - returns current user's credit balance
- Admin can add credits: `POST /api/admin/credits` - requires admin role
- Video generation is blocked when credits reach 0 (deduction happens before generation starts)
- Credit deduction is handled by `service.creditService.DeductCredits()` called by video generation service
- Credits are stored in `users.credits` column as integer

### Credit Flow

1. User calls `POST /api/videos/generate`
2. Video Handler validates user has sufficient credits
3. Credit Service deducts 1 credit from user's account
4. Video variant generation is enqueued as a `GenerationJob`
5. If generation completes, credits remain deducted
6. If generation fails, credits are **not** refunded (by design)
7. Admin can manually add credits for failed jobs via `POST /api/admin/credits`

## Soft Delete & Recovery

- When account is deleted, `deleted_at` is set to current timestamp
- User cannot login after deletion
- Can restore account using refresh_token via `/api/auth/restore` endpoint
- After restoration, `deleted_at` is set to null and account is fully active again
- Business briefs, creative briefs, projects, and other resources also support soft delete

## API Documentation

For comprehensive API documentation with full request/response schemas, see:
- [API Documentation](docs/API_DOCUMENTATION.md)
- [API Collection](docs/API_collection.json) — import into Postman/Bruno for ready-to-use requests with auto-token-saving scripts

## Testing & Validation

### Testing Tools

The project includes comprehensive testing tools in the `/testing/` folder:

**PowerShell Scripts:**
- `complete_test.ps1` - Full endpoint test suite (tests all major API flows)
- `add_credits.ps1` - Utility script for adding credits during testing
- `INDEX.md` - Testing documentation index
- `QUICKSTART.md` - Quick start guide for running tests

**Reports:**
The `/reports/` folder contains:
- `ENDPOINT_TEST_REPORT_*.txt` - Detailed endpoint test execution reports with timestamps
- `COMPLETE_TESTING_SUITE.md` - Comprehensive testing documentation
- `CONSOLIDATED_TESTING_REFERENCE.md` - Testing reference guide for all endpoints
- `INDEX.md` - Reports navigation and summary

### Running Tests

**Quick Start (Windows PowerShell):**
```powershell
# Navigate to testing directory
cd testing/

# Run complete test suite
.\complete_test.ps1

# Add credits to a test user
.\add_credits.ps1
```

See `/testing/QUICKSTART.md` for detailed instructions and prerequisites.

### Postman Testing

Import `/docs/API_COLLECTION.json` into Postman for ready-to-use requests with:
- All endpoints pre-configured
- Example request/response bodies
- Authentication token variables
- Environment setup documentation

## Troubleshooting

**Problem: "AI service is unreachable"**
- Ensure Python AI service is running on port 8000
- Check `AI_SERVICE_URL` in .env configuration
- Run: `cd ai-service && python main.py`
- Verify network connectivity between Go and Python services

**Problem: "Invalid token / Unauthorized"**
- Token may have expired, use refresh endpoint: `POST /api/auth/refresh`
- Ensure Authorization header format is correct: `Authorization: Bearer {token}`
- Check `JWT_SECRET` and `JWT_REFRESH_SECRET` match in .env
- Verify token is from `/api/auth/login` or `/api/auth/register`, not refresh token

**Problem: "Connection refused / Cannot connect to database"**
- Go backend listening on port 5000
- PostgreSQL running on configured `DB_HOST:DB_PORT` (default localhost:5432)
- Check database credentials in .env file
- Verify PostgreSQL service is running: verify in services (Windows) or `sudo systemctl status postgresql` (Linux)

**Problem: "Database migration errors"**
- Tables should auto-create via `AutoMigrate()` on startup
- Check PostgreSQL has proper permissions for the configured DB_USER
- Drop and recreate database if schema is corrupted: 
  ```bash
  dropdb go_auth
  createdb go_auth
  # Then restart Go server to auto-migrate
  ```

**Problem: "Credits not deducting / Video generation blocked"**
- User must have sufficient credits (initial: 10 credits per user)
- Each video generation costs 1 credit
- Admin can add credits via `POST /api/admin/credits`
- Check user credit balance: `GET /api/credits`
