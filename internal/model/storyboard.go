package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== Storyboard ====================

type Storyboard struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title         string         `gorm:"not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	TotalDuration int            `json:"total_duration"` // in seconds (30, 45, 60, etc)
	Style         string         `json:"style"`          // template style name (e.g., "dynamic", "narrative", "energetic")
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relations
	User     User                `gorm:"foreignKey:UserID" json:"-"`
	Project  *Project             `gorm:"foreignKey:ProjectID" json:"-"`
	Sections []StoryboardSection `gorm:"foreignKey:StoryboardID" json:"sections,omitempty"`
}

func (s *Storyboard) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// ==================== StoryboardSection ====================
// Manual 3-part storyboard structure: Hook/Intro, Value/Highlight, CTA

type StoryboardSection struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StoryboardID uuid.UUID      `gorm:"type:uuid;not null;index" json:"storyboard_id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	SectionType  string         `gorm:"type:varchar(50);not null" json:"section_type"` // "hook", "value", "cta"
	Content      string         `gorm:"type:text;not null" json:"content"`             // user-written content
	Duration     int            `json:"duration"`                                      // in seconds for this section
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Storyboard *Storyboard `gorm:"foreignKey:StoryboardID" json:"-"`
}

func (ss *StoryboardSection) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return nil
}

// ==================== Request Types ====================

type GenerateStoryboardTemplatesRequest struct {
	ProjectID     string `json:"project_id" validate:"required"`
	VideoDuration int    `json:"video_duration" validate:"required,min=15,max=300"` // in seconds
}

type CreateManualStoryboardRequest struct {
	ProjectID   string                         `json:"project_id" validate:"required"`
	Title       string                         `json:"title" validate:"required"`
	Description string                         `json:"description"`
	Style       string                         `json:"style"`
	Duration    int                            `json:"duration"` // total video duration
	Sections    []CreateStoryboardSectionInput `json:"sections" validate:"required,len=3"`
}

type CreateStoryboardSectionInput struct {
	SectionType string `json:"section_type" validate:"required,oneof=hook value cta"` // "hook", "value", "cta"
	Content     string `json:"content" validate:"required"`
	Duration    int    `json:"duration"` // duration for this section
}

type StoryboardTemplate struct {
	TemplateID  string                      `json:"template_id"` // internal id for this template variant
	Style       string                      `json:"style"`       // e.g., "Dynamic", "Narrative", "Energetic"
	Description string                      `json:"description"` // human-readable description
	Duration    int                         `json:"duration"`    // total duration in seconds
	Sections    []TemplateStoryboardSection `json:"sections"`    // 3 sections with suggested content
}

type TemplateStoryboardSection struct {
	SectionType       string `json:"section_type"`       // "hook", "value", "cta"
	Title             string `json:"title"`              // section title for display
	SuggestedDuration int    `json:"suggested_duration"` // recommended duration for this section
	Content           string `json:"content"`            // suggested content based on project data
	Tips              string `json:"tips"`               // tips for this section
}

type UpdateManualStoryboardRequest struct {
	Title       *string                        `json:"title"`
	Description *string                        `json:"description"`
	Style       *string                        `json:"style"`
	Sections    []UpdateStoryboardSectionInput `json:"sections"`
}

type UpdateStoryboardSectionInput struct {
	SectionType string `json:"section_type"` // "hook", "value", "cta"
	Content     string `json:"content"`
	Duration    int    `json:"duration"`
}

type SelectStoryboardRequest struct {
	StoryboardID string `json:"storyboard_id" validate:"required"`
}

type Veo3TestPayload struct {
	Prompt          string   `json:"prompt"`

}
