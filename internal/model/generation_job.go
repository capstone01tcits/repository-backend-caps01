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
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProjectID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	StoryboardID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"storyboard_id"`
	VideoID         *uuid.UUID     `gorm:"type:uuid;index" json:"video_id,omitempty"`
	JobType         string         `gorm:"not null;index" json:"job_type"` // generate, regenerate, regenerate_scene
	Status          string         `gorm:"default:queued;index" json:"status"` // queued, processing, completed, failed
	Priority        int            `gorm:"default:0" json:"priority"`
	Prompt          datatypes.JSON `gorm:"type:jsonb" json:"prompt"`
	SceneCount      int            `gorm:"default:2" json:"scene_count"`
	VideoDuration   int            `gorm:"default:10" json:"video_duration"` // in seconds
	Provider        string         `json:"provider"`    // wavespeed, wan, open_source
	Model           string         `json:"model"`       // gen4.5, wan2.1, etc
	Resolution      string         `gorm:"default:1080p" json:"resolution"`
	ProcessingNotes datatypes.JSON `gorm:"type:jsonb" json:"processing_notes"`
	ErrorMessage    string         `gorm:"type:text" json:"error_message,omitempty"`
	CreditsRequired int            `json:"credits_required"`
	CreditsUsed     int            `json:"credits_used"`
	RetryCount      int            `gorm:"default:0" json:"retry_count"`
	MaxRetries      int            `gorm:"default:3" json:"max_retries"`
	StartedAt       *time.Time     `gorm:"type:timestamp" json:"started_at,omitempty"`
	CompletedAt     *time.Time     `gorm:"type:timestamp" json:"completed_at,omitempty"`
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
	// VideoID is left NULL initially since videos are created during generation process
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



type VideoStatusResponse struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	VideoURL      string    `json:"video_url,omitempty"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	Duration      int       `json:"duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GenerationJobResponse struct {
	ID              string              `json:"id"`
	Status          string              `json:"status"`
	JobType         string              `json:"job_type"`
	Progress        int                 `json:"progress"` // percentage
	Video           *VideoStatusResponse `json:"video,omitempty"`
	ErrorMessage    string              `json:"error_message,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
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
