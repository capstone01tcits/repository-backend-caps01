package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"default:'draft'" json:"status"` // draft, in_progress, completed
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User            User             `gorm:"foreignKey:UserID" json:"-"`
	BusinessBriefs  []BusinessBrief  `gorm:"foreignKey:ProjectID" json:"business_briefs,omitempty"`
	ContentPillars  []ContentPillar  `gorm:"foreignKey:ProjectID" json:"content_pillars,omitempty"`
	Storyboards     []Storyboard     `gorm:"foreignKey:ProjectID" json:"storyboards,omitempty"`
	Videos          []Video          `gorm:"foreignKey:ProjectID" json:"videos,omitempty"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

// ==================== Request Types ====================

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}
