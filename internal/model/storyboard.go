package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== Storyboard ====================

type Storyboard struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	IsSelected  bool           `gorm:"default:false" json:"is_selected"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
	Scenes  []Scene `gorm:"foreignKey:StoryboardID" json:"scenes,omitempty"`
}

func (s *Storyboard) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

// ==================== Scene ====================

type Scene struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StoryboardID uuid.UUID      `gorm:"type:uuid;not null;index" json:"storyboard_id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	SceneNumber  int            `gorm:"not null" json:"scene_number"`
	Title        string         `json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	VisualDesc   string         `gorm:"type:text" json:"visual_description"`
	Duration     int            `json:"duration"` // in seconds
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Storyboard Storyboard `gorm:"foreignKey:StoryboardID" json:"-"`
}

func (sc *Scene) BeforeCreate(tx *gorm.DB) error {
	sc.ID = uuid.New()
	return nil
}

// ==================== Request Types ====================

type GenerateStoryboardRequest struct {
	ProjectID      string `json:"project_id" validate:"required"`
	ContentThemeID string `json:"content_theme_id" validate:"required"`
}

type SelectStoryboardRequest struct {
	StoryboardID string `json:"storyboard_id" validate:"required"`
}
