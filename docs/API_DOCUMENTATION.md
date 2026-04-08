# API Documentation — AI Video Generation Platform

> **Version:** 2.1.0 (Updated Sprint 3 - April 2026)
> **Base URL:** `http://localhost:3000`  
> **AI Service URL:** `http://localhost:8000`
> **Collection Format:** Bruno API Collection (docs/bruno_api_collection.json)

---

## Sprint 3 Updates (April 2026) ✓

**Field Changes in Business Brief:**
- ✓ `company_name` → `institute_name` (for educational institutions)
- ✓ `industry` → `education` (changed semantics for education domain)
- ✓ `target_audience` string → boolean (yes/no representation)
- ✓ `budget` + `timeline` → `deadline` (merged into single datetime field)

**New Fields Added:**
- ✓ `theme` field in Project model
- ✓ `prompt` & `video_url` fields in ContentPillar
- ✓ `prompt` field in Storyboard
- ✓ `caption` field in Scene (for generated captions)
- ✓ `regenerate_count` in Video, VideoVariant, Scene (max 3 times)

**Regenerate Limits:**
- ✓ Max 3 regenerate attempts per video
- ✓ Max 3 regenerate attempts per scene
- ✗ Error returned when limit exceeded

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

#### Restore Account

```http
POST /api/auth/restore
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

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
    "created_at": "2026-03-13T10:00:00Z",
    "updated_at": "2026-03-13T10:00:00Z"
  }
}
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
  "description": "Marketing video campaign for Q1",
  "theme": "Corporate Branding"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Project created successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Brand Campaign Q1 2026",
    "description": "Marketing video campaign for Q1",
    "theme": "Corporate Branding",
    "status": "draft",
    "created_at": "2026-03-13T10:05:00Z",
    "updated_at": "2026-03-13T10:05:00Z"
  }
}
```

**✓ NEW in Sprint 3:** Added `theme` field to projects for storing theme information

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
  "description": "Updated description",
  "theme": "Updated Theme"
}
```

**✓ NEW in Sprint 3:** `theme` field can be updated

#### Delete Project

```http
DELETE /api/projects/{project_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Project deleted successfully",
  "data": null
}
```

---

### Business Brief APIs

#### Create Business Brief

```http
POST /api/briefs/business
Authorization: Bearer {access_token}
Content-Type: application/json

{institute_name": "Universitas XYZ",
  "education": "Higher Education",
  "target_audience": true,
  "project_objective": "Increase awareness",
  "key_message": "Innovation",
  "deadline": "2026-05-15T00:00:00Z",
  "competitors": "Kompetitor lainnya",
  "additional_notes": "Catatan tambahan"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Business brief created successfully",
  "data": {
    "id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_name": "Q1 Campaign",
    "institute_name": "Universitas XYZ",
    "education": "Higher Education",
    "target_audience": true,
    "deadline": "2026-05-15T00:00:00Z41d4-a716-446655440000",
    "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_name": "Q1 Campaign",
    "company_name": "Tech Corp",
    "status": "draft",
    "created_at": "2026-03-13T10:10:00Z",
    "updated_at": "2026-03-13T10:10:00Z"
  }
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
  "status": "submitted"
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
  "success": true,
  "message": "Business brief deleted successfully",
  "data": null
}
```

#### List Creative Briefs by Business Brief

```http
GET /api/briefs/business/{business_brief_id}/creative
Authorization: Bearer {access_token}
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
  "success": true,
  "message": "Creative brief created successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "business_brief_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Product Launch Video",
    "duration": 60,
    "status": "draft",
    "created_at": "2026-03-13T10:40:00Z",
    "updated_at": "2026-03-13T10:40:00Z"
  }
}
```

#### List Creative Briefs

```http
GET /api/briefs/creative
Authorization: Bearer {access_token}
```

#### Get Creative Brief

```http
GET /api/briefs/creative/{creative_brief_id}
Authorization: Bearer {access_token}
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

#### Delete Creative Brief

```http
DELETE /api/briefs/creative/{creative_brief_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Creative brief deleted successfully",
  "data": null
}
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

Cost: 0 credits

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Content pillars generated successfully",
  "data": [
    {
      "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Brand Awareness",
      "description": "Content focused on increasing brand visibility and recognition",
      "is_selected": false,
      "content_themes": [
        {
          "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
          "title": "Brand Awareness - Theme A",
          "description": "First theme variation for Brand Awareness",
          "is_selected": false
        },
        {
          "id": "0ea7b810-9dad-11d1-80b4-00c04fd430c8",
          "title": "Brand Awareness - Theme B",
          "description": "Second theme variation for Brand Awareness",
          "is_selected": false
        }
      ],
      "created_at": "2026-03-13T10:20:00Z",
      "updated_at": "2026-03-13T10:20:00Z"
    }
  ]
}
```

#### List Content Pillars

```http
GET /api/projects/{project_id}/content-pillars
Authorization: Bearer {access_token}
```

#### Get Content Pillar

```http
GET /api/content-pillars/{pillar_id}
Authorization: Bearer {access_token}
```

#### Select Content Pillar

```http
POST /api/content-pillars/{pillar_id}/select
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Content pillar selected",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "updated_at": "2026-03-13T10:21:00Z"
  }
}
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

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Content theme selected",
  "data": {
    "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "updated_at": "2026-03-13T10:22:00Z"
  }
}
```

---

### Storyboard APIs

#### Generate Storyboards

```http
POST /api/projects/{project_id}/storyboards/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "content_theme_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

Cost: 0 credits

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Storyboards generated successfully",
  "data": [
    {
      "id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Storyboard A - Dynamic",
      "description": "A dynamic, fast-paced storyboard",
      "is_selected": false,
      "scenes": [
        {
          "id": "6ka7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 1,
          "title": "Opening Hook",
          "description": "Attention-grabbing opening sequence",
          "visual_description": "Wide shot with dramatic lighting and motion graphics",
          "duration": 5
        },
        {
          "id": "7la7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 2,
          "title": "Problem Statement",
          "description": "Present the challenge or pain point",
          "visual_description": "Close-up shots with text overlay",
          "duration": 8
        },
        {
          "id": "8ma7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 3,
          "title": "Solution Reveal",
          "description": "Introduce the product as solution",
          "visual_description": "Product showcase with smooth transitions",
          "duration": 10
        },
        {
          "id": "9na7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 4,
          "title": "Call to Action",
          "description": "End with clear CTA",
          "visual_description": "Logo animation with contact info",
          "duration": 7
        }
      ],
      "created_at": "2026-03-13T10:25:00Z",
      "updated_at": "2026-03-13T10:25:00Z"
    }
  ]
}
```

#### List Storyboards

```http
GET /api/projects/{project_id}/storyboards
Authorization: Bearer {access_token}
```

#### Get Storyboard

```http
GET /api/storyboards/{storyboard_id}
Authorization: Bearer {access_token}
```

#### Select Storyboard

```http
POST /api/storyboards/{storyboard_id}/select
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Storyboard selected",
  "data": {
    "id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "updated_at": "2026-03-13T10:26:00Z"
  }
}
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

The video generation system creates **3 automatic variations** of videos from a single briefing using **scene-based generation**.

### Core Concept

- **Input**: One project briefing with storyboard
- **Output**: 3 video variations (cinematic, vibrant, professional)
- **Composition**: 2-3 scenes per video, 4-6 seconds each = 8-12 seconds total
- **Processing**: Asynchronous queue with background workers
- **Cost**: Only charged for video generation, not for content pillar or storyboard generation

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
- `custom_prompt` (optional): Custom prompt override

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

#### Get Generation Job Status

```http
GET /api/videos/generation/{job_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Job status retrieved",
  "data": {
    "id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "job_type": "generate",
    "status": "processing",
    "scene_count": 2,
    "video_duration": 10,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "credits_required": 120,
    "credits_used": 0,
    "retry_count": 0,
    "started_at": "2026-03-13T10:30:15Z",
    "created_at": "2026-03-13T10:30:00Z",
    "updated_at": "2026-03-13T10:30:15Z"
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
  "success": true,
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
      "scenes": [
        {
          "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 1,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-1.mp4",
          "duration": 5,
          "updated_at": "2026-03-13T10:35:30Z"
        },
        {
          "id": "0ea7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 2,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-2.mp4",
          "duration": 5,
          "updated_at": "2026-03-13T10:35:45Z"
        }
      ],
      "created_at": "2026-03-13T10:30:00Z",
      "updated_at": "2026-03-13T10:35:45Z"
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

#### Regenerate Video Variant (Sprint 3 ✓ NEW)

✓ **NEW in Sprint 3** – Regenerate with `regenerate_count` tracking and max 3 limit

```http
POST /api/videos/{variant_id}/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "provider": "Runway-Gen3",
  "regenerate_all_scenes": true,
  "new_prompt": "More emphasis on campus facilities with modern aesthetic"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Video regeneration job created",
  "data": {
    "generation_job_id": "3fa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "variant_id": "{variant_id}",
    "regenerate_count": 1,
    "max_regenerate": 3,
    "remaining_attempts": 2,
    "credit_cost": 40,
    "current_balance": 120
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
