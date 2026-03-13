package repository

import (
	"Sevima-AI-Content-Creator/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentRepository interface {
	// Content Pillar
	CreateContentPillar(pillar *model.ContentPillar) error
	FindContentPillarByID(id string) (*model.ContentPillar, error)
	FindContentPillarsByProjectID(projectID string) ([]model.ContentPillar, error)
	UpdateContentPillar(pillar *model.ContentPillar) error
	DeleteContentPillar(id string) error
	DeselectAllPillarsByProjectID(projectID string) error

	// Content Theme
	CreateContentTheme(theme *model.ContentTheme) error
	FindContentThemeByID(id string) (*model.ContentTheme, error)
	FindContentThemesByPillarID(pillarID string) ([]model.ContentTheme, error)
	UpdateContentTheme(theme *model.ContentTheme) error
	DeleteContentTheme(id string) error
	DeselectAllThemesByPillarID(pillarID string) error
}

type contentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db}
}

// ==================== Content Pillar ====================

func (r *contentRepository) CreateContentPillar(pillar *model.ContentPillar) error {
	return r.db.Create(pillar).Error
}

func (r *contentRepository) FindContentPillarByID(id string) (*model.ContentPillar, error) {
	var pillar model.ContentPillar
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uid).First(&pillar).Error
	if err != nil {
		return nil, err
	}
	return &pillar, nil
}

func (r *contentRepository) FindContentPillarsByProjectID(projectID string) ([]model.ContentPillar, error) {
	var pillars []model.ContentPillar
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return pillars, err
	}
	err = r.db.Where("project_id = ?", pid).Preload("ContentThemes").Order("created_at ASC").Find(&pillars).Error
	return pillars, err
}

func (r *contentRepository) UpdateContentPillar(pillar *model.ContentPillar) error {
	return r.db.Save(pillar).Error
}

func (r *contentRepository) DeleteContentPillar(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.ContentPillar{}).Error
}

func (r *contentRepository) DeselectAllPillarsByProjectID(projectID string) error {
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return err
	}
	return r.db.Model(&model.ContentPillar{}).Where("project_id = ?", pid).Update("is_selected", false).Error
}

// ==================== Content Theme ====================

func (r *contentRepository) CreateContentTheme(theme *model.ContentTheme) error {
	return r.db.Create(theme).Error
}

func (r *contentRepository) FindContentThemeByID(id string) (*model.ContentTheme, error) {
	var theme model.ContentTheme
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uid).First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *contentRepository) FindContentThemesByPillarID(pillarID string) ([]model.ContentTheme, error) {
	var themes []model.ContentTheme
	pilID, err := uuid.Parse(pillarID)
	if err != nil {
		return themes, err
	}
	err = r.db.Where("content_pillar_id = ?", pilID).Order("created_at ASC").Find(&themes).Error
	return themes, err
}

func (r *contentRepository) UpdateContentTheme(theme *model.ContentTheme) error {
	return r.db.Save(theme).Error
}

func (r *contentRepository) DeleteContentTheme(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.ContentTheme{}).Error
}

func (r *contentRepository) DeselectAllThemesByPillarID(pillarID string) error {
	pilID, err := uuid.Parse(pillarID)
	if err != nil {
		return err
	}
	return r.db.Model(&model.ContentTheme{}).Where("content_pillar_id = ?", pilID).Update("is_selected", false).Error
}
