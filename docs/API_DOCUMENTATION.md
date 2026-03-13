# API Documentation — AI Video Generation Platform Backend Services

> **Version:** 2.0.0  
> **Base URL:** `http://localhost:3000`  
> **AI Service URL:** `http://localhost:8000`

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Authentication & Security](#authentication--security)
4. [Standard Response Format](#standard-response-format)
5. [Error Handling](#error-handling)
6. [API Endpoints](#api-endpoints)
   - [Health Check](#health-check)
   - [Authentication APIs](#authentication-apis)
   - [Project APIs](#project-apis)
   - [Business Brief APIs](#business-brief-apis)
   - [Creative Brief APIs](#creative-brief-apis)
   - [Content Pillar APIs](#content-pillar-apis)
   - [Content Theme APIs](#content-theme-apis)
   - [Storyboard APIs](#storyboard-apis)
   - [Video APIs](#video-apis)
   - [Credit APIs](#credit-apis)
   - [Admin APIs](#admin-apis)
   - [AI Gateway APIs](#ai-gateway-apis)
7. [Video Generation Workflow](#video-generation-workflow)
8. [Database Schema](#database-schema)
9. [Interactive Documentation](#interactive-documentation)

---

## Overview

The AI Video Generation Platform Backend provides RESTful APIs for:

- **User authentication** with role-based access (user / admin)
- **Project management** (dashboard for video projects)
- **Brief management** (business briefs & creative briefs)
- **AI content pillar generation** with selection workflow
- **AI storyboard generation** with scene timeline
- **AI video generation** with credit system
- **Credit management** (user balance & admin top-up)
- **AI service gateway** (proxy to Python AI service)

The backend is composed of two services:

| Service | Port | Stack | Purpose |
|---------|------|-------|---------|
| Go Backend | 3000 | Fiber + GORM + PostgreSQL | Auth, Projects, Briefs, Content, Storyboards, Videos, Credits |
| Python AI Service | 8000 | FastAPI + Uvicorn | AI Processing Stubs |

---

## Architecture

```
Client / Frontend
       │
       ▼
┌──────────────────────────────────────────┐
│   Go Backend (Port 3000)                 │
│                                          │
│  ┌──────────┐  ┌──────────────────────┐  │
│  │ Auth API │  │  Project/Brief API   │  │
│  └────┬─────┘  └──────┬───────────────┘  │
│       │               │                  │
│  ┌────┴───────────────┴───────────────┐  │
│  │  Content Pillar / Storyboard API   │  │
│  └────────────────┬───────────────────┘  │
│                   │                      │
│  ┌────────────────┴───────────────────┐  │
│  │  Video Generation + Credit API     │  │
│  └────────────────┬───────────────────┘  │
│                   │                      │
│  ┌────────────────┴───────────────────┐  │
│  │    JWT Middleware + Role Check     │  │
│  └────────────────┬───────────────────┘  │
│                   │                      │
│  ┌────────────────┴───────────────────┐  │
│  │    PostgreSQL (GORM) — 9 Tables    │  │
│  └────────────────────────────────────┘  │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  AI Gateway (Proxy)                │──┼──► Python AI Service (Port 8000)
│  └────────────────────────────────────┘  │
└──────────────────────────────────────────┘
```

---

## Authentication & Security

### Authentication Method

All protected endpoints use **JWT Bearer Token** authentication.

### Headers

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

### Token Flow

1. **Register** or **Login** → Receive `access_token` + `refresh_token`
2. Use `access_token` in `Authorization` header for protected endpoints
3. When `access_token` expires → call `/api/auth/refresh` with `refresh_token`
4. `access_token` expiry: configurable (default 24 hours)
5. `refresh_token` expiry: configurable (default 168 hours / 7 days)

### Roles

| Role | Default Credits | Capabilities |
|------|----------------|--------------|
| `user` | 10 | All standard endpoints |
| `admin` | 10 | All user capabilities + credit management |

### Token Validation Errors

| Status | Message | Cause |
|--------|---------|-------|
| 401 | Authorization header required | Missing `Authorization` header |
| 401 | Invalid authorization format | Not using `Bearer <token>` format |
| 401 | Invalid or expired token | Token expired or tampered |

---

## Standard Response Format

### Success Response

```json
{
  "success": true,
  "message": "Descriptive success message",
  "data": { }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Descriptive error message"
}
```

---

## Error Handling

| HTTP Status | Meaning | Example |
|-------------|---------|---------|
| 400 | Bad Request | Invalid request body, missing required fields |
| 401 | Unauthorized | Invalid/expired token, wrong credentials, insufficient role |
| 403 | Forbidden | Not owner of resource |
| 404 | Not Found | Resource does not exist |
| 500 | Internal Server Error | Server-side processing failure |

---

## API Endpoints

### Health Check

#### `GET /health`

Check Go backend status.

- **Auth Required:** No

**Response:**
```json
{ "status": "ok" }
```

---

#### `GET /api/ai/health`

Check Python AI service connectivity.

- **Auth Required:** No

**Response (success):**
```json
{
  "success": true,
  "status": "ok",
  "message": "AI service is running"
}
```

---

### Authentication APIs

#### `POST /api/auth/register`

Register a new user account. New users get `role: "user"` and `credits: 10` by default.

- **Auth Required:** No

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | User's full name (min 2 chars) |
| email | string | Yes | User's email address |
| password | string | Yes | Password (min 6 chars) |

**Example Request:**
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure123"
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 1741564800,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-10T10:00:00Z",
      "updated_at": "2026-03-10T10:00:00Z"
    }
  }
}
```

---

#### `POST /api/auth/login`

Authenticate user and receive tokens.

- **Auth Required:** No

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | Registered email |
| password | string | Yes | Account password |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 1741564800,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-10T10:00:00Z",
      "updated_at": "2026-03-10T10:00:00Z"
    }
  }
}
```

---

#### `POST /api/auth/refresh`

Refresh an expired access token.

- **Auth Required:** No

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| refresh_token | string | Yes | Valid refresh token |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token refreshed",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 1741564800,
    "user": { "...user object..." }
  }
}
```

---

#### `GET /api/auth/me`

Get the currently authenticated user's profile.

- **Auth Required:** Yes

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
    "created_at": "2026-03-10T10:00:00Z",
    "updated_at": "2026-03-10T10:00:00Z"
  }
}
```

---

#### `POST /api/auth/change-password`

Change the authenticated user's password.

- **Auth Required:** Yes

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| old_password | string | Yes | Current password (min 6 chars) |
| new_password | string | Yes | New password (min 6 chars) |

**Response (200 OK):**
```json
{ "success": true, "message": "Password changed successfully" }
```

---

#### `DELETE /api/auth/account`

Soft-delete the authenticated user's account.

- **Auth Required:** Yes

**Response (200 OK):**
```json
{ "success": true, "message": "Account deleted successfully" }
```

---

#### `POST /api/auth/restore`

Restore a soft-deleted account using a refresh token.

- **Auth Required:** No (uses refresh_token)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| refresh_token | string | Yes | Refresh token from before deletion |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Account restored successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "credits": 10,
    "created_at": "2026-03-10T10:00:00Z",
    "updated_at": "2026-03-10T10:00:00Z"
  }
}
```

---

### Project APIs

All Project endpoints require authentication. Projects serve as the top-level container for the complete video generation workflow.

#### `POST /api/projects`

Create a new project.

- **Auth Required:** Yes

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Project name |
| description | string | No | Project description |

**Example Request:**
```bash
curl -X POST http://localhost:3000/api/projects \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Brand Video Campaign Q3",
    "description": "A promotional video campaign for Q3 2026"
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Project created successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Brand Video Campaign Q3",
    "description": "A promotional video campaign for Q3 2026",
    "status": "draft",
    "created_at": "2026-03-10T10:00:00Z",
    "updated_at": "2026-03-10T10:00:00Z"
  }
}
```

---

#### `GET /api/projects`

List all projects for the authenticated user.

- **Auth Required:** Yes

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Projects retrieved successfully",
  "data": [ { "...project objects..." } ]
}
```

---

#### `GET /api/projects/:id`

Get a single project by ID.

- **Auth Required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Project ID |

---

#### `PUT /api/projects/:id`

Update an existing project. Only the project owner can update.

- **Auth Required:** Yes

**Request Body:** (all fields optional)

| Field | Type | Description |
|-------|------|-------------|
| name | string | Project name |
| description | string | Project description |
| status | string | draft / in_progress / completed |

---

#### `DELETE /api/projects/:id`

Soft-delete a project. Only the project owner can delete.

- **Auth Required:** Yes

---

### Business Brief APIs

All Business Brief endpoints require authentication.

#### `POST /api/briefs/business`

Create a new business brief, optionally linked to a project.

- **Auth Required:** Yes

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| project_id | string (UUID) | No | ID of the parent project |
| project_name | string | **Yes** | Name of the project |
| company_name | string | No | Company or brand name |
| industry | string | No | Industry sector |
| target_audience | string | No | Target audience description |
| project_objective | string | No | Objective / goals of the project |
| key_message | string | No | Core message to convey |
| budget | string | No | Budget range or amount |
| timeline | string | No | Expected timeline |
| competitors | string | No | Competitor references |
| additional_notes | string | No | Any additional information |

---

#### `GET /api/briefs/business`

List all business briefs for the authenticated user.

- **Auth Required:** Yes

---

#### `GET /api/briefs/business/:id`

Get a single business brief by ID.

- **Auth Required:** Yes

---

#### `PUT /api/briefs/business/:id`

Update an existing business brief.

- **Auth Required:** Yes

**Request Body:** (all fields optional — only send fields to update)

| Field | Type | Description |
|-------|------|-------------|
| project_name | string | Name of the project |
| company_name | string | Company or brand name |
| industry | string | Industry sector |
| target_audience | string | Target audience |
| project_objective | string | Objective |
| key_message | string | Core message |
| budget | string | Budget |
| timeline | string | Timeline |
| competitors | string | Competitors |
| additional_notes | string | Additional info |
| status | string | draft / submitted / approved / rejected |

---

#### `DELETE /api/briefs/business/:id`

Soft-delete a business brief.

- **Auth Required:** Yes

---

#### `GET /api/briefs/business/:id/creative`

Get all creative briefs linked to a specific business brief.

- **Auth Required:** Yes

---

### Creative Brief APIs

All Creative Brief endpoints require authentication.

#### `POST /api/briefs/creative`

Create a new creative brief (must be linked to an existing business brief).

- **Auth Required:** Yes

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| business_brief_id | string (UUID) | **Yes** | ID of the parent business brief |
| title | string | **Yes** | Creative brief title |
| video_type | string | No | promotional / educational / testimonial / explainer |
| duration | integer | No | Video duration in seconds |
| style | string | No | cinematic / animated / minimalist |
| tone | string | No | professional / casual / energetic |
| script | string | No | Video script text |
| storyboard | string | No | Storyboard description |
| visual_references | string | No | Visual reference URLs or descriptions |
| music_preference | string | No | Music style preference |
| call_to_action | string | No | CTA text |
| output_format | string | No | mp4 / webm |
| resolution | string | No | 1080p / 4K |
| additional_notes | string | No | Additional notes |

---

#### `GET /api/briefs/creative`

List all creative briefs for the authenticated user.

- **Auth Required:** Yes

---

#### `GET /api/briefs/creative/:id`

Get a single creative brief by ID.

- **Auth Required:** Yes

---

#### `PUT /api/briefs/creative/:id`

Update an existing creative brief.

- **Auth Required:** Yes

---

#### `DELETE /api/briefs/creative/:id`

Soft-delete a creative brief.

- **Auth Required:** Yes

---

### Content Pillar APIs

Content pillars are AI-generated topic categories for a project's video content.

#### `POST /api/projects/:id/content-pillars/generate`

Generate content pillars for a project using AI. Currently returns stub data: 3 pillars with 2 themes each.

- **Auth Required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Project ID |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| project_id | string (UUID) | Yes | Project ID |

**Example Request:**
```bash
curl -X POST http://localhost:3000/api/projects/{project_id}/content-pillars/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "{project_id}"}'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Content pillars generated successfully",
  "data": [
    {
      "id": "uuid",
      "project_id": "uuid",
      "user_id": "uuid",
      "title": "Educational Content",
      "description": "Informative content that educates your audience about your product or industry",
      "is_selected": false,
      "content_themes": [
        {
          "id": "uuid",
          "content_pillar_id": "uuid",
          "title": "How-To Tutorials",
          "description": "Step-by-step guides showing how to use your product",
          "is_selected": false
        },
        {
          "id": "uuid",
          "content_pillar_id": "uuid",
          "title": "Industry Insights",
          "description": "Share knowledge about trends and developments",
          "is_selected": false
        }
      ]
    },
    { "...2 more pillars..." }
  ]
}
```

---

#### `GET /api/projects/:id/content-pillars`

List all content pillars for a project (with their themes).

- **Auth Required:** Yes

---

#### `GET /api/content-pillars/:id`

Get a single content pillar by ID (with its themes).

- **Auth Required:** Yes

---

#### `POST /api/content-pillars/:id/select`

Select a content pillar. Deselects all other pillars in the same project first.

- **Auth Required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Content pillar ID |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Content pillar selected",
  "data": { "...updated content pillar..." }
}
```

---

#### `GET /api/content-pillars/:id/themes`

Get all content themes under a specific content pillar.

- **Auth Required:** Yes

---

### Content Theme APIs

#### `POST /api/content-themes/:id/select`

Select a content theme. Deselects all other themes in the same pillar first.

- **Auth Required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Content theme ID |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Content theme selected",
  "data": { "...updated content theme..." }
}
```

---

### Storyboard APIs

Storyboards represent visual narrative structures for video content, each containing multiple scenes.

#### `POST /api/projects/:id/storyboards/generate`

Generate storyboard variations for a project using AI. Currently returns stub data: 2 storyboard variations with 4 scenes each.

- **Auth Required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Project ID |

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| project_id | string (UUID) | Yes | Project ID |
| content_theme_id | string (UUID) | Yes | Selected content theme ID |

**Example Request:**
```bash
curl -X POST http://localhost:3000/api/projects/{project_id}/storyboards/generate \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "{project_id}",
    "content_theme_id": "{theme_id}"
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Storyboards generated successfully",
  "data": [
    {
      "id": "uuid",
      "project_id": "uuid",
      "title": "Storyboard Variation 1",
      "description": "A dynamic narrative approach...",
      "is_selected": false,
      "scenes": [
        {
          "id": "uuid",
          "storyboard_id": "uuid",
          "scene_number": 1,
          "title": "Opening Hook",
          "description": "Start with an attention-grabbing visual",
          "visual_description": "Wide aerial shot transitioning to close-up",
          "duration": 5
        },
        { "...3 more scenes..." }
      ]
    },
    { "...1 more storyboard variation..." }
  ]
}
```

---

#### `GET /api/projects/:id/storyboards`

List all storyboards for a project.

- **Auth Required:** Yes

---

#### `GET /api/storyboards/:id`

Get a single storyboard by ID (with all scenes ordered by scene_number).

- **Auth Required:** Yes

---

#### `POST /api/storyboards/:id/select`

Select a storyboard. Deselects all other storyboards in the same project first.

- **Auth Required:** Yes

---

#### `GET /api/storyboards/:id/scenes`

Get all scenes for a storyboard, ordered by scene_number.

- **Auth Required:** Yes

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Scenes retrieved",
  "data": [
    {
      "id": "uuid",
      "storyboard_id": "uuid",
      "scene_number": 1,
      "title": "Opening Hook",
      "description": "Start with an attention-grabbing visual",
      "visual_description": "Wide aerial shot transitioning to close-up",
      "duration": 5,
      "created_at": "2026-03-10T10:00:00Z"
    }
  ]
}
```

---

### Video APIs

Video generation requires credits. Each generation costs 1 credit.

#### `POST /api/videos/generate`

Generate a video from a selected storyboard. Deducts 1 credit from the user.

- **Auth Required:** Yes

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| project_id | string (UUID) | Yes | Project ID |
| storyboard_id | string (UUID) | Yes | Selected storyboard ID |
| title | string | No | Video title (defaults to "Generated Video") |
| format | string | No | mp4 / webm (default: mp4) |
| resolution | string | No | 1080p / 4K (default: 1080p) |

**Example Request:**
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

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Video generated successfully",
  "data": {
    "id": "uuid",
    "project_id": "uuid",
    "user_id": "uuid",
    "storyboard_id": "uuid",
    "title": "My Brand Video",
    "status": "completed",
    "video_url": "https://storage.example.com/videos/uuid.mp4",
    "thumbnail_url": "https://storage.example.com/thumbnails/uuid.jpg",
    "duration": 120,
    "format": "mp4",
    "resolution": "1080p",
    "file_size": 15728640,
    "credits_used": 1,
    "created_at": "2026-03-10T10:00:00Z"
  }
}
```

**Error Responses:**

| Status | Message |
|--------|---------|
| 400 | project_id and storyboard_id are required |
| 400 | insufficient credits |
| 400 | project not found |
| 400 | storyboard not found |

---

#### `GET /api/videos`

List all videos generated by the authenticated user.

- **Auth Required:** Yes

---

#### `GET /api/videos/:id`

Get a single video by ID.

- **Auth Required:** Yes

---

#### `GET /api/videos/:id/download`

Get download information for a video (stub — returns download metadata).

- **Auth Required:** Yes

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Download info retrieved",
  "data": {
    "video_id": "uuid",
    "video_url": "https://storage.example.com/videos/uuid.mp4",
    "format": "mp4",
    "resolution": "1080p",
    "file_size": 15728640,
    "download_url": "https://storage.example.com/download/uuid.mp4"
  }
}
```

---

#### `GET /api/projects/:id/videos`

List all videos for a specific project.

- **Auth Required:** Yes

---

### Credit APIs

#### `GET /api/credits`

Get the authenticated user's current credit balance.

- **Auth Required:** Yes

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

Admin endpoints require the user to have `role: "admin"`.

#### `POST /api/admin/credits`

Add credits to a user account. Only accessible by admin users.

- **Auth Required:** Yes (admin role)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| user_id | string (UUID) | Yes | Target user ID |
| amount | integer | Yes | Number of credits to add (must be positive) |

**Example Request:**
```bash
curl -X POST http://localhost:3000/api/admin/credits \
  -H "Authorization: Bearer {admin_access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 10
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Credits added successfully",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "credits_added": 10,
    "new_balance": 20
  }
}
```

**Error Responses:**

| Status | Message |
|--------|---------|
| 400 | user_id is required |
| 400 | amount must be positive |
| 400 | unauthorized: only admins can add credits |
| 400 | target user not found |

---

### AI Gateway APIs

#### `GET /api/ai/health`

Check if the AI service is running.

- **Auth Required:** No

---

#### `ANY /api/ai/*`

All requests under `/api/ai/` are proxied to the Python AI service with user context injected.

- **Auth Required:** Yes
- **Injected Headers:** `X-User-ID`, `X-User-Email`

Available AI service stub endpoints:

| Method | AI Service Endpoint | Description |
|--------|-------------------|-------------|
| POST | `/generate-content-pillars` | Generate content pillar suggestions |
| POST | `/generate-storyboard` | Generate storyboard with scenes |
| POST | `/generate-video` | Generate video from storyboard |
| GET | `/health` | AI service health check |
| GET | `/status` | AI service status |

---

## Video Generation Workflow

### Complete Flow

```
1. User registers/logs in (gets 10 credits)
        │
        ▼
2. Create a Project
   POST /api/projects
        │
        ▼
3. Create Business Brief (linked to Project)
   POST /api/briefs/business { project_id: "..." }
        │
        ▼
4. Generate Content Pillars (AI returns 3 pillars with themes)
   POST /api/projects/:id/content-pillars/generate
        │
        ▼
5. Browse pillars & select one
   GET /api/projects/:id/content-pillars
   POST /api/content-pillars/:id/select
        │
        ▼
6. Browse themes & select one
   GET /api/content-pillars/:id/themes
   POST /api/content-themes/:id/select
        │
        ▼
7. Generate Storyboards (AI returns 2 variations with scenes)
   POST /api/projects/:id/storyboards/generate { content_theme_id: "..." }
        │
        ▼
8. Review scenes & select a storyboard
   GET /api/storyboards/:id/scenes
   POST /api/storyboards/:id/select
        │
        ▼
9. Generate Video (costs 1 credit)
   POST /api/videos/generate { project_id, storyboard_id }
        │
        ▼
10. Preview & download
    GET /api/videos/:id
    GET /api/videos/:id/download
```

### Credit Flow

1. User starts with 10 credits on registration
2. Each video generation deducts 1 credit
3. User can check balance via `GET /api/credits`
4. Admin can top-up credits via `POST /api/admin/credits`
5. Generation is blocked when credits reach 0

---

## Database Schema

### Users Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | User identifier |
| name | VARCHAR(255) | NOT NULL | Full name |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Email address |
| password | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| role | VARCHAR(255) | DEFAULT 'user' | user / admin |
| credits | INTEGER | DEFAULT 10 | AI credit balance |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Projects Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Project identifier |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| name | VARCHAR(255) | NOT NULL | Project name |
| description | TEXT | | Project description |
| status | VARCHAR(255) | DEFAULT 'draft' | draft / in_progress / completed |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Business Briefs Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Brief identifier |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| project_id | UUID | FK→projects.id | Parent project |
| project_name | VARCHAR(255) | NOT NULL | Project name |
| company_name | VARCHAR(255) | | Company name |
| industry | VARCHAR(255) | | Industry sector |
| target_audience | VARCHAR(255) | | Target audience |
| project_objective | TEXT | | Project objectives |
| key_message | TEXT | | Core message |
| budget | VARCHAR(255) | | Budget info |
| timeline | VARCHAR(255) | | Timeline info |
| competitors | TEXT | | Competitor info |
| additional_notes | TEXT | | Additional notes |
| status | VARCHAR(255) | DEFAULT 'draft' | draft/submitted/approved/rejected |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Creative Briefs Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Brief identifier |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| business_brief_id | UUID | NOT NULL, FK→business_briefs.id | Parent business brief |
| title | VARCHAR(255) | NOT NULL | Brief title |
| video_type | VARCHAR(255) | | promotional/educational/testimonial/explainer |
| duration | INTEGER | | Duration in seconds |
| style | VARCHAR(255) | | cinematic/animated/minimalist |
| tone | VARCHAR(255) | | professional/casual/energetic |
| script | TEXT | | Video script |
| storyboard | TEXT | | Storyboard description |
| visual_references | TEXT | | Visual references |
| music_preference | VARCHAR(255) | | Music style |
| call_to_action | VARCHAR(255) | | CTA text |
| output_format | VARCHAR(255) | | mp4/webm |
| resolution | VARCHAR(255) | | 1080p/4K |
| additional_notes | TEXT | | Additional notes |
| status | VARCHAR(255) | DEFAULT 'draft' | draft/submitted/in_production/completed |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Content Pillars Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Pillar identifier |
| project_id | UUID | NOT NULL, FK→projects.id | Parent project |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| title | VARCHAR(255) | NOT NULL | Pillar title |
| description | TEXT | | Pillar description |
| is_selected | BOOLEAN | DEFAULT FALSE | Selection status |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Content Themes Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Theme identifier |
| content_pillar_id | UUID | NOT NULL, FK→content_pillars.id | Parent pillar |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| title | VARCHAR(255) | NOT NULL | Theme title |
| description | TEXT | | Theme description |
| is_selected | BOOLEAN | DEFAULT FALSE | Selection status |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Storyboards Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Storyboard identifier |
| project_id | UUID | NOT NULL, FK→projects.id | Parent project |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| title | VARCHAR(255) | NOT NULL | Storyboard title |
| description | TEXT | | Storyboard description |
| is_selected | BOOLEAN | DEFAULT FALSE | Selection status |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Scenes Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Scene identifier |
| storyboard_id | UUID | NOT NULL, FK→storyboards.id | Parent storyboard |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| scene_number | INTEGER | NOT NULL | Order in storyboard |
| title | VARCHAR(255) | | Scene title |
| description | TEXT | | Scene description |
| visual_description | TEXT | | Visual direction notes |
| duration | INTEGER | | Scene duration in seconds |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Videos Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Video identifier |
| project_id | UUID | NOT NULL, FK→projects.id | Parent project |
| user_id | UUID | NOT NULL, FK→users.id | Owner user |
| storyboard_id | UUID | NOT NULL, FK→storyboards.id | Source storyboard |
| title | VARCHAR(255) | NOT NULL | Video title |
| status | VARCHAR(255) | DEFAULT 'pending' | pending/processing/completed/failed |
| video_url | TEXT | | Video file URL |
| thumbnail_url | TEXT | | Thumbnail URL |
| duration | INTEGER | | Duration in seconds |
| format | VARCHAR(255) | DEFAULT 'mp4' | mp4/webm |
| resolution | VARCHAR(255) | DEFAULT '1080p' | 1080p/4K |
| file_size | BIGINT | | File size in bytes |
| credits_used | INTEGER | DEFAULT 1 | Credits consumed |
| error_message | TEXT | | Error details (if failed) |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

---

## Interactive Documentation

### FastAPI Auto-generated Docs

The Python AI service (FastAPI) automatically exposes interactive documentation:

| URL | Description |
|-----|-------------|
| `http://localhost:8000/docs` | Swagger UI — interactive API testing |
| `http://localhost:8000/redoc` | ReDoc — readable API documentation |
| `http://localhost:8000/openapi.json` | OpenAPI 3.0 JSON schema |

### Postman / Bruno Collection

A ready-to-import API collection is available at:

- `docs/postman_collection.json` — Postman v2.1 format (also compatible with Bruno)

Import into Postman or Bruno for interactive endpoint testing with pre-configured requests and environment variables.
