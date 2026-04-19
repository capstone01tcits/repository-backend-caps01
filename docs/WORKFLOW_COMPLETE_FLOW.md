# Complete Workflow Flow - Register to Video Generation

Simplified complete flow from Register until video ready for download.

**Note:** Use Bruno API Collection (`docs/API_COLLECTION.json`) for automatic variable population. Variables (access_token, project_id, etc.) auto-set after each endpoint—no manual copy-paste needed!

---

## Table of Contents
1. [Step 1: Register User](#step-1-register-user)
2. [Step 2: Login](#step-2-login)
3. [Step 3: Initialize Project (from FE Wizard)](#step-3-initialize-project-from-fe-wizard)
4. [Step 4: Generate Storyboard](#step-4-generate-storyboard)
5. [Step 5: Generate Video](#step-5-generate-video)
6. [Step 6: Check Video Status](#step-6-check-video-status)
7. [Timeline](#complete-visual-timeline)

---

## STEP 1: Register User

### Request
```bash
POST /api/auth/register
Content-Type: application/json

{
  "name": "Budi Wiranto",
  "email": "budi@company.id",
  "password": "secure123456"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Budi Wiranto",
      "email": "budi@company.id",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-13T10:00:00Z",
      "updated_at": "2026-03-13T10:00:00Z"
    }
  }
}
```

### Key Info to Save
- `access_token` - Use for all subsequent requests
- `user.id` - User ID
- `user.credits` - 10 credits (default new user)

---

## STEP 2: Login

### Request
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "budi@company.id",
  "password": "secure123456"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Budi Wiranto",
      "email": "budi@company.id",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-13T10:00:00Z",
      "updated_at": "2026-03-13T10:00:00Z"
    }
  }
}
```

### Headers for Next Requests
```
Authorization: Bearer {access_token}
```

---

## STEP 3: Initialize Project (from FE Wizard)

**✓ NEW in Sprint 4** - Unified atomic endpoint that creates Project + Brief data in one call!

This endpoint accepts all 12 fields from the Frontend wizard (9 required + 3 optional) and automatically creates project with all necessary metadata.

### Request
```bash
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
  "editable_copywriting": "Halo generasi masa depan! Bergabunglah dengan keluarga besar kami dan wujudkan impianmu bersama SMA Negeri 1 Jakarta.",
  "editable_hashtags": "#SMANegeri1Jakarta #PenerimaanMahasiswaBaru #Pendidikan #KampusImpian"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Project initialized successfully with briefs",
  "data": {
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "project_name": "Penerimaan Mahasiswa Baru for SMA Negeri 1 Jakarta",
    "institution_name": "SMA Negeri 1 Jakarta",
    "school_level": "Senior High School",
    "theme": "Tur Kampus Sinematik",
    "tone": "Santai & Ramah",
    "duration": 15,
    "status": "draft",
    "created_at": "2026-03-13T10:05:00Z"
  }
}
```

### What Gets Auto-Created
- ✅ **Project** - with name, theme, institution info
- ✅ **Brief Data** - auto-filled from FE input + sensible defaults
- ✅ **Auto-Mappings:**
  - `event_content` → video_type
  - `tone_of_voice` → music_preference
  - `selected_theme` → visual style
  - Auto-defaults: industry="Education", target_audience="Students", format="mp4", resolution="1080p"

### Credits Used
- Cost: 0 credits
- Current Balance: 10 credits

### Save Variables
```
project_id = 550e8400-e29b-41d4-a716-446655440001
```

---

## STEP 4: Generate Storyboard

### Request
```bash
POST /api/storyboard/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Storyboard generated successfully",
  "data": {
    "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "status": "ready",
    "total_duration": 19,
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
    "created_at": "2026-03-13T10:10:00Z"
  }
}
```

### Credits Used
- Cost: 0 credits
- Current Balance: 10 credits

### Save Variables
```
storyboard_id = 5ja7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 5: Generate Video

### Request
```bash
POST /api/videos/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440001",
  "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Video generation job created",
  "data": {
    "generation_job_id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
    "credits_required": 100,
    "created_at": "2026-03-13T10:15:00Z"
  }
}
```

### Credits Used
- Cost: 100 credits per video
- Current Balance: **-90 credits (Insufficient!)**
- **Status:** Job queued but will fail if credits insufficient

### Important: Check Credits Before Generating

If you don't have enough credits, use admin endpoint to add credits:
```bash
POST /api/credits/add
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 200
}
```

### Save Variables
```
generation_job_id = 7ca7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 6: Check Video Status

### Request
```bash
GET /api/videos/7ca7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer {access_token}
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "Video status retrieved",
  "data": {
    "id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "completed",
    "video_url": "http://127.0.0.1:8000/video/output_7ca7b810.mp4",
    "duration": 19,
    "provider": "ltx",
    "model": "ltx-2-3-fast",
    "created_at": "2026-03-13T10:15:00Z",
    "updated_at": "2026-03-13T10:25:00Z"
  }
}
```

### Video Statuses
- `queued` - Waiting in queue
- `processing` - Currently generating
- `completed` - ✅ Ready for download
- `failed` - ❌ Generation failed

### Download Video
```bash
GET {video_url}
# http://127.0.0.1:8000/video/output_7ca7b810.mp4
```

---

## Complete Visual Timeline

```
Timeline (Estimated)
├─ STEP 1: Register (0s)
│  └─ Response: access_token + 10 credits
│
├─ STEP 2: Login (5s)
│  └─ Response: access_token
│
├─ STEP 3: Initialize Project (1-2s)
│  └─ Response: project_id, brief data auto-created
│  └─ Credits: 10 (no cost)
│
├─ STEP 4: Generate Storyboard (2-3s)
│  └─ Response: storyboard_id + 3 scenes
│  └─ Credits: 10 (no cost)
│
└─ STEP 5-6: Generate Video (15-60s depending on provider)
   ├─ Response: generation_job_id
   ├─ Status: queued → processing → completed
   ├─ Credits: 10 - 100 = -90 (FAIL if insufficient)
   └─ Download: video_url ready

Total Time: ~2-5 minutes for complete video generation
Total Credits Used: 100 (video generation only)
```

---

## Key Changes from Previous Version

- ✅ **Unified Initialization** - Single endpoint creates project + briefs atomically
- ✅ **Simplified Workflow** - Removed Content Pillars, Content Themes, manual selections
- ✅ **Auto-Fill** - 18+ fields automatically populated from FE wizard
- ✅ **Cleaner Flow** - 5 steps instead of 10+
- ✅ **FE-aligned** - Workflow matches exactly what FE wizard provides
- ✅ **Updated Credits** - Only charged for video generation, not briefs or storyboard
