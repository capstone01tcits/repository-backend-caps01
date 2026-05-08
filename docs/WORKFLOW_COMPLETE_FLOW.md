# Complete Workflow - Register to Video Generation

Updated: May 2026 (Automated Storyboard Generation)
Complete flow from Registration through Video Generation

Note: Use Bruno/Postman API Collection (docs/API_COLLECTION.json) for automatic variable population.
Variables (access_token, project_id, etc.) auto-set after each endpoint.

---

## Table of Contents

1. [Step 1: Register User](#step-1-register-user)
2. [Step 2: Login](#step-2-login)
3. [Step 3: Initialize Project](#step-3-initialize-project)
4. [Step 4: Generate Storyboard](#step-4-generate-storyboard)
5. [Step 5: Generate Video](#step-5-generate-video)
6. [Step 6: Check Video Status](#step-6-check-video-status)
7. [Step 7: Download Video](#step-7-download-video)
8. [Complete Timeline](#complete-visual-timeline)

---

## STEP 1: Register User

Create a new user account and receive initial JWT tokens and 10 free credits.

Request:
POST /api/auth/register
Content-Type: application/json

{
  "name": "Budi Wiranto",
  "email": "budi@example.id",
  "password": "secure123456"
}

Response (201 Created):
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
      "email": "budi@example.id",
      "role": "user",
      "credits": 10,
      "created_at": "2026-03-13T10:00:00Z"
    }
  }
}

Key Info to Save:
- access_token: Use for all subsequent requests
- user.id: User ID
- user.credits: Starting with 10 credits

---

## STEP 2: Login

Authenticate existing user and receive new JWT tokens.

Request:
POST /api/auth/login
Content-Type: application/json

{
  "email": "budi@example.id",
  "password": "secure123456"
}

Response (200 OK):
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
      "email": "budi@example.id",
      "role": "user",
      "credits": 10
    }
  }
}

Headers for Next Requests:
Authorization: Bearer {access_token}

---

## STEP 3: Initialize Project

Updated in April 2026: Unified atomic endpoint that creates Project, BusinessBrief, and CreativeBrief in one call.

This endpoint accepts all fields from the frontend wizard form and handles image uploads as base64 strings.

Request:
POST /api/projects/initialize
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "institution_name": "SMA Negeri 1 Jakarta",
  "institution_history": "Sekolah terkemuka dengan program pendidikan berkualitas tinggi sejak 1950",
  "school_level": "SMA",
  "offered_degrees": "IPA, IPS, Bahasa",
  "event_content": "Penerimaan Siswa Baru",
  "tone_of_voice": "Santai & Ramah",
  "selected_key_message": "Bergabunglah dengan komunitas pembelajar yang dinamis",
  "video_duration": "15 detik",
  "prompt": "Video yang engaging dan inspiring untuk Gen Z",
  "selected_theme": "Tur Kampus Sinematik",
  "editable_copywriting": "Halo calon siswa baru! SMA Negeri 1 Jakarta menghadirkan pengalaman belajar yang berbeda.",
  "editable_hashtags": "#SMANegeri1Jakarta #PenerimaanSiswaBaru #Sekolah",
  "logo_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "env_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFAABAAAAAAAAAAAAAAAAAAAAA/8AAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCwAA8A/9k=",
  "document_base64": ""
}

Response (201 Created):
{
  "success": true,
  "message": "Project created successfully with briefs",
  "data": {
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "business_brief_id": "550e8400-e29b-41d4-a716-446655440002",
    "creative_brief_id": "550e8400-e29b-41d4-a716-446655440003",
    "project_name": "Penerimaan Siswa Baru for SMA Negeri 1 Jakarta",
    "theme": "Tur Kampus Sinematik",
    "tone": "Santai & Ramah",
    "duration": 15,
    "school_level": "SMA",
    "institution_name": "SMA Negeri 1 Jakarta",
    "key_message": "Bergabunglah dengan komunitas pembelajar yang dinamis",
    "copywriting": "Halo calon siswa baru! SMA Negeri 1 Jakarta menghadirkan pengalaman belajar yang berbeda.",
    "hashtags": "#SMANegeri1Jakarta #PenerimaanSiswaBaru #Sekolah",
    "created_at": "2026-03-13T10:15:00Z"
  }
}

Key Info to Save:
- project_id: Use for next steps
- business_brief_id: Reference to business context
- creative_brief_id: Reference to creative direction
- Project now has all briefs and images automatically created

---

## STEP 4: Generate Storyboard

Backend automatically generates storyboard by combining project data with pre-prepared content templates.

Request:
POST /api/storyboard/templates/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440001",
  "video_duration": 30
}

Response (201 Created):
{
  "success": true,
  "message": "Templates generated successfully",
  "data": {
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "video_duration": 30,
    "templates": [
      {
        "template_id": "dynamic_template",
        "style": "Dynamic",
        "description": "Fast-paced, attention-grabbing with motion and energy. Ideal for attracting Gen Z students.",
        "duration": 30,
        "sections": [
          {
            "section_type": "hook",
            "title": "Hook - Grab Attention",
            "suggested_duration": 10,
            "content": "START WITH IMPACT: Capture attention in first 10s with striking campus visuals.\n\nOpening ideas:\n- Student quote: \"I came to study, I found my future\"\n- Question: \"What if your college changed everything?\"\n- Statement: \"Universitas Indonesia is reshaping education\"\n\nVisuals: Campus landmarks, students collaborating, quick campus cuts, aerial shots",
            "tips": "Use fast transitions, trending audio, relatable student voices"
          },
          {
            "section_type": "value",
            "title": "Value - Highlight Excellence",
            "suggested_duration": 10,
            "content": "SHOWCASE INSTITUTIONAL STRENGTHS:\n- Ranked #1 Engineering in Southeast Asia\n- 95% employment rate within 6 months\n- 50+ active student clubs\n- State-of-the-art research labs\n- Industry partnerships: Google, Microsoft, Gojek\n- Alumni: 2000+ at Fortune 500 companies",
            "tips": "Use concrete data and visual proof"
          },
          {
            "section_type": "cta",
            "title": "CTA - Drive Enrollment",
            "suggested_duration": 10,
            "content": "DRIVE ENROLLMENT:\n- \"Apply now at ui-admissions.ac.id\"\n- \"Application deadline: March 31, 2026\"\n- \"Early bird scholarships available\"\n- \"Join 50,000+ students at Universitas Indonesia\"\n\nDisplay: Website link, QR code, phone number, campus image with logo",
            "tips": "Make action obvious. Include deadline to create urgency"
          }
        ]
      },
      {
        "template_id": "narrative_template",
        "style": "Narrative",
        "description": "Story-driven approach following student journey from admission through graduation and career success.",
        "duration": 30,
        "sections": [
          {
            "section_type": "hook",
            "title": "Hook - Student Story Begins",
            "suggested_duration": 10,
            "content": "ESTABLISH STUDENT JOURNEY:\n- \"Meet Rini - she came with a dream of becoming an engineer\"\n- \"Three years ago, Budi stood where you are now\"\n- \"This is the story of how Universitas Indonesia changed lives\"\n\nShow diverse characters: first-generation student, career-changer, athlete, international student",
            "tips": "Make it personal and relatable. Show diverse backgrounds"
          },
          {
            "section_type": "value",
            "title": "Value - The Transformation",
            "suggested_duration": 10,
            "content": "SHOW THE TRANSFORMATION:\n- First day in class with inspiring professors\n- Collaborative project building smart systems\n- Research at national conference\n- Campus clubs and lifelong friendships\n- Real internship at tech startup\n- Student gaining confidence and expertise\n\nHighlight: Faculty mentorship, hands-on learning, industry partnerships, campus community",
            "tips": "Show genuine transformation. Use real testimonials"
          },
          {
            "section_type": "cta",
            "title": "CTA - Your Turn Starts Now",
            "suggested_duration": 10,
            "content": "INSPIRE FUTURE STUDENTS:\n- \"Rini graduated and was hired by Google\"\n- \"95% placement rate within 6 months\"\n- \"Alumni network: 200,000+ professionals globally\"\n- \"Your story at Universitas Indonesia starts here\"\n\nCTA: Apply now at ui-admissions.ac.id | Campus tour: tours.ui.ac.id | Deadline: March 31",
            "tips": "End on high note showing success and achievement"
          }
        ]
      },
      {
        "template_id": "energetic_template",
        "style": "Energetic",
        "description": "High-energy, youth-focused showcasing campus culture, student life, and social vibrancy."
      },
      {
        "template_id": "minimalist_template",
        "style": "Minimalist",
        "description": "Clean, professional emphasizing institutional prestige and academic excellence."
      }
    ],
    "count": 4
  }
}

Key Info to Save:
- templates: Review 4 style options with AI-suggested content specific to your institution
- Each template has complete hook/value/cta with realistic examples
- Choose the template that best matches your vision and target audience
- Note the template_id to reference if you want to use it for video generation

---

### Option B: Create Manual Storyboard (Full Control)

Compose your own storyboard with complete control over content and messaging.

Request:
POST /api/storyboard/create
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440001",
  "title": "Engineering Excellence Campaign - 2026",
  "description": "Target freshman students interested in engineering programs",
  "total_duration": 30,
  "style": "Dynamic",
  "sections": [
    {
      "section_type": "hook",
      "content": "Open with student testimonial: 'I came to study engineering, I found my future' - cut to 3-second campus skyline with dynamic music, quick cuts of students working on robots in labs and prototyping",
      "duration": 10
    },
    {
      "section_type": "value",
      "content": "Showcase: state-of-the-art engineering labs with 3D printers and IoT equipment, student project showcase (smart city project won national award), industry partnerships: Google, Gojek, Traveloka plaques visible. Graphics overlay: '#1 Ranked Engineering Program in SE Asia', '2000+ alumni at Fortune 500 companies', '95% employment rate', 'Full scholarships available for top students', '50+ research clubs'",
      "duration": 10
    },
    {
      "section_type": "cta",
      "content": "Text overlay: 'Transform Your Future in Engineering' with bold font. Voice-over: 'Join 5000+ engineering students at Universitas Indonesia. Early application closes March 31, 2026. Full scholarships for top 100 students. Apply now at ui-admissions.ac.id or call +62-21-7886-1234. Open house: Feb 20, 27, Mar 6'",
      "duration": 10
    }
  ]
}

Response (201 Created):
{
  "success": true,
  "message": "Storyboard created successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "Engineering Excellence Campaign - 2026",
    "description": "Target freshman students interested in engineering programs",
    "total_duration": 30,
    "style": "Dynamic",
    "is_selected": false,
    "created_at": "2026-05-08T10:25:00Z",
    "sections": [
      {
        "id": "section-1",
        "section_type": "hook",
        "content": "Open with student testimonial: 'I came to study engineering, I found my future' - cut to 3-second campus skyline with dynamic music, quick cuts of students working on robots in labs",
        "duration": 10
      },
      {
        "id": "section-2",
        "section_type": "value",
        "content": "Showcase: state-of-the-art engineering labs with 3D printers and IoT equipment, student project showcase (smart city project won national award), industry partnerships: Google, Gojek, Traveloka. Graphics: '#1 Ranked Engineering Program in SE Asia', '95% employment rate', 'Full scholarships for top 100 students'",
        "duration": 10
      },
      {
        "id": "section-3",
        "section_type": "cta",
        "content": "Text overlay: 'Transform Your Future in Engineering'. CTA: 'Apply now at ui-admissions.ac.id | Call +62-21-7886-1234 | Deadline: March 31, 2026 | Full scholarships available | Open house: Feb 20, 27, Mar 6'",
        "duration": 10
      }
    ]
  }
}

Key Info to Save:
- storyboard_id: "6ba7b810-9dad-11d1-80b4-00c04fd430c8" (Use this for video generation)
- All 3 sections ready for AI video generation
- Your custom messaging and branding will be used in the final video

---

### Which Option Should You Choose?

**Choose Auto-Generate Templates (Option A) if:**
- You want AI-assisted suggestions based on your project data
- You need inspiration for structuring your video
- You want to compare 4 different professional styles
- You want to save time on content planning

**Choose Manual Storyboard (Option B) if:**
- You have a specific creative vision already defined
- You want complete control over every message and detail
- You have existing scripts or content requirements
---

## STEP 4: Generate Storyboard

Backend automatically generates storyboard by combining project data with pre-prepared content templates.

The system uses your project information (institution name, tone, key message, etc.) and automatically fills in pre-prepared educational content structures.

Request:
POST /api/storyboard/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440001"
}

Response (201 Created):
{
  "success": true,
  "message": "Storyboard generated successfully",
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "project_id": "550e8400-e29b-41d4-a716-446655440001",
    "institution_name": "Universitas Indonesia",
    "total_duration": 30,
    "created_at": "2026-05-08T10:25:00Z",
    "sections": [
      {
        "id": "section-1",
        "section_type": "hook",
        "title": "Hook - Grab Attention",
        "content": "START WITH IMPACT: Capture attention in first 10s with striking campus visuals. Student testimonial: 'I came to study, I found my future'. Question: 'What if your college changed everything?'. Visuals: Campus landmarks, students collaborating, quick cuts, aerial shots.",
        "duration": 10
      },
      {
        "id": "section-2",
        "section_type": "value",
        "title": "Value - Institutional Excellence",
        "content": "SHOWCASE STRENGTHS: Ranked #1 Engineering in Southeast Asia. 95% employment rate. 50+ student clubs. State-of-the-art research labs. Industry partnerships: Google, Microsoft, Gojek. Alumni: 2000+ at Fortune 500 companies. Scholarship opportunities available.",
        "duration": 10
      },
      {
        "id": "section-3",
        "section_type": "cta",
        "title": "CTA - Call to Action",
        "content": "DRIVE ENROLLMENT: Apply now at ui-admissions.ac.id. Application deadline: March 31, 2026. Early bird scholarships available. Join 50,000+ students. Contact: +62-21-7886-1234. Open house: Feb 20, 27, Mar 6.",
        "duration": 10
      }
    ]
  }
}

Key Info to Save:
- storyboard_id: "6ba7b810-9dad-11d1-80b4-00c04fd430c8" (Use this for video generation)
- All 3 sections auto-generated and ready for video creation
- Content is automatically combined with your project data

---

## STEP 5: Generate Video

Initiate video generation from the auto-generated storyboard.

Request:
POST /api/videos/generate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "project_id": "550e8400-e29b-41d4-a716-446655440001",
  "storyboard_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}

Response (201 Created):
{
  "success": true,
  "message": "Video generation job created",
  "data": {
    "job_id": "7ca7b810-9dad-11d1-80b4-00c04fd430c8",
    "status": "queued",
    "created_at": "2026-03-13T10:25:00Z"
  }
}

Expected Processing Time:
- Queued: 0-2 minutes (waiting for available worker)
- Processing: 30-60 seconds per scene (3 scenes = 1.5-3 minutes)
- Total: 2-5 minutes from job creation to completion

Credits Used:
- Video generation: 1 credit per job
- Remaining after video generation: 9 credits

---

## STEP 6: Check Video Status

Poll the video endpoint to check generation status. Repeat until status is "completed".

Request:
GET /api/videos/{video_id}
Authorization: Bearer {access_token}

Response (200 OK - Processing):
{
  "success": true,
  "message": "Video retrieved",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "variant_number": 1,
    "status": "processing",
    "video_url": null,
    "progress": "Generating scene 2 of 3",
    "created_at": "2026-03-13T10:25:00Z",
    "updated_at": "2026-03-13T10:27:30Z"
  }
}

Response (200 OK - Completed):
{
  "success": true,
  "message": "Video retrieved",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "variant_number": 1,
    "status": "completed",
    "video_url": "https://storage.example.com/videos/8da7b810.mp4",
    "thumbnail_url": "https://storage.example.com/thumbnails/8da7b810.jpg",
    "duration": 15,
    "provider": "ltx",
    "model": "ltx-2-fast",
    "resolution": "1920x1080",
    "file_size": 51234567,
    "created_at": "2026-03-13T10:25:00Z",
    "updated_at": "2026-03-13T10:30:15Z",
    "scenes": [
      {
        "scene_number": 1,
        "status": "completed",
        "video_url": "https://storage.example.com/scenes/scene1.mp4",
        "duration": 5
      },
      {
        "scene_number": 2,
        "status": "completed",
        "video_url": "https://storage.example.com/scenes/scene2.mp4",
        "duration": 7
      },
      {
        "scene_number": 3,
        "status": "completed",
        "video_url": "https://storage.example.com/scenes/scene3.mp4",
        "duration": 3
      }
    ]
  }
}

---

## STEP 7: Download Video

Once video generation is complete, download using the download endpoint.

Request:
GET /api/videos/download/{video_id}
Authorization: Bearer {access_token}

Response (200 OK):
{
  "success": true,
  "message": "Video download ready",
  "data": {
    "id": "8da7b810-9dad-11d1-80b4-00c04fd430c8",
    "download_url": "https://storage.example.com/videos/8da7b810.mp4",
    "file_size": 51234567,
    "format": "mp4",
    "resolution": "1920x1080",
    "duration": 15
  }
}

Now you can download the generated video for sharing on social media or other platforms.

---

## Complete Visual Timeline

Typical workflow progression from start to finish:

Time    Action                              Status                  Details
00:00   User Registration                   COMPLETE               Access token issued, 10 credits
00:10   User Login                          COMPLETE               Token refreshed
00:15   Project Initialization              COMPLETE               Project + Briefs created with images
00:20   Storyboard Generation               COMPLETE               3 scenes generated
00:25   Video Generation Started            QUEUED                 Job submitted
01:30   Video Generation                    PROCESSING             Worker started, generating scenes
02:45   Video Generation                    PROCESSING             Scene 2/3 in progress
03:30   Video Generation Completed          COMPLETE               All scenes merged
03:35   Video Ready for Download            DOWNLOADABLE           Final video available
03:45   Video Downloaded                    COMPLETE               User has MP4 file (51 MB)

Total Time: Approximately 3.75 minutes from start to downloadable video

Credit Usage Summary:
- Initial credits: 10
- Video generation: 1 credit
- Final balance: 9 credits

---

## Database Tables Used Per Workflow Step

This section details which database tables are modified/queried at each stage of the workflow:

### Step 1: User Registration
```
Register
  └─> INSERT users
      ├─ Create user row with email (unique), password (bcrypt hash)
      ├─ Set role = 'user'
      └─ Set credits = 10 (initial balance)
```

### Step 2: User Login
```
Login
  └─> SELECT users (WHERE email = ?)
      ├─ Validate password hash
      ├─ Generate JWT access_token (24h expiry)
      └─ Generate refresh_token (7 day expiry)
```

### Step 3: Initialize Project
```
Initialize Project (FE Wizard)
  └─> INSERT projects
      ├─ project_name, theme, tone, duration
      ├─ institution_name, school_level
      ├─ status = 'draft'
      └─ user_id (FK to users)
  
  └─> INSERT business_briefs
      ├─ institution_history, offered_degrees
      ├─ logo_path (save base64 image)
      ├─ environment_path (save base64 image)
      ├─ document_path (save base64 document)
      ├─ target_audience, key_message, budget, timeline
      ├─ status = 'draft'
      ├─ project_id (FK to projects)
      └─ user_id (FK to users)
  
  └─> INSERT creative_briefs
      ├─ video_type (promotional/educational/testimonial)
      ├─ duration, style, tone
      ├─ script, copywriting, hashtags
      ├─ color_palette, music_preference, mood
      ├─ business_brief_id (FK to business_briefs)
      └─ user_id (FK to users)
```

### Step 4: Generate Storyboard
```
Generate Storyboard (AI Service call)
  └─> INSERT storyboards
      ├─ title, description
      ├─ prompt (AI prompt used)
      ├─ is_selected = false
      ├─ project_id (FK to projects)
      └─ user_id (FK to users)
  
  └─> INSERT scenes (3-5 rows, one per scene)
      ├─ scene_number (1, 2, 3, ...)
      ├─ title, description, visual_description
      ├─ duration (calculated from total video_duration)
      ├─ caption (narration/voiceover text)
      ├─ regenerate_count = 0
      ├─ storyboard_id (FK to storyboards)
      └─ user_id (FK to users)
```

### Step 5: Generate Video (Enqueue)
```
Generate Video (Async Job Submission)
  └─> INSERT generation_jobs
      ├─ job_type = 'generate'
      ├─ status = 'queued'
      ├─ priority = 0
      ├─ prompt (JSONB with all context)
      ├─ scene_count, video_duration
      ├─ provider = pending (to be assigned)
      ├─ model = pending (to be assigned)
      ├─ retry_count = 0, max_retries = 3
      ├─ created_at = NOW()
      ├─ project_id (FK to projects)
      ├─ storyboard_id (FK to storyboards)
      └─ user_id (FK to users)
  
  └─> INSERT video_variants (3 rows)
      ├─ Row 1: variant_number = 1, style = 'cinematic'
      ├─ Row 2: variant_number = 2, style = 'vibrant'
      ├─ Row 3: variant_number = 3, style = 'professional'
      ├─ status = 'pending' (for all)
      ├─ provider = pending (for all)
      ├─ project_id (FK to projects)
      ├─ storyboard_id (FK to storyboards)
      └─ user_id (FK to users)
  
  └─> UPDATE users
      ├─ credits = credits - 1
      └─ WHERE id = user_id
  
  └─> Enqueue job to background worker channel (non-blocking return to FE)
```

### Step 6: Background Worker Processing (Async)
```
Background Worker (3 concurrent goroutines)
  
  └─> UPDATE generation_jobs
      ├─ status = 'processing'
      ├─ started_at = NOW()
      └─ WHERE id = job_id
  
  ├─> For each scene (loop):
  │   │
  │   ├─> INSERT scene_generations
  │   │   ├─ scene_number (from scene)
  │   │   ├─ scene_index
  │   │   ├─ prompt (text description for LTX/Runway)
  │   │   ├─ duration (in seconds)
  │   │   ├─ status = 'queued'
  │   │   ├─ external_job_id = null (to be filled)
  │   │   └─ variant_id (FK to video_variants)
  │   │
  │   ├─> Call AI Service (LTX/Runway API)
  │   │   ├─ POST /generate/video/scene
  │   │   ├─ Send: scene_number, description, duration, theme, style
  │   │   └─ Receive: video_url, external_job_id, metadata
  │   │
  │   ├─> UPDATE scene_generations
  │   │   ├─ status = 'completed'
  │   │   ├─ video_url (path to scene MP4)
  │   │   ├─ external_job_id (from AI Service)
  │   │   └─ completed_at = NOW()
  │   │
  │   └─> UPDATE generation_jobs
  │       ├─ progress = (scene_index / total_scenes * 100)%
  │       ├─ processing_notes = JSON update with progress
  │       └─ updated_at = NOW()
  │
  ├─> Merge all scene videos (FFmpeg concatenation)
  │   └─ Output: /outputs/videos/merged_{variant_style}.mp4
  │
  ├─> INSERT videos (final merged video)
  │   ├─ title (auto-generated from project)
  │   ├─ status = 'completed'
  │   ├─ video_url (path to final MP4)
  │   ├─ thumbnail_url (first frame thumbnail)
  │   ├─ duration = total_duration
  │   ├─ resolution = '1920x1080'
  │   ├─ file_size (bytes of MP4)
  │   ├─ credits_used = 1
  │   ├─ project_id (FK to projects)
  │   ├─ storyboard_id (FK to storyboards)
  │   └─ user_id (FK to users)
  │
  ├─> UPDATE generation_jobs
  │   ├─ status = 'completed'
  │   ├─ progress = 100%
  │   ├─ completed_at = NOW()
  │   └─ WHERE id = job_id
  │
  └─> UPDATE video_variants
      ├─ status = 'completed'
      ├─ video_url (path to final MP4)
      ├─ file_size (bytes)
      ├─ provider = 'ltx' (or actual provider used)
      ├─ model = 'ltx-2-fast' (or actual model used)
      ├─ resolution = '1920x1080'
      └─ updated_at = NOW()
```

### Step 7: Check Video Status (Poll)
```
Check Video Status (Client polling every 5 seconds)
  
  └─> SELECT video_variants
      ├─ WHERE user_id = ? AND project_id = ?
      ├─ Returns: status, progress (from generation_jobs)
      ├─ Returns: video_url, file_size, provider, model
      └─ Response shows processing/completed status to FE
  
  └─> SELECT scene_generations (for detailed progress)
      ├─ WHERE variant_id IN (...)
      ├─ Returns: scene_number, status, progress
      └─ Allows display of "Generating scene 2 of 3"
```

### Step 8: Download Video
```
Download Video
  
  └─> SELECT videos
      ├─ WHERE id = video_id AND user_id = user_id
      ├─ Verify status = 'completed'
      ├─ Get video_url
      └─ Return download_url to FE
  
  └─> Log download event (optional analytics)
      └─ INSERT INTO download_logs (video_id, user_id, timestamp)
```

---

## Error Handling During Workflow

If you encounter errors at any step:

1. Check Bearer token validity in Authorization header
2. Verify user_id, project_id, storyboard_id are valid UUIDs
3. Ensure required fields are provided based on endpoint specs
4. Check response message for specific error details
5. Retry failed operations (most are idempotent)

Common Status Codes:
- 200/201: Success
- 400: Invalid request parameters
- 401: Missing/invalid authentication token
- 403: Insufficient permissions or credits
- 404: Resource not found
- 500: Server error (contact support)

---

## Using Bruno API Collection

For automated testing without manual variable management:

1. Install Bruno: https://www.usebruno.com
2. Open docs/API_COLLECTION.json
3. Set base_url variable to your backend URL
4. Run requests in sequence - variables auto-populate
5. Each endpoint automatically extracts and stores IDs for next steps

This eliminates manual copy-paste of IDs and tokens throughout the workflow.

---

# Complete Integrated Architecture - AI Service + Backend + Frontend

Updated: April 2026
Full end-to-end system flow showing how Frontend, Backend, and AI Service work together

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         SEVIMA AI Video Gen System                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │  FRONTEND (Next.js - Port 3000)                                  │   │
│  │  - Login/Register forms                                          │   │
│  │  - Project wizard (institution info, creative brief)             │   │
│  │  - Storyboard preview & editor                                   │   │
│  │  - Video gallery & download manager                              │   │
│  │  - Credit balance display                                        │   │
│  └──────────────────┬───────────────────────────────────────────────┘   │
│                     │ HTTPS Requests (JSON API)                         │
│                     │                                                   │
│  ┌──────────────────▼───────────────────────────────────────────────┐   │
│  │  BACKEND (Go Fiber - Port 5000)                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐     │   │
│  │  │ Handlers (7 active)                                     │     │   │
│  │  │ - Auth: register, login, refresh, change-password       │     │   │
│  │  │ - Projects: initialize (atomic: project + briefs)       │     │   │
│  │  │ - Storyboard: generate                                  │     │   │
│  │  │ - Videos: generate, get, list, download                 │     │   │
│  │  │ - Credits: get balance, admin add                       │     │   │
│  │  └───────────────┬─────────────────────────────────────────┘     │   │
│  │                  │                                               │   │
│  │  ┌───────────────▼──────────┐  ┌─────────────────────────┐       │   │
│  │  │ Business Logic (Services)│  │ Background Job Queue    │       │   │
│  │  ├──────────────────────────┤  ├─────────────────────────┤       │   │
│  │  │ - AuthService            │  │ 3 Worker Goroutines     │       │   │
│  │  │ - BriefService (unified) │  │ - Process video jobs    │       │   │
│  │  │ - ProjectService         │  │ - Call AI Service       │       │   │
│  │  │ - StoryboardService      │  │ - Store results in DB   │       │   │
│  │  │ - VideoGenerationService │  │ - Auto-retry failures   │       │   │
│  │  │ - CreditService          │  │ - Status tracking       │       │   │
│  │  └──────────────┬───────────┘  └─────────────────────────┘       │   │
│  │                 │                                                │   │
│  │  ┌──────────────▼─────────────────────────────────────────┐      │   │
│  │  │ Data Persistence (Repositories + GORM)                 │      │   │
│  │  │ PostgreSQL Tables:                                     │      │   │
│  │  │ - users (auth & credits)                               │      │   │
│  │  │ - projects, business_briefs, creative_briefs           │      │   │
│  │  │ - storyboards, scenes                                  │      │   │
│  │  │ - videos, generation_jobs, video_variants, scenes      │      │   │
│  │  └────────────────────────────────────────────────────────┘      │   │
│  │                 │                                                │   │
│  └─────────────────┼────────────────────────────────────────────────┘   │
│                    │ HTTP (with X-User-ID, X-User-Email headers)        │
│  ┌─────────────────▼──────────────────────────────────────────────┐     │
│  │  AI SERVICE (FastAPI - Port 8000)                              │     │
│  │  ┌─────────────────────────────────────────────────────────┐   │     │
│  │  │ Processing Pipelines                                    │   │     │
│  │  │ - Prompt optimization & content generation              │   │     │
│  │  │ - Storyboard scene generation with AI models            │   │     │
│  │  │ - Video generation via multiple providers:              │   │     │
│  │  │   * LTX Video (Fast, high-quality)                      │   │     │
│  │  │   * Runway Gen-3 (Alternative provider)                 │   │     │
│  │  │ - Asset management (logos, backgrounds, overlays)       │   │     │
│  │  │ - Scene-by-scene video compilation                      │   │     │
│  │  └─────────────────────────────────────────────────────────┘   │     │
│  │                                                                │     │
│  │  Outputs:                                                      │     │
│  │  - ai-service/outputs/videos/ (final MP4s)                     │     │
│  │  - ai-service/outputs/reports/ (generation logs)               │     │
│  │  - Video URLs stored in PostgreSQL                             │     │
│  └────────────────────────────────────────────────────────────────┘     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Component Interactions - Detailed Data Flow

### Phase 1: User Registration & Authentication (FE ↔ BE)

FE: User enters name, email, password in login form
  ↓
BE POST /api/auth/register
  ├─ Hash password with bcrypt
  ├─ Create user in DB with role=user, credits=10
  ├─ Generate JWT access_token (24h expiry)
  ├─ Generate refresh_token (7 day expiry)
  └─ Return tokens to FE
FE: Store access_token in session/localStorage
  ├─ Redirect to dashboard
  └─ Include Authorization: Bearer {access_token} in all subsequent requests

---

### Phase 2: Project Initialization (FE ↔ BE ↔ PostgreSQL)

FE: User fills out wizard form with:
  ├─ Institution info (name, history, school level, degrees)
  ├─ Event/Campaign info (penerimaan siswa baru, tone, key message)
  ├─ Creative direction (theme, duration, prompt)
  ├─ Copywriting (editable text, hashtags)
  └─ Uploads images (logo, environment, document) as base64
  ↓
BE POST /api/projects/initialize (atomic operation)
  ├─ Validate JWT token from header
  ├─ Create Project row with:
  │  ├─ project_name, theme, tone, duration
  │  ├─ institution_name, school_level
  │  └─ user_id (extracted from JWT)
  │
  ├─ Create BusinessBrief row with:
  │  ├─ institution_history, offered_degrees
  │  ├─ logo_path (save base64 to /outputs/)
  │  ├─ environment_path (save base64 to /outputs/)
  │  └─ project_id (FK)
  │
  ├─ Create CreativeBrief row with:
  │  ├─ event_content, tone_of_voice, key_message
  │  ├─ selected_theme, video_duration
  │  ├─ prompt, copywriting, hashtags
  │  └─ project_id (FK)
  │
  └─ Return project_id, brief_ids to FE
FE: Store project_id locally and redirect to storyboard view

---

### Phase 3: Storyboard Generation (FE ↔ BE → AI Service → DB)

FE User clicks "Generate Storyboard" button
  ↓
BE POST /api/storyboard/generate (user authenticated)
  ├─ Fetch Project + BusinessBrief + CreativeBrief from DB
  ├─ Build prompt for AI with all context:
  │  ├─ Institution info, key message, tone
  │  ├─ Event details, copywriting, hashtags
  │  ├─ Theme requirements, video duration
  │  └─ Examples of good storytelling
  │
  ├─ Call AI Service: POST /generate/storyboard
  │  └─ AI Service:
  │     ├─ Process prompt with GPT-4/Claude
  │     ├─ Generate 3-5 scenes with:
  │     │  ├─ Scene descriptions (for video generation)
  │     │  ├─ Narration/voiceover text
  │     │  └─ Visual direction (lighting, mood, composition)
  │     └─ Return scenes array to Backend
  │
  ├─ Create Storyboard row with status=ready
  ├─ Create Scene rows (1 per scene) with:
  │  ├─ scene_number, description, narration
  │  ├─ duration (calculated from video_duration)
  │  └─ storyboard_id (FK)
  │
  └─ Return storyboard_id + scenes to FE
FE: Display storyboard in preview
  ├─ Show each scene with description & narration
  └─ User can edit descriptions/narration before video generation

---

### Phase 4: Video Generation (FE ↔ BE → Job Queue → AI Service → Storage)

FE User clicks "Generate Video" button
  ↓
BE POST /api/videos/generate
  ├─ Validate JWT & check user credits >= 1
  ├─ Fetch Storyboard + Scenes from DB
  ├─ Deduct 1 credit from user (UPDATE users.credits)
  │
  ├─ Create GenerationJob row with:
  │  ├─ status=queued, progress=0%
  │  ├─ project_id, storyboard_id, user_id
  │  └─ created_at timestamp
  │
  ├─ Create VideoVariant rows (3 variants):
  │  ├─ Variant 1: style=cinematic (dramatic, high-contrast)
  │  ├─ Variant 2: style=vibrant (colorful, energetic)
  │  ├─ Variant 3: style=professional (corporate, clean)
  │  └─ Each with status=queued, provider=pending
  │
  ├─ Enqueue job to Channel (shared across 3 worker goroutines)
  └─ Return job_id to FE immediately (async processing)

FE: Redirect to "Video in Progress" page
  ├─ Show spinning loader
  ├─ Display estimated time (2-5 minutes)
  └─ Start polling GET /api/videos/{video_id} every 5 seconds

---

### Phase 4B: Async Video Processing (Background Workers)

BE Job Queue Worker (3 concurrent goroutines):
  ├─ Receive job from channel (blocking until available)
  │
  ├─ Update GenerationJob: status=processing
  │
  ├─ For Each Scene in Storyboard:
  │  ├─ Create SceneGeneration row (status=queued, scene_id)
  │  │
  │  ├─ Call AI Service: POST /generate/video/scene
  │  │  │  Request:
  │  │  │  {
  │  │  │    "scene_number": 1,
  │  │  │    "description": "Wide shot of campus...",
  │  │  │    "narration": "Halo generasi masa depan!",
  │  │  │    "duration": 5,
  │  │  │    "theme": "Tur Kampus Sinematik",
  │  │  │    "style": "cinematic"
  │  │  │  }
  │  │  │
  │  │  │  AI Service:
  │  │  │  ├─ Convert text description to visual prompt for LTX/Runway
  │  │  │  ├─ Generate video via LTX API call
  │  │  │  │  ├─ Model: ltx-2-fast
  │  │  │  │  ├─ Duration: 5 seconds
  │  │  │  │  └─ Resolution: 1920x1080
  │  │  │  ├─ Optional: Add narration audio overlay
  │  │  │  ├─ Save video to: ai-service/outputs/videos/scene_1_cinematic.mp4
  │  │  │  └─ Return video_url + metadata
  │  │  │
  │  │  └─ Backend receives video_url
  │  │
  │  ├─ Update SceneGeneration:
  │  │  ├─ status=completed
  │  │  ├─ video_url (local path or S3 URL)
  │  │  ├─ duration_actual
  │  │  └─ completed_at
  │  │
  │  └─ Update GenerationJob: progress = (scene_number / total_scenes * 100)%
  │
  ├─ After ALL scenes processed: Merge videos
  │  ├─ Use FFmpeg to concatenate scene videos
  │  ├─ Output: ai-service/outputs/videos/merged_cinematic.mp4
  │  ├─ Generate thumbnail from first frame
  │  └─ Calculate final file_size
  │
  ├─ Update GenerationJob: status=completed
  │
  ├─ Update each VideoVariant:
  │  ├─ status=completed
  │  ├─ video_url = final_merged_video_url
  │  ├─ provider=ltx
  │  ├─ model=ltx-2-fast
  │  ├─ resolution=1920x1080
  │  ├─ file_size = actual_size_in_bytes
  │  └─ completed_at = timestamp
  │
  └─ Worker releases goroutine (returns to waiting for next job)

---

### Phase 5: Video Status Polling (FE ↔ BE ↔ DB)

FE Polling Loop (every 5 seconds):
  ├─ GET /api/videos/{video_id}
  ├─ BE Response shows:
  │  ├─ Initially: status=queued or processing, progress=0%
  │  ├─ During: status=processing, progress=33%, "Generating scene 2 of 3"
  │  └─ Finally: status=completed, video_url provided, all scenes available
  │
  └─ When status=completed:
     ├─ FE displays video preview
     ├─ Shows all 3 variants (cinematic, vibrant, professional)
     ├─ User can watch/download each variant
     └─ Stop polling

---

### Phase 6: Video Download (FE ↔ BE → Storage)

FE User clicks "Download Video" button
  ↓
BE GET /api/videos/download/{video_id}
  ├─ Validate JWT token
  ├─ Fetch VideoVariant from DB
  ├─ Verify video_url exists
  ├─ Log download event (for analytics)
  └─ Return response with download_url
FE: Trigger browser download
  ├─ Create <a> element with href={download_url}
  ├─ Click to initiate download
  └─ File saved as: SMA_Negeri_1_Jakarta_cinematic.mp4

---

## Request Flow Timing & Performance

```
Timeline (seconds)  Frontend State              Backend Action              AI Service Action
00:00              Loading login              POST /api/auth/register
00:05              Dashboard loaded           User in DB, tokens issued
00:10              Filling wizard form
00:15              Form submitted             POST /api/projects/initialize
00:16              Storyboard loading         Calling AI Service          POST /generate/storyboard
00:18              Storyboard displayed       Scenes stored in DB         Returning scenes
00:25              User reviews scenes
00:35              Clicks generate video      POST /api/videos/generate
00:36              Shows "generating..."      Job queued                  Waiting for worker
00:40              Polling status             Job status=processing       Worker started
00:50              Polling status             Progress 33%                Generating scene 1
01:20              Polling status             Progress 66%                Generating scene 2
02:00              Polling status             Progress 99%                Merging scenes
02:30              Polling status             status=completed            Completed
02:35              Video preview shown        Sending video_url           FFmpeg output saved
03:00              User clicks download       GET /api/videos/download
03:05              Video downloaded (51 MB)   Log event                   N/A
```

---

## Data Models Relationships

```
USER (PK: id)
├─ 1 ← → N PROJECTS
├─ 1 ← → N GENERATION_JOBS
├─ 1 ← → N VIDEO_VARIANTS
└─ credits: int (deducted per video generation)

PROJECT (PK: id, FK: user_id)
├─ 1 ← → 1 BUSINESS_BRIEF (FK: project_id)
├─ 1 ← → 1 CREATIVE_BRIEF (FK: business_brief_id)
├─ 1 ← → N STORYBOARDS (FK: project_id)
└─ 1 ← → N GENERATION_JOBS (FK: project_id)

STORYBOARD (PK: id, FK: project_id)
└─ 1 ← → N SCENES (FK: storyboard_id)

SCENE (PK: id, FK: storyboard_id)
└─ 1 ← → N SCENE_GENERATIONS (FK: scene_id)

GENERATION_JOB (PK: id, FK: project_id, storyboard_id, user_id)
├─ status: queued → processing → completed → failed
├─ progress: 0-100%
└─ 1 ← → N VIDEO_VARIANTS (FK: generation_job_id)

(Note: ContentPillar & ContentTheme removed in April 2026 audit - not used in active workflow)

VIDEO_VARIANT (PK: id, FK: generation_job_id, user_id)
├─ style: cinematic | vibrant | professional
├─ status: queued → processing → completed
├─ video_url: path to final MP4
├─ provider: ltx | runway
└─ 1 ← → N SCENE_GENERATIONS (FK: variant_id)

SCENE_GENERATION (PK: id, FK: scene_id, variant_id)
├─ status: queued → processing → completed
├─ video_url: path to scene MP4
└─ provider: ltx | runway
```

---

## Error Recovery & Retry Logic

### Scenario 1: AI Service Timeout During Scene Generation
```
Worker processing scene 2:
├─ Call to AI Service times out (30s)
├─ Catch error: connection timeout
├─ Update SceneGeneration: status=failed, error_message
├─ Update GenerationJob: retry_count++
├─ If retry_count < 3:
│  └─ Re-queue job to front of channel (priority retry)
└─ If retry_count >= 3:
   ├─ Update GenerationJob: status=failed
   ├─ Notify user via FE polling response
   └─ User can click "Retry" to start fresh
```

### Scenario 2: Insufficient Credits
```
FE POST /api/videos/generate:
├─ BE checks: user.credits >= 1
├─ If insufficient:
│  ├─ Return 403 Forbidden
│  ├─ Response: "Insufficient credits. Please buy more credits."
│  └─ FE shows upgrade prompt
```

### Scenario 3: AI Service Returns Empty/Invalid Video
```
Worker receives video_url from AI Service:
├─ Validate: file exists and size > 0
├─ If invalid:
│  ├─ Log error with detailed context
│  ├─ Attempt retry (up to 2 times)
│  └─ If persistent: mark as failed, notify user
```

---

## Credit System Integration

```
User Registration
├─ Initial credits: 10
│
Video Generation Flow
├─ FE shows: "This will use 1 credit"
├─ User clicks "Generate Video"
├─ BE: Deduct 1 credit atomically
│  ├─ UPDATE users.credits = credits - 1 WHERE id = user_id
│  ├─ If credits < 0, ROLLBACK (prevent negative credits)
│  └─ Create GenerationJob only if deduction succeeds
│
└─ Return JobId to FE immediately
   ├─ FE shows: "Credit used. Your balance: 9/10"
   └─ Job processes in background

Admin Add Credits
├─ Admin POST /api/credits/ with user_id + amount
├─ BE: UPDATE users.credits = credits + amount WHERE id = user_id
└─ User sees updated balance next login
```

---

## Key Performance Characteristics

### Request Latency (end-to-end, client perspective)
- Registration: 200ms (DB write)
- Login: 150ms (DB query + JWT generation)
- Project Initialize: 300ms (atomic multi-table insert + file uploads)
- Storyboard Generate: 2-3 seconds (calls AI Service, waits for response)
- Video Generate: 100ms (immediate response, actual processing async)
- Video Status Check: 50-100ms (DB query only)
- Video Download: 30ms (redirect response)

### AI Service Processing Time (per video)
- Storyboard generation (3-5 scenes): 2-3 seconds (GPT-4 inference)
- Scene video generation per scene: 30-60 seconds (LTX model inference)
- 3 scenes = 90-180 seconds total generation time
- Video merging: 10-20 seconds (FFmpeg)
- Total queue + processing: 2-5 minutes typical

### Database Performance
- User lookup (by email): <5ms (indexed)
- Project + Briefs atomic insert: 15-25ms
- Scene generation batch insert: 10-15ms
- Status polling query: <5ms (indexed by user_id + created_at)

### Worker Concurrency
- 3 worker goroutines running concurrently
- Allows 3 videos to generate simultaneously
- If 4th job arrives, it waits in channel until worker available
- Typical queue wait: 0-2 minutes depending on current load

---

## Deployment Architecture

### Development (localhost)
```
Frontend:   http://localhost:3000
Backend:    http://localhost:5000
AI Service: http://localhost:8000
Database:   localhost:5432 (PostgreSQL)
```

### Production
```
Frontend:   Deployed on Vercel/AWS S3 + CloudFront
Backend:    Docker container on AWS ECS/Kubernetes
AI Service: Docker container on GPU instance
Database:   Managed PostgreSQL (AWS RDS)
Video Storage: S3 bucket with CDN (CloudFront)

Environment Variables:
├─ BACKEND_URL = https://api.sevima.example.com
├─ AI_SERVICE_URL = http://ai-service.internal:8000 (internal)
├─ DATABASE_URL = postgresql://user:pass@rds-endpoint:5432/db
├─ JWT_SECRET = [secure random string]
└─ AWS_S3_BUCKET = sevima-videos-prod
```

---

## Summary

The integrated system works as follows:

1. **Frontend** - User-facing interface built with Next.js
   - Handles UI/UX for form submission
   - Manages JWT tokens in session storage
   - Polls backend for async video progress
   - Displays results and enables downloads

2. **Backend** - Orchestration layer built with Go Fiber
   - Validates all requests & JWT tokens
   - Manages database transactions atomically
   - Coordinates with AI Service for content generation
   - Implements background job queue for video processing
   - Handles credit deduction and balance management

3. **AI Service** - Processing layer built with FastAPI
   - Generates storyboard scenes using LLMs
   - Creates videos using LTX/Runway APIs
   - Saves outputs to local storage or S3
   - Returns URLs back to Backend

4. **Database** - Persistent storage (PostgreSQL)
   - 12 tables modeling complete workflow
   - Tracks all generation jobs and status
   - Maintains user credit balances
   - Enables status polling and historical data

The workflow is designed for scalability: Frontend is stateless, Backend coordinates work asynchronously via job queues, and AI Service handles compute-intensive tasks on specialized hardware. Users experience fast initial responses even though video generation takes minutes in the background.
