# Complete Workflow Flow - Register to Video Generation

Complete flow from start (Register) until video ready for download.

**Note:** Use Bruno API Collection (`docs/API_COLLECTION.json`) for automatic variable population. Variables (access_token, project_id, etc.) auto-set after each Create/Generate/List/Get endpoint—no manual copy-paste needed!

---

## Table of Contents
1. [Step 1: Register User](#step-1-register-user)
2. [Step 2: Login](#step-2-login)
3. [Step 3: Create Project](#step-3-create-project)
4. [Step 4: Create Business Brief](#step-4-create-business-brief)
5. [Step 5: Generate Content Pillars](#step-5-generate-content-pillars)
6. [Step 6: Select Content Pillar & Theme](#step-6-select-content-pillar--theme)
7. [Step 7: Generate Storyboards](#step-7-generate-storyboards)
8. [Step 7.5: Check Credits & Top-up if Needed](#step-75-check-credits--top-up-if-needed)
9. [Step 8: Generate Videos (3 Variants)](#step-8-generate-videos-3-variants)
10. [Step 9: Poll Job Status](#step-9-poll-job-status)
11. [Step 10: Download Videos](#step-10-download-videos)
12. [Timeline](#complete-visual-timeline)

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
- `access_token` - Gunakan untuk semua request berikutnya
- `refresh_token` - Untuk refresh token ketika expired
- `user.id` - User ID
- `user.credits` - 10 credits (default new user)

### Credits Tracking
**Starting Balance: 10 credits**
```
Initial Credits: 10
Available for use: 10
```

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
Use this header for all subsequent requests:
```
Authorization: Bearer {access_token}
```

### Save Variables
```
access_token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
user_id = 550e8400-e29b-41d4-a716-446655440000
credits_balance = 10
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
  "description": "Video promosi produk terbaru untuk Q2 2026",
  "theme": "Corporate Branding"
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
    "theme": "Corporate Branding",
    "status": "draft",
    "created_at": "2026-03-13T10:05:00Z",
    "updated_at": "2026-03-13T10:05:00Z"
  }
}
```

### Credits Used
- Cost: 0 credits
- Current Balance: 10 credits

### Save Variable
```
project_id = 6ba7b810-9dad-11d1-80b4-00c04fd430c8
```

**✓ NEW in Sprint 3:** Theme field added for project theme

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
  "institute_name": "PT Teknologi Indonesia",
  "education": "Higher Education",
  "target_audience": true,
  "project_objective": "Meningkatkan brand awareness sebesar 50% dan lead generation",
  "key_message": "Inovasi teknologi untuk masa depan yang lebih cerah",
  "deadline": "2026-05-15T00:00:00Z",
  "competitors": "Kompetitor lainnya",
  "additional_notes": "Catatan tambahan"
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
    "institute_name": "PT Teknologi Indonesia",
    "education": "Higher Education",
    "target_audience": true,
    "project_objective": "Meningkatkan brand awareness sebesar 50% dan lead generation",
    "key_message": "Inovasi teknologi untuk masa depan yang lebih cerah",
    "deadline": "2026-05-15T00:00:00Z",
    "competitors": "Kompetitor lainnya",
    "additional_notes": "Catatan tambahan",
    "status": "draft",
    "created_at": "2026-03-13T10:10:00Z",
    "updated_at": "2026-03-13T10:10:00Z"
  }
}
```

### Field Updates (Sprint 3)
- ✓ `company_name` changed to `institute_name`
- ✓ `industry` changed to `education`
- ✓ `target_audience` is now boolean (true/false instead of string)
- ✓ `budget` + `timeline` merged into `deadline` (time.Time format: ISO 8601)
    "key_message": "Inovasi teknologi untuk masa depan yang lebih cerah",
    "budget": "500000000",
    "timeline": "8 minggu",
    "status": "draft",
    "created_at": "2026-03-13T10:10:00Z",
    "updated_at": "2026-03-13T10:10:00Z"
  }
}
```

### Credits Used
- Cost: 0 credits
- Current Balance: 10 credits

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
      "prompt": "Custom AI prompt untuk video generation (optional)",
      "video_url": "https://..../generated_video.mp4",
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
      ],
      "created_at": "2026-03-13T10:15:00Z",
      "updated_at": "2026-03-13T10:15:00Z"
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
      ],
      "created_at": "2026-03-13T10:15:00Z",
      "updated_at": "2026-03-13T10:15:00Z"
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
      ],
      "created_at": "2026-03-13T10:15:00Z",
      "updated_at": "2026-03-13T10:15:00Z"
    }
  ]
}
```

### Credits Used
- Cost: 0 credits
- Current Balance: 10 credits

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
  "message": "Content pillar selected",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "updated_at": "2026-03-13T10:16:00Z"
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
  "message": "Content theme selected",
  "data": {
    "id": "9da7b810-9dad-11d1-80b4-00c04fd430c8",
    "is_selected": true,
    "updated_at": "2026-03-13T10:17:00Z"
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
      "prompt": "Custom AI prompt untuk caption generation (optional)",
      "is_selected": false,
      "scenes": [
        {
          "id": "6ka7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 1,
          "title": "Pembukaan - Introduksi Produk",
          "description": "Wide shot dari kantor modern dengan produk di tengah",
          "duration": 5,
          "visual_style": "cinematic",
          "caption": "Memperkenalkan solusi inovatif untuk kebutuhan bisnis modern",
          "regenerate_count": 0
        },
        {
          "id": "7la7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 2,
          "title": "Demo Fitur Utama",
          "description": "Close-up product in action, menunjukkan UI dan features",
          "duration": 4,
          "visual_style": "cinematic",
          "caption": "Fitur-fitur unggulan dirancang khusus untuk kemudahan penggunaan",
          "regenerate_count": 0
        },
        {
          "id": "8ma7b810-9dad-11d1-80b4-00c04fd430c8",
          "scene_number": 3,
          "title": "Call to Action",
          "description": "Website dan CTA button dengan text 'Kunjungi Sekarang'",
          "duration": 3,
          "visual_style": "cinematic",
          "caption": "Bergabunglah dengan ribuan pengguna puas hari ini",
          "regenerate_count": 0
        }
      ],
      "total_duration": 12,
      "created_at": "2026-03-13T10:20:00Z",
      "updated_at": "2026-03-13T10:20:00Z"
    },
    {
      "id": "9na7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Storyboard 2 - Vibrant Dynamic",
      "description": "Energetik dan vibrant dengan banyak transisi",
      "is_selected": false,
      "scenes": [...]
    },
    {
      "id": "0oa7b810-9dad-11d1-80b4-00c04fd430c8",
      "title": "Storyboard 3 - Professional Corporate",
      "description": "Corporate style dengan fokus pada kredibilitas",
      "is_selected": false,
      "scenes": [...]
    }
  ]
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

## STEP 7.5: Check Credits & Top-up if Needed

### Current Status
```
Current Balance: 10 credits
Required for Video Generation (3 variants): 120 credits
Deficit: -110 credits
status: INSUFFICIENT CREDITS
```

### Top-up Credits - Admin Request

#### Request
```bash
POST /api/admin/credits
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 150
}
```

#### Response (200 OK)
```json
{
  "success": true,
  "message": "Credits added successfully",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "credits_added": 150,
    "total_credits": 160,
    "created_at": "2026-03-13T10:28:00Z"
  }
}
```

### Credits After Top-up
```
Previous Balance: 10 credits
Top-up Amount: +150 credits
New Balance: 160 credits
Ready for: Video Generation
```

### Save Variable
```
credits_balance = 160
```

---

## STEP 8: Generate Videos (3 Variants)

#### Pre-Flight Check
```
Current Credits: 160
Video Generation Cost: 120 credits (40 per variant x 3)
After Generation: 40 credits
Check: APPROVED - Sufficient credits
```

#### Request (User)
```bash
POST /api/videos/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "storyboard_id": "5ja7b810-9dad-11d1-80b4-00c04fd430c8"
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
    "credits_to_deduct": 120,
    "created_at": "2026-03-13T10:30:00Z",
    "estimated_completion": "2026-03-13T11:00:00Z"
  }
}
```

### Credits Deducted
```
Previous Balance: 160 credits
Cost Breakdown:
  - Variant 1 (Cinematic): 40 credits
  - Variant 2 (Vibrant): 40 credits
  - Variant 3 (Professional): 40 credits
Total Deducted: -120 credits

Current Balance: 40 credits
```

### What Happens Behind the Scenes
```
1. Creates 3 VideoVariants:
   - Variant 1: Cinematic
   - Variant 2: Vibrant
   - Variant 3: Professional

2. Each variant gets 3 SceneGenerations:
   - Scene 1 (5 sec)
   - Scene 2 (4 sec)
   - Scene 3 (3 sec)

3. Job queued -> Status: "queued"
   Waiting for worker to pick up

4. Worker processes all scenes asynchronously with AI provider
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

Copy the download_url and download the video.

#### Request 2: Download Vibrant Version
```bash
GET /api/videos/7va7b810-9dad-11d1-80b4-00c04fd430c8/download
Authorization: Bearer {access_token}
```

#### Request 3: Download Professional Version
```bash
GET /api/videos/8wa7b810-9dad-11d1-80b4-00c04fd430c8/download
Authorization: Bearer {access_token}
```

---

## Complete Visual Timeline

```
+-----------------------------------------------------------------------------+
| COMPLETE WORKFLOW TIMELINE - Register hingga Download Video                 |
+-----------------------------------------------------------------------------+

T=0:00     Register User
           |- User ID: 550e8400-...
           |- Initial Credits: 10
           +- Access Token: eyJ...

T=0:05     Login
           |- Access Token: eyJ...
           |- Refresh Token: eyJ...
           +- Ready for API calls

T=0:10     Create Project
           |- Project ID: 6ba7b810-...
           |- Name: Kampanye Brand Video Q2 2026
           +- Status: draft

T=0:15     Create Business Brief
           |- Brief ID: 7ca7b810-...
           |- Target Audience: Profesional muda 25-40
           +- Company: PT Teknologi Indonesia

T=0:20     Generate Content Pillars (0 credits)
           |- Starting: 10 credits
           |- Cost: 0 credits
           |- Balance: 10 credits
           |- Pillars: 3 variants
           |  |- Inovasi Produk
           |  |- Kepercayaan & Kredibilitas
           |  +- Call to Action
           +- Themes: 3-4 per pillar

T=0:30     Select Content Pillar & Theme (0 credits)
           |- Starting: 10 credits
           |- Cost: 0 credits
           |- Balance: 10 credits
           |- Selected: Inovasi Produk
           +- Theme: Product Features Showcase

T=0:40     Generate Storyboards (0 credits)
           |- Starting: 10 credits
           |- Cost: 0 credits
           |- Balance: 10 credits
           |- Storyboards: 3 variations
           |  |- Cinematic
           |  |- Vibrant
           |  +- Professional
           |- Scenes per storyboard: 3
           +- Total duration per SB: 12 seconds

T=0:50     Check & Top-up Credits (0 credits)
           |- Previous: 10 credits
           |- Required: 120 credits
           |- Top-up: +150 credits
           +- New Balance: 160 credits

T=1:00     Generate Videos (120 credits)
           |- Starting: 160 credits
           |- Cost: 40 credits/variant x 3 = 120 credits
           |- Balance: 40 credits
           |- Job started: queued
           |- 3 Variants x 3 Scenes = 9 scenes total
           |- Provider: LTX-2-Fast
           +- Estimated time: ~90 seconds

T=1:05     Poll Status (Every 5 sec)
           |- Status: queued to processing
           +- Updated at each poll

T=1:30     [ VIDEO GENERATION IN PROGRESS ]
           |- Variant 1 (Cinematic): 50% complete
           |- Variant 2 (Vibrant): 20% complete
           |- Variant 3 (Professional): queued
           +- Credits Allocated: 120/120

T=1:60     [ VIDEOS NEARLY COMPLETE ]
           |- Variant 1 & 2: DONE
           |- Variant 3: 80% complete
           |- Credits Allocated: 120/120
           +- Remaining Balance: 40 credits (after generation)

T=1:90     VIDEOS COMPLETE
           |- Status: completed
           |- 3 Full Videos Ready
           |  |- Cinematic (52 MB) - 40 credits
           |  |- Vibrant (52 MB) - 40 credits
           |  +- Professional (52 MB) - 40 credits
           |- Total Used: 120 credits
           +- Current Balance: 40 credits

T=2:00     Get All Variants
           |- Endpoint: /api/videos/storyboard/{id}
           |- Response: 3 videos with URLs
           +- Each video: 12 sec at 1080p

T=2:05     Download Cinematic
           |- URL: https://storage.example.com/videos/...mp4
           |- Size: 52 MB
           +- Duration: 12 seconds

T=2:10     Download Vibrant
           +- 3 videos successfully downloaded

T=2:15     Download Professional
           +- Ready for use

T=2:20     [COMPLETE - READY TO USE]
           |- 3 video variants downloaded
           |- Credits remaining: 40
           +- Can regenerate or generate more
```

---

## API Steps Summary

| Step | Action | Credit Cost | Time |
|------|--------|-------------|------|
| 1 | Register User | 0 | Instant |
| 2 | Login | 0 | Instant |
| 3 | Create Project | 0 | Instant |
| 4 | Create Business Brief | 0 | Instant |
| 5 | Generate Content Pillars | 0 | Instant |
| 6 | Select Pillar & Theme | 0 | Instant |
| 7 | Generate Storyboards | 0 | Instant |
| 7.5 | Top-up Credits (if needed) | - | Instant |
| 8 | Generate Videos | 120 | ~60-90 seconds |
| 9 | Poll Status | 0 | Variable |
| 10 | Download Videos | 0 | Instant |

---

## Step 11: Regenerate Video with Limit Tracking (SPRINT 3 ✓ NEW)

### Video Regenerate Limit Rules
- **✓ NEW in Sprint 3:** RegenerateCount field added to Video model
- **Max Regenerate:** 3 times per video (after regenerate_count reaches 3, no more regeneration allowed)
- **Cost:** Same as original (40 credits per variant)
- **Use Case:** Update style, fix quality issues, apply different AI provider

### Request: Regenerate Full Video
```bash
POST /api/videos/3ra7b810-9dad-11d1-80b4-00c04fd430c8/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "provider": "Runway-Gen3",
  "regenerate_all_scenes": true,
  "new_prompt": "Lebih vibrant dengan warna-warna cerah dan energi tinggi"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Video regeneration job created",
  "data": {
    "generation_job_id": "11ga7b810-new-job-id-12345678",
    "status": "queued",
    "regenerate_count": 1,
    "max_regenerate": 3,
    "remaining_attempts": 2,
    "credit_cost": 40,
    "current_balance": 0
  }
}
```

### Status Tracking (Regenerate)
```bash
GET /api/videos/3ra7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer {access_token}
```

```json
{
  "success": true,
  "data": {
    "id": "3ra7b810-9dad-11d1-80b4-00c04fd430c8",
    "title": "Cinematic Version v2",
    "status": "completed",
    "regenerate_count": 1,
    "max_regenerate": 3,
    "video_url": "https://storage.example.com/videos/3ra7b810-cinematic-v2.mp4",
    "variant_number": 1,
    "updated_at": "2026-03-13T11:00:00Z",
    "scenes": [
      {
        "id": "4sa7b810-9dad-11d1-80b4-00c04fd430c8",
        "scene_number": 1,
        "status": "completed",
        "caption": "Tampilan produk dengan cahaya sinematik",
        "regenerate_count": 1,
        "max_regenerate": 3,
        "updated_at": "2026-03-13T11:00:30Z"
      }
    ]
  }
}
```

---

## Step 12: Regenerate Single Scene with Limit Tracking (SPRINT 3 ✓ NEW)

### Scene Regenerate Limit Rules
- **✓ NEW in Sprint 3:** RegenerateCount & Caption fields added to Scene model
- **Max Regenerate:** 3 times per scene (independent from video limit)
- **Cost:** 15 credits per scene
- **Use Case:** Fix specific scene, update caption, improve specific part

### Request: Regenerate Scene Only
```bash
POST /api/videos/scene/4sa7b810-9dad-11d1-80b4-00c04fd430c8/regenerate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "new_prompt": "Fokus pada fitur produk dengan close-up detail dan pencahayaan lebih terang",
  "update_caption": true,
  "provider": "LTX-2-Fast"
}
```

### Response (201 Created) - Scene Regenerate
```json
{
  "success": true,
  "message": "Scene regeneration job created",
  "data": {
    "generation_job_id": "12ha7b810-new-scene-regen-56789",
    "scene_id": "4sa7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "regenerate_count": 1,
    "max_regenerate": 3,
    "remaining_attempts": 2,
    "credit_cost": 15,
    "current_balance": 25,
    "estimated_time": "30 seconds"
  }
}
```

### Scene After Regeneration
```bash
GET /api/storyboards/2ba7b810-9dad-11d1-80b4-00c04fd430c8
Authorization: Bearer {access_token}
```

```json
{
  "success": true,
  "data": {
    "id": "2ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "scenes": [
      {
        "id": "4sa7b810-9dad-11d1-80b4-00c04fd430c8",
        "scene_number": 1,
        "status": "completed",
        "caption": "Tampilan close-up produk dengan detail fitur utama",
        "regenerate_count": 1,
        "max_regenerate": 3,
        "video_url": "https://storage.example.com/scenes/scene-1-cinematic-v2.mp4",
        "duration": 5,
        "updated_at": "2026-03-13T11:05:00Z"
      }
    ]
  }
}
```

### Regenerate Limit Exceeded
```json
{
  "success": false,
  "error": "REGENERATE_LIMIT_EXCEEDED",
  "message": "Scene has reached maximum regenerate attempts (3/3). Cannot regenerate further.",
  "data": {
    "scene_id": "4sa7b810-9dad-11d1-80b4-00c04fd430c8",
    "regenerate_count": 3,
    "max_regenerate": 3
  }
}
```

---

## Regenerate Tracking Summary

| Entity | Field | Current | Max | Increment On |
|--------|-------|---------|-----|--------------|
| Video | `regenerate_count` | 0-3 | 3 | Each full video regenerate |
| Scene | `regenerate_count` | 0-3 | 3 | Each scene regenerate |
| VideoVariant | `regenerate_count` | 0-3 | 3 | Each variant regenerate |

**✓ NEW in Sprint 3:** All regenerate counts now tracked independently

---

## Complete Regenerate Credit Costs

| Action | Credit Cost | Limit | Notes |
|--------|-------------|-------|-------|
| Generate Video (full) | 40 | Unlimited | Initial generation |
| Regenerate Video (full) | 40 | 3x | Entire video with all scenes |
| Regenerate Scene | 15 | 3x per scene | Single scene in storyboard |
| Change Variant Style | 40 | 3x per variant | Different AI provider/prompt |

**Total Credits Available:** 10 (new user) + top-ups
**Smart Decision:** Choose top-up carefully based on regenerate needs

---

## Error Handling: Regenerate Scenarios

### Scenario 1: Insufficient Credits
```json
{
  "success": false,
  "error": "INSUFFICIENT_CREDITS",
  "message": "Not enough credits. Required: 40, Available: 25",
  "data": {
    "required_credits": 40,
    "available_credits": 25,
    "shortfall": 15
  }
}
```

### Scenario 2: Multiple Regenerate Attempts Exceeded
```json
{
  "success": false,
  "error": "REGENERATE_LIMIT_EXCEEDED",
  "message": "Cannot regenerate. Maximum attempts reached (3/3)",
  "data": {
    "entity_id": "video-id-here",
    "regenerate_count": 3,
    "max_regenerate": 3,
    "solutions": ["Create new video variant", "Start new project"]
  }
}
```

### Scenario 3: Video Still Processing
```json
{
  "success": false,
  "error": "VIDEO_PROCESSING",
  "message": "Video is currently being regenerated. Try again in 60 seconds.",
  "data": {
    "current_status": "processing",
    "estimated_completion": "2026-03-13T11:35:00Z"
  }
}
```

---

## Recommended Workflow with Regenerate

```
START
  │
  ├─→ Generate Initial Video (40 credits)
  │     │
  │     ├─→ Not satisfied?
  │     │     ├─→ Regenerate Full Video (40 credits, regenerate_count=1)
  │     │     │
  │     │     └─→ Still not satisfied?
  │     │           ├─→ Regenerate Specific Scene (15 credits, regenerate_count=1)
  │     │           │
  │     │           └─→ Repeat until satisfied (max 3 per scene/video)
  │     │
  │     └─→ Satisfied!
  │           ├─→ Download Video
  │           └─→ Ready for use
  │
  └─→ Create New Project (For different campaign)
```

---

## Next Steps After Video Download

Videos ready to use for:
- Social media (TikTok, Instagram, YouTube)
- Website marketing
- Email campaigns
- Advertisements

Other options:
- **✓ NEW:** Regenerate with max 3 attempts per video/scene
- Create additional storyboards
- Start new project for different campaign
- Top-up credits for more generation/regenerate cycles
