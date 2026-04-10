# AI Video Gen - Backend Service

Go REST API backend untuk AI Video Content Creator platform — Capstone project dengan SEVIMA. Meliputi authentication, project management, creative brief, content pillar generation, storyboard generation, video generation, credit system, dan admin management.

## Documentation

**Comprehensive API & Workflow References:**

- **[API Documentation](docs/API_DOCUMENTATION.md)** — Complete API reference dengan semua endpoints, request/response examples, dan sistem video generation
- **[Postman Collection](docs/postman_collection.json)** — Import ke Postman untuk testing (termasuk semua endpoints dengan pre-configured variables)
- **[Complete Workflow](docs/WORKFLOW_COMPLETE_FLOW.md)** — End-to-end flow dari register hingga video download dengan HTTP requests, database changes, dan timeline

## Tech Stack

**Go Backend (Port 3000) - Module: `Sevima-AI-Content-Creator`**
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
- Purpose: AI processing microservice for content pillars, storyboard, and video generation
- Integrated with: Go backend via HTTP proxy gateway

**Database**
- PostgreSQL (auto-created tables via GORM AutoMigrate on startup)
- 12 tables managed: users, projects, business_briefs, creative_briefs, content_pillars, content_themes, storyboards, scenes, videos, generation_jobs, video_variants, scene_generations

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
│   │   └── providers.go                     # Provider implementations (LTX, Runway, etc)
│   ├── handler/
│   │   ├── ai_handler.go                   # AI proxy handler
│   │   ├── auth_handler.go                 # Auth handlers (register, login, refresh, restore)
│   │   ├── brief_handler.go                # Brief CRUD handlers (business & creative briefs)
│   │   ├── content_handler.go              # Content Pillar & Theme handlers (generation, selection)
│   │   ├── credit_handler.go               # Credit management handlers (balance, add)
│   │   ├── project_handler.go              # Project CRUD handlers (dashboard)
│   │   ├── storyboard_handler.go           # Storyboard & Scene handlers (generation, selection)
│   │   └── video_handler.go                # Video generation handlers (variants, scenes, regeneration)
│   ├── middleware/
│   │   └── auth.go                         # JWT validation & role-based access control
│   ├── migration/
│   │   └── video_generation.go             # Database migration helpers
│   ├── model/
│   │   ├── brief.go                        # BusinessBrief & CreativeBrief models
│   │   ├── content.go                      # ContentPillar & ContentTheme models
│   │   ├── generation_job.go               # GenerationJob, VideoVariant, SceneGeneration models
│   │   ├── project.go                      # Project model
│   │   ├── storyboard.go                   # Storyboard & Scene models
│   │   ├── user.go                         # User struct & request types
│   │   └── video.go                        # Video model
│   ├── queue/
│   │   └── job_queue.go                    # Asynchronous job queue system
│   ├── repository/
│   │   ├── brief_repository.go             # Brief database queries (business & creative)
│   │   ├── content_repository.go           # Content pillar & theme database queries
│   │   ├── generation_job_repository.go    # Generation job tracking queries
│   │   ├── project_repository.go           # Project database queries
│   │   ├── scene_generation_repository.go  # Scene generation tracking queries
│   │   ├── storyboard_repository.go        # Storyboard & scene database queries
│   │   ├── user_repository.go              # User database queries
│   │   ├── video_repository.go             # Video database queries
│   │   └── video_variant_repository.go     # Video variant database queries
│   └── service/
│       ├── auth_service.go                 # Auth business logic
│       ├── brief_service.go                # Brief business logic
│       ├── content_service.go              # Content pillar/theme generation logic
│       ├── credit_service.go               # Credit deduction & balance management
│       ├── project_service.go              # Project business logic
│       ├── storyboard_service.go           # Storyboard generation & scene logic
│       └── video_generation_service.go     # Video variant generation & regeneration logic
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
└── app.exe                                 # Compiled Go executable (Windows)
```

## Database Initialization

Database tables are auto-created on startup via GORM AutoMigrate in `cmd/main.go`. Twelve tables are managed:

| Table | Model | Description |
|-------|-------|-------------|
| `users` | `model.User` | User accounts with role & credits, soft delete |
| `projects` | `model.Project` | Video projects (container for full workflow) |
| `business_briefs` | `model.BusinessBrief` | Business brief forms (FK → users, projects) |
| `creative_briefs` | `model.CreativeBrief` | Creative brief forms (FK → users, business_briefs) |
| `content_pillars` | `model.ContentPillar` | AI-generated content pillars (FK → projects, users) |
| `content_themes` | `model.ContentTheme` | Content themes per pillar (FK → content_pillars, users) |
| `storyboards` | `model.Storyboard` | Storyboard variations (FK → projects, users) |
| `scenes` | `model.Scene` | Individual scenes in storyboards (FK → storyboards, users) |
| `videos` | `model.Video` | Generated videos (FK → projects, storyboards, users) |
| `generation_jobs` | `model.GenerationJob` | Video generation queue jobs (track status & progress) |
| `video_variants` | `model.VideoVariant` | 3 video variations per storyboard (cinematic, vibrant, professional) |
| `scene_generations` | `model.SceneGeneration` | Individual scene generation tracking (scene-level progress) |

No manual migration needed — just start the server and tables will be created/updated automatically via `AutoMigrate()` in main.go.

## Environment Configuration

Create `.env` file:
```bash
cp .env.example .env
```

Configure in `.env`:
```
APP_PORT=3000
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

Server berjalan di `http://localhost:3000`

**Python AI Service:**
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

Service berjalan di `http://localhost:8000`

## API Endpoints

### Health Check (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check Go backend status |
| GET | `/api/ai/health` | Check Python AI service connectivity |

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user (default role=user, credits=10) |
| POST | `/api/auth/login` | Login and get access/refresh tokens |
| POST | `/api/auth/refresh` | Refresh access token using refresh token |
| POST | `/api/auth/restore` | Restore deleted account (soft delete recovery) |

### Profile (Protected - requires Bearer token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Get current user profile (includes role & credits) |
| GET | `/api/auth/users/:user_id` | Get another user's profile (returns limited info) |
| POST | `/api/auth/change-password` | Change user password |
| DELETE | `/api/auth/account` | Delete user account (soft delete) |

### Projects / Dashboard (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects` | Create new project |
| GET | `/api/projects` | List current user's projects |
| GET | `/api/projects/:id` | Get project by ID (with details) |
| PUT | `/api/projects/:id` | Update project (name, description, status) |
| DELETE | `/api/projects/:id` | Delete project (soft delete) |

### Business Brief (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/briefs/business` | Create new business brief (with optional project_id) |
| GET | `/api/briefs/business` | List current user's business briefs |
| GET | `/api/briefs/business/:id` | Get business brief by ID |
| PUT | `/api/briefs/business/:id` | Update business brief |
| DELETE | `/api/briefs/business/:id` | Delete business brief (soft delete) |
| GET | `/api/briefs/business/:id/creative` | List creative briefs under a business brief |

### Creative Brief (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/briefs/creative` | Create new creative brief (linked to business brief) |
| GET | `/api/briefs/creative` | List current user's creative briefs |
| GET | `/api/briefs/creative/:id` | Get creative brief by ID |
| PUT | `/api/briefs/creative/:id` | Update creative brief |
| DELETE | `/api/briefs/creative/:id` | Delete creative brief (soft delete) |

### Content Pillars (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/:id/content-pillars/generate` | AI-generate content pillars (3 variations with themes) for project |
| GET | `/api/projects/:id/content-pillars` | List all content pillars for a project |
| GET | `/api/content-pillars/:id` | Get single content pillar by ID |
| PUT | `/api/content-pillars/:id` | Update content pillar |
| POST | `/api/content-pillars/:id/select` | Select a content pillar as active |
| GET | `/api/content-pillars/:id/themes` | List all themes under a content pillar |

### Content Themes (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/content-themes/:id/select` | Select a content theme as active |

### Storyboards (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/:id/storyboards/generate` | AI-generate storyboard variations (2-3 variations) for project |
| GET | `/api/projects/:id/storyboards` | List all storyboards for a project |
| GET | `/api/storyboards/:id` | Get storyboard by ID (with scenes) |
| PUT | `/api/storyboards/:id` | Update storyboard |
| POST | `/api/storyboards/:id/select` | Select a storyboard as active |
| GET | `/api/storyboards/:id/scenes` | List all scenes in a storyboard |

### Videos (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/videos/generate` | Generate 3 video variants (deducts 1 credit) |
| GET | `/api/videos/generation/:jobId` | Check video generation job status (polling) |
| GET | `/api/videos/storyboard/:storyboardId` | Get all 3 video variants for a storyboard |
| GET | `/api/videos/:variantId` | Get single video variant by ID (with all scenes) |
| GET | `/api/videos/:variantId/download` | Get video download URL/info |
| POST | `/api/videos/:variantId/regenerate` | Regenerate entire video variant with new prompt |
| POST | `/api/videos/scene/:sceneId/regenerate` | Regenerate individual scene within a video |

### Credits (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/credits` | Get current user's credit balance |

### Admin (Protected - requires admin role)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/credits` | Add credits to a user (admin only) |

### AI Gateway (Protected - routes to Python AI Service)

| Method | Endpoint | Description |
|--------|----------|-------------|
| ANY | `/api/ai/*` | Proxy any request to Python AI service with user context headers |

All AI gateway endpoints inject user context headers (`X-User-ID`, `X-User-Email`) automatically.

## Video Generation System

The backend includes a **scene-based video generation system** that creates 3 video variations from a single storyboard:

### Video Generation Overview

**Key Features:**
- **3 Video Variants**: Automatically generates 3 variations per storyboard (e.g., cinematic, vibrant, professional)
- **Scene-Based Architecture**: Each video composed of individual scenes (4-6 sec each), total 8-12 seconds per video
- **Multiple Providers**: Pluggable provider architecture supporting:
  - **LTX-2-Fast** (standard tier) - 2 credits/second
  - **LTX-2-Pro** (premium tier) - 3 credits/second
  - **Runway** (gen4.5 model) - Variable credits
  - **Wan 2.1** (research tier) - Variable credits
  - **Open-source models** - 1 credit/second
- **Queue System**: Background workers process videos asynchronously via `GenerationJob`
- **Job Tracking**: Full job lifecycle: queued → processing → completed/failed
- **Regeneration**: Regenerate full videos or individual scenes with new prompts
- **Retry Logic**: Automatic retry mechanism with configurable max retries per job

### Models & Components

**Core Models** (in `internal/model/`):
- `GenerationJob` - Video generation task (tracks status, provider, credits, retries)
- `VideoVariant` - Individual video variant (one of 3 per storyboard)
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
- `LTXStandardProvider` - LTX-2-Fast implementation
- `LTXPremiumProvider` - LTX-2-Pro implementation
- `RunwayProvider` - Runway gen4.5 implementation
- `Wan2Provider` - Wan 2.1 implementation
- `ProviderFactory` - Factory for creating appropriate provider instances

### Video Generation Endpoints & Flow

Complete endpoints for video generation:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/videos/generate` | Generate 3 video variants from storyboard (deducts 1 credit) |
| GET | `/api/videos/generation/:jobId` | Check generation job status (polling) |
| GET | `/api/videos/storyboard/:storyboardId` | Get all 3 video variants for a storyboard |
| GET | `/api/videos/:variantId` | Get single video variant details (with all scenes) |
| GET | `/api/videos/:variantId/download` | Download video file |
| POST | `/api/videos/:variantId/regenerate` | Regenerate entire video variant |
| POST | `/api/videos/scene/:sceneId/regenerate` | Regenerate individual scene |

### Quick Video Generation CLI Example

```bash
# 1. Register & Login
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"pass123"}'

# Save token from response: export TOKEN="eyJ..."

# 2. Create Project
curl -X POST http://localhost:3000/api/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Video Project"}' # Save project_id

# 3. Create Business Brief
curl -X POST http://localhost:3000/api/briefs/business \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "project_name": "My Video",
    "company_name": "Acme Corp",
    "industry": "Tech",
    "target_audience": "Developers"
  }' # Save business_brief_id

# 4. Generate Content Pillars
curl -X POST http://localhost:3000/api/projects/{project_id}/content-pillars/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "{project_id}"}' # Get pillar & theme IDs

# 5. Select Content Pillar & Theme
curl -X POST http://localhost:3000/api/content-pillars/{pillar_id}/select \
  -H "Authorization: Bearer $TOKEN"

curl -X POST http://localhost:3000/api/content-themes/{theme_id}/select \
  -H "Authorization: Bearer $TOKEN"

# 6. Generate Storyboards
curl -X POST http://localhost:3000/api/projects/{project_id}/storyboards/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "{project_id}","content_theme_id": "{theme_id}"}' # Get storyboard_id

# 7. Select Storyboard
curl -X POST http://localhost:3000/api/storyboards/{storyboard_id}/select \
  -H "Authorization: Bearer $TOKEN"

# 8. Generate Video (3 variants)
curl -X POST http://localhost:3000/api/videos/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "storyboard_id": "{storyboard_id}",
    "title": "My Generated Video"
  }' # Get generation_job_id

# 9. Poll Generation Status
curl http://localhost:3000/api/videos/generation/{generation_job_id} \
  -H "Authorization: Bearer $TOKEN"

# 10. Get Variants When Completed
curl http://localhost:3000/api/videos/storyboard/{storyboard_id} \
  -H "Authorization: Bearer $TOKEN"

# 11. Download Video
curl http://localhost:3000/api/videos/{variant_id}/download \
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
4. Creates 3 VideoVariant records (one per variant type)
5. Enqueues 3 jobs in JobQueue

    |
    v
Background Workers start processing
    |
    +-- For each VideoVariant:
    |   |
    |   v
    |   1. Fetch SceneGenerations for variant
    |   2. For each scene: call VideoProvider.GenerateScene()
    |   3. Update SceneGeneration status as scenes complete
    |   4. JobQueue polls provider for scene status (60-sec intervals)
    |   5. Max polling duration: 2 hours
    |   |
    |   v
    |
    +-- When all scenes complete:
        1. VideoVariant marked as completed
        2. Download/store video file
        3. GenerationJob marked as completed
        
    |
    v
GET /api/videos/generation/{jobId} returns status: completed
```

### Credit Calculation

**Standard Costs (single scene):**
- LTX-2-Fast: 2 credits/second × duration
- LTX-2-Pro: 3 credits/second × duration
- Runway: Variable (model-dependent)
- Open-source: 1 credit/second × duration

**Example (3 variants of 10-second video):**
- 10 sec × 2 credits/sec (LTX-2-Fast) × 3 variants = 60 credits
- Total deducted: 1 credit (fixed per generation call)

## Quick Copy-Paste Commands

### 1. Health Check (No Auth)
```bash
curl http://localhost:3000/health
```

### 2. Register User
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure123"
  }'
```

### 3. Login (save tokens)
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secure123"
  }'
```

### 4. Get Profile (Protected)
```bash
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer {access_token}"
```

### 5. Create Project
```bash
curl -X POST http://localhost:3000/api/projects \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Brand Video Campaign Q3",
    "description": "A promotional video campaign for Q3 2026"
  }'
```

### 6. List Projects
```bash
curl http://localhost:3000/api/projects \
  -H "Authorization: Bearer {access_token}"
```

### 7. Create Business Brief
```bash
curl -X POST http://localhost:3000/api/briefs/business \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "project_name": "Brand Video Q3",
    "company_name": "Acme Corp",
    "industry": "Technology",
    "target_audience": "Developers aged 25-40",
    "project_objective": "Increase brand awareness by 30%",
    "key_message": "Innovation made simple",
    "budget": "50000000",
    "timeline": "4 weeks"
  }'
```

### 8. Generate Content Pillars
```bash
curl -X POST http://localhost:3000/api/projects/{project_id}/content-pillars/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "{project_id}"}'
```

### 9. List Content Pillars & Select One
```bash
curl http://localhost:3000/api/projects/{project_id}/content-pillars \
  -H "Authorization: Bearer {access_token}"

curl -X POST http://localhost:3000/api/content-pillars/{pillar_id}/select \
  -H "Authorization: Bearer {access_token}"
```

### 10. Select Content Theme
```bash
curl http://localhost:3000/api/content-pillars/{pillar_id}/themes \
  -H "Authorization: Bearer {access_token}"

curl -X POST http://localhost:3000/api/content-themes/{theme_id}/select \
  -H "Authorization: Bearer {access_token}"
```

### 11. Generate Storyboards
```bash
curl -X POST http://localhost:3000/api/projects/{project_id}/storyboards/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "{project_id}", "content_theme_id": "{theme_id}"}'
```

### 12. Select Storyboard & View Scenes
```bash
curl -X POST http://localhost:3000/api/storyboards/{storyboard_id}/select \
  -H "Authorization: Bearer {access_token}"

curl http://localhost:3000/api/storyboards/{storyboard_id}/scenes \
  -H "Authorization: Bearer {access_token}"
```

### 13. Generate Video (costs 1 credit)
```bash
curl -X POST http://localhost:3000/api/videos/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "storyboard_id": "{storyboard_id}",
    "title": "My Brand Video",
    "format": "mp4",
    "resolution": "1080p"
  }'
```

### 14. View & Download Video
```bash
curl http://localhost:3000/api/videos/{video_id} \
  -H "Authorization: Bearer {access_token}"

curl http://localhost:3000/api/videos/{video_id}/download \
  -H "Authorization: Bearer {access_token}"
```

### 15. Check Credits
```bash
curl http://localhost:3000/api/credits \
  -H "Authorization: Bearer {access_token}"
```

### 16. Admin: Add Credits to User
```bash
curl -X POST http://localhost:3000/api/admin/credits \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "{target_user_id}", "amount": 10}'
```

## Usage Examples with Responses

### Register User

```bash
curl -X POST http://localhost:3000/api/auth/register \
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
curl -X POST http://localhost:3000/api/auth/login \
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
curl -X POST http://localhost:3000/api/projects \
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

### Generate Content Pillars

```bash
curl -X POST http://localhost:3000/api/projects/{project_id}/content-pillars/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "{project_id}"}'
```

Response:
```json
{
  "success": true,
  "message": "Content pillars generated successfully",
  "data": [
    {
      "id": "uuid",
      "project_id": "uuid",
      "title": "Educational Content",
      "description": "Informative content that educates...",
      "is_selected": false,
      "content_themes": [
        {"id": "uuid", "title": "How-To Tutorials", "is_selected": false},
        {"id": "uuid", "title": "Industry Insights", "is_selected": false}
      ]
    }
  ]
}
```

### Generate Video

```bash
curl -X POST http://localhost:3000/api/videos/generate \
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

### Users Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(255) DEFAULT 'user',       -- user, admin
  credits INTEGER DEFAULT 10,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Projects Table
```sql
CREATE TABLE projects (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  status VARCHAR(255) DEFAULT 'draft',    -- draft, in_progress, completed
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Business Briefs Table
```sql
CREATE TABLE business_briefs (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  project_id UUID REFERENCES projects(id),
  project_name VARCHAR(255) NOT NULL,
  company_name VARCHAR(255),
  industry VARCHAR(255),
  target_audience VARCHAR(255),
  project_objective TEXT,
  key_message TEXT,
  budget VARCHAR(255),
  timeline VARCHAR(255),
  competitors TEXT,
  additional_notes TEXT,
  status VARCHAR(255) DEFAULT 'draft',
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Creative Briefs Table
```sql
CREATE TABLE creative_briefs (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  business_brief_id UUID NOT NULL REFERENCES business_briefs(id),
  title VARCHAR(255) NOT NULL,
  video_type VARCHAR(255),
  duration INTEGER,
  style VARCHAR(255),
  tone VARCHAR(255),
  script TEXT,
  storyboard TEXT,
  visual_references TEXT,
  music_preference VARCHAR(255),
  call_to_action VARCHAR(255),
  output_format VARCHAR(255),
  resolution VARCHAR(255),
  additional_notes TEXT,
  status VARCHAR(255) DEFAULT 'draft',
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Content Pillars Table
```sql
CREATE TABLE content_pillars (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL REFERENCES projects(id),
  user_id UUID NOT NULL REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  is_selected BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Content Themes Table
```sql
CREATE TABLE content_themes (
  id UUID PRIMARY KEY,
  content_pillar_id UUID NOT NULL REFERENCES content_pillars(id),
  user_id UUID NOT NULL REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  is_selected BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Storyboards Table
```sql
CREATE TABLE storyboards (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL REFERENCES projects(id),
  user_id UUID NOT NULL REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  is_selected BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Scenes Table
```sql
CREATE TABLE scenes (
  id UUID PRIMARY KEY,
  storyboard_id UUID NOT NULL REFERENCES storyboards(id),
  user_id UUID NOT NULL REFERENCES users(id),
  scene_number INTEGER NOT NULL,
  title VARCHAR(255),
  description TEXT,
  visual_description TEXT,
  duration INTEGER,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

### Videos Table
```sql
CREATE TABLE videos (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL REFERENCES projects(id),
  user_id UUID NOT NULL REFERENCES users(id),
  storyboard_id UUID NOT NULL REFERENCES storyboards(id),
  title VARCHAR(255) NOT NULL,
  status VARCHAR(255) DEFAULT 'pending',  -- pending, processing, completed, failed
  video_url TEXT,
  thumbnail_url TEXT,
  duration INTEGER,
  format VARCHAR(255) DEFAULT 'mp4',
  resolution VARCHAR(255) DEFAULT '1080p',
  file_size BIGINT,
  credits_used INTEGER DEFAULT 1,
  error_message TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
```

## Project Architecture

```
Client/Frontend (Web/Mobile)
    |
    v
Go Backend (Port 3000) - Sevima-AI-Content-Creator
    |
    +-- Middleware Layer
    |   ├── JWT Validation (Protected routes)
    |   ├── Role-Based Access Control (user/admin)
    |   ├── CORS, Logger, Panic Recovery
    |
    +-- Handler Layer (8 handlers)
    |   ├── Auth Handler        --> Auth Service --> User Repository --> PostgreSQL
    |   ├── Project Handler     --> Project Service --> Project Repository --> PostgreSQL
    |   ├── Brief Handler       --> Brief Service --> Brief Repository --> PostgreSQL
    |   ├── Content Handler     --> Content Service --> Content Repository --> PostgreSQL
    |   ├── Storyboard Handler  --> Storyboard Service --> Storyboard Repository --> PostgreSQL
    |   ├── Video Handler       --> Video Gen Service --> Job Queue --> Generation Workers
    |   ├── Credit Handler      --> Credit Service --> User Repository --> PostgreSQL
    |   └── AI Handler (Proxy)  --> Python AI Service (Port 8000)
    |
    +-- Service Layer (Business Logic)
    |   ├── Auth Service
    |   ├── Project Service
    |   ├── Brief Service
    |   ├── Content Service
    |   ├── Storyboard Service
    |   ├── Video Generation Service
    |   └── Credit Service
    |
    +-- Repository Layer (Data Access)
    |   ├── User Repository
    |   ├── Project Repository
    |   ├── Brief Repository (Business + Creative)
    |   ├── Content Repository (Pillars + Themes)
    |   ├── Storyboard Repository (Storyboards + Scenes)
    |   ├── Generation Job Repository
    |   ├── Video Variant Repository
    |   ├── Scene Generation Repository
    |   └── Video Repository
    |
    +-- Queue System (Async Processing)
    |   └── Job Queue --> Background Workers --> Video Generation Providers
    |
    +-- AI Provider Layer
    |   ├── LTX-2-Fast Provider
    |   ├── LTX-2-Pro Provider
    |   ├── Runway Provider
    |   └── Open-source Models Provider
    |
    v
PostgreSQL Database (12 tables)
    |
    v
Python AI Service (Port 8000) - FastAPI
    |
    ├── Content Pillar Generation
    ├── Storyboard Generation
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
2. Create a Project (Dashboard)
   POST /api/projects
        │
        ▼
3. Create Business Brief (project context, linked to Project)
   POST /api/briefs/business
        │
        ▼
4. Generate Content Pillars (AI stub generates 3 pillars with themes)
   POST /api/projects/:id/content-pillars/generate
        │
        ▼
5. Browse & Select Content Pillar
   POST /api/content-pillars/:id/select
        │
        ▼
6. Browse & Select Content Theme
   POST /api/content-themes/:id/select
        │
        ▼
7. Generate Storyboards (AI stub generates 2 variations with scenes)
   POST /api/projects/:id/storyboards/generate
        │
        ▼
8. Review Scenes & Select Storyboard
   GET /api/storyboards/:id/scenes
   POST /api/storyboards/:id/select
        │
        ▼
9. Generate Video (deducts 1 credit)
   POST /api/videos/generate
        │
        ▼
10. Preview & Download Video
    GET /api/videos/:id
    GET /api/videos/:id/download
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
| `user` | 10 | All standard endpoints (projects, briefs, content, storyboards, videos, credits) |
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
- [Postman Collection](docs/postman_collection.json) — import into Postman for ready-to-use requests with auto-token-saving scripts

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
- Go backend listening on port 3000
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
