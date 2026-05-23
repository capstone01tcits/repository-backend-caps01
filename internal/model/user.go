package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role          string         `gorm:"default:user" json:"role"`   // user, admin
	Credits       int            `gorm:"default:1000" json:"credits"`
	EmailAlerts   bool           `gorm:"default:true" json:"email_alerts"`
	Newsletter    bool           `gorm:"default:false" json:"newsletter"`
	PublicProfile bool           `gorm:"default:false" json:"public_profile"`
	DataTraining  bool           `gorm:"default:true" json:"data_training"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

// Request & Response types

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	User         UserInfo  `json:"user"`
}

type UserInfo struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	Credits       int       `json:"credits"`
	EmailAlerts   bool      `json:"email_alerts"`
	Newsletter    bool      `json:"newsletter"`
	PublicProfile bool      `json:"public_profile"`
	DataTraining  bool      `json:"data_training"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdatePreferencesRequest struct {
	EmailAlerts   *bool `json:"email_alerts"`
	Newsletter    *bool `json:"newsletter"`
	PublicProfile *bool `json:"public_profile"`
	DataTraining  *bool `json:"data_training"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
