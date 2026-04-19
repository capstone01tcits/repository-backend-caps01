package repository

import (
	"Sevima-AI-Content-Creator/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BriefRepository interface {
	// Business Brief
	CreateBusinessBrief(brief *model.BusinessBrief) error
	FindBusinessBriefByID(id string) (*model.BusinessBrief, error)
	FindBusinessBriefsByUserID(userID string) ([]model.BusinessBrief, error)
	FindBusinessBriefByProjectID(projectID string) (*model.BusinessBrief, error)
	UpdateBusinessBrief(brief *model.BusinessBrief) error
	DeleteBusinessBrief(id string) error

	// Creative Brief
	CreateCreativeBrief(brief *model.CreativeBrief) error
	FindCreativeBriefByID(id string) (*model.CreativeBrief, error)
	FindCreativeBriefsByUserID(userID string) ([]model.CreativeBrief, error)
	FindCreativeBriefsByBusinessBriefID(businessBriefID string) ([]model.CreativeBrief, error)
	UpdateCreativeBrief(brief *model.CreativeBrief) error
	DeleteCreativeBrief(id string) error
}

type briefRepository struct {
	db *gorm.DB
}

func NewBriefRepository(db *gorm.DB) BriefRepository {
	return &briefRepository{db}
}

// ==================== Business Brief ====================

func (r *briefRepository) CreateBusinessBrief(brief *model.BusinessBrief) error {
	return r.db.Create(brief).Error
}

func (r *briefRepository) FindBusinessBriefByID(id string) (*model.BusinessBrief, error) {
	var brief model.BusinessBrief
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uid).First(&brief).Error
	if err != nil {
		return nil, err
	}
	return &brief, nil
}

func (r *briefRepository) FindBusinessBriefsByUserID(userID string) ([]model.BusinessBrief, error) {
	var briefs []model.BusinessBrief
	uid, err := uuid.Parse(userID)
	if err != nil {
		return briefs, err
	}
	err = r.db.Where("user_id = ?", uid).Order("created_at DESC").Find(&briefs).Error
	return briefs, err
}

func (r *briefRepository) FindBusinessBriefByProjectID(projectID string) (*model.BusinessBrief, error) {
	var brief model.BusinessBrief
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("project_id = ?", pid).First(&brief).Error
	if err != nil {
		return nil, err
	}
	return &brief, nil
}

func (r *briefRepository) UpdateBusinessBrief(brief *model.BusinessBrief) error {
	return r.db.Save(brief).Error
}

func (r *briefRepository) DeleteBusinessBrief(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.BusinessBrief{}).Error
}

// ==================== Creative Brief ====================

func (r *briefRepository) CreateCreativeBrief(brief *model.CreativeBrief) error {
	return r.db.Create(brief).Error
}

func (r *briefRepository) FindCreativeBriefByID(id string) (*model.CreativeBrief, error) {
	var brief model.CreativeBrief
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uid).First(&brief).Error
	if err != nil {
		return nil, err
	}
	return &brief, nil
}

func (r *briefRepository) FindCreativeBriefsByUserID(userID string) ([]model.CreativeBrief, error) {
	var briefs []model.CreativeBrief
	uid, err := uuid.Parse(userID)
	if err != nil {
		return briefs, err
	}
	err = r.db.Where("user_id = ?", uid).Order("created_at DESC").Find(&briefs).Error
	return briefs, err
}

func (r *briefRepository) FindCreativeBriefsByBusinessBriefID(businessBriefID string) ([]model.CreativeBrief, error) {
	var briefs []model.CreativeBrief
	bid, err := uuid.Parse(businessBriefID)
	if err != nil {
		return briefs, err
	}
	err = r.db.Where("business_brief_id = ?", bid).Order("created_at DESC").Find(&briefs).Error
	return briefs, err
}

func (r *briefRepository) UpdateCreativeBrief(brief *model.CreativeBrief) error {
	return r.db.Save(brief).Error
}

func (r *briefRepository) DeleteCreativeBrief(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.CreativeBrief{}).Error
}
