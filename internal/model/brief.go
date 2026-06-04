package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== Business Brief ====================

type BusinessBrief struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProjectID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	InstitutionName    string         `gorm:"default:'-'" json:"institution_name"`
	InstitutionHistory string         `gorm:"type:text" json:"institution_history"`
	SchoolLevel        string         `json:"school_level"`
	OfferedDegrees     string         `gorm:"type:text" json:"offered_degrees"`
	LogoPath           string         `json:"logo_path"`
	EnvironmentPath    string         `json:"environment_path"`
	DocumentPath       string         `json:"document_path"`
	Status             string         `gorm:"default:draft" json:"status"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User           User            `gorm:"foreignKey:UserID" json:"-"`
	Project        Project         `gorm:"foreignKey:ProjectID" json:"-"`
	CreativeBriefs []CreativeBrief `gorm:"foreignKey:BusinessBriefID" json:"creative_briefs,omitempty"`
}

func (b *BusinessBrief) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// ==================== Creative Brief ====================

type CreativeBrief struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	BusinessBriefID uuid.UUID      `gorm:"type:uuid;not null;index" json:"business_brief_id"`
	Title           string         `json:"title"` // Internal name
	EventContent    string         `gorm:"not null" json:"event_content"`
	VideoDuration   string         `json:"video_duration"`
	ToneOfVoice     string         `json:"tone_of_voice"`
	KeyMessage      string         `gorm:"type:text" json:"key_message"`
	Prompt          string         `gorm:"type:text" json:"prompt"`
	Theme           string         `json:"theme"`
	Copywriting     string         `gorm:"type:text" json:"copywriting"`
	Hashtags        string         `gorm:"type:text" json:"hashtags"`
	Status          string         `gorm:"default:draft" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User          User          `gorm:"foreignKey:UserID" json:"-"`
	BusinessBrief BusinessBrief `gorm:"foreignKey:BusinessBriefID" json:"-"`
}

func (c *CreativeBrief) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// ==================== Request Types ====================

type CreateBusinessBriefRequest struct {
	ProjectID          string `json:"project_id" validate:"required"`
	InstitutionName    string `json:"institution_name" validate:"required"`
	InstitutionHistory string `json:"institution_history"`
	SchoolLevel        string `json:"school_level"`
	OfferedDegrees     string `json:"offered_degrees"`
}

type UpdateBusinessBriefRequest struct {
	InstitutionName    *string `json:"institution_name"`
	InstitutionHistory *string `json:"institution_history"`
	SchoolLevel        *string `json:"school_level"`
	OfferedDegrees     *string `json:"offered_degrees"`
	Status             *string `json:"status"`
}

type CreateCreativeBriefRequest struct {
	BusinessBriefID string `json:"business_brief_id" validate:"required"`
	EventContent    string `json:"event_content" validate:"required"`
	VideoDuration   string `json:"video_duration"`
	ToneOfVoice     string `json:"tone_of_voice"`
	KeyMessage      string `json:"key_message"`
	Prompt          string `json:"prompt"`
	Theme           string `json:"theme"`
	Copywriting     string `json:"copywriting"`
	Hashtags        string `json:"hashtags"`
}

type UpdateCreativeBriefRequest struct {
	EventContent  *string `json:"event_content"`
	VideoDuration *string `json:"video_duration"`
	ToneOfVoice   *string `json:"tone_of_voice"`
	KeyMessage    *string `json:"key_message"`
	Prompt        *string `json:"prompt"`
	Theme         *string `json:"theme"`
	Copywriting   *string `json:"copywriting"`
	Hashtags      *string `json:"hashtags"`
	Status        *string `json:"status"`
}

// ==================== Simplified FE Request (Matches Frontend Exactly) ====================

type CreateProjectFromFERequest struct {
	ProjectID          string `json:"project_id"` // optional - for updating existing
	ProjectName        string `json:"project_name"` // optional - custom name
	// Step 1: Business Brief
	InstitutionName    string `json:"institution_name" validate:"required"`
	InstitutionHistory string `json:"institution_history"` // optional - not always sent by FE
	SchoolLevel        string `json:"school_level"`        // optional - not always sent by FE
	OfferedDegrees     string `json:"offered_degrees"`     // optional

	// Step 2: Creative Brief
	EventContent       string `json:"event_content" validate:"required"`
	ToneOfVoice        string `json:"tone_of_voice" validate:"required"`
	SelectedKeyMessage string `json:"selected_key_message" validate:"required"`
	VideoDuration      string `json:"video_duration"` // optional - not always sent by FE (e.g., "15 detik", "30 detik", "60 detik")
	Prompt             string `json:"prompt"`         // optional

	// Step 3: Theme
	SelectedTheme string `json:"selected_theme" validate:"required"`

	// Step 4: Summary
	EditableCopywriting string `json:"editable_copywriting"` // optional
	EditableHashtags    string `json:"editable_hashtags"`    // optional - not always sent by FE

	// Images (Base64 encoded)
	LogoBase64     string `json:"logo_base64"`     // optional - institution logo
	EnvBase64      string `json:"env_base64"`      // optional - environment photo
	DocumentBase64 string `json:"document_base64"` // optional - pdf/doc about institution
}
