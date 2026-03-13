package repository

import (
	"context"
	"errors"

	"app/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenerationJobRepository interface {
	Create(ctx context.Context, job *model.GenerationJob) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.GenerationJob, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]model.GenerationJob, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, limit int, offset int) ([]model.GenerationJob, error)
	GetPendingJobs(ctx context.Context, limit int) ([]model.GenerationJob, error)
	GetProcessingJobs(ctx context.Context, limit int) ([]model.GenerationJob, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMsg string) error
	UpdateProgress(ctx context.Context, id uuid.UUID, notes map[string]interface{}) error
	Update(ctx context.Context, job *model.GenerationJob) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetJobsByStatus(ctx context.Context, status string, limit int) ([]model.GenerationJob, error)
}

type generationJobRepository struct {
	db *gorm.DB
}

func NewGenerationJobRepository(db *gorm.DB) GenerationJobRepository {
	return &generationJobRepository{db: db}
}

func (r *generationJobRepository) Create(ctx context.Context, job *model.GenerationJob) error {
	result := r.db.WithContext(ctx).Create(job)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *generationJobRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.GenerationJob, error) {
	var job model.GenerationJob
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&job)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("generation job not found")
		}
		return nil, result.Error
	}
	return &job, nil
}

func (r *generationJobRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]model.GenerationJob, error) {
	var jobs []model.GenerationJob
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

func (r *generationJobRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID, limit int, offset int) ([]model.GenerationJob, error) {
	var jobs []model.GenerationJob
	result := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

func (r *generationJobRepository) GetPendingJobs(ctx context.Context, limit int) ([]model.GenerationJob, error) {
	var jobs []model.GenerationJob
	result := r.db.WithContext(ctx).
		Where("status = ?", "queued").
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

func (r *generationJobRepository) GetProcessingJobs(ctx context.Context, limit int) ([]model.GenerationJob, error) {
	var jobs []model.GenerationJob
	result := r.db.WithContext(ctx).
		Where("status = ?", "processing").
		Limit(limit).
		Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

func (r *generationJobRepository) GetJobsByStatus(ctx context.Context, status string, limit int) ([]model.GenerationJob, error) {
	var jobs []model.GenerationJob
	result := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at ASC").
		Limit(limit).
		Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

func (r *generationJobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMsg string) error {
	updates := map[string]interface{}{"status": status}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}
	if status == "completed" {
		now := gorm.NowFunc()
		updates["completed_at"] = now
	}
	if status == "processing" {
		now := gorm.NowFunc()
		updates["started_at"] = now
	}

	result := r.db.WithContext(ctx).Model(&model.GenerationJob{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *generationJobRepository) UpdateProgress(ctx context.Context, id uuid.UUID, notes map[string]interface{}) error {
	data, _ := gorm.NowFunc().MarshalJSON()
	notes["updated_at"] = data

	result := r.db.WithContext(ctx).
		Model(&model.GenerationJob{}).
		Where("id = ?", id).
		Update("processing_notes", gorm.Expr("jsonb_set(COALESCE(processing_notes, '{}'), '{0}', ?)", notes))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *generationJobRepository) Update(ctx context.Context, job *model.GenerationJob) error {
	result := r.db.WithContext(ctx).Save(job)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *generationJobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.GenerationJob{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
