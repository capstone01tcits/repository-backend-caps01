# API Documentation — AI Video Generation Platform

> **Version:** 2.2.0 (Updated Sprint 4 - April 2026)
> **Base URL:** `http://localhost:5000`  
> **AI Service URL:** `http://localhost:8000`
> **Collection Format:** Bruno API Collection (`docs/API_COLLECTION.json`)
> **Auto-Set Feature:** Variables automatically populate during workflow execution via test scripts

---

## Sprint 4 Updates (April 2026) ✓

**Unified FE-BE Integration Endpoint:**
- ✓ Endpoint: `POST /api/projects/initialize` (formerly `/api/projects/from-fe`)
- ✓ Accepts 12 fields from FE wizard (9 required + 3 optional)
- ✓ Atomically creates Project + Brief data in single call
- ✓ Auto-fills 18+ missing backend fields with sensible defaults
- ✓ Backend port: 5000

**Auto-Fill Mappings (FE → BE):**
- `institution_name` → `company_name`, `institute_name`
- `event_content` → `video_type` (smart enum mapping)
- `selected_theme` → `style` (theme → visual style mapping)
- `tone_of_voice` → `music_preference` (tone → music mapping)
- `video_duration` → integer duration (string parsing)
- Auto-defaults: industry="Education", target_audience="Students", format="mp4", resolution="1080p"

**Workflow Simplified:**
- FE Wizard (4 steps) → Backend initialization → Storyboard generation → Video generation
- Register → Login → Initialize Project (from FE) → Generate Storyboard → Generate Video (5 steps)
- Briefs auto-created during project initialization

**15 Total Endpoints (Simplified Core Workflow):**
- 6 Auth endpoints (register, login, profile, change-password, refresh, delete account)
- 3 Project endpoints (initialize with 12 fields, list, get)
- 1 Storyboard endpoint (generate)
- 4 Video endpoints (generate, get, list, download)
- 1 Health check endpoint

**Plus Support Endpoints:**
- 2 Credit endpoints (get balance, admin add credits)

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Authentication & Security](#authentication--security)
4. [Standard Response Format](#standard-response-format)
5. [Error Handling](#error-handling)
6. [Core API Endpoints](#core-api-endpoints)
7. [Video Generation API](#video-generation-api)

---

## Overview

The AI Video Generation Platform Backend provides RESTful APIs for:

- **User authentication** with role-based access (user / admin)
- **Project management** (unified initialization from FE wizard)
- **AI storyboard generation** with automatic scene generation
- **AI video generation** with provider routing (LTX, Runway)
- **Credit management** (user balance & admin top-up)
- **AI service gateway** (proxy to Python AI service at port 8000)

### Service Stack

| Service | Port | Stack | Purpose |
|---------|------|-------|---------|
| Go Backend | 5000 | Fiber + GORM + PostgreSQL | Auth, Projects, Briefs, Video Generation |
| Python AI Service | 8000 | FastAPI + Uvicorn | AI Processing |

---

## Architecture

```
Frontend (Next.js)
       │
       ▼
┌──────────────────────────────────────────┐
│   Go Backend (Port 5000)                 │
│                                          │
│  ┌──────────────────────────────────┐    │
│  │  API Layer (Handlers)            │    │
│  │  - Auth, Projects, Briefs        │    │
│  │  - Video Generation              │    │
│  └────────────────┬─────────────────┘    │
│                   │                      │
│  ┌────────────────┴─────────────────┐    │
│  │  Business Logic Layer (Services) │    │
│  │  - VideoGenerationService        │    │
│  │  - Credit Management             │    │
│  └────────────────┬─────────────────┘    │
│                   │                      │
│  ┌────────────────┴─────────────────┐    │
│  │  Data Layer (Repositories)       │    │
│  │  - Video Repositories            │    │
│  │  - Job Queue Management          │    │
│  └────────────────┬─────────────────┘    │
│                   │                      │
│  ┌────────────────┴─────────────────┐    │
│  │  PostgreSQL + 12 Tables          │    │
│  │  - Original 9 tables             │    │
│  │  - GenerationJob, VideoVariant   │    │
│  │  - SceneGeneration               │    │
│  └────────────────┬─────────────────┘    │
│                   │                      │
│  ┌────────────────┴───────────────────┐  │
│  │  Background Job Queue              │  │
│  │  - 3 Worker Goroutines             │  │
│  │  - Status Polling (60s interval)   │  │
│  │  - Automatic Retry Logic           │  │
│  └────────────────┬───────────────────┘  │
│                   │                      │
│  ┌────────────────┴───────────────────┐  │
│  │  AI Provider Abstraction           │  │
│  │  - LTX, Runway, Wan2.1             │  │
│  │  - Mock implementations            │  │
│  └────────────────┬───────────────────┘  │
│                   │                      │
│  ┌────────────────┴───────────────────┐  │
│  │  AI Gateway (Proxy)                │──┼──► Python AI Service (Port 8000)
│  └────────────────────────────────────┘  │
└──────────────────────────────────────────┘
```

---

## Authentication & Security

### Authentication Method

All protected endpoints use **JWT Bearer Token** authentication.

### Request Headers

```
Authorization: Bearer {access_token}
Content-Type: application/json
```

### Token Structure

**Access Token:**
- Valid for 24 hours
- Contains: user_id, email, role, exp
- Used for API requests

**Refresh Token:**
- Valid for 7 days
- Used to obtain new access tokens
- Can restore deleted accounts

### Role-Based Access

| Role | Permissions |
|------|-------------|
| `user` | Create/edit own projects, access own briefs and content, generate videos (credit-limited), view own credits |
| `admin` | All user permissions, add credits to users |

---

## Standard Response Format

### Success Response (2xx)

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Error Response (4xx, 5xx)

```json
{
  "success": false,
  "message": "Error description"
}
```

---

## Error Handling

### HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | OK | Request successful |
| 201 | Created | Resource created |
| 400 | Bad Request | Invalid parameters |
| 401 | Unauthorized | Missing/invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error |

---

## Core API Endpoints

### Health Check (Public)

```http
GET /health
```

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

### Authentication APIs

#### Register

```http
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-13T10:00:00Z",
      "updated_at": "2026-03-13T10:00:00Z"
    }
  }
}
```

#### Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-13T10:00:00Z",
      "updated_at": "2026-03-13T10:00:00Z"
    }
  }
}
```

#### Get Profile

```http
GET /api/auth/me
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "credits": 10,
    "created_at": "2026-03-13T10:00:00Z",
    "updated_at": "2026-03-13T10:00:00Z"
  }
}
```

#### Change Password

```http
POST /api/auth/change-password
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "old_password": "secure123",
  "new_password": "newsecure123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

#### Refresh Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token refreshed",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-13T10:00:00Z",
      "updated_at": "2026-03-13T10:00:00Z"
    }
  }
}
```

#### Delete Account

```http
DELETE /api/auth/account
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Account deleted successfully",
  "data": null
}
```

---

### Project APIs

#### Initialize Project (from FE Wizard)

This endpoint accepts 12 fields from the frontend wizard (9 required + 3 optional) and atomically creates:
- **Project** with metadata
- **BusinessBrief** with 10 auto-filled fields
- **CreativeBrief** with 8 auto-filled fields

```http
POST /api/projects/initialize
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "institution_name": "SMA Negeri 1 Jakarta",
  "institution_history": "Sekolah terkemuka dengan program pendidikan berkualitas tinggi",
  "school_level": "Senior High School",
  "offered_degrees": "",
  "event_content": "Penerimaan Mahasiswa Baru",
  "tone_of_voice": "Santai & Ramah",
  "selected_key_message": "Bergabunglah dengan keluarga besar kami",
  "video_duration": "15 detik",
  "prompt": "",
  "selected_theme": "Tur Kampus Sinematik",
  "editable_copywriting": "Halo generasi masa depan! Bergabunglah dengan keluarga besar kami.",
  "editable_hashtags": "#SMANegeri1Jakarta #PenerimaanMahasiswaBaru #Pendidikan"
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| institution_name | string | Yes | Name of institution |
| institution_history | string | Yes | Background/history info |
| school_level | string | Yes | Education level (e.g., "Senior High School") |
| offered_degrees | string | No | Offered degree programs |
| event_content | string | Yes | Event/program being promoted |
| tone_of_voice | string | Yes | Tone style (Santai & Ramah, Profesional & Formal, etc.) |
| selected_key_message | string | Yes | Main message for video |
| video_duration | string | Yes | Duration (e.g., "15 detik", "30 detik") |
| prompt | string | No | Additional custom prompt |
| selected_theme | string | Yes | Visual theme (Tur Kampus Sinematik, etc.) |
| editable_copywriting | string | No | Custom copywriting content |
| editable_hashtags | string | No | Hashtags for promotion |

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Project created successfully with briefs",
  "data": {
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "business_brief_id": "550e8400-e29b-41d4-a716-446655440002",
    "creative_brief_id": "550e8400-e29b-41d4-a716-446655440003",
    "project_name": "Penerimaan Mahasiswa Baru for SMA Negeri 1 Jakarta",
    "theme": "Tur Kampus Sinematik",
    "tone": "Santai & Ramah",
    "duration": 15,
    "school_level": "Senior High School",
    "institution_name": "SMA Negeri 1 Jakarta"
  }
}
```

**Auto-fill Behavior:**
- `project_name` = `event_content` + " for " + `institution_name`
- `company_name` = `institution_name`
- `industry` = "Education" (default)
- `target_audience` = "Students" (default)
- `style` = Auto-mapped from `selected_theme`
- `music_preference` = Auto-mapped from `tone_of_voice`
- `video_type` = Auto-mapped from `event_content`
- `output_format` = "mp4" (default)
- `resolution` = "1080p" (default)
- `status` = "draft" (all resources)

#### List Projects

```http
GET /api/projects
Authorization: Bearer {access_token}
```

#### Get Project

```http
GET /api/projects/{project_id}
Authorization: Bearer {access_token}
```



### Storyboard API

#### Generate Storyboard

```http
POST /api/storyboard/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Storyboard generated successfully",
  "data": {
    "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "ready",
    "scenes": [
      {
        "scene_id": "6ka7b810-9dad-11d1-80b4-00c04fd430c8",
        "scene_number": 1,
        "title": "Opening Hook",
        "narration": "Halo generasi masa depan!",
        "visual_description": "Wide shot with dramatic lighting showing campus landmark",
        "duration": 5
      },
      {
        "scene_id": "7la7b810-9dad-11d1-80b4-00c04fd430c8",
        "scene_number": 2,
        "title": "Campus Life",
        "narration": "Di sini, kami siap membantu mewujudkan impian Anda",
        "visual_description": "Students in activities, modern facilities, interactive learning",
        "duration": 8
      },
      {
        "scene_id": "8ma7b810-9dad-11d1-80b4-00c04fd430c8",
        "scene_number": 3,
        "title": "Call to Action",
        "narration": "Jangan lewatkan kesempatan ini. Raih mimpimu bersama kami!",
        "visual_description": "Campus logo with CTA text and animated graphics",
        "duration": 6
      }
    ],
    "total_duration": 19,
    "created_at": "2026-03-13T10:25:00Z"
  }
}
```

---

### Credit APIs

#### Get Credit Balance

```http
GET /api/credits
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Credits retrieved",
  "data": {
    "credits": 10
  }
}
```

---

### Admin APIs

#### Add Credits to User

```http
POST /api/admin/credits
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 100
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Credits added successfully",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "credits_added": 100,
    "total_credits": 110
  }
}
```

---

## Video Generation API

The video generation system creates videos from storyboards in the background using the job queue system.

### Video Generation Endpoints

#### Generate Video

```http
POST /api/videos/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

**Parameters:**
- `project_id` (required): UUID of project
- `storyboard_id` (required): UUID of storyboard

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Video generation job created",
  "data": {
    "job_id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "created_at": "2026-03-13T10:30:00Z"
  }
}
```

---

#### Get Video by ID

```http
GET /api/videos/{id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Video retrieved",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "completed",
    "video_url": "https://storage.example.com/videos/8da7b810.mp4",
    "thumbnail_url": "https://storage.example.com/thumbnails/8da7b810.jpg",
    "duration": 10,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "created_at": "2026-03-13T10:30:00Z",
    "updated_at": "2026-03-13T10:35:45Z"
  }
}
```

---

#### List User Videos

```http
GET /api/videos
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Videos retrieved",
  "data": [
    {
      "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
      "status": "completed",
      "video_url": "https://storage.example.com/videos/8da7b810.mp4",
      "duration": 10,
      "created_at": "2026-03-13T10:30:00Z"
    }
  ]
}
```

---

#### Download Video

```http
GET /api/videos/download/{id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Video download ready",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "download_url": "https://storage.example.com/videos/8da7b810.mp4",
    "file_size": 51234567,
    "format": "mp4",
    "resolution": "1920x1080"
  }
}
```

**Regenerate Limit Exceeded (400):**
```json
{
  "success": false,
  "error": "REGENERATE_LIMIT_EXCEEDED",
  "message": "Video has reached maximum regenerate attempts (3/3)",
  "data": {
    "regenerate_count": 3,
    "max_regenerate": 3
  }
}
```

---

#### Regenerate Single Scene (Sprint 3 ✓ NEW)

✓ **NEW in Sprint 3** – Scene regenerate with independent `regenerate_count` tracking (max 3 per scene)

```http
POST /api/videos/scene/{scene_id}/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "new_prompt": "Show aerial campus views with dynamic panning",
  "update_caption": true,
  "provider": "LTX-2-Fast"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Scene regeneration job created",
  "data": {
    "generation_job_id": "4ga7b810-9dad-11d1-80b4-00c04fd430c8",
    "scene_id": "{scene_id}",
    "status": "queued",
    "regenerate_count": 1,
    "max_regenerate": 3,
    "remaining_attempts": 2,
    "credit_cost": 15,
    "current_balance": 105,
    "estimated_time": "30 seconds"
  }
}
```

**Scene Regenerate Limit Exceeded (400):**
```json
{
  "success": false,
  "error": "REGENERATE_LIMIT_EXCEEDED",
  "message": "Scene has reached maximum regenerate attempts (3/3)",
  "data": {
    "scene_id": "{scene_id}",
    "regenerate_count": 3,
    "max_regenerate": 3
  }
}
```

**Credit Cost for Regenerate:**
- Full video regenerate: 40 credits (same as initial generation)
- Scene regenerate: 15 credits per scene

---

#### Download Video Variant

```http
GET /api/videos/{variant_id}/download
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Video download ready",
  "data": {
    "variant_id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "download_url": "https://storage.example.com/videos/8da7b810.mp4",
    "file_size": 52428800,
    "format": "mp4",
    "resolution": "1080p"
  }
}
```

**Requirements:**
- Video status must be "completed"

---

## Video Generation System Details

### Architecture Components

**1. Models**
- **GenerationJob**: Tracks video generation tasks in the queue
- **VideoVariant**: Represents one of 3 video variations
- **SceneGeneration**: Individual scene generation tracking

**2. Repositories**
- Full CRUD operations for all video models
- Database indexing for performance optimization

**3. Service Layer**
- VideoGenerationService: Business logic for generation and regeneration
- Credit validation and deduction
- Scene planning and prompt variation

**4. Queue System**
- Background workers (configurable count, default 3)
- Job polling mechanism (60-second intervals)
- Automatic retry logic (up to 3 attempts)
- Priority-based job dequeuing

**5. AI Providers**
- LTX-2-Fast (Standard): 4 credits/sec
- LTX-2-Pro (Premium): 6 credits/sec
- Runway Gen4.5: 12 credits/sec
- Runway Turbo: 5 credits/sec
- Wan2.1 (Open Source): 1 credit/sec
- LTX Open Source: 1 credit/sec

### Job Lifecycle

```
User Request
    |
Validate credits
    |
Create GenerationJob (status: queued)
Create 3 VideoVariants
Create 2-3 SceneGenerations per variant
Deduct credits
    |
Enqueue to JobQueue
    |
Worker picks up job (status: processing)
    |
For each scene:
  - Call VideoProvider.GenerateScene()
  - Store external job ID
    |
Start polling routine (every 60 seconds)
    |
Provider returns scene status
    |
When all scenes complete:
  - Update variant status
  - Update video URLs
    |
When all variants complete:
  - Update job status (completed)
  - Videos ready for download
```

### Credit System

Credit hanya dikenakan pada operasi video generation. Generate content pillar dan storyboard tidak mengurangi credit.

**Provider Costs (per second):**
- LTX-2-Fast: 4 credits/sec
- LTX-2-Pro: 6 credits/sec
- Runway Gen4.5: 12 credits/sec
- Runway Turbo: 5 credits/sec
- Open Source: 1 credit/sec

### Polling Recommendations

**Frontend Polling Strategy:**
1. Poll every 2-5 seconds while `status: processing`
2. After 30 seconds, reduce to 10-second intervals
3. After 1 minute, use 60-second intervals
4. Stop polling after 2 hours (timeout)

---

## Database Schema

### Tables (12 total)

| Table | Purpose | Key Fields |
|-------|---------|-----------|
| users | User accounts | id, email, password, role, credits |
| projects | Video projects | id, user_id, name, status |
| business_briefs | Business input | id, project_id, company_name, target_audience |
| creative_briefs | Creative direction | id, business_brief_id, tone, style |
| content_pillars | AI-generated pillars | id, project_id, title, is_selected |
| content_themes | Pillar themes | id, content_pillar_id, title, is_selected |
| storyboards | Storyboard variations | id, project_id, title, is_selected |
| scenes | Individual scenes | id, storyboard_id, scene_number, duration |
| videos | Generated videos (legacy) | id, storyboard_id, video_url |
| generation_jobs | Video job queue | id, status, provider, credits_required |
| video_variants | 3 video variations | id, variant_number, video_url, status |
| scene_generations | Individual scene tracking | id, variant_id, scene_number, video_url |

---

## Setup & Configuration

### Environment Variables

```bash
# Server
APP_PORT=3000
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=go_auth

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRE_HOURS=24
JWT_REFRESH_SECRET=your_refresh_secret
JWT_REFRESH_EXPIRE_HOURS=168

# AI Service
AI_SERVICE_URL=http://localhost:8000
```

### Running the Services

**Docker (Recommended):**
```bash
docker-compose up -d
```

**Manual:**
```bash
# Go backend
go mod tidy
go run cmd/main.go

# Python AI service (separate terminal)
cd ai-service
pip install -r requirements.txt
python main.py
```

### Testing the API

See [Postman Collection](./postman_collection.json) for all endpoints with pre-configured requests.

Quick test:
```bash
# Health check
curl http://localhost:3000/health

# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"123456"}'

# Login and get token
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# Generate videos (requires valid token)
curl -X POST http://localhost:3000/api/videos/generate \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id":"{project_id}",
    "storyboard_id":"{storyboard_id}"
  }'
```

---

## Additional Resources

- **Postman Collection**: [postman_collection.json](./postman_collection.json)
- **Main README**: [../README.md](../README.md)
- **Complete Workflow**: [WORKFLOW_COMPLETE_FLOW.md](./WORKFLOW_COMPLETE_FLOW.md)
