package repository

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	UpdatePassword(id string, newPassword string) error
	UpdateCredits(id string, credits int) error
	DeductCredits(id string, amount int) error
	UpdatePreferences(id string, req *model.UpdatePreferencesRequest) error
	Delete(id string) error
	Restore(id string) error
	FindByIDIncludeDeleted(id string) (*model.User, error)
	FindAll() ([]model.User, error)
	Count() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	// Parse string to UUID for proper database query
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.db.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(id string, newPassword string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Model(&model.User{}).Where("id = ?", uid).Update("password", newPassword).Error
}

func (r *userRepository) UpdateCredits(id string, credits int) error {
	if credits < 0 {
		return errors.New("credits cannot be negative")
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Model(&model.User{}).Where("id = ?", uid).Update("credits", credits).Error
}

func (r *userRepository) DeductCredits(id string, amount int) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	result := r.db.Model(&model.User{}).
		Where("id = ? AND credits >= ?", uid, amount).
		Update("credits", gorm.Expr("credits - ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient credits or user not found")
	}
	return nil
}

func (r *userRepository) UpdatePreferences(id string, req *model.UpdatePreferencesRequest) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{}
	if req.EmailAlerts != nil {
		updates["email_alerts"] = *req.EmailAlerts
	}
	if req.Newsletter != nil {
		updates["newsletter"] = *req.Newsletter
	}
	if req.PublicProfile != nil {
		updates["public_profile"] = *req.PublicProfile
	}
	if req.DataTraining != nil {
		updates["data_training"] = *req.DataTraining
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&model.User{}).Where("id = ?", uid).Updates(updates).Error
}

func (r *userRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.User{}).Error
}

func (r *userRepository) Restore(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Unscoped().Model(&model.User{}).Where("id = ?", uid).Update("deleted_at", nil).Error
}

func (r *userRepository) FindByIDIncludeDeleted(id string) (*model.User, error) {
	// Parse string to UUID for proper database query
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.db.Unscoped().Where("id = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}
