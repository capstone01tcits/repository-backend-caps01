package repository

import (
	"Sevima-AI-Content-Creator/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRepository interface {
	Create(video *model.Video) error
	FindByID(id string) (*model.Video, error)
	FindByProjectID(projectID string) ([]model.Video, error)
	FindByStoryboardID(storyboardID string) (*model.Video, error)
	FindAllByStoryboardID(storyboardID string) ([]model.Video, error)
	FindByUserID(userID string) ([]model.Video, error)
	Update(video *model.Video) error
	Delete(id string) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db}
}

func (r *videoRepository) Create(video *model.Video) error {
	return r.db.Create(video).Error
}

func (r *videoRepository) FindByID(id string) (*model.Video, error) {
	var video model.Video
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uid).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindByProjectID(projectID string) ([]model.Video, error) {
	var videos []model.Video
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return videos, err
	}
	err = r.db.Where("project_id = ?", pid).Order("created_at DESC").Find(&videos).Error
	return videos, err
}

func (r *videoRepository) FindByStoryboardID(storyboardID string) (*model.Video, error) {
	var video model.Video
	sid, err := uuid.Parse(storyboardID)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("storyboard_id = ?", sid).Order("created_at DESC").First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindAllByStoryboardID(storyboardID string) ([]model.Video, error) {
	var videos []model.Video
	sid, err := uuid.Parse(storyboardID)
	if err != nil {
		return videos, err
	}
	err = r.db.Where("storyboard_id = ?", sid).Order("created_at ASC").Find(&videos).Error
	return videos, err
}

func (r *videoRepository) FindByUserID(userID string) ([]model.Video, error) {
	var videos []model.Video
	uid, err := uuid.Parse(userID)
	if err != nil {
		return videos, err
	}
	err = r.db.Where("user_id = ?", uid).Order("created_at DESC").Find(&videos).Error
	return videos, err
}

func (r *videoRepository) Update(video *model.Video) error {
	return r.db.Save(video).Error
}

func (r *videoRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ?", uid).Delete(&model.Video{}).Error
}
