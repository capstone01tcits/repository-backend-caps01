package repository

import (
	"context"
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SceneGenerationRepository interface {
	Create(ctx context.Context, scene *model.SceneGeneration) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SceneGeneration, error)
	GetByVariantID(ctx context.Context, variantID uuid.UUID) ([]model.SceneGeneration, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateWithVideoURL(ctx context.Context, id uuid.UUID, videoURL string) error
	Update(ctx context.Context, scene *model.SceneGeneration) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByExternalJobID(ctx context.Context, externalJobID string) (*model.SceneGeneration, error)
	GetPendingScenes(ctx context.Context, limit int) ([]model.SceneGeneration, error)
}

type sceneGenerationRepository struct {
	db *gorm.DB
}

func NewSceneGenerationRepository(db *gorm.DB) SceneGenerationRepository {
	return &sceneGenerationRepository{db: db}
}

func (r *sceneGenerationRepository) Create(ctx context.Context, scene *model.SceneGeneration) error {
	result := r.db.WithContext(ctx).Create(scene)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *sceneGenerationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SceneGeneration, error) {
	var scene model.SceneGeneration
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&scene)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("scene generation not found")
		}
		return nil, result.Error
	}
	return &scene, nil
}

func (r *sceneGenerationRepository) GetByVariantID(ctx context.Context, variantID uuid.UUID) ([]model.SceneGeneration, error) {
	var scenes []model.SceneGeneration
	result := r.db.WithContext(ctx).
		Where("variant_id = ?", variantID).
		Order("scene_index ASC").
		Find(&scenes)
	if result.Error != nil {
		return nil, result.Error
	}
	return scenes, nil
}

func (r *sceneGenerationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	result := r.db.WithContext(ctx).Model(&model.SceneGeneration{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *sceneGenerationRepository) UpdateWithVideoURL(ctx context.Context, id uuid.UUID, videoURL string) error {
	result := r.db.WithContext(ctx).
		Model(&model.SceneGeneration{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"video_url": videoURL,
			"status":    "completed",
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *sceneGenerationRepository) Update(ctx context.Context, scene *model.SceneGeneration) error {
	result := r.db.WithContext(ctx).Save(scene)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *sceneGenerationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.SceneGeneration{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *sceneGenerationRepository) GetByExternalJobID(ctx context.Context, externalJobID string) (*model.SceneGeneration, error) {
	var scene model.SceneGeneration
	result := r.db.WithContext(ctx).Where("external_job_id = ?", externalJobID).First(&scene)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("scene generation not found")
		}
		return nil, result.Error
	}
	return &scene, nil
}

func (r *sceneGenerationRepository) GetPendingScenes(ctx context.Context, limit int) ([]model.SceneGeneration, error) {
	var scenes []model.SceneGeneration
	result := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Limit(limit).
		Find(&scenes)
	if result.Error != nil {
		return nil, result.Error
	}
	return scenes, nil
}
