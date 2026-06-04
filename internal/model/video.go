package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	StoryboardID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"storyboard_id"`
	Title           string         `gorm:"not null" json:"title"`
	Status          string         `gorm:"default:pending" json:"status"` // pending, processing, completed, failed
	VideoURL        string         `json:"video_url"`
	ThumbnailURL    string         `json:"thumbnail_url"`
	Duration        int            `json:"duration"`      // total duration in seconds
	Format          string         `gorm:"default:mp4" json:"format"`
	Resolution      string         `gorm:"default:1080p" json:"resolution"`
	FileSize        int64          `json:"file_size"`     // in bytes
	CreditsUsed     int            `gorm:"default:1" json:"credits_used"`
	RegenerateCount int            `gorm:"default:0" json:"regenerate_count"` // max 3
	ErrorMessage    string         `gorm:"type:text" json:"error_message,omitempty"`
	ExternalJobID   string         `json:"external_job_id,omitempty"` // job ID from provider (Wavespeed)
	SectionType     string         `gorm:"default:''" json:"section_type"` // hook, value, cta
	SceneIndex      int            `gorm:"default:0" json:"scene_index"`   // 1, 2, 3
	VideoMode       string         `gorm:"default:'text-to-video'" json:"video_mode"` // text-to-video | image-to-video | start-end-to-video
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Project    Project    `gorm:"foreignKey:ProjectID" json:"-"`
	Storyboard Storyboard `gorm:"foreignKey:StoryboardID" json:"-"`
}

func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// ==================== Request Types ====================

// Video requests are defined in generation_job.go
