# AI Video Generator Backend - Implementation Guide

## Overview

This document describes the backend implementation of the AI Video Generator system. The system generates 3 variations of marketing videos from a single project briefing using scene-based video generation.

## Architecture

### Core Components

#### 1. **Models** (`internal/model/`)

- **GenerationJob**: Tracks video generation tasks in the queue
  - Manages job lifecycle: queued → processing → completed/failed
  - Stores provider and model selection
  - Tracks retry attempts and credit usage

- **VideoVariant**: Represents one of 3 video variations per storyboard
  - Variant number (1, 2, or 3)
  - Scene plan and prompt used
  - Status tracking and metadata
  - Supports revision tracking

- **SceneGeneration**: Individual scene within a video variant
  - Tracks scene-level generation status
  - Stores external job IDs from providers
  - Links to video URLs per scene

#### 2. **Repositories** (`internal/repository/`)

- **GenerationJobRepository**: CRUD operations for generation jobs
- **VideoVariantRepository**: CRUD operations for video variants
- **SceneGenerationRepository**: CRUD operations for scene generation tasks

#### 3. **AI Provider System** (`internal/ai/`)

Abstracts video generation providers:

- **LTX-2-Fast**: Standard tier, $0.04/sec
- **LTX-2-Pro**: Premium tier, $0.06/sec
- **Runway Gen4.5**: Premium, 12 credits/sec
- **Runway Gen4 Turbo**: Fast premium, 5 credits/sec
- **Wan2.1**: Open-source, minimal credits
- **LTX-Video Open Source**: Internal/research, minimal credits

Each provider:
- Generates video scenes from prompts
- Polls job status asynchronously
- Calculates credits consumed
- Provides mock responses for development

#### 4. **Video Generation Service** (`internal/service/video_generation_service.go`)

Core business logic:

```go
GenerateVideoVariants()    // Creates 3 video variants with scene-based generation
RegenerateVideoVariant()   // Regenerate individual video with new prompt
RegenerateScene()          // Regenerate individual scene
GetJobStatus()             // Track generation progress
ProcessGenerationJob()     // Worker: execute job (scene generation)
PollJobStatus()            // Worker: poll provider for completion
```

**Features**:
- Automatic scene planning (2-3 scenes per video, 4-6 sec each)
- 8-12 sec total duration per video
- Credit validation and deduction
- Automatic provider selection based on tier

#### 5. **Job Queue System** (`internal/queue/`)

**SimpleJobQueue**: Manages background processing

- Enqueues generation jobs
- Tracks job status through lifecycle
- Worker pool for parallel processing
- Automatic polling for job completion
- Retry logic with configurable max attempts

**Worker Process**:
1. Dequeue pending job (ordered by priority)
2. Mark job as processing
3. Process each video variant's scenes
4. Call video provider to start generation
5. Start polling routine for job completion
6. Update database with results

**Polling Mechanism**:
- Polls every 60 seconds for job status
- Continues for up to 2 hours
- Updates scene/variant status when complete
- Handles timeout and failure scenarios

#### 6. **API Handlers** (`internal/handler/video_handler.go`)

**Endpoints**:

```
POST   /api/videos/generate                  # Generate 3 video variants
GET    /api/videos/generation/:jobId         # Get generation job status
GET    /api/videos/storyboard/:storyboardId  # Get all variants for storyboard
GET    /api/videos/:variantId                # Get single variant with scenes
POST   /api/videos/:variantId/regenerate     # Regenerate variant
POST   /api/videos/scene/:sceneId/regenerate # Regenerate single scene
GET    /api/videos/:variantId/download       # Get download URL
```

### Data Flow

```
Frontend Request
    ↓
Handler validates request
    ↓
Service creates GenerationJob + 3 VideoVariants + SceneGenerations
    ↓
Job enqueued in JobQueue
    ↓
Worker processes job:
    - Calls VideoProvider.GenerateScene() for each scene
    - Stores external job ID
    ↓
Polling routine starts:
    - Calls VideoProvider.GetJobStatus() periodically
    - Updates scene status when complete
    - Updates variant status when all scenes done
    ↓
Frontend polls /api/videos/generation/:jobId
    ↓
Response with video URLs when ready
```

## Usage

### 1. Generate Videos

**Request**:
```bash
POST /api/videos/generate
{
  "project_id": "uuid",
  "storyboard_id": "uuid",
  "custom_prompt": "optional custom direction"
}
```

**Response**:
```json
{
  "generation_job_id": "uuid",
  "status": "queued",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 2. Poll Generation Status

**Request**:
```bash
GET /api/videos/generation/{jobId}
```

**Response**:
```json
{
  "id": "uuid",
  "job_type": "generate",
  "status": "processing",
  "processing_notes": {...},
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 3. Get Video Variants

**Request**:
```bash
GET /api/videos/storyboard/{storyboardId}
```

**Response**:
```json
[
  {
    "id": "uuid",
    "variant_number": 1,
    "status": "completed",
    "video_url": "https://storage.example.com/...",
    "thumbnail_url": "https://storage.example.com/...",
    "prompt_used": "cinematic...",
    "duration": 10,
    "scenes": [
      {
        "id": "uuid",
        "scene_number": 1,
        "status": "completed",
        "video_url": "...",
        "duration": 5
      }
    ]
  }
]
```

### 4. Regenerate Video

**Request**:
```bash
POST /api/videos/{variantId}/regenerate
{
  "new_prompt": "different premise..."
}
```

**Response**:
```json
{
  "generation_job_id": "uuid",
  "status": "queued"
}
```

## Credit System

**Credit Calculation**:
- Base cost: 2 credits per second per video
- Formula: `duration_sec × scene_count × video_count × 2`
- Example: 10 sec × 2 scenes × 3 videos × 2 = 120 credits

**Regeneration Costs**:
- Full video regeneration: Same as original
- Single scene regeneration: Lower cost (1 credit × duration)

**Provider-Specific**:
- LTX-2-Fast: 4 credits/sec
- LTX-2-Pro: 6 credits/sec
- Runway Gen4.5: 12 credits/sec
- Open-source: 1 credit/sec (minimal internal tracking)

## Configuration

### Environment Variables

```
VIDEO_GENERATION_WORKERS=3      # Number of worker goroutines
VIDEO_POLLING_INTERVAL=60s      # Poll frequency in seconds
VIDEO_MAX_RETRIES=3             # Max retry attempts per job
VIDEO_SCENE_DURATION_MIN=4      # Minimum scene duration in seconds
VIDEO_SCENE_DURATION_MAX=6      # Maximum scene duration in seconds
VIDEO_TOTAL_DURATION_MIN=8      # Minimum total video duration
VIDEO_TOTAL_DURATION_MAX=12     # Maximum total video duration
```

### Initializing Services

```go
// In your main.go or dependency injection setup:

// Create repositories
jobRepo := repository.NewGenerationJobRepository(db)
variantRepo := repository.NewVideoVariantRepository(db)
sceneRepo := repository.NewSceneGenerationRepository(db)

// Create service
videoGenService := service.NewVideoGenerationService(
    jobRepo,
    variantRepo,
    sceneRepo,
    creditService,
)

// Create queue
jobQueue := queue.NewJobQueue(jobRepo, videoGenService)

// Start workers
jobQueue.Start(context.Background(), 3) // 3 workers

// Create handler
videoHandler := handler.NewVideoHandler(videoGenService)
```

## Testing

### Mock Providers

All providers include mock implementations that:
- Simulate API delays (400-1000ms)
- Return random status transitions (processing → completed)
- Generate mock video URLs
- Calculate realistic credit costs

### Manual Testing

```bash
# Start server
go run cmd/main.go

# Generate videos
curl -X POST http://localhost:8000/api/videos/generate \
  -H "Content-Type: application/json" \
  -d '{"project_id":"...", "storyboard_id":"..."}'

# Check status (returns soon - mock processing)
curl http://localhost:8000/api/videos/generation/{jobId}

# View variants (may show "processing" initially)
curl http://localhost:8000/api/videos/storyboard/{storyboardId}
```

## Process Flow Diagrams

### Generation Flow

```
User submits briefing
    ↓
Create GenerationJob (status: queued)
Create 3 VideoVariants (status: pending)
Create 2-3 SceneGenerations per variant (status: pending)
    ↓
Deduct credits from user account
    ↓
Job enqueued in JobQueue
    ↓
Background Worker picks up job
    ↓
For each scene:
    - Provider.GenerateScene(prompt, duration)
    - Store external job ID
    - Update status to "processing"
    ↓
Start polling routine (fires every 60 sec)
    ↓
Provider.GetJobStatus(external_job_id)
    ↓
If completed:
    - Update scene with video URL
    - Update variant status
    - When all scenes done, variant = completed
    ↓
If all variants completed:
    - Update job status to "completed"
```

### Regeneration Flow

```
User selects variant to regenerate
    ↓
Submit new prompt (optional)
    ↓
Create new GenerationJob (type: regenerate)
Deduct regeneration credits
    ↓
Create new SceneGenerations with updated prompts
    ↓
Queue job for processing
    ↓
Same as Generation Flow above
```

## Database Schema

### GenerationJob Table
```sql
- id UUID PRIMARY KEY
- user_id UUID (FOREIGN KEY)
- project_id UUID (FOREIGN KEY)
- storyboard_id UUID (FOREIGN KEY)
- video_id UUID (FOREIGN KEY)
- job_type VARCHAR (generate, regenerate, regenerate_scene)
- status VARCHAR (queued, processing, completed, failed) - INDEXED
- priority INT
- prompt JSONB
- scene_count INT
- video_duration INT
- provider VARCHAR
- model VARCHAR
- credits_required INT
- credits_used INT
- retry_count INT
- max_retries INT
- started_at TIMESTAMP
- completed_at TIMESTAMP
- error_message TEXT
- created_at TIMESTAMP
- updated_at TIMESTAMP
- deleted_at TIMESTAMP
```

### VideoVariant Table
```sql
- id UUID PRIMARY KEY
- user_id UUID
- project_id UUID
- storyboard_id UUID
- variant_number INT
- scene_plan JSONB
- prompt_used TEXT
- provider VARCHAR
- model VARCHAR
- duration INT
- resolution VARCHAR
- status VARCHAR (pending, processing, completed, failed)
- video_url VARCHAR
- thumbnail_url VARCHAR
- file_size BIGINT
- credits_used INT
- external_job_id VARCHAR
- revision_of UUID (self-reference for regenerations)
- error_message TEXT
- created_at TIMESTAMP
- updated_at TIMESTAMP
- deleted_at TIMESTAMP
```

### SceneGeneration Table
```sql
- id UUID PRIMARY KEY
- variant_id UUID (FOREIGN KEY)
- scene_number INT
- scene_index INT
- prompt TEXT
- duration INT
- status VARCHAR (pending, processing, completed, failed)
- external_job_id VARCHAR
- video_url VARCHAR
- error_message TEXT
- processing_notes JSONB
- created_at TIMESTAMP
- updated_at TIMESTAMP
- deleted_at TIMESTAMP
```

## Error Handling

### Common Errors

| Error | Cause | Handler Response |
|-------|-------|------------------|
| Insufficient Credits | User balance < required | 400 Bad Request |
| Job Not Found | Invalid job ID | 404 Not Found |
| Invalid Format | Request parsing fails | 400 Bad Request |
| Generation Timeout | Job > 2 hours | Auto-fail, update DB |
| Provider Error | API call fails | Retry up to 3 times |
| Database Error | DB operation fails | 500 Internal Error |

### Retry Strategy

- Automatic retry up to 3 times (configurable)
- Exponential backoff recommended in production
- Failed jobs marked with error message
- User can manually retry or regenerate

## Performance Considerations

### Optimization Tips

1. **Database Indexing**:
   - Index on `job.status` for queue queries
   - Index on `job.project_id`, `job.user_id` for lookups
   - Index on `variant.storyboard_id` for batch retrieval

2. **Worker Tuning**:
   - Default: 3 workers for ~100 concurrent users
   - Increase workers for higher concurrency
   - Monitor polling latency and adjust poll_interval

3. **Scene Settings**:
   - 2 scenes: faster generation, lower cost
   - 3 scenes: longer videos, higher cost
   - Adjustable per request

4. **Caching**:
   - Cache provider responses
   - Cache processed videos in CDN
   - TTL-based stale file cleanup

## Future Enhancements

- [ ] Multi-provider load balancing
- [ ] Intelligence in provider selection (based on speed/cost)
- [ ] Media composition pipeline (merging scenes into final video)
- [ ] Audio & music integration
- [ ] Subtitle generation
- [ ] Advanced analytics and reporting
- [ ] Webhook notifications for job completion
- [ ] Batch generation for multiple projects

## Related Files

- [API Documentation](../docs/API_DOCUMENTATION.md)
- [Model Definitions](../internal/model/generation_job.go)
- [Service Implementation](../internal/service/video_generation_service.go)
- [Queue Implementation](../internal/queue/job_queue.go)
- [Provider Implementations](../internal/ai/providers.go)
