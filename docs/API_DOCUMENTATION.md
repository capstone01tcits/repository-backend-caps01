# API Documentation — AI Video Generation Platform

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
6. [Core API Endpoints](#core-api-endpoints)
7. [Video Generation API](#video-generation-api)
8. [Video Generation System Details](#video-generation-system-details)
9. [Database Schema](#database-schema)
10. [Setup & Configuration](#setup--configuration)

---

## Overview

The AI Video Generation Platform Backend provides RESTful APIs for:

- **User authentication** with role-based access (user / admin)
- **Project management** (dashboard for video projects)
- **Brief management** (business briefs & creative briefs)
- **AI content pillar generation** with selection workflow
- **AI storyboard generation** with scene timeline
- **AI video generation** (3 variations per briefing with scene-based generation)
- **Credit management** (user balance & admin top-up)
- **AI service gateway** (proxy to Python AI service)

### Service Stack

| Service | Port | Stack | Purpose |
|---------|------|-------|---------|
| Go Backend | 3000 | Fiber + GORM + PostgreSQL | Auth, Projects, Briefs, Video Generation |
| Python AI Service | 8000 | FastAPI + Uvicorn | AI Processing |

---

## Architecture

```
Frontend (Ember JS)
       │
       ▼
┌──────────────────────────────────────────┐
│   Go Backend (Port 3000)                 │
│                                          │
│  ┌──────────────────────────────────┐    │
│  │  API Layer (Handlers)            │    │
│  │  - Auth, Projects, Briefs        │    │
│  │  - Video Generation (NEW)        │    │
│  └────────────┬─────────────────────┘    │
│               │                          │
│  ┌────────────┴─────────────────────┐    │
│  │  Business Logic Layer (Services) │    │
│  │  - VideoGenerationService (NEW)  │    │
│  │  - Credit Management             │    │
│  └────────────┬─────────────────────┘    │
│               │                          │
│  ┌────────────┴─────────────────────┐    │
│  │  Data Layer (Repositories)       │    │
│  │  - Video Repositories (NEW)      │    │
│  │  - Job Queue Management (NEW)    │    │
│  └────────────┬─────────────────────┘    │
│               │                          │
│  ┌────────────┴─────────────────────┐    │
│  │  PostgreSQL + 12 Tables          │    │
│  │  - Original 9 tables             │    │
│  │  - GenerationJob, VideoVariant   │    │
│  │  - SceneGeneration (NEW)         │    │
│  └──────────────────────────────────┘    │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  Background Job Queue              │  │
│  │  - 3 Worker Goroutines             │  │
│  │  - Status Polling (60s interval)   │  │
│  │  - Automatic Retry Logic           │  │
│  └────────────┬───────────────────────┘  │
│               │                          │
│  ┌────────────┴───────────────────────┐  │
│  │  AI Provider Abstraction (NEW)     │  │
│  │  - LTX, Runway, Wan2.1             │  │
│  │  - Mock implementations            │  │
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
| `user` | • Create/edit own projects<br/>• Access own briefs and content<br/>• Generate videos (credit-limited)<br/>• View own credits |
| `admin` | • All user permissions<br/>• Add credits to users<br/>• View system statistics |

---

## Standard Response Format

### Success Response (2xx)

```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    ...
  },
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "2.0.0"
  }
}
```

### Error Response (4xx, 5xx)

```json
{
  "status": "error",
  "message": "Invalid request",
  "error": "Email already registered",
  "code": "DUPLICATE_EMAIL"
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
| 409 | Conflict | Duplicate resource |
| 500 | Internal Server Error | Server error |

### Common Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| INVALID_CREDENTIALS | 401 | Email/password incorrect |
| INSUFFICIENT_CREDITS | 400 | User doesn't have enough credits |
| DUPLICATE_EMAIL | 409 | Email already registered |
| NOT_FOUND | 404 | Resource doesn't exist |
| UNAUTHORIZED | 401 | No/invalid authorization |
| FORBIDDEN | 403 | No permission for action |

---

## Core API Endpoints

### Health Check (Public)

```http
GET /health
```

Check backend and AI service status.

**Response (200 OK):**
```json
{
  "status": "ok",
  "message": "Backend service is running",
  "ai_service": "connected"
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
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "credits": 10,
    "created_at": "2024-01-15T10:00:00Z"
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
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "credits": 10
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
  "status": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "credits": 10,
    "created_at": "2024-01-15T10:00:00Z"
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

#### Refresh Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Delete Account

```http
DELETE /api/auth/account
Authorization: Bearer {access_token}
```

---

### Project APIs

#### Create Project

```http
POST /api/projects
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Brand Campaign Q1 2026",
  "description": "Marketing video campaign for Q1"
}
```

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

#### Update Project

```http
PUT /api/projects/{project_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Project Name",
  "description": "Updated description"
}
```

#### Delete Project

```http
DELETE /api/projects/{project_id}
Authorization: Bearer {access_token}
```

---

### Business Brief APIs

#### Create Business Brief

```http
POST /api/briefs/business
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "project_name": "Q1 Campaign",
  "company_name": "Tech Corp",
  "industry": "Technology",
  "target_audience": "Developers 25-40",
  "project_objective": "Increase awareness",
  "key_message": "Innovation",
  "budget": "50000000",
  "timeline": "4 weeks"
}
```

#### List Business Briefs

```http
GET /api/briefs/business
Authorization: Bearer {access_token}
```

#### Get Business Brief

```http
GET /api/briefs/business/{brief_id}
Authorization: Bearer {access_token}
```

#### Update Business Brief

```http
PUT /api/briefs/business/{brief_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_name": "Updated Brand Video Q3",
  "company_name": "Updated Company Name",
  "industry": "Updated Industry",
  "target_audience": "Updated audience",
  "project_objective": "Updated objective",
  "key_message": "Updated message",
  "budget": "75000000",
  "timeline": "6 weeks",
  "status": "submitted"
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Business brief updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "project_name": "Updated Brand Video Q3",
    "status": "submitted",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

#### Delete Business Brief

```http
DELETE /api/briefs/business/{brief_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Business brief deleted successfully"
}
```

---

### Creative Brief APIs

#### Create Creative Brief

```http
POST /api/briefs/creative
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "business_brief_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Product Launch Video",
  "video_type": "promotional",
  "duration": 60,
  "style": "cinematic",
  "tone": "professional",
  "script": "Opening shot: wide aerial view...",
  "call_to_action": "Visit our website",
  "output_format": "mp4",
  "resolution": "1080p"
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Creative brief created successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "business_brief_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Product Launch Video",
    "duration": 60,
    "created_at": "2024-01-15T10:40:00Z"
  }
}
```

#### List Creative Briefs

```http
GET /api/briefs/creative
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "data": [
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "business_brief_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Product Launch Video",
      "video_type": "promotional",
      "duration": 60,
      "created_at": "2024-01-15T10:40:00Z"
    }
  ]
}
```

#### Get Single Creative Brief

```http
GET /api/briefs/creative/{creative_brief_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "business_brief_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Product Launch Video",
    "video_type": "promotional",
    "duration": 60,
    "style": "cinematic",
    "tone": "professional",
    "script": "Opening shot: wide aerial view...",
    "call_to_action": "Visit our website",
    "output_format": "mp4",
    "resolution": "1080p",
    "created_at": "2024-01-15T10:40:00Z"
  }
}
```

#### Update Creative Brief

```http
PUT /api/briefs/creative/{creative_brief_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "Updated Product Launch Video",
  "duration": 90,
  "tone": "uplifting",
  "status": "submitted"
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Creative brief updated successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "title": "Updated Product Launch Video",
    "duration": 90,
    "tone": "uplifting",
    "status": "submitted"
  }
}
```

####Delete Creative Brief

```http
DELETE /api/briefs/creative/{creative_brief_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Creative brief deleted successfully"
}
```

#### List Creative Briefs by Business Brief

```http
GET /api/briefs/business/{business_brief_id}/creative
Authorization: Bearer {access_token}
```

---

### Content Pillar APIs

#### Generate Content Pillars

```http
POST /api/projects/{project_id}/content-pillars/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

Cost: 1 credit per generation

#### List Content Pillars

```http
GET /api/projects/{project_id}/content-pillars
Authorization: Bearer {access_token}
```

#### Select Content Pillar

```http
POST /api/content-pillars/{pillar_id}/select
Authorization: Bearer {access_token}
```

---

### Content Theme APIs

#### List Themes for Pillar

```http
GET /api/content-pillars/{pillar_id}/themes
Authorization: Bearer {access_token}
```

#### Select Theme

```http
POST /api/content-themes/{theme_id}/select
Authorization: Bearer {access_token}
```

---

### Storyboard APIs

#### Generate Storyboards

```http
POST /api/projects/{project_id}/storyboards/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "content_theme_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

Cost: 1 credit per generation

#### List Storyboards

```http
GET /api/projects/{project_id}/storyboards
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "data": [
    {
      "id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Storyboard 1 - Premium Cinematic",
      "description": "Professional product showcase",
      "version": 1,
      "total_duration": 12,
      "created_at": "2024-01-15T10:45:00Z"
    }
  ]
}
```

#### Get Single Storyboard

```http
GET /api/storyboards/{storyboard_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Storyboard 1 - Premium Cinematic",
    "description": "Professional product showcase",
    "version": 1,
    "total_duration": 12,
    "created_at": "2024-01-15T10:45:00Z"
  }
}
```

#### Select Storyboard

```http
POST /api/storyboards/{storyboard_id}/select
Authorization: Bearer {access_token}
```

#### Get Scenes

```http
GET /api/storyboards/{storyboard_id}/scenes
Authorization: Bearer {access_token}
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
  "status": "success",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "balance": 25,
    "total_earned": 50,
    "total_used": 25
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
  "amount": 100,
  "reason": "Promotion bonus"
}
```

---

## Video Generation API

The video generation system creates **3 automatic variations** of videos from a single briefing using **scene-based generation**.

### Core Concept

- **Input**: One project briefing with storyboard
- **Output**: 3 video variations (cinematic, vibrant, professional)
- **Composition**: 2-3 scenes per video, 4-6 seconds each = 8-12 sec total
- **Processing**: Asynchronous queue with background workers
- **Cost**: Calculated by provider and operation type

### Video Generation Endpoints

#### Generate Video

```http
POST /api/videos/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "title": "My Brand Video",
  "format": "mp4",
  "resolution": "1080p"
}
```

**Parameters:**
- `project_id` (required): UUID of project
- `storyboard_id` (required): UUID of storyboard  
- `title` (required): Video title
- `format` (optional): Video format (default: "mp4")
- `resolution` (optional): Video resolution (default: "1080p")

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Video generation job created",
  "data": {
    "generation_job_id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "title": "My Brand Video",
    "format": "mp4",
    "resolution": "1080p",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

#### Get Generation Job Status

```http
GET /api/videos/generation/{job_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Job status retrieved",
  "data": {
    "id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "job_type": "generate",
    "status": "processing",
    "priority": 1,
    "scene_count": 2,
    "video_duration": 10,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "credits_required": 120,
    "credits_used": 0,
    "retry_count": 0,
    "started_at": "2024-01-15T10:30:15Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:15Z"
  }
}
```

**Job Statuses:**
- `queued`: Waiting for worker to process
- `processing`: Currently generating videos
- `completed`: All videos generated successfully
- `failed`: Generation failed (see error_message)

---

#### Get All Video Variants for Storyboard

```http
GET /api/videos/storyboard/{storyboard_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Video variants retrieved",
  "data": [
    {
      "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
      "variant_number": 1,
      "status": "completed",
      "video_url": "https://storage.example.com/videos/8da7b810.mp4",
      "thumbnail_url": "https://storage.example.com/thumbnails/8da7b810.jpg",
      "prompt_used": "Professional marketing video cinematic style",
      "duration": 10,
      "provider": "ltx",
      "model": "ltx-2-fast",
      "resolution": "1080p",
      "credits_used": 40,
      "scenes": [
        {
          "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 1,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-1.mp4",
          "duration": 5,
          "updated_at": "2024-01-15T10:35:30Z"
        },
        {
          "id": "0ea7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 2,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-2.mp4",
          "duration": 5,
          "updated_at": "2024-01-15T10:35:45Z"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:35:45Z"
    }
  ]
}
```

---

#### Get Single Video Variant

```http
GET /api/videos/{variant_id}
Authorization: Bearer {access_token}
```

Returns variant details with all scenes.

---

#### List My Videos

```http
GET /api/videos
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "data": [
    {
      "id": "video_1",
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "My Brand Video",
      "status": "completed",
      "format": "mp4",
      "resolution": "1080p",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

#### Get Single Video

```http
GET /api/videos/{video_id}
Authorization: Bearer {access_token}
```

---

#### Get Videos by Project

```http
GET /api/projects/{project_id}/videos
Authorization: Bearer {access_token}
```

---

#### Download Video

```http
POST /api/videos/{variant_id}/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "new_prompt": "More emphasis on campus facilities"
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Video regeneration job created",
  "data": {
    "generation_job_id": "3fa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued"
  }
}
```

**Cost**: Same as original video generation (~120 credits)

---

#### Regenerate Single Scene

```http
POST /api/videos/scene/{scene_id}/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "new_prompt": "Show aerial campus views"
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Scene regeneration job created",
  "data": {
    "generation_job_id": "4ga7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued"
  }
}
```

**Cost**: Lower than full video (~10 credits for single scene)

---

#### Download Video Variant

```http
GET /api/videos/{variant_id}/download
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "status": "success",
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
- Download URL is signed and expires after 1 hour
- Can be downloaded multiple times

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
- LTX-2-Fast (Standard): $0.04/sec = 4 credits/sec
- LTX-2-Pro (Premium): $0.06/sec = 6 credits/sec
- Runway Gen4.5: 12 credits/sec
- Runway Turbo: 5 credits/sec
- Wan2.1 (Open Source): 1 credit/sec
- LTX Open Source: 1 credit/sec

### Job Lifecycle

```
User Request
    ↓
Validate credits
    ↓
Create GenerationJob (status: queued)
Create 3 VideoVariants
Create 2-3 SceneGenerations per variant
Deduct credits
    ↓
Enqueue to JobQueue
    ↓
Worker picks up job (status: processing)
    ↓
For each scene:
  - Call VideoProvider.GenerateScene()
  - Store external job ID
    ↓
Start polling routine (every 60 seconds)
    ↓
Provider returns scene status
    ↓
When all scenes complete:
  - Update variant status
  - Generate thumbnail
  - Update video URLs
    ↓
When all variants complete:
  - Update job status (completed)
  - Videos ready for download
```

### Credit System

**Generation Costs:**
```
Cost = video_duration × scene_count × variant_count × provider_multiplier

Example (Standard Tier):
10 sec × 2 scenes × 3 variants × 2 = 120 credits
```

**Provider Costs (per second):**
- LTX-2-Fast: 2 credits/sec
- LTX-2-Pro: 3 credits/sec
- Runway Gen4.5: 12 credits/sec
- Runway Turbo: 5 credits/sec
- Open Source: 1 credit/sec

**Regeneration:**
- Full video: Same as original generation
- Single scene: Significantly cheaper (~1/3 cost)

### Polling Recommendations

**Frontend Polling Strategy:**
1. Poll every 2-5 seconds while `status: processing`
2. After 30 seconds, reduce to 10-second intervals
3. After 1 minute, use 60-second intervals
4. Stop polling after 2 hours (timeout)

**Example Poll Sequence:**
```
0s:   POST /api/videos/generate → job_id
2s:   GET /api/videos/generation/{job_id} → status: queued
5s:   GET /api/videos/generation/{job_id} → status: processing
10s:  GET /api/videos/storyboard/{id} → variants (mixed statuses)
20s:  GET /api/videos/storyboard/{id} → variant 1 complete
45s:  GET /api/videos/storyboard/{id} → all variants complete
```

---

## Database Schema

### Tables (12 total)

| Table | Purpose | Key Fields |
|-------|---------|-----------|
| users | User accounts | id, email, password_hash, role, credits |
| projects | Video projects | id, user_id, name, status |
| business_briefs | Business input | id, project_id, company_name, target_audience |
| creative_briefs | Creative direction | id, business_brief_id, tone, style |
| content_pillars | AI-generated pillars | id, project_id, title |
| content_themes | Pillar themes | id, pillar_id, title |
| storyboards | Storyboard variations | id, project_id, title |
| scenes | Individual scenes | id, storyboard_id, scene_number, duration |
| videos | Generated videos | id, storyboard_id, video_url (legacy) |
| generation_jobs | Video job queue | id, status, provider, credits_required |
| video_variants | 3 video variations | id, variant_number, video_url, status |
| scene_generations | Individual scenes | id, variant_id, scene_number, video_url |

### Primary Indexes

- `generation_jobs.status` (queue queries)
- `generation_jobs.user_id`, `generation_jobs.project_id`
- `video_variants.storyboard_id` (batch retrieval)
- `scene_generations.variant_id` (composition)

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

# Video Generation
VIDEO_GENERATION_WORKERS=3
VIDEO_POLLING_INTERVAL=60
VIDEO_MAX_RETRIES=3
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
- **Architecture Details**: See README.md Video Generation System section
