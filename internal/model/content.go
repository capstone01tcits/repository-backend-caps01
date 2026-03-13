package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== Content Pillar ====================

type ContentPillar struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	IsSelected  bool           `gorm:"default:false" json:"is_selected"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	Project       Project        `gorm:"foreignKey:ProjectID" json:"-"`
	ContentThemes []ContentTheme `gorm:"foreignKey:ContentPillarID" json:"content_themes,omitempty"`
}

func (cp *ContentPillar) BeforeCreate(tx *gorm.DB) error {
	cp.ID = uuid.New()
	return nil
}

// ==================== Content Theme ====================

type ContentTheme struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ContentPillarID uuid.UUID      `gorm:"type:uuid;not null;index" json:"content_pillar_id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title           string         `gorm:"not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	IsSelected      bool           `gorm:"default:false" json:"is_selected"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User          User          `gorm:"foreignKey:UserID" json:"-"`
	ContentPillar ContentPillar `gorm:"foreignKey:ContentPillarID" json:"-"`
}

func (ct *ContentTheme) BeforeCreate(tx *gorm.DB) error {
	ct.ID = uuid.New()
	return nil
}

// ==================== Request Types ====================

type GenerateContentPillarsRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
}

type SelectContentPillarRequest struct {
	ContentPillarID string `json:"content_pillar_id" validate:"required"`
}

type SelectContentThemeRequest struct {
	ContentThemeID string `json:"content_theme_id" validate:"required"`
}
