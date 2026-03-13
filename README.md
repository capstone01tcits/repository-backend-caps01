# AI Video Gen - Backend Service

Go REST API backend untuk AI Video Content Creator platform — Capstone project dengan SEVIMA. Meliputi authentication, project management, creative brief, content pillar generation, storyboard generation, video generation, credit system, dan admin management.

## Documentation

**Comprehensive API & Workflow References:**

- **[API Documentation](docs/API_DOCUMENTATION.md)** — Complete API reference dengan semua endpoints, request/response examples, dan sistem video generation
- **[Postman Collection](docs/postman_collection.json)** — Import ke Postman untuk testing (termasuk semua endpoints dengan pre-configured variables)
- **[Complete Workflow](docs/WORKFLOW_COMPLETE_FLOW.md)** — End-to-end flow dari register hingga video download dengan HTTP requests, database changes, dan timeline

## Tech Stack

**Go Backend (Port 3000)**
- Framework: Fiber v2
- Database: PostgreSQL + GORM (AutoMigrate)
- Authentication: JWT (Access + Refresh Token)
- Password Hashing: bcrypt
- Role-based Access Control (user / admin)
- AI Gateway: Proxy ke Python AI Service

**Python AI Service (Port 8000)**
- Framework: FastAPI + Uvicorn
- Purpose: AI processing stubs (content pillars, storyboard, video generation)

## Struktur Folder

```
Sevima-BackEnd Ai Video Gen/
├── cmd/
│   └── main.go                        # Go server entry point
├── config/
│   └── config.go                      # Config loader & DB connection
├── internal/
│   ├── handler/
│   │   ├── ai_handler.go             # AI proxy handler
│   │   ├── auth_handler.go           # Auth handlers
│   │   ├── brief_handler.go          # Brief CRUD handlers
│   │   ├── content_handler.go        # Content Pillar & Theme handlers
│   │   ├── credit_handler.go         # Credit management handlers
│   │   ├── project_handler.go        # Project CRUD handlers
│   │   ├── storyboard_handler.go     # Storyboard & Scene handlers
│   │   └── video_handler.go          # Video generation handlers
│   ├── middleware/
│   │   └── auth.go                   # JWT validation middleware
│   ├── model/
│   │   ├── brief.go                  # BusinessBrief & CreativeBrief models
│   │   ├── content.go                # ContentPillar & ContentTheme models
│   │   ├── project.go                # Project model
│   │   ├── storyboard.go             # Storyboard & Scene models
│   │   ├── user.go                   # User struct & request types
│   │   └── video.go                  # Video model
│   ├── repository/
│   │   ├── brief_repository.go       # Brief database queries
│   │   ├── content_repository.go     # Content pillar/theme queries
│   │   ├── project_repository.go     # Project database queries
│   │   ├── storyboard_repository.go  # Storyboard/scene queries
│   │   ├── user_repository.go        # User database queries
│   │   └── video_repository.go       # Video database queries
│   └── service/
│       ├── auth_service.go           # Auth business logic
│       ├── brief_service.go          # Brief business logic
│       ├── content_service.go        # Content pillar/theme logic
│       ├── credit_service.go         # Credit management logic
│       ├── project_service.go        # Project business logic
│       ├── storyboard_service.go     # Storyboard generation logic
│       └── video_service.go          # Video generation logic
├── pkg/utils/
│   ├── jwt.go                        # JWT token generation
│   └── response.go                   # Standardized response format
├── ai-service/
│   ├── main.py                       # Python AI service entry
│   └── requirements.txt              # Python dependencies
├── docs/
│   ├── API_DOCUMENTATION.md          # Comprehensive API documentation
│   └── postman_collection.json       # Postman v2.1 importable collection
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Database Initialization

Database tables are auto-created on startup via GORM AutoMigrate. Nine tables are managed:

| Table | Model | Description |
|-------|-------|-------------|
| `users` | `model.User` | User accounts with role & credits, soft delete |
| `projects` | `model.Project` | Video projects (container for full workflow) |
| `business_briefs` | `model.BusinessBrief` | Business brief forms (FK → users, projects) |
| `creative_briefs` | `model.CreativeBrief` | Creative brief forms (FK → users, business_briefs) |
| `content_pillars` | `model.ContentPillar` | AI-generated content pillars (FK → projects) |
| `content_themes` | `model.ContentTheme` | Content themes per pillar (FK → content_pillars) |
| `storyboards` | `model.Storyboard` | Storyboard variations (FK → projects) |
| `scenes` | `model.Scene` | Individual scenes in storyboards (FK → storyboards) |
| `videos` | `model.Video` | Generated videos (FK → projects, storyboards) |
| `generation_jobs` | `model.GenerationJob` | Video generation queue jobs |
| `video_variants` | `model.VideoVariant` | 3 video variations per storyboard |
| `scene_generations` | `model.SceneGeneration` | Individual scene generation tracking |

No manual migration needed — just start the server and tables will be created/updated automatically.

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
| POST | `/api/auth/login` | Login and get tokens |
| POST | `/api/auth/refresh` | Refresh access token |
| POST | `/api/auth/restore` | Restore deleted account using refresh token |

### Profile (Protected - requires Bearer token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Get current user profile (includes role & credits) |
| POST | `/api/auth/change-password` | Change user password |
| DELETE | `/api/auth/account` | Delete user account (soft delete) |

### Projects / Dashboard (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects` | Create new project |
| GET | `/api/projects` | List user's projects |
| GET | `/api/projects/:id` | Get project by ID |
| PUT | `/api/projects/:id` | Update project |
| DELETE | `/api/projects/:id` | Delete project (soft delete) |

### Business Brief (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/briefs/business` | Create new business brief (with optional project_id) |
| GET | `/api/briefs/business` | List user's business briefs |
| GET | `/api/briefs/business/:id` | Get business brief by ID |
| PUT | `/api/briefs/business/:id` | Update business brief |
| DELETE | `/api/briefs/business/:id` | Delete business brief (soft delete) |
| GET | `/api/briefs/business/:id/creative` | List creative briefs under a business brief |

### Creative Brief (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/briefs/creative` | Create new creative brief |
| GET | `/api/briefs/creative` | List user's creative briefs |
| GET | `/api/briefs/creative/:id` | Get creative brief by ID |
| PUT | `/api/briefs/creative/:id` | Update creative brief |
| DELETE | `/api/briefs/creative/:id` | Delete creative brief (soft delete) |

### Content Pillars (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/:id/content-pillars/generate` | AI-generate content pillars for project |
| GET | `/api/projects/:id/content-pillars` | List content pillars for project |
| GET | `/api/content-pillars/:id` | Get content pillar by ID |
| POST | `/api/content-pillars/:id/select` | Select a content pillar |
| GET | `/api/content-pillars/:id/themes` | List themes under a content pillar |

### Content Themes (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/content-themes/:id/select` | Select a content theme |

### Storyboards (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/:id/storyboards/generate` | AI-generate storyboard variations for project |
| GET | `/api/projects/:id/storyboards` | List storyboards for project |
| GET | `/api/storyboards/:id` | Get storyboard by ID (with scenes) |
| POST | `/api/storyboards/:id/select` | Select a storyboard |
| GET | `/api/storyboards/:id/scenes` | List scenes in a storyboard |

### Videos (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/videos/generate` | Generate video (deducts 1 credit) |
| GET | `/api/videos` | List user's generated videos |
| GET | `/api/videos/:id` | Get video by ID |
| GET | `/api/videos/:id/download` | Get video download info |
| GET | `/api/projects/:id/videos` | List videos for a project |

### Credits (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/credits` | Get current user's credit balance |

### Admin (Protected - requires admin role)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/credits` | Add credits to a user |

### AI Gateway (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| ANY | `/api/ai/*` | Proxy to Python AI service |

All AI gateway endpoints inject user context headers (`X-User-ID`, `X-User-Email`) automatically.

## Video Generation System

The backend includes a **scene-based video generation system** that creates 3 video variations from a single briefing:

### Features
- **3 Video Variants**: Automatically generates cinematic, vibrant, and professional variations
- **Scene-Based**: Each video composed of 2-3 scenes (4-6 sec each), total 8-12 seconds
- **Multiple Providers**: LTX-2-Fast, LTX-2-Pro, Runway, and open-source models
- **Queue System**: Background workers process videos asynchronously
- **Regeneration**: Regenerate full videos or individual scenes with new prompts
- **Credit System**: Configurable costs per provider and operation type

### Video Generation Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/videos/generate` | Generate 3 video variants from storyboard |
| GET | `/api/videos/generation/:jobId` | Check generation job status |
| GET | `/api/videos/storyboard/:storyboardId` | Get all 3 variants for storyboard |
| GET | `/api/videos/:variantId` | Get single variant with all scenes |
| POST | `/api/videos/:variantId/regenerate` | Regenerate video variant |
| POST | `/api/videos/scene/:sceneId/regenerate` | Regenerate individual scene |
| GET | `/api/videos/:variantId/download` | Get video download URL |

### Quick Video Generation Flow

```bash
# 1. Generate 3 videos (default: LTX-2-Fast)
curl -X POST http://localhost:3000/api/videos/generate \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "storyboard_id": "{storyboard_id}"
  }'
# Returns: generation_job_id

# 2. Poll status (videos being generated)
curl http://localhost:3000/api/videos/generation/{generation_job_id} \
  -H "Authorization: Bearer {token}"

# 3. Get variants when ready (status: completed)
curl http://localhost:3000/api/videos/storyboard/{storyboard_id} \
  -H "Authorization: Bearer {token}"

# 4. Download individual video
curl http://localhost:3000/api/videos/{variant_id}/download \
  -H "Authorization: Bearer {token}"
```

### Architecture

**Component Stack:**
- **Models**: GenerationJob, VideoVariant, SceneGeneration
- **Repositories**: Full CRUD operations for all video models
- **Service Layer**: Business logic for generation and regeneration
- **Queue System**: Background workers with polling mechanism
- **Providers**: Abstract interface supporting LTX, Runway, and open-source models

**Job Lifecycle:**
```
queued → processing → completed/failed
         ↓
    Workers process scenes
    ↓
    Providers generate video
    ↓
    Polling tracks progress (60-sec intervals, 2-hour max)
    ↓
    Database updated with results
```

### Credit Costs

- LTX-2-Fast: 2 credits/second (standard)
- LTX-2-Pro: 3 credits/second (premium)
- Runway: Variable by model
- Open-source: 1 credit/second (internal)

Example: 10-second video with 2 scenes × 3 variants = 120 credits

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
Client/Frontend
    |
    v
Go Backend (Port 3000)
    |
    +-- Auth Service       --> PostgreSQL (Users)
    |
    +-- Project Service    --> PostgreSQL (Projects)
    |
    +-- Brief Service      --> PostgreSQL (BusinessBriefs, CreativeBriefs)
    |
    +-- Content Service    --> PostgreSQL (ContentPillars, ContentThemes)
    |
    +-- Storyboard Service --> PostgreSQL (Storyboards, Scenes)
    |
    +-- Video Service      --> PostgreSQL (Videos) + Credit deduction
    |
    +-- Credit Service     --> PostgreSQL (Users.credits)
    |
    +-- Middleware (JWT + Role Validation)
    |
    +-- AI Gateway (Proxy) --> Python AI Service (Port 8000)
                                   |
                                   v
                            AI Processing Stubs
                            (Content Pillars, Storyboard, Video)
```

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

## Authorization

Protected endpoints require Authorization header:

```
Authorization: Bearer {access_token}
```

Token is obtained from `/api/auth/login` response. Access tokens expire after `JWT_EXPIRE_HOURS` (default 24 hours). Use refresh token to get new access token via `/api/auth/refresh`.

### Roles

| Role | Default Credits | Capabilities |
|------|----------------|--------------|
| `user` | 10 | All standard endpoints, video generation |
| `admin` | 10 | All user capabilities + credit management (`POST /api/admin/credits`) |

## Credit System

- New users start with **10 credits** by default
- Each video generation costs **1 credit**
- Check balance: `GET /api/credits`
- Admin can add credits: `POST /api/admin/credits`
- Video generation is blocked when credits reach 0

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

## Troubleshooting

**Problem: "AI service is unreachable"**
- Ensure Python AI service is running on port 8000
- Check `AI_SERVICE_URL` in .env configuration
- Run: `cd ai-service && python main.py`

**Problem: "Invalid token"**
- Token may have expired, use refresh endpoint
- Ensure Authorization header format is correct: `Bearer {token}`
- Check `JWT_SECRET` matches between service and token

**Problem: "Connection refused"**
- Go backend listening on port 3000
- PostgreSQL running on configured host/port
- Both services have network connectivity
