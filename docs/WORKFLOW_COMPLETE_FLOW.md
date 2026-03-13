# Complete Workflow Flow - Register to Video Generation

**Complete flow from start (Register) until video ready for download!**

---

## Table of Contents
1. [Step 1: Register User](#step-1-register-user)
2. [Step 2: Login](#step-2-login)
3. [Step 3: Create Project](#step-3-create-project)
4. [Step 4: Create Business Brief](#step-4-create-business-brief)
5. [Step 5: Generate Content Pillars](#step-5-generate-content-pillars)
6. [Step 6: Select Content Pillar & Theme](#step-6-select-content-pillar--theme)
7. [Step 7: Generate Storyboards](#step-7-generate-storyboards)
8. [Step 8: Generate Videos (3 Variants)](#step-8-generate-videos-3-variants)
9. [Step 9: Poll Job Status](#step-9-poll-job-status)
10. [Step 10: Download Videos](#step-10-download-videos)
11. [Timeline](#timeline)

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
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJleHAiOjE3NDEyNzI2MDB9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJleHAiOjE3NDE5Mjc2MDB9...",
    "expires_in": 86400,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Budi Wiranto",
      "email": "budi@company.id",
      "role": "user",
      "created_at": "2026-03-13T10:00:00Z"
    }
  }
}
```

### Key Info to Save
-  `access_token` → Gunakan untuk semua request berikutnya
-  `refresh_token` → Untuk refresh token ketika expired
-  `user.id` → User ID
-  `user.credits` → 10 credits (default new user)

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
      "role": "user"
    }
  }
}
```

### Headers for Next Requests
**Use this header for all subsequent requests:**
```
Authorization: Bearer {access_token}
```

### Save Variables
```
access_token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
user_id = 550e8400-e29b-41d4-a716-446655440000
```

---

## STEP 3: Create Project

### Request
```bash
POST /api/projects
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Kampanye Brand Video Q2 2026",
  "description": "Video promosi produk terbaru untuk Q2 2026"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Project created successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Kampanye Brand Video Q2 2026",
    "description": "Video promosi produk terbaru untuk Q2 2026",
    "status": "draft",
    "created_at": "2026-03-13T10:05:00Z",
    "updated_at": "2026-03-13T10:05:00Z"
  }
}
```

### Save Variable
```
project_id = 6ba7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 4: Create Business Brief

### Request
```bash
POST /api/briefs/business
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "project_name": "Kampanye Brand Video Q2 2026",
  "company_name": "PT Teknologi Indonesia",
  "industry": "Technology",
  "target_audience": "Profesional muda 25-40 tahun, tech-savvy, urban",
  "project_objective": "Meningkatkan brand awareness sebesar 50% dan lead generation",
  "key_message": "Inovasi teknologi untuk masa depan yang lebih cerah",
  "budget": "500000000",
  "timeline": "8 minggu"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Business brief created successfully",
  "data": {
    "id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "project_name": "Kampanye Brand Video Q2 2026",
    "company_name": "PT Teknologi Indonesia",
    "industry": "Technology",
    "target_audience": "Profesional muda 25-40 tahun, tech-savvy, urban",
    "project_objective": "Meningkatkan brand awareness sebesar 50% dan lead generation",
    "key_message": "Inovasi teknologi untuk masa depan yang lebih cerah",
    "budget": "500000000",
    "timeline": "8 minggu",
    "status": "draft",
    "created_at": "2026-03-13T10:10:00Z"
  }
}
```

### Save Variable
```
business_brief_id = 7ca7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 5: Generate Content Pillars

### Request
```bash
POST /api/projects/6ba7b810-9dad-11d1-80b4-00c04fd430c8/content-pillars/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Content pillars generated successfully",
  "data": [
    {
      "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Inovasi Produk",
      "description": "Menghadirkan fitur-fitur revolusioner yang mengubah industri",
      "is_selected": false,
      "content_themes": [
        {
          "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
          "title": "Product Features Showcase",
          "description": "Demo fitur unggulan produk",
          "is_selected": false
        },
        {
          "id": "0ea7b810-9dad-11d1-80b4-00c04fd430c8",
          "title": "Problem Solving",
          "description": "Bagaimana produk menyelesaikan masalah pelanggan",
          "is_selected": false
        }
      ]
    },
    {
      "id": "1fa7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Kepercayaan & Kredibilitas",
      "description": "Membangun kepercayaan melalui testimoni dan track record",
      "is_selected": false,
      "content_themes": [
        {
          "id": "2ga7b810-9dad-11d1-80b4-00c04fd430c8",
          "title": "Customer Testimonials",
          "description": "Kisah sukses pelanggan yang puas",
          "is_selected": false
        }
      ]
    },
    {
      "id": "3ha7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Call to Action",
      "description": "Mendorong pelanggan untuk ambil tindakan",
      "is_selected": false,
      "content_themes": [
        {
          "id": "4ia7b810-9dad-11d1-80b4-00c04fd430c8",
          "title": "Limited Time Offer",
          "description": "Penawaran terbatas waktu untuk urgency",
          "is_selected": false
        }
      ]
    }
  ]
}
```

### Credits Used
- Cost: 1 credit
- User now has: 9 credits

### Save Variables
```
content_pillar_id = 8da7b810-9dad-11d1-80b4-00c04fd430c8  (Pillar 1)
content_theme_id = 9da7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 6: Select Content Pillar & Theme

### Request 6A: Select Content Pillar
```bash
POST /api/content-pillars/8da7b810-9dad-11d1-80b4-00c04fd430c8/select
Authorization: Bearer {access_token}
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "Content pillar selected successfully",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "selected_at": "2026-03-13T10:15:00Z"
  }
}
```

### Request 6B: Select Content Theme
```bash
POST /api/content-themes/9da7b810-9dad-11d1-80b4-00c04fd430c8/select
Authorization: Bearer {access_token}
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "Content theme selected successfully",
  "data": {
    "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "selected_at": "2026-03-13T10:16:00Z"
  }
}
```

---

## STEP 7: Generate Storyboards

### Request
```bash
POST /api/projects/6ba7b810-9dad-11d1-80b4-00c04fd430c8/storyboards/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "content_theme_id": "9da7b810-9dad-11d1-80b4-00c04fd430c8"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Storyboards generated successfully",
  "data": [
    {
      "id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "content_theme_id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Storyboard 1 - Premium Cinematic",
      "description": "Professional product showcase dengan gaya cinematic",
      "version": 1,
      "is_selected": false,
      "scenes": [
        {
          "id": "6ka7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 1,
          "title": "Pembukaan - Introduksi Produk",
          "description": "Wide shot dari kantor modern dengan produk di tengah",
          "duration": 5,
          "visual_style": "cinematic"
        },
        {
          "id": "7la7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 2,
          "title": "Demo Fitur Utama",
          "description": "Close-up product in action, menunjukkan UI dan features",
          "duration": 4,
          "visual_style": "cinematic"
        },
        {
          "id": "8ma7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 3,
          "title": "Call to Action",
          "description": "Website dan CTA button dengan text 'Kunjungi Sekarang'",
          "duration": 3,
          "visual_style": "cinematic"
        }
      ],
      "total_duration": 12,
      "created_at": "2026-03-13T10:20:00Z"
    },
    {
      "id": "9na7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "content_theme_id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Storyboard 2 - Vibrant Dynamic",
      "description": "Energetik dan vibrant dengan banyak transisi",
      "version": 2,
      "is_selected": false,
      "scenes": [...]
    },
    {
      "id": "0oa7b810-9dad-11d1-80b4-00c04fd430c8",
      "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "content_theme_id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Storyboard 3 - Professional Corporate",
      "description": "Corporate style dengan fokus pada kredibilitas",
      "version": 3,
      "is_selected": false,
      "scenes": [...]
    }
  ]
}
```

### Credits Used
- Cost: 1 credit
- User now has: 8 credits

### Save Variables
```
storyboard_id = 5ja7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 8: Generate Videos (3 Variants)

#### Request (User)
```bash
POST /api/videos/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
  "custom_prompt": "Fokus pada inovasi produk, profesional namun modern",
  "scene_count": 3,
  "video_duration": 12
}
```

#### Response (201 Created)
```json
{
  "success": true,
  "message": "Video generation job created",
  "data": {
    "generation_job_id": "2qa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "scene_count": 3,
    "video_duration": 12,
    "total_variants": 3,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "created_at": "2026-03-13T10:30:00Z",
    "estimated_completion": "2026-03-13T11:00:00Z"
  }
}
```

### What Happens Behind the Scenes
```
1. Creates 3 VideoVariants:
   - Variant 1: Cinematic
   - Variant 2: Vibrant
   - Variant 3: Professional

3. Each variant gets 3 SceneGenerations:
   - Scene 1 (5 sec)
   - Scene 2 (4 sec)
   - Scene 3 (3 sec)

2. Job queued → Status: "queued"
   Waiting for worker to pick up

3. Worker processes all scenes asynchronously with AI provider
```

### Save Variables
```
generation_job_id = 2qa7b810-9dad-11d1-80b4-00c04fd430c8
```

---

## STEP 9: Poll Job Status

### Request (Poll Every 5 Seconds)
```bash
GET /api/videos/generation/2qa7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer {access_token}
```

### Response Timeline

#### T=0s (Queued)
```json
{
  "success": true,
  "data": {
    "id": "2qa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "job_type": "generate",
    "priority": 1,
    "scene_count": 3,
    "video_duration": 12,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "created_at": "2026-03-13T10:30:00Z",
    "updated_at": "2026-03-13T10:30:05Z"
  }
}
```

#### T=5s (Processing)
```json
{
  "success": true,
  "data": {
    "id": "2qa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "processing",
    "job_type": "generate",
    "priority": 1,
    "scene_count": 3,
    "video_duration": 12,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "started_at": "2026-03-13T10:30:08Z",
    "created_at": "2026-03-13T10:30:00Z",
    "updated_at": "2026-03-13T10:30:35Z"
  }
}
```

#### T=20s (30% Complete)
```json
{
  "success": true,
  "data": {
    "status": "processing",
    "progress": {
      "total_variants": 3,
      "completed_variants": 0,
      "total_scenes": 9,
      "completed_scenes": 3,
      "percentage": 33
    },
    "details": "Variant 1 (Cinematic): Scene 1-3 processing, Variant 2: queued"
  }
}
```

#### T=60s (80% Complete)
```json
{
  "success": true,
  "data": {
    "status": "processing",
    "progress": {
      "total_variants": 3,
      "completed_variants": 2,
      "total_scenes": 9,
      "completed_scenes": 8,
      "percentage": 88
    },
    "details": "Variant 1 & 2 completed, Variant 3 in progress"
  }
}
```

#### T=90s (COMPLETED)
```json
{
  "success": true,
  "message": "Video generation completed",
  "data": {
    "id": "2qa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "completed",
    "job_type": "generate",
    "priority": 1,
    "scene_count": 3,
    "video_duration": 12,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "started_at": "2026-03-13T10:30:08Z",
    "completed_at": "2026-03-13T10:31:30Z",
    "created_at": "2026-03-13T10:30:00Z",
    "updated_at": "2026-03-13T10:31:30Z",
    "generated_videos": 3
  }
}
```

---

## STEP 10: Get & Download Videos

### Request: Get All Variants for Storyboard
```bash
GET /api/videos/storyboard/5ja7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer {access_token}
```

### Response: 3 Video Variants
```json
{
  "success": true,
  "message": "Video variants retrieved",
  "data": [
    {
      "id": "3ra7b810-9dad-11d1-80b4-00c04fd430c8",
      "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8",
      "variant_number": 1,
      "title": "Cinematic Version",
      "status": "completed",
      "video_url": "https://storage.example.com/videos/3ra7b810-cinematic.mp4",
      "thumbnail_url": "https://storage.example.com/thumbnails/3ra7b810-thumb.jpg",
      "prompt_used": "Fokus pada inovasi produk, profesional namun modern - Cinematic style",
      "duration": 12,
      "provider": "ltx",
      "model": "ltx-2-fast",
      "resolution": "1080p",
      "file_size": 52428800,
      "scenes": [
        {
          "id": "4sa7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 1,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-1-cinematic.mp4",
          "duration": 5,
          "updated_at": "2026-03-13T10:31:10Z"
        },
        {
          "id": "5ta7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 2,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-2-cinematic.mp4",
          "duration": 4,
          "updated_at": "2026-03-13T10:31:18Z"
        },
        {
          "id": "6ua7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 3,
          "status": "completed",
          "video_url": "https://storage.example.com/scenes/scene-3-cinematic.mp4",
          "duration": 3,
          "updated_at": "2026-03-13T10:31:25Z"
        }
      ],
      "created_at": "2026-03-13T10:30:00Z",
      "updated_at": "2026-03-13T10:31:25Z"
    },
    {
      "id": "7va7b810-9dad-11d1-80b4-00c04fd430c8",
      "variant_number": 2,
      "title": "Vibrant Version",
      "status": "completed",
      "video_url": "https://storage.example.com/videos/7va7b810-vibrant.mp4",
      "thumbnail_url": "https://storage.example.com/thumbnails/7va7b810-thumb.jpg",
      "prompt_used": "Fokus pada inovasi produk, profesional namun modern - Vibrant, energetic style",
      "duration": 12,
      "resolution": "1080p",
      "credits_used": 40,
      "scenes": [...]
    },
    {
      "id": "8wa7b810-9dad-11d1-80b4-00c04fd430c8",
      "variant_number": 3,
      "title": "Professional Version",
      "status": "completed",
      "video_url": "https://storage.example.com/videos/8wa7b810-professional.mp4",
      "thumbnail_url": "https://storage.example.com/thumbnails/8wa7b810-thumb.jpg",
      "prompt_used": "Fokus pada inovasi produk, profesional namun modern - Corporate professional",
      "duration": 12,
      "resolution": "1080p",
      "credits_used": 40,
      "scenes": [...]
    }
  ]
}
```

### DOWNLOAD EACH VIDEO

#### Request 1: Download Cinematic Version
```bash
GET /api/videos/3ra7b810-9dad-11d1-80b4-00c04fd430c8/download
Authorization: Bearer {access_token}
```

#### Response
```json
{
  "success": true,
  "message": "Video download ready",
  "data": {
    "variant_id": "3ra7b810-9dad-11d1-80b4-00c04fd430c8",
    "title": "Cinematic Version",
    "download_url": "https://storage.example.com/videos/3ra7b810-cinematic.mp4",
    "file_size": 52428800,
    "format": "mp4",
    "resolution": "1080p",
    "duration": 12,
    "expires_at": "2026-03-13T11:31:30Z"
  }
}
```

Copy the download_url and download the video!

#### Request 2: Download Vibrant Version
```bash
GET /api/videos/7va7b810-9dad-11d1-80b4-00c04fd430c8/download
```

#### Request 3: Download Professional Version
```bash
GET /api/videos/8wa7b810-9dad-11d1-80b4-00c04fd430c8/download
```

---

## Complete Visual Timeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ COMPLETE WORKFLOW TIMELINE - Register hingga Download Video                 │
└─────────────────────────────────────────────────────────────────────────────┘

T=0:00     Register User
           ├─ User ID: 550e8400-...
           ├─ Initial Credits: 10
           └─ Access Token: eyJ...

T=0:05     Login
           ├─ Access Token: eyJ...
           ├─ Refresh Token: eyJ...
           └─ Ready for API calls

T=0:10     Create Project
           ├─ Project ID: 6ba7b810-...
           ├─ Name: Kampanye Brand Video Q2 2026
           └─ Status: draft

T=0:15     Create Business Brief
           ├─ Brief ID: 7ca7b810-...
           ├─ Target Audience: Profesional muda 25-40
           └─ Company: PT Teknologi Indonesia

T=0:20     Generate Content Pillars (1 credit)
           ├─ Credits: 10 to 9
           ├─ Pillars: 3 variants
           │  ├─ Inovasi Produk
           │  ├─ Kepercayaan & Kredibilitas
           │  └─ Call to Action
           └─ Themes: 3-4 per pillar

T=0:30     Select Content Pillar & Theme (0 credits)
           ├─ Selected: Inovasi Produk
           ├─ Theme: Product Features Showcase
           └─ Credits: 9 (no change)

T=0:40     Generate Storyboards (1 credit)
           ├─ Credits: 9 to 8
           ├─ Storyboards: 3 variations
           │  ├─ Cinematic
           │  ├─ Vibrant
           │  └─ Professional
           ├─ Scenes per storyboard: 3
           └─ Total duration per SB: 12 seconds

T=0:50     Admin Top-up Credits (0 credits)
           ├─ Added: 150 credits
           └─ Credits: 8 to 158

T=1:00     Generate Videos (120 credits)
           ├─ Credits: 158 to 38
           ├─ Job started: queued
           ├─ 3 Variants x 3 Scenes = 9 scenes total
           ├─ Provider: LTX-2-Fast
           └─ Estimated time: ~90 seconds

T=1:05     Poll Status (Every 5 sec)
           ├─ Status: queued to processing
           └─ Updated at each poll

T=1:30     [ VIDEO GENERATION IN PROGRESS ]
           ├─ Variant 1 (Cinematic): 50% complete
           ├─ Variant 2 (Vibrant): 20% complete
           ├─ Variant 3 (Professional): queued
           └─ Credits Used: 80/120

T=1:60     [ VIDEOS NEARLY COMPLETE ]
           ├─ Variant 1 & 2: DONE
           ├─ Variant 3: 80% complete
           └─ Credits Used: 110/120

T=1:90     VIDEOS COMPLETE
           ├─ Status: completed
           ├─ 3 Full Videos Ready
           │  ├─ Cinematic (52 MB)
           │  ├─ Vibrant (52 MB)
           │  └─ Professional (52 MB)
           └─ Credits Final: 38

T=2:00     Get All Variants
           ├─ Endpoint: /api/videos/storyboard/{id}
           ├─ Response: 3 videos with URLs
           └─ Each video: 12 sec at 1080p

T=2:05     Download Cinematic
           ├─ URL: https://storage.example.com/videos/...mp4
           ├─ Size: 52 MB
           └─ Duration: 12 seconds

T=2:10     Download Vibrant
           └─ 3 videos successfully downloaded

T=2:15     Download Professional
           └─ Ready for use

T=2:20     [COMPLETE - READY TO USE]
           ├─ 3 video variants downloaded
           ├─ Credits remaining: 38
           └─ Can regenerate or generate more
```

---

## API Steps Summary

| Step | Action | Time |
|------|--------|------|
| 1 | Register User | Instant |
| 2 | Login | Instant |
| 3 | Create Project | Instant |
| 4 | Create Business Brief | Instant |
| 5 | Generate Content Pillars | Instant |
| 6 | Select Pillar & Theme | Instant |
| 7 | Generate Storyboards | Instant |
| 8 | Generate Videos | ~60-90 seconds |
| 9 | Poll Status | Variable |
| 10 | Download Videos | Instant |

---

## Optional: Regenerate Single Scene

Jika ingin edit hanya satu scene:

### Request
```bash
POST /api/videos/scene/4sa7b810-9dad-11d1-80b4-00c04fd430c8/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "new_prompt": "Fokus pada layar produk dengan close-up yang lebih detail"
}
```



---

## Next Steps After Video Download

**Videos ready to use for:**
- Social media (TikTok, Instagram, YouTube)
- Website marketing
- Email campaigns
- Advertisements

**Other options:**
- Regenerate full video with different prompt
- Regenerate individual scene
- Create additional storyboards
- Start new project for different campaign

