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
	UnscopedFindByID(id string) (*model.Storyboard, error)
	FindByProjectID(projectID string) (*model.Storyboard, error)
	Update(storyboard *model.Storyboard) error
	Delete(id string) error
	Restore(id string) error

	// StoryboardSection
	CreateSection(section *model.StoryboardSection) error
	FindSectionsByStoryboardID(storyboardID string) ([]model.StoryboardSection, error)
	UpdateSection(section *model.StoryboardSection) error
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
	err = r.db.Where("id = ?", uid).Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&storyboard).Error
	if err != nil {
		return nil, err
	}
	return &storyboard, nil
}

func (r *storyboardRepository) UnscopedFindByID(id string) (*model.Storyboard, error) {
	var storyboard model.Storyboard
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Unscoped().Where("id = ?", uid).Preload("Sections").First(&storyboard).Error
	if err != nil {
		return nil, err
	}
	return &storyboard, nil
}


func (r *storyboardRepository) FindByProjectID(projectID string) (*model.Storyboard, error) {
	var storyboard model.Storyboard
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("project_id = ?", pid).Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&storyboard).Error
	if err != nil {
		return nil, err
	}
	return &storyboard, nil
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

func (r *storyboardRepository) Restore(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Unscoped().Model(&model.Storyboard{}).Where("id = ?", uid).Update("deleted_at", nil).Error
}



func (r *storyboardRepository) CreateSection(section *model.StoryboardSection) error {
	return r.db.Create(section).Error
}

func (r *storyboardRepository) FindSectionsByStoryboardID(storyboardID string) ([]model.StoryboardSection, error) {
	var sections []model.StoryboardSection
	sid, err := uuid.Parse(storyboardID)
	if err != nil {
		return sections, err
	}
	err = r.db.Where("storyboard_id = ?", sid).Order("created_at ASC").Find(&sections).Error
	return sections, err
}

func (r *storyboardRepository) UpdateSection(section *model.StoryboardSection) error {
	return r.db.Save(section).Error
}
