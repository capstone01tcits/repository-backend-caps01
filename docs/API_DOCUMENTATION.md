# API Documentation - AI Video Generation Platform

Version: 4.1.0 (May 2026 - Veo 3 Integration & Linear Automated Flow)
Base URL: http://localhost:5000
AI Service URL: http://localhost:8000
Collection Format: Bruno/Postman API Collection (docs/API_collection.json)

---

## Backend Status - Template Generator Ready

LATEST UPDATE: May 2026
- Veo 3 Integration: Real HTTP requests with specialized `Veo3Payload`, automated FFmpeg stitching, or direct aggregator integration via **Wavespeed**.
- Linear Automated Flow: Auto-generates a 3-scene Storyboard (Hook, Value, CTA) directly during Project Initialization.
- Soft-delete and restore for projects and storyboards
- Build Status: Compiles successfully with zero errors

Core Workflow (4 Steps):
1. User Registration/Login
2. Project Initialization (Creates Briefs & Auto-Generates Storyboard)
3. Video Generation (Sends payload to Veo 3)
4. Video Retrieval and Download (polling `generating_assets` and `stitching_video` statuses)

Total Active Endpoints: 27
- Authentication: 6 endpoints
- Projects: 5 endpoints
- Storyboard: 6 endpoints (create, get by project, get detail, update, delete, restore)
- Videos: 6 endpoints
- Credits: 1 endpoint
- Admin: 1 endpoint
- Health/AI Gateway: 2 endpoints

**Database Tables (10 via AutoMigrate):**
users, projects, business_briefs, creative_briefs, storyboards, storyboard_sections, videos, generation_jobs, video_variants, scene_generations

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
- **Project management** (unified initialization from FE wizard, soft-delete/restore)
- **Storyboard management** (manual create, CRUD with sections, linked 1:1 to Project)
- **AI video generation** with Veo 3 / Wavespeed integration
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
│  │  PostgreSQL + 10 Tables          │    │
│  │  - User, Project, Briefs (2)     │    │
│  │  - Storyboard, StoryboardSection │    │
│  │  - Video, GenerationJob          │    │
│  │  - VideoVariant, SceneGeneration │    │
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
│  │  AI Provider Implementation        │  │
│  │  - Veo 3 (Internal / Simulation)   │  │
│  │  - Wavespeed (External Aggregator) │  │
│  │  - FFmpeg Video Stitching          │  │
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
| `user` | Create/edit own projects, access own briefs and content, generate videos (10 credits/sec), view own credits |
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
- **BusinessBrief** with auto-filled fields
- **CreativeBrief** with auto-filled fields
- **Storyboard** with 3 auto-generated scenes (Hook, Value, CTA) in draft status

```http
POST /api/projects/initialize
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "institution_name": "SMA Negeri 1 Jakarta",
  "institution_history": "Sekolah terkemuka dengan program pendidikan berkualitas tinggi",
  "school_level": "SMA",
  "offered_degrees": "IPA, IPS, Bahasa",
  "event_content": "Penerimaan Mahasiswa Baru",
  "tone_of_voice": "Santai & Ramah",
  "selected_key_message": "Bergabunglah dengan keluarga besar kami",
  "video_duration": "15 detik",
  "prompt": "Buat video yang eye-catching dan engaging untuk Gen Z",
  "selected_theme": "Tur Kampus Sinematik",
  "editable_copywriting": "Halo generasi masa depan! Bergabunglah dengan keluarga besar kami.",
  "editable_hashtags": "#SMANegeri1Jakarta #PenerimaanMahasiswaBaru #Pendidikan",
  "logo_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "env_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFAABAAAAAAAAAAAAAAAAAAAAA//EAAUEAQEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCwAA8A/9k=",
  "document_base64": ""
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| institution_name | string | Yes | Name of institution |
| institution_history | string | No | Background/history info (auto-fills if empty) |
| school_level | string | No | Education level (PreSchool, TK, SD, SMP, SMA, SMK, Perguruan Tinggi). Defaults to "Perguruan Tinggi" if not provided |
| offered_degrees | string | No | Offered degree programs |
| event_content | string | Yes | Event/program being promoted |
| tone_of_voice | string | Yes | Tone style (Santai & Ramah, Profesional & Formal, Kreatif & Inovatif, Berwibawa & Meyakinkan) |
| selected_key_message | string | Yes | Main message for video |
| video_duration | string | No | Duration (e.g., "15 detik", "30 detik", "60 detik"). Defaults to 30 seconds if not provided |
| prompt | string | No | Additional custom prompt |
| selected_theme | string | Yes | Visual theme (Tur Kampus Sinematik, Cerita Kehidupan Mahasiswa, Keunggulan Akademik, Tren & Gaya Hidup Cepat) |
| editable_copywriting | string | No | Custom copywriting content for social media caption |
| editable_hashtags | string | No | Hashtags for social media promotion |
| logo_base64 | string | No | Institution logo as base64 encoded image (PNG/JPEG) |
| env_base64 | string | No | Environment photo as base64 encoded image (PNG/JPEG) |
| document_base64 | string | No | Optional PDF/document about institution as base64 encoded string |

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
    "school_level": "SMA",
    "institution_name": "SMA Negeri 1 Jakarta",
    "event_content": "Penerimaan Mahasiswa Baru",
    "key_message": "Bergabunglah dengan keluarga besar kami",
    "copywriting": "Halo generasi masa depan! Bergabunglah dengan keluarga besar kami.",
    "hashtags": "#SMANegeri1Jakarta #PenerimaanMahasiswaBaru #Pendidikan",
    "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
  }
}
```

**Auto-fill & Optional Field Behavior:**
- Fields marked as optional will not cause validation errors if missing
- `institution_history` - Auto-generates if empty: "Video production project for {institution_name}"
- `school_level` - Defaults to "Perguruan Tinggi" if not provided
- `video_duration` - Defaults to 30 seconds if not provided or invalid
- `logo_base64`, `env_base64`, `document_base64` - Stored directly in database for future file storage implementation
- `editable_hashtags` - Stored in creative brief for social media use
- Auto-mapped fields:
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

#### Soft Delete Project

```http
DELETE /api/projects/{project_id}
Authorization: Bearer {access_token}
```

#### Restore Project

```http
POST /api/projects/{project_id}/restore
Authorization: Bearer {access_token}
```



### Storyboard API

#### Get Storyboard by Project

```http
GET /api/storyboard/{project_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Storyboard retrieved successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Dynamic Campaign",
    "total_duration": 30,
    "style": "Dynamic",
    "sections": [
      { "id": "...", "section_type": "hook", "content": "...", "duration": 10 },
      { "id": "...", "section_type": "value", "content": "...", "duration": 10 },
      { "id": "...", "section_type": "cta", "content": "...", "duration": 10 }
    ]
  }
}
```

#### Get Storyboard Detail

```http
GET /api/storyboard/detail/{storyboard_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Storyboard retrieved successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "title": "Dynamic Campaign",
    "total_duration": 30,
    "style": "Dynamic",
    "is_selected": true,
    "sections": [
      { "id": "...", "section_type": "hook", "content": "...", "duration": 10 },
      { "id": "...", "section_type": "value", "content": "...", "duration": 10 },
      { "id": "...", "section_type": "cta", "content": "...", "duration": 10 }
    ]
  }
}
```


#### Update Storyboard

```http
PUT /api/storyboard/{storyboard_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description",
  "sections": [
    { "section_type": "hook", "content": "New hook", "duration": 10 },
    { "section_type": "value", "content": "New value", "duration": 10 },
    { "section_type": "cta", "content": "New CTA", "duration": 10 }
  ]
}
```

#### Get Storyboard Sections

```http
GET /api/storyboard/{storyboard_id}/sections
Authorization: Bearer {access_token}
```

---

### Storyboard Data Flow

Shows how storyboard data flows from project initialization through video generation.

**Flow Diagram:**

```
Step 1: Initialize Project (Auto-Generates Storyboard)
(POST /api/projects/initialize)
        ↓
Returns: project_id, storyboard_id
        ↓
Step 2: Generate Video FROM Storyboard
(POST /api/videos/generate)
Input: project_id + storyboard_id (from Step 1)
        ↓
Returns: job_id, status: "queued"
        ↓
Step 3: Poll Video Status
(GET /api/videos/{id})
        ↓
Returns: status (generating_assets → stitching_video → completed)
        ↓
Step 4: Download Generated Video
(GET /api/videos/download/{id})
```

**Data Mapping:**

| Source Endpoint | Data Used | Target Endpoint |
|-----------------|-----------|-----------------|
| POST /api/projects/initialize | Returns `project_id` | Used in storyboard endpoints |
| POST /api/storyboard/templates/generate | Returns 4 templates with `sections[]` content | User reviews and creates custom storyboard or uses as template |
| POST /api/storyboard/create | Returns `storyboard_id` and `sections[].id` | Used in POST /api/videos/generate |
| POST /api/videos/generate | Returns `job_id` | Used in GET /api/videos/{id} for polling |

**Critical Fields for Workflow:**

1. **From Project Initialization:**
   - `project_id` - Required for all subsequent endpoints

2. **From Template Generation (if using Option A):**
   - `sections[].content` - AI-suggested content for each section
   - `sections[].suggested_duration` - Recommended timing
   - `style` - Selected template style

3. **From Manual Storyboard (Option B):**
   - `id` - Storyboard ID for video generation
   - `sections[].id` - Individual section IDs
   - `sections[].content` - Your custom content
   - `total_duration` - Total video length

4. **For Video Generation:**
   - `project_id` + `storyboard_id` - The essential pair for creating videos

**Complete Request-Response Chain:**

```json
Request 1: POST /api/projects/initialize
Response 1 → project_id: "abc-123", storyboard_id: "xyz-789"

Request 2: POST /api/videos/generate
Body: { project_id: "abc-123", storyboard_id: "xyz-789" }
Response 2 → job_id: "job-456", status: "queued"

Request 3: GET /api/videos/{job_id}
Response 3 → status: "completed", video_url: "..."
```

---

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
    "generation_job_id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
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

---

#### Regenerate Video Variant

```http
POST /api/videos/{variantId}/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "new_prompt": "Alternative video with different visual style"
}
```

**Parameters:**
- `variantId` (required): UUID of video variant to regenerate
- `new_prompt` (optional): Custom prompt for regeneration

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Video regeneration job created",
  "data": {
    "generation_job_id": "3ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "created_at": "2026-03-13T10:40:00Z"
  }
}
```

**Insufficient Credits (400):**
```json
{
  "success": false,
  "message": "insufficient credits for regeneration",
  "data": null
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

### Tables (10 total via AutoMigrate)

| Table | Purpose | Key Fields |
|-------|---------|-----------|
| users | User accounts | id, email, password, role, credits |
| projects | Video projects | id, user_id, name, status |
| business_briefs | Business input | id, project_id, institute_name, school_level, target_audience |
| creative_briefs | Creative direction | id, business_brief_id, tone, style, copywriting, hashtags |
| storyboards | Storyboard blueprints | id, project_id, title, style, total_duration, is_selected |
| storyboard_sections | 3-part sections (hook/value/cta) | id, storyboard_id, section_type, content, duration |
| videos | Generated videos | id, storyboard_id, video_url |
| generation_jobs | Video job queue | id, status, provider, credits_required |
| video_variants | 3 video variations | id, variant_number, video_url, status |
| scene_generations | Individual scene tracking | id, variant_id, scene_number, video_url |

---

## Setup & Configuration

### Environment Variables

```bash
# Server
APP_PORT=5000
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
curl http://localhost:5000/health

# Register
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"123456"}'

# Login and get token
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# Generate videos (requires valid token)
curl -X POST http://localhost:5000/api/videos/generate \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id":"{project_id}",
    "storyboard_id":"{storyboard_id}"
  }'
```

---

## Additional Resources

- **API Collection**: [API_collection.json](./API_collection.json)
- **Main README**: [../README.md](../README.md)
- **Complete Workflow**: [WORKFLOW_COMPLETE_FLOW.md](./WORKFLOW_COMPLETE_FLOW.md)
