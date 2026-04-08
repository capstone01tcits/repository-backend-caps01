package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== Business Brief ====================

type BusinessBrief struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProjectID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	ProjectName      string         `gorm:"not null" json:"project_name"`
	CompanyName      string         `json:"company_name"`
	InstituteName    string         `json:"institute_name"`
	Education        string         `json:"education"`
	Industry         string         `json:"industry"`
	TargetAudience   string         `json:"target_audience"`
	ProjectObjective string         `gorm:"type:text" json:"project_objective"`
	KeyMessage       string         `gorm:"type:text" json:"key_message"`
	Budget           string         `json:"budget"`
	Timeline         string         `json:"timeline"`
	Deadline         time.Time      `json:"deadline"`
	Competitors      string         `gorm:"type:text" json:"competitors"`
	AdditionalNotes  string         `gorm:"type:text" json:"additional_notes"`
	Status           string         `gorm:"default:draft" json:"status"` // draft, submitted, approved, rejected
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User           User            `gorm:"foreignKey:UserID" json:"-"`
	Project        Project         `gorm:"foreignKey:ProjectID" json:"-"`
	CreativeBriefs []CreativeBrief `gorm:"foreignKey:BusinessBriefID" json:"creative_briefs,omitempty"`
}

func (b *BusinessBrief) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}

// ==================== Creative Brief ====================

type CreativeBrief struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	BusinessBriefID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"business_brief_id"`
	Title            string         `gorm:"not null" json:"title"`
	VideoType        string         `json:"video_type"` // promotional, educational, testimonial, explainer, etc.
	Duration         int            `json:"duration"`   // in seconds
	Style            string         `json:"style"`      // cinematic, animated, minimalist, etc.
	Tone             string         `json:"tone"`       // professional, casual, energetic, etc.
	Script           string         `gorm:"type:text" json:"script"`
	Storyboard       string         `gorm:"type:text" json:"storyboard"`
	VisualReferences string         `gorm:"type:text" json:"visual_references"`
	MusicPreference  string         `json:"music_preference"`
	CallToAction     string         `json:"call_to_action"`
	OutputFormat     string         `json:"output_format"` // mp4, webm, etc.
	Resolution       string         `json:"resolution"`    // 1080p, 4K, etc.
	AdditionalNotes  string         `gorm:"type:text" json:"additional_notes"`
	Status           string         `gorm:"default:draft" json:"status"` // draft, submitted, in_production, completed
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User          User          `gorm:"foreignKey:UserID" json:"-"`
	BusinessBrief BusinessBrief `gorm:"foreignKey:BusinessBriefID" json:"-"`
}

func (c *CreativeBrief) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}

// ==================== Request Types ====================

type CreateBusinessBriefRequest struct {
	ProjectID        string    `json:"project_id" validate:"required"`
	ProjectName      string    `json:"project_name" validate:"required"`
	CompanyName      string    `json:"company_name"`
	InstituteName    string    `json:"institute_name"`
	Education        string    `json:"education"`
	Industry         string    `json:"industry"`
	TargetAudience   string    `json:"target_audience"`
	ProjectObjective string    `json:"project_objective"`
	KeyMessage       string    `json:"key_message"`
	Budget           string    `json:"budget"`
	Timeline         string    `json:"timeline"`
	Deadline         time.Time `json:"deadline"`
	Competitors      string    `json:"competitors"`
	AdditionalNotes  string    `json:"additional_notes"`
}

type UpdateBusinessBriefRequest struct {
	ProjectName      *string    `json:"project_name"`
	CompanyName      *string    `json:"company_name"`
	InstituteName    *string    `json:"institute_name"`
	Education        *string    `json:"education"`
	Industry         *string    `json:"industry"`
	TargetAudience   *string    `json:"target_audience"`
	ProjectObjective *string    `json:"project_objective"`
	KeyMessage       *string    `json:"key_message"`
	Budget           *string    `json:"budget"`
	Timeline         *string    `json:"timeline"`
	Deadline         *time.Time `json:"deadline"`
	Competitors      *string    `json:"competitors"`
	AdditionalNotes  *string    `json:"additional_notes"`
	Status           *string    `json:"status"`
}

type CreateCreativeBriefRequest struct {
	BusinessBriefID  string `json:"business_brief_id" validate:"required"`
	Title            string `json:"title" validate:"required"`
	VideoType        string `json:"video_type"`
	Duration         int    `json:"duration"`
	Style            string `json:"style"`
	Tone             string `json:"tone"`
	Script           string `json:"script"`
	Storyboard       string `json:"storyboard"`
	VisualReferences string `json:"visual_references"`
	MusicPreference  string `json:"music_preference"`
	CallToAction     string `json:"call_to_action"`
	OutputFormat     string `json:"output_format"`
	Resolution       string `json:"resolution"`
	AdditionalNotes  string `json:"additional_notes"`
}

type UpdateCreativeBriefRequest struct {
	Title            *string `json:"title"`
	VideoType        *string `json:"video_type"`
	Duration         *int    `json:"duration"`
	Style            *string `json:"style"`
	Tone             *string `json:"tone"`
	Script           *string `json:"script"`
	Storyboard       *string `json:"storyboard"`
	VisualReferences *string `json:"visual_references"`
	MusicPreference  *string `json:"music_preference"`
	CallToAction     *string `json:"call_to_action"`
	OutputFormat     *string `json:"output_format"`
	Resolution       *string `json:"resolution"`
	AdditionalNotes  *string `json:"additional_notes"`
	Status           *string `json:"status"`
}
