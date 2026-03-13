package migrations

import "gorm.io/gorm"

// MigrateVideoGeneration creates tables for video generation system
func MigrateVideoGeneration(db *gorm.DB) error {
	return db.AutoMigrate(
		&GenerationJob{},
		&VideoVariant{},
		&SceneGeneration{},
	)
}

// ============================================================================
// Models for migration (duplicated from internal/model for migration context)
// ============================================================================

type GenerationJob struct {
	ID              string `gorm:"type:uuid;primary_key"`
	UserID          string `gorm:"type:uuid;not null;index"`
	ProjectID       string `gorm:"type:uuid;not null;index"`
	StoryboardID    string `gorm:"type:uuid;not null;index"`
	VideoID         string `gorm:"type:uuid;index"`
	JobType         string `gorm:"not null;index"` // generate, regenerate, regenerate_scene
	Status          string `gorm:"default:'queued';index"` // queued, processing, completed, failed
	Priority        int    `gorm:"default:0"`
	Prompt          interface{} `gorm:"type:jsonb"`
	SceneCount      int    `gorm:"default:2"`
	VideoDuration   int    `gorm:"default:10"` // in seconds
	Provider        string // ltx, runway, wan, open_source
	Model           string // specific model name
	Resolution      string `gorm:"default:'1080p'"`
	ProcessingNotes interface{} `gorm:"type:jsonb"`
	ErrorMessage    string `gorm:"type:text"`
	CreditsRequired int
	CreditsUsed     int
	RetryCount      int `gorm:"default:0"`
	MaxRetries      int `gorm:"default:3"`
	StartedAt       *int64
	CompletedAt     *int64
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       *int64 `gorm:"index"`
}

type VideoVariant struct {
	ID               string `gorm:"type:uuid;primary_key"`
	UserID           string `gorm:"type:uuid;not null;index"`
	ProjectID        string `gorm:"type:uuid;not null;index"`
	StoryboardID     string `gorm:"type:uuid;not null;index"`
	VariantNumber    int    `gorm:"not null"` // 1, 2, or 3
	ScenePlan        interface{} `gorm:"type:jsonb"`
	PromptUsed       string `gorm:"type:text"`
	Provider         string
	Model            string
	Duration         int    `json:"duration"`
	Resolution       string `gorm:"default:'1080p'"`
	Status           string `gorm:"default:'pending'"` // pending, processing, completed, failed
	VideoURL         string
	ThumbnailURL     string
	FileSize         int64
	CreditsUsed      int
	ErrorMessage     string `gorm:"type:text"`
	RevisionOf       *string // uuid reference
	ExternalJobID    string
	CreatedAt        int64
	UpdatedAt        int64
	DeletedAt        *int64 `gorm:"index"`
}

type SceneGeneration struct {
	ID               string `gorm:"type:uuid;primary_key"`
	VariantID        string `gorm:"type:uuid;not null;index"`
	SceneNumber      int
	SceneIndex       int
	Prompt           string `gorm:"type:text"`
	Duration         int
	Status           string `gorm:"default:'pending'"` // pending, processing, completed, failed
	ExternalJobID    string
	VideoURL         string
	ErrorMessage     string `gorm:"type:text"`
	ProcessingNotes  interface{} `gorm:"type:jsonb"`
	CreatedAt        int64
	UpdatedAt        int64
	DeletedAt        *int64 `gorm:"index"`
}

// ============================================================================
// Raw SQL Migration (Alternative approach)
// ============================================================================

const createGenerationJobTable = `
CREATE TABLE IF NOT EXISTS generation_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    project_id UUID NOT NULL,
    storyboard_id UUID NOT NULL,
    video_id UUID,
    job_type VARCHAR NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'queued',
    priority INTEGER DEFAULT 0,
    prompt JSONB,
    scene_count INTEGER DEFAULT 2,
    video_duration INTEGER DEFAULT 10,
    provider VARCHAR,
    model VARCHAR,
    resolution VARCHAR DEFAULT '1080p',
    processing_notes JSONB,
    error_message TEXT,
    credits_required INTEGER,
    credits_used INTEGER,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_generation_job_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_generation_job_project FOREIGN KEY(project_id) REFERENCES projects(id),
    CONSTRAINT fk_generation_job_storyboard FOREIGN KEY(storyboard_id) REFERENCES storyboards(id)
);

CREATE INDEX idx_generation_jobs_user_id ON generation_jobs(user_id);
CREATE INDEX idx_generation_jobs_project_id ON generation_jobs(project_id);
CREATE INDEX idx_generation_jobs_storyboard_id ON generation_jobs(storyboard_id);
CREATE INDEX idx_generation_jobs_status ON generation_jobs(status);
CREATE INDEX idx_generation_jobs_job_type ON generation_jobs(job_type);
CREATE INDEX idx_generation_jobs_deleted_at ON generation_jobs(deleted_at);
`

const createVideoVariantTable = `
CREATE TABLE IF NOT EXISTS video_variants (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    project_id UUID NOT NULL,
    storyboard_id UUID NOT NULL,
    variant_number INTEGER NOT NULL,
    scene_plan JSONB,
    prompt_used TEXT,
    provider VARCHAR,
    model VARCHAR,
    duration INTEGER,
    resolution VARCHAR DEFAULT '1080p',
    status VARCHAR DEFAULT 'pending',
    video_url VARCHAR,
    thumbnail_url VARCHAR,
    file_size BIGINT,
    credits_used INTEGER,
    error_message TEXT,
    revision_of UUID,
    external_job_id VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_video_variant_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_video_variant_project FOREIGN KEY(project_id) REFERENCES projects(id),
    CONSTRAINT fk_video_variant_storyboard FOREIGN KEY(storyboard_id) REFERENCES storyboards(id),
    CONSTRAINT fk_video_variant_revision FOREIGN KEY(revision_of) REFERENCES video_variants(id)
);

CREATE INDEX idx_video_variants_user_id ON video_variants(user_id);
CREATE INDEX idx_video_variants_project_id ON video_variants(project_id);
CREATE INDEX idx_video_variants_storyboard_id ON video_variants(storyboard_id);
CREATE INDEX idx_video_variants_status ON video_variants(status);
CREATE INDEX idx_video_variants_deleted_at ON video_variants(deleted_at);
CREATE UNIQUE INDEX idx_video_variants_storyboard_variant ON video_variants(storyboard_id, variant_number) 
    WHERE deleted_at IS NULL;
`

const createSceneGenerationTable = `
CREATE TABLE IF NOT EXISTS scene_generations (
    id UUID PRIMARY KEY,
    variant_id UUID NOT NULL,
    scene_number INTEGER,
    scene_index INTEGER,
    prompt TEXT,
    duration INTEGER,
    status VARCHAR DEFAULT 'pending',
    external_job_id VARCHAR,
    video_url VARCHAR,
    error_message TEXT,
    processing_notes JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_scene_generation_variant FOREIGN KEY(variant_id) REFERENCES video_variants(id)
);

CREATE INDEX idx_scene_generations_variant_id ON scene_generations(variant_id);
CREATE INDEX idx_scene_generations_status ON scene_generations(status);
CREATE INDEX idx_scene_generations_deleted_at ON scene_generations(deleted_at);
CREATE INDEX idx_scene_generations_scene_number ON scene_generations(variant_id, scene_number);
`

// RunRawMigrations executes raw SQL migrations
func RunRawMigrations(db *gorm.DB) error {
	if err := db.Exec(createGenerationJobTable).Error; err != nil {
		return err
	}
	if err := db.Exec(createVideoVariantTable).Error; err != nil {
		return err
	}
	if err := db.Exec(createSceneGenerationTable).Error; err != nil {
		return err
	}
	return nil
}

// Rollback functions for migration reversal

const dropGenerationJobTable = `DROP TABLE IF EXISTS generation_jobs CASCADE;`
const dropVideoVariantTable = `DROP TABLE IF EXISTS video_variants CASCADE;`
const dropSceneGenerationTable = `DROP TABLE IF EXISTS scene_generations CASCADE;`

func RollbackVideoGeneration(db *gorm.DB) error {
	if err := db.Exec(dropSceneGenerationTable).Error; err != nil {
		return err
	}
	if err := db.Exec(dropVideoVariantTable).Error; err != nil {
		return err
	}
	if err := db.Exec(dropGenerationJobTable).Error; err != nil {
		return err
	}
	return nil
}
