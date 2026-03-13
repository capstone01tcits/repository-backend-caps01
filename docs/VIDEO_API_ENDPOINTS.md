# Video Generation API Endpoints

## Overview

This document describes all API endpoints for the AI Video Generator backend.

## Base URL
```
http://localhost:8000/api
```

## Authentication
All endpoints require JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## Endpoints

### 1. Generate Video Variants
Generate 3 video variations from a storyboard briefing.

**Request**
```http
POST /videos/generate
Content-Type: application/json
Authorization: Bearer <token>

{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "custom_prompt": "Focus on campus beauty and student life",
  "scene_count": 2,
  "video_duration": 10
}
```

**Parameters**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| project_id | uuid | Yes | Project containing the storyboard |
| storyboard_id | uuid | Yes | Storyboard to generate videos from |
| custom_prompt | string | No | Custom direction for video generation |
| scene_count | int | No | Number of scenes (2-3, default: 2) |
| video_duration | int | No | Total duration in seconds (8-12, default: 10) |

**Response** (201 Created)
```json
{
  "status": "success",
  "message": "Video generation job created",
  "data": {
    "generation_job_id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Status Codes**
- 201: Job created successfully
- 400: Invalid request or insufficient credits
- 401: Unauthorized
- 500: Server error

---

### 2. Get Generation Job Status
Check the progress of a video generation job.

**Request**
```http
GET /videos/generation/7ca7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer <token>
```

**Response** (200 OK)
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
    "completed_at": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:15Z"
  }
}
```

**Status Values**
| Status | Description |
|--------|-------------|
| queued | Waiting to be processed |
| processing | Currently being generated |
| completed | Successfully generated all videos |
| failed | Generation failed |

---

### 3. Get Video Variants for Storyboard
Retrieve all 3 video variants generated for a storyboard.

**Request**
```http
GET /videos/storyboard/6ba7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer <token>
```

**Response** (200 OK)
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
      "prompt_used": "Generate a professional marketing video with cinematic style. This is variation 1.",
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
    },
    {
      "id": "1ea7b810-9dad-11d1-80b4-00c04fd430c8",
      "variant_number": 2,
      "status": "processing",
      "video_url": null,
      "thumbnail_url": null,
      "prompt_used": "Generate a professional marketing video with vibrant and dynamic style. This is variation 2.",
      "duration": 10,
      "provider": "ltx",
      "model": "ltx-2-fast",
      "scenes": [...]
    },
    {
      "id": "2ea7b810-9dad-11d1-80b4-00c04fd430c8",
      "variant_number": 3,
      "status": "queued",
      "video_url": null,
      "thumbnail_url": null,
      "prompt_used": "Generate a professional marketing video with professional and polished style. This is variation 3.",
      "duration": 10,
      "provider": "ltx",
      "model": "ltx-2-fast",
      "scenes": []
    }
  ]
}
```

---

### 4. Get Single Video Variant
Retrieve details of a specific video variant with all its scenes.

**Request**
```http
GET /videos/8da7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer <token>
```

**Response** (200 OK)
```json
{
  "status": "success",
  "message": "Video retrieved",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "variant_number": 1,
    "status": "completed",
    "video_url": "https://storage.example.com/videos/8da7b810.mp4",
    "thumbnail_url": "https://storage.example.com/thumbnails/8da7b810.jpg",
    "prompt_used": "Generate a professional marketing video with cinematic style. This is variation 1.",
    "duration": 10,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "resolution": "1080p",
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
}
```

---

### 5. Regenerate Video Variant
Regenerate a single video variant with optional new prompt.

**Request**
```http
POST /videos/8da7b810-9dad-11d1-80b4-00c04fd430c8/regenerate
Content-Type: application/json
Authorization: Bearer <token>

{
  "new_prompt": "Focus more on student activities and campus facilities"
}
```

**Parameters**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| new_prompt | string | No | Updated prompt for all scenes in the variant |

**Response** (201 Created)
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

**Notes**
- Creates new GenerationJob with type "regenerate"
- Deducts regeneration credits from user account
- Original video variant remains unchanged until regeneration completes
- New variant returned on completion

---

### 6. Regenerate Scene
Regenerate a single scene within a video variant.

**Request**
```http
POST /videos/scene/9da7b810-9dad-11d1-80b4-00c04fd430c8/regenerate
Content-Type: application/json
Authorization: Bearer <token>

{
  "new_prompt": "Show more campus buildings in this scene"
}
```

**Parameters**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| new_prompt | string | No | Updated prompt for this specific scene |

**Response** (201 Created)
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

**Notes**
- Creates new GenerationJob with type "regenerate_scene"
- Lower credit cost than full video regeneration
- Only specified scene is regenerated
- Download costs are cheaper than full video

---

### 7. Download Video Variant
Get download URL for a completed video variant.

**Request**
```http
GET /videos/8da7b810-9dad-11d1-80b4-00c04fd430c8/download
Authorization: Bearer <token>
```

**Response** (200 OK)
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

**Status Codes**
- 200: Download URL ready
- 400: Video not completed yet or not found
- 401: Unauthorized

**Notes**
- Only available if variant status is "completed"
- URL is signed and has temporary expiration (typically 1 hour)
- Can be downloaded multiple times
- For streaming, use video_url from variant details

---

## Error Responses

### 400 Bad Request
```json
{
  "status": "error",
  "message": "Invalid request",
  "error": "project_id is required"
}
```

### 401 Unauthorized
```json
{
  "status": "error",
  "message": "Unauthorized",
  "error": "Invalid or missing authentication token"
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "Not found",
  "error": "Video variant not found"
}
```

### 500 Internal Server Error
```json
{
  "status": "error",
  "message": "Internal server error",
  "error": "Database operation failed"
}
```

---

## Status Transitions

### Video Variant Status Flow
```
pending → processing → completed
              ↓
            failed
```

scenes individually:
```
pending → processing → completed
              ↓
            failed
```

### Generation Job Status Flow
```
queued → processing → completed
            ↓
          failed
```

---

## Rate Limiting

Currently no rate limiting implemented. Recommended:
- 100 requests per minute per user for reading
- 10 concurrent generation jobs per user
- 3 regenerations per video per 24 hours

---

## Polling Recommendations

### For Frontend
1. Initially poll every 2-5 seconds while status is "processing"
2. After 30 seconds, reduce to 10-second intervals
3. Implement exponential backoff up to 60-second intervals
4. Set maximum poll duration of 2 hours

### Example Poll Sequence
```
0s:    POST /videos/generate → job_id
2s:    GET /videos/generation/{job_id} → status: queued
5s:    GET /videos/generation/{job_id} → status: processing
10s:   GET /videos/storyboard/{storyboard_id} → variants (mixed statuses)
20s:   GET /videos/storyboard/{storyboard_id} → variant 1 complete
30s:   GET /videos/storyboard/{storyboard_id} → variants 1,2 complete
45s:   GET /videos/storyboard/{storyboard_id} → all variants complete
```

---

## WebSocket Events (Future)

Planned for real-time updates:
```
- video:generation:started
- video:generation:processing
- video:variant:completed
- video:scene:completed
- video:generation:failed
```
