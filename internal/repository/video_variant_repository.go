package repository

import (
	"context"
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoVariantRepository interface {
	Create(ctx context.Context, variant *model.VideoVariant) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.VideoVariant, error)
	GetByStoryboardID(ctx context.Context, storyboardID uuid.UUID) ([]model.VideoVariant, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, limit int, offset int) ([]model.VideoVariant, error)
	GetByVariantNumber(ctx context.Context, storyboardID uuid.UUID, variantNumber int) (*model.VideoVariant, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateWithVideoURL(ctx context.Context, id uuid.UUID, videoURL, thumbnailURL string, fileSize int64) error
	Update(ctx context.Context, variant *model.VideoVariant) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByExternalJobID(ctx context.Context, externalJobID string) (*model.VideoVariant, error)
}

type videoVariantRepository struct {
	db *gorm.DB
}

func NewVideoVariantRepository(db *gorm.DB) VideoVariantRepository {
	return &videoVariantRepository{db: db}
}

func (r *videoVariantRepository) Create(ctx context.Context, variant *model.VideoVariant) error {
	result := r.db.WithContext(ctx).Create(variant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *videoVariantRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.VideoVariant, error) {
	var variant model.VideoVariant
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&variant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("video variant not found")
		}
		return nil, result.Error
	}
	return &variant, nil
}

func (r *videoVariantRepository) GetByStoryboardID(ctx context.Context, storyboardID uuid.UUID) ([]model.VideoVariant, error) {
	var variants []model.VideoVariant
	result := r.db.WithContext(ctx).
		Where("storyboard_id = ?", storyboardID).
		Order("variant_number ASC, created_at DESC").
		Find(&variants)
	if result.Error != nil {
		return nil, result.Error
	}
	return variants, nil
}

func (r *videoVariantRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID, limit int, offset int) ([]model.VideoVariant, error) {
	var variants []model.VideoVariant
	result := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&variants)
	if result.Error != nil {
		return nil, result.Error
	}
	return variants, nil
}

func (r *videoVariantRepository) GetByVariantNumber(ctx context.Context, storyboardID uuid.UUID, variantNumber int) (*model.VideoVariant, error) {
	var variant model.VideoVariant
	result := r.db.WithContext(ctx).
		Where("storyboard_id = ? AND variant_number = ?", storyboardID, variantNumber).
		First(&variant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("video variant not found")
		}
		return nil, result.Error
	}
	return &variant, nil
}

func (r *videoVariantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	result := r.db.WithContext(ctx).Model(&model.VideoVariant{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *videoVariantRepository) UpdateWithVideoURL(ctx context.Context, id uuid.UUID, videoURL, thumbnailURL string, fileSize int64) error {
	updates := map[string]interface{}{
		"video_url":     videoURL,
		"thumbnail_url": thumbnailURL,
		"file_size":     fileSize,
		"status":        "completed",
	}
	result := r.db.WithContext(ctx).Model(&model.VideoVariant{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *videoVariantRepository) Update(ctx context.Context, variant *model.VideoVariant) error {
	result := r.db.WithContext(ctx).Save(variant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *videoVariantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.VideoVariant{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *videoVariantRepository) GetByExternalJobID(ctx context.Context, externalJobID string) (*model.VideoVariant, error) {
	var variant model.VideoVariant
	result := r.db.WithContext(ctx).Where("external_job_id = ?", externalJobID).First(&variant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("video variant not found")
		}
		return nil, result.Error
	}
	return &variant, nil
}
