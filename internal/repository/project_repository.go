package repository

import (
	"Sevima-AI-Content-Creator/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *model.Project) error
	FindByID(id string) (*model.Project, error)
	UnscopedFindByID(id string) (*model.Project, error)
	FindByUserID(userID string) ([]model.Project, error)
	Update(project *model.Project) error
	Delete(id string) error
	Restore(id string) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db}
}

func (r *projectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id string) (*model.Project, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	
	var project model.Project
	err = r.db.
		Preload("BusinessBriefs").
		Preload("BusinessBriefs.CreativeBriefs").
		Preload("Storyboard").
		Preload("Storyboard.Sections").
		Preload("Videos").
		Where("id = ?", uid).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) UnscopedFindByID(id string) (*model.Project, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	
	var project model.Project
	err = r.db.Unscoped().Where("id = ?", uid).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindByUserID(userID string) ([]model.Project, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	
	var projects []model.Project
	err = r.db.Preload("Videos").Where("user_id = ?", uid).Order("created_at DESC").Find(&projects).Error
	return projects, err
}

func (r *projectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.Project{}).Error
}

func (r *projectRepository) Restore(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Unscoped().Model(&model.Project{}).Where("id = ?", uid).Update("deleted_at", nil).Error
}
