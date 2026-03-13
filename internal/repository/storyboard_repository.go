package repository

import (
	"Sevima-AI-Content-Creator/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoryboardRepository interface {
	// Storyboard
	Create(storyboard *model.Storyboard) error
	FindByID(id string) (*model.Storyboard, error)
	FindByProjectID(projectID string) ([]model.Storyboard, error)
	Update(storyboard *model.Storyboard) error
	Delete(id string) error
	DeselectAllByProjectID(projectID string) error

	// Scene
	CreateScene(scene *model.Scene) error
	FindScenesByStoryboardID(storyboardID string) ([]model.Scene, error)
}

type storyboardRepository struct {
	db *gorm.DB
}

func NewStoryboardRepository(db *gorm.DB) StoryboardRepository {
	return &storyboardRepository{db}
}

func (r *storyboardRepository) Create(storyboard *model.Storyboard) error {
	return r.db.Create(storyboard).Error
}

func (r *storyboardRepository) FindByID(id string) (*model.Storyboard, error) {
	var storyboard model.Storyboard
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uid).Preload("Scenes", func(db *gorm.DB) *gorm.DB {
		return db.Order("scene_number ASC")
	}).First(&storyboard).Error
	if err != nil {
		return nil, err
	}
	return &storyboard, nil
}

func (r *storyboardRepository) FindByProjectID(projectID string) ([]model.Storyboard, error) {
	var storyboards []model.Storyboard
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return storyboards, err
	}
	err = r.db.Where("project_id = ?", pid).Preload("Scenes", func(db *gorm.DB) *gorm.DB {
		return db.Order("scene_number ASC")
	}).Order("created_at ASC").Find(&storyboards).Error
	return storyboards, err
}

func (r *storyboardRepository) Update(storyboard *model.Storyboard) error {
	return r.db.Save(storyboard).Error
}

func (r *storyboardRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.Storyboard{}).Error
}

func (r *storyboardRepository) DeselectAllByProjectID(projectID string) error {
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Storyboard{}).Where("project_id = ?", pid).Update("is_selected", false).Error
}

func (r *storyboardRepository) CreateScene(scene *model.Scene) error {
	return r.db.Create(scene).Error
}

func (r *storyboardRepository) FindScenesByStoryboardID(storyboardID string) ([]model.Scene, error) {
	var scenes []model.Scene
	sid, err := uuid.Parse(storyboardID)
	if err != nil {
		return scenes, err
	}
	err = r.db.Where("storyboard_id = ?", sid).Order("scene_number ASC").Find(&scenes).Error
	return scenes, err
}
