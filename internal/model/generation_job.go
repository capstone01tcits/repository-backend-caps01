package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ==================== GenerationJob ====================
// Represents a video generation task in the queue

type GenerationJob struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProjectID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	StoryboardID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"storyboard_id"`
	VideoID         uuid.UUID      `gorm:"type:uuid;index" json:"video_id"`
	JobType         string         `gorm:"not null;index" json:"job_type"` // generate, regenerate, regenerate_scene
	Status          string         `gorm:"default:'queued';index" json:"status"` // queued, processing, completed, failed
	Priority        int            `gorm:"default:0" json:"priority"`
	Prompt          datatypes.JSON `gorm:"type:jsonb" json:"prompt"`
	SceneCount      int            `gorm:"default:2" json:"scene_count"`
	VideoDuration   int            `gorm:"default:10" json:"video_duration"` // in seconds
	Provider        string         `json:"provider"`    // ltx, runway, wan, open_source
	Model           string         `json:"model"`       // ltx-2-fast, gen4.5, wan2.1, etc
	Resolution      string         `gorm:"default:'1080p'" json:"resolution"`
	ProcessingNotes datatypes.JSON `gorm:"type:jsonb" json:"processing_notes"`
	ErrorMessage    string         `gorm:"type:text" json:"error_message,omitempty"`
	CreditsRequired int            `json:"credits_required"`
	CreditsUsed     int            `json:"credits_used"`
	RetryCount      int            `gorm:"default:0" json:"retry_count"`
	MaxRetries      int            `gorm:"default:3" json:"max_retries"`
	StartedAt       *time.Time     `json:"started_at,omitempty"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Project    Project    `gorm:"foreignKey:ProjectID" json:"-"`
	Storyboard Storyboard `gorm:"foreignKey:StoryboardID" json:"-"`
	Video      *Video     `gorm:"foreignKey:VideoID" json:"-"`
}

func (gj *GenerationJob) BeforeCreate(tx *gorm.DB) error {
	gj.ID = uuid.New()
	if gj.VideoID == uuid.Nil {
		gj.VideoID = uuid.New()
	}
	return nil
}

// ==================== VideoVariant ====================
// Represents one of 3 video variants generated from a briefing

type VideoVariant struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	StoryboardID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"storyboard_id"`
	VariantNumber  int            `gorm:"not null" json:"variant_number"` // 1, 2, or 3
	ScenePlan      datatypes.JSON `gorm:"type:jsonb" json:"scene_plan"`
	PromptUsed     string         `gorm:"type:text" json:"prompt_used"`
	Provider       string         `json:"provider"`    // ltx, runway, wan, open_source
	Model          string         `json:"model"`       // specific model name
	Duration       int            `json:"duration"`    // in seconds
	Resolution     string         `gorm:"default:'1080p'" json:"resolution"`
	Status         string         `gorm:"default:'pending'" json:"status"` // pending, processing, completed, failed
	VideoURL       string         `json:"video_url"`
	ThumbnailURL   string         `json:"thumbnail_url"`
	FileSize       int64          `json:"file_size"`
	CreditsUsed    int            `json:"credits_used"`
	ErrorMessage   string         `gorm:"type:text" json:"error_message,omitempty"`
	RevisionOf     *uuid.UUID     `json:"revision_of,omitempty"` // if this is a regenerated version
	ExternalJobID  string         `json:"external_job_id"`        // job ID from provider (LTX, Runway, etc)
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Project    Project    `gorm:"foreignKey:ProjectID" json:"-"`
	Storyboard Storyboard `gorm:"foreignKey:StoryboardID" json:"-"`
}

func (vv *VideoVariant) BeforeCreate(tx *gorm.DB) error {
	vv.ID = uuid.New()
	return nil
}

// ==================== SceneGeneration ====================
// Tracks individual scene generation within a video

type SceneGeneration struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	VariantID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"variant_id"`
	SceneNumber     int            `json:"scene_number"`
	SceneIndex      int            `json:"scene_index"`
	Prompt          string         `gorm:"type:text" json:"prompt"`
	Duration        int            `json:"duration"` // in seconds
	Status          string         `gorm:"default:'pending'" json:"status"` // pending, processing, completed, failed
	ExternalJobID   string         `json:"external_job_id"`
	VideoURL        string         `json:"video_url"`
	ErrorMessage    string         `gorm:"type:text" json:"error_message,omitempty"`
	ProcessingNotes datatypes.JSON `gorm:"type:jsonb" json:"processing_notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Variant VideoVariant `gorm:"foreignKey:VariantID" json:"-"`
}

func (sg *SceneGeneration) BeforeCreate(tx *gorm.DB) error {
	sg.ID = uuid.New()
	return nil
}

// ==================== Request Types ====================

type GenerateVideoRequest struct {
	ProjectID    string `json:"project_id" validate:"required"`
	StoryboardID string `json:"storyboard_id" validate:"required"`
	CustomPrompt string `json:"custom_prompt,omitempty"`
	SceneCount   int    `json:"scene_count,omitempty"` // default 2-3
	VideoDuration int   `json:"video_duration,omitempty"` // default 8-12
}

type RegenerateVideoRequest struct {
	VideoVariantID string `json:"video_variant_id" validate:"required"`
	NewPrompt      string `json:"new_prompt,omitempty"`
}

type RegenerateSceneRequest struct {
	SceneGenerationID string `json:"scene_generation_id" validate:"required"`
	NewPrompt         string `json:"new_prompt,omitempty"`
}

type VideoStatusResponse struct {
	ID            string                    `json:"id"`
	VariantNumber int                       `json:"variant_number"`
	Status        string                    `json:"status"`
	VideoURL      string                    `json:"video_url,omitempty"`
	ThumbnailURL  string                    `json:"thumbnail_url,omitempty"`
	PromptUsed    string                    `json:"prompt_used"`
	Duration      int                       `json:"duration"`
	Scenes        []SceneStatusResponse     `json:"scenes,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type SceneStatusResponse struct {
	ID            string    `json:"id"`
	SceneNumber   int       `json:"scene_number"`
	Status        string    `json:"status"`
	VideoURL      string    `json:"video_url,omitempty"`
	Duration      int       `json:"duration"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GenerationJobResponse struct {
	ID              string            `json:"id"`
	Status          string            `json:"status"`
	JobType         string            `json:"job_type"`
	Progress        int               `json:"progress"` // percentage
	VideoVariants   []VideoStatusResponse `json:"video_variants,omitempty"`
	ErrorMessage    string            `json:"error_message,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// ==================== Prompt Templates ====================

type PromptGenerationRequest struct {
	StoryboardID string `json:"storyboard_id" validate:"required"`
	Variation    int    `json:"variation"` // 1, 2, or 3
}

type PromptData struct {
	SceneDescriptions []string `json:"scene_descriptions"`
	KeyMessages       []string `json:"key_messages"`
	TargetAudience    string   `json:"target_audience"`
	CampaignTheme     string   `json:"campaign_theme"`
}

func (pd *PromptData) MarshalJSON() ([]byte, error) {
	return json.Marshal(pd)
}

func (pd *PromptData) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, pd)
}

// Scan implements sql.Scanner interface
func (pd *PromptData) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &pd)
}

// Value implements driver.Valuer interface
func (pd *PromptData) Value() (driver.Value, error) {
	return json.Marshal(pd)
}
