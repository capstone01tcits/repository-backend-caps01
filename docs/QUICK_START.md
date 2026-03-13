# Quick Start Guide - Video Generation Backend

## 5-Minute Setup

### 1. Database Migration

Run migrations to create tables:

```go
// In your main.go
import "app/internal/migration"

// During app initialization
db := initializeDB()
if err := migration.MigrateVideoGeneration(db); err != nil {
    log.Fatal("Migration failed:", err)
}
```

Or use raw SQL migration:
```go
if err := migration.RunRawMigrations(db); err != nil {
    log.Fatal("SQL Migration failed:", err)
}
```

### 2. Initialize Services & Queue

```go
// In dependency injection setup
import (
    "app/internal/repository"
    "app/internal/service"
    "app/internal/queue"
    "app/internal/handler"
)

// Create repositories
jobRepo := repository.NewGenerationJobRepository(db)
variantRepo := repository.NewVideoVariantRepository(db)
sceneRepo := repository.NewSceneGenerationRepository(db)

// Create service
creditService := service.NewCreditService(creditRepo) // Assume exists
videoGenService := service.NewVideoGenerationService(
    jobRepo,
    variantRepo,
    sceneRepo,
    creditService,
)

// Create queue
jobQueue := queue.NewJobQueue(jobRepo, videoGenService)

// Start queue workers (3 workers for concurrent processing)
ctx := context.Background()
jobQueue.Start(ctx, 3)

// Create handler
videoHandler := handler.NewVideoHandler(videoGenService)
```

### 3. Register Routes

```go
// In your router setup
api := app.Group("/api")

videoRoutes := api.Group("/videos")
{
    // Generate videos
    videoRoutes.Post("/generate", videoHandler.GenerateVideoVariants)
    
    // Check generation status
    videoRoutes.Get("/generation/:jobId", videoHandler.GetGenerationJobStatus)
    
    // Get variants for storyboard
    videoRoutes.Get("/storyboard/:storyboardId", videoHandler.GetVideoVariants)
    
    // Get single variant
    videoRoutes.Get("/:variantId", videoHandler.GetVideoVariant)
    
    // Regenerate video
    videoRoutes.Post("/:variantId/regenerate", videoHandler.RegenerateVideoVariant)
    
    // Regenerate scene
    videoRoutes.Post("/scene/:sceneId/regenerate", videoHandler.RegenerateScene)
    
    // Download video
    videoRoutes.Get("/:variantId/download", videoHandler.DownloadVideo)
}
```

### 4. Test the Flow

**Terminal 1 - Start the server:**
```bash
go run cmd/main.go
```

**Terminal 2 - Generate videos:**
```bash
# 1. Create a project (assume you have existing endpoints)
curl -X POST http://localhost:8000/api/projects \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Campus Marketing"}'
# Save project_id

# 2. Create a storyboard (or use existing)
# Save storyboard_id

# 3. Generate videos
curl -X POST http://localhost:8000/api/videos/generate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id":"PROJECT_ID",
    "storyboard_id":"STORYBOARD_ID",
    "custom_prompt":"Beautiful campus scenes"
  }'

# Response will include generation_job_id
# Save it as JOB_ID
```

**Terminal 2 - Check status:**
```bash
# Check every 5 seconds
while true; do
  curl -X GET "http://localhost:8000/api/videos/generation/JOB_ID" \
    -H "Authorization: Bearer YOUR_TOKEN"
  echo ""
  sleep 5
done

# Or get all variants
curl -X GET "http://localhost:8000/api/videos/storyboard/STORYBOARD_ID" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Terminal 2 - Download completed video:**
```bash
# Once status is "completed"
curl -X GET "http://localhost:8000/api/videos/VARIANT_ID/download" \
  -H "Authorization: Bearer YOUR_TOKEN" | jq .data.download_url
```

## Common Tasks

### Generate Videos with Custom Direction

```bash
curl -X POST http://localhost:8000/api/videos/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "uuid1",
    "storyboard_id": "uuid2",
    "custom_prompt": "Focus on student success stories",
    "scene_count": 3,
    "video_duration": 12
  }'
```

### Regenerate a Single Video

```bash
# Get the variant_id from /videos/storyboard/{storyboardId}
curl -X POST http://localhost:8000/api/videos/{variant_id}/regenerate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "new_prompt": "More emphasis on campus facilities"
  }'
```

### Regenerate a Single Scene

```bash
# Get the scene_id from /videos/{variant_id}
curl -X POST http://localhost:8000/api/videos/scene/{scene_id}/regenerate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "new_prompt": "Show aerial campus views"
  }'
```

### Monitor Queue Status

```go
// Add a debug endpoint to check queue stats
api.Get("/admin/queue-stats", func(c *fiber.Ctx) error {
    stats := jobQueue.GetStats(c.Context())
    return c.JSON(stats)
})
```

## Behind the Scenes

### What Happens When You Generate Videos

1. **Request Validation**: Checks project_id, storyboard_id, user credits
2. **Job Creation**: Creates GenerationJob with status "queued"
3. **Variant Creation**: Creates 3 VideoVariant records
4. **Scene Planning**: Creates 2-3 SceneGeneration records per variant
5. **Credit Deduction**: Deducts credits from user account
6. **Queue Enqueue**: Adds job to processing queue
7. **Worker Processing**: Background worker picks up job
8. **Provider Calls**: Calls video provider API for each scene
9. **Polling**: Continuously polls for job completion
10. **Database Updates**: Updates scenes/variants/job as they complete
11. **Frontend Notified**: Videos ready for download

### Status Transitions

```
User Request
    ↓
GenerationJob created (status: queued)
    ↓
Worker picks up (status: processing)
    ↓
Scenes being generated (VideoVariant status: processing)
    ↓
Polling begins (checking provider status every 60 seconds)
    ↓
Scenes complete (VideoVariant status: completed)
    ↓
Job complete (GenerationJob status: completed)
    ↓
Videos ready for download
```

## Monitoring & Debugging

### Check Queue Status

```go
stats := jobQueue.GetStats(c.Context())
// Returns:
// {
//   "is_running": true,
//   "processing_workers": 3,
//   "poll_interval": "60s"
// }
```

### View Pending Jobs

```bash
# Direct database query
SELECT * FROM generation_jobs 
WHERE status = 'queued' 
ORDER BY priority DESC, created_at ASC;
```

### View Processing Jobs

```bash
SELECT * FROM generation_jobs 
WHERE status = 'processing' 
ORDER BY started_at DESC;
```

### View Failed Jobs

```bash
SELECT * FROM generation_jobs 
WHERE status = 'failed' 
ORDER BY created_at DESC;
```

### Check Logs

The background workers log important events:
```
[Worker 0 started]
[Worker 0 processing job abc123 (type: generate)]
[Worker 0 error processing job abc123: ...]
[Job abc123 polling timeout after 120 attempts]
```

## Configuration Variables

Set these in your `.env` or config:

```
VIDEO_GENERATION_WORKERS=3
VIDEO_POLLING_INTERVAL=60
VIDEO_MAX_RETRIES=3
VIDEO_SCENE_DURATION_MIN=4
VIDEO_SCENE_DURATION_MAX=6
VIDEO_TOTAL_DURATION_MIN=8
VIDEO_TOTAL_DURATION_MAX=12
VIDEO_MIN_SCENE_COUNT=2
VIDEO_MAX_SCENE_COUNT=3
DEFAULT_VIDEO_PROVIDER=ltx
DEFAULT_VIDEO_MODEL=ltx-2-fast
```

## Performance Tips

### For Development
- Start with 1-2 workers
- Use shorter polling intervals (10-30 seconds)
- Enable debug logging

### For Production
- Use 3-5 workers (scale based on load)
- Use 60-second polling interval
- Implement database connection pooling
- Add database indexes (already included in migration)
- Use CDN for video storage/delivery

### Database Optimization
- Indexes are created for common queries
- Consider partitioning generation_jobs by date for large volumes
- Archive old jobs after 30 days

## Troubleshooting

### Videos not generating
1. Check job status: `GET /api/videos/generation/{jobId}`
2. Look for error_message field
3. Verify user has sufficient credits
4. Check worker logs for errors

### Stuck in "processing"
1. Check if workers are running: `GET /admin/queue-stats`
2. Verify polling is working
3. Check provider mock implementations (they simulate delays)
4. Increase polling attempts or timeout

### High credit usage
1. Check video duration configured
2. Check scene count
3. Review provider-specific costs
4. Consider using cheaper tier (LTX-2-Fast vs Pro)

### Database errors
1. Ensure migration ran successfully
2. Check database connection
3. Verify user/project/storyboard records exist
4. Check foreign key constraints

## Next Steps

1. **Integrate Frontend**: Connect frontend to these endpoints
2. **Implement Variants UI**: Show 3 video options to user
3. **Add Real Storage**: Replace mock storage with S3/GCS
4. **Setup Webhooks**: Notify frontend when videos complete
5. **Analytics**: Track video generation success rates
6. **Billing Dashboard**: Show credit usage to users

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (Ember JS)                  │
│         (User submits briefing, views videos)           │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP Request
                     ↓
┌─────────────────────────────────────────────────────────┐
│                   API Layer                             │
│  ┌──────────────────────────────────────────────────┐  │
│  │  VideoHandler                                    │  │
│  │  - POST /generate                                │  │
│  │  - GET /generation/{jobId}                       │  │
│  │  - POST /{variantId}/regenerate                  │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│              Service Layer                              │
│  ┌──────────────────────────────────────────────────┐  │
│  │  VideoGenerationService                          │  │
│  │  - Business logic for video generation           │  │
│  │  - Credit management                             │  │
│  │  - Job orchestration                             │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │
          ┌──────────┴──────────┐
          ↓                     ↓
    ┌──────────────┐     ┌──────────────┐
    │ Repository   │     │ Queue        │
    │ Layer        │     │ System       │
    ├──────────────┤     ├──────────────┤
    │- Job Repo    │     │- Worker 1    │
    │- Variant     │     │- Worker 2    │
    │  Repo        │     │- Worker 3    │
    │- Scene Repo  │     │- Polling     │
    └──────────────┘     └──────────────┘
          │                     │
          ↓                     ↓
    ┌──────────────┐     ┌──────────────┐
    │  Database    │     │ AI Providers │
    │  (PostgreSQL)│     │ (Mock)       │
    │              │     ├──────────────┤
    │- Jobs        │     │- LTX-2-Fast  │
    │- Variants    │     │- Runway      │
    │- Scenes      │     │- Wan2.1      │
    └──────────────┘     └──────────────┘
```

## File Structure

```
internal/
├── model/
│   └── generation_job.go      # Data models
├── repository/
│   ├── generation_job_repository.go
│   ├── video_variant_repository.go
│   └── scene_generation_repository.go
├── service/
│   └── video_generation_service.go
├── handler/
│   └── video_handler.go
├── queue/
│   └── job_queue.go
├── ai/
│   ├── provider.go            # Provider interface
│   └── providers.go           # Mock implementations
└── migration/
    └── video_generation.go

docs/
├── VIDEO_GENERATION_BACKEND.md  # Full documentation
├── VIDEO_API_ENDPOINTS.md       # API reference
└── QUICK_START.md              # This file
```

---

**Ready to generate some videos!** 🎬
