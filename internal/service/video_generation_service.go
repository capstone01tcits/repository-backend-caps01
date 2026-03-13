package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Sevima-AI-Content-Creator/internal/ai"
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"
	"github.com/google/uuid"
)

type VideoGenerationService interface {
	// Generate 3 video variants from a storyboard
	GenerateVideoVariants(ctx context.Context, userID, projectID, storyboardID uuid.UUID, customPrompt string) (*model.GenerationJob, error)
	
	// Regenerate a specific video variant
	RegenerateVideoVariant(ctx context.Context, variantID uuid.UUID, newPrompt string) (*model.GenerationJob, error)
	
	// Regenerate a specific scene
	RegenerateScene(ctx context.Context, sceneID uuid.UUID, newPrompt string) (*model.GenerationJob, error)
	
	// Get generation job status
	GetJobStatus(ctx context.Context, jobID uuid.UUID) (*model.GenerationJob, error)
	
	// Get video variants for a storyboard
	GetVideoVariants(ctx context.Context, storyboardID uuid.UUID) ([]model.VideoVariant, error)
	
	// Get single video variant
	GetVideoVariant(ctx context.Context, variantID uuid.UUID) (*model.VideoVariant, error)
	
	// Get video variant with scenes
	GetVideoVariantWithScenes(ctx context.Context, variantID uuid.UUID) (*model.VideoVariant, []model.SceneGeneration, error)
	
	// Calculate credits for generation
	CalculateCreditsForGeneration(duration int, sceneCount int, videoCount int) int
	
	// Process generation job (called by worker)
	ProcessGenerationJob(ctx context.Context, jobID uuid.UUID) error
	
	// Poll and update job status (called by worker)
	PollJobStatus(ctx context.Context, jobID uuid.UUID) error
}

type videoGenerationService struct {
	jobRepo          repository.GenerationJobRepository
	variantRepo      repository.VideoVariantRepository
	sceneRepo        repository.SceneGenerationRepository
	creditService    CreditService
	providerFactory  *ai.ProviderFactory
}

func NewVideoGenerationService(
	jobRepo repository.GenerationJobRepository,
	variantRepo repository.VideoVariantRepository,
	sceneRepo repository.SceneGenerationRepository,
	creditService CreditService,
) VideoGenerationService {
	return &videoGenerationService{
		jobRepo:         jobRepo,
		variantRepo:     variantRepo,
		sceneRepo:       sceneRepo,
		creditService:   creditService,
		providerFactory: &ai.ProviderFactory{},
	}
}

func (s *videoGenerationService) GenerateVideoVariants(ctx context.Context, userID, projectID, storyboardID uuid.UUID, customPrompt string) (*model.GenerationJob, error) {
	// Validate user has credits
	credits, err := s.creditService.GetUserCredits(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %w", err)
	}

	// Calculate credits needed for 3 videos: 8-12 sec each, 2-3 scenes
	sceneDuration := 5 // middle of 4-6 sec
	videoDuration := 10 // middle of 8-12 sec
	sceneCount := 2 // middle of 2-3
	creditsNeeded := s.CalculateCreditsForGeneration(videoDuration, sceneCount, 3)

	if credits < creditsNeeded {
		return nil, errors.New("insufficient credits for video generation")
	}

	// Create generation job
	job := &model.GenerationJob{
		UserID:          userID,
		ProjectID:       projectID,
		StoryboardID:    storyboardID,
		JobType:         "generate",
		Status:          "queued",
		Priority:        1,
		SceneCount:      sceneCount,
		VideoDuration:   videoDuration,
		Provider:        "ltx",
		Model:           "ltx-2-fast",
		Resolution:      "1080p",
		CreditsRequired: creditsNeeded,
		MaxRetries:      3,
	}

	// Store prompt if provided
	if customPrompt != "" {
		promptData := map[string]interface{}{
			"custom_prompt": customPrompt,
			"timestamp":     time.Now(),
		}
		promptJSON, _ := json.Marshal(promptData)
		job.Prompt = promptJSON
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create generation job: %w", err)
	}

	// Create 3 video variants
	for i := 1; i <= 3; i++ {
		variant := &model.VideoVariant{
			UserID:       userID,
			ProjectID:    projectID,
			StoryboardID: storyboardID,
			VariantNumber: i,
			Status:       "pending",
			Duration:     videoDuration,
			Resolution:   "1080p",
			Provider:     "ltx",
			Model:        "ltx-2-fast",
		}

		// Generate prompt based on variant
		prompt := s.generatePromptForVariant(ctx, storyboardID, i, customPrompt)
		variant.PromptUsed = prompt

		// Create scene plan (2-3 scenes)
		scenePlan := s.generateScenePlan(sceneDuration, sceneCount)
		scenePlanJSON, _ := json.Marshal(scenePlan)
		variant.ScenePlan = scenePlanJSON

		if err := s.variantRepo.Create(ctx, variant); err != nil {
			return nil, fmt.Errorf("failed to create video variant: %w", err)
		}

		// Create individual scene generation tasks
		for j, scene := range scenePlan {
			// Handle both int and float64 for duration since JSON unmarshaling may change types
			var duration int
			switch d := scene["duration"].(type) {
			case int:
				duration = d
			case float64:
				duration = int(d)
			}

			sceneGen := &model.SceneGeneration{
				VariantID:   variant.ID,
				SceneNumber: j + 1,
				SceneIndex:  j,
				Prompt:      scene["prompt"].(string),
				Duration:    duration,
				Status:      "pending",
			}

			if err := s.sceneRepo.Create(ctx, sceneGen); err != nil {
				return nil, fmt.Errorf("failed to create scene generation task: %w", err)
			}
		}
	}

	// Deduct credits
	if err := s.creditService.DeductCredits(ctx, userID, creditsNeeded, "video_generation"); err != nil {
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}

	return job, nil
}

func (s *videoGenerationService) RegenerateVideoVariant(ctx context.Context, variantID uuid.UUID, newPrompt string) (*model.GenerationJob, error) {
	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video variant: %w", err)
	}

	// Deduct credits for regeneration
	credits, err := s.creditService.GetUserCredits(ctx, variant.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %w", err)
	}

	creditsNeeded := s.CalculateCreditsForGeneration(variant.Duration, 2, 1)
	if credits < creditsNeeded {
		return nil, errors.New("insufficient credits for regeneration")
	}

	// Create new generation job for regeneration
	job := &model.GenerationJob{
		UserID:          variant.UserID,
		ProjectID:       variant.ProjectID,
		StoryboardID:    variant.StoryboardID,
		JobType:         "regenerate",
		Status:          "queued",
		Priority:        2,
		SceneCount:      2,
		VideoDuration:   variant.Duration,
		Provider:        variant.Provider,
		Model:           variant.Model,
		Resolution:      variant.Resolution,
		CreditsRequired: creditsNeeded,
		MaxRetries:      3,
	}

	if newPrompt != "" {
		variant.PromptUsed = newPrompt
	}

	promptJSON, _ := json.Marshal(map[string]interface{}{"prompt": variant.PromptUsed})
	job.Prompt = promptJSON

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create regeneration job: %w", err)
	}

	// Mark old scenes as superseded and create new ones
	oldScenes, err := s.sceneRepo.GetByVariantID(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get old scenes: %w", err)
	}

	for _, oldScene := range oldScenes {
		// Create new scene with updated prompt
		newScene := &model.SceneGeneration{
			VariantID:   variant.ID,
			SceneNumber: oldScene.SceneNumber,
			SceneIndex:  oldScene.SceneIndex,
			Prompt:      newPrompt,
			Duration:    oldScene.Duration,
			Status:      "pending",
		}

		if err := s.sceneRepo.Create(ctx, newScene); err != nil {
			return nil, fmt.Errorf("failed to create new scene: %w", err)
		}
	}

	// Deduct credits
	if err := s.creditService.DeductCredits(ctx, variant.UserID, creditsNeeded, "video_regeneration"); err != nil {
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}

	return job, nil
}

func (s *videoGenerationService) RegenerateScene(ctx context.Context, sceneID uuid.UUID, newPrompt string) (*model.GenerationJob, error) {
	scene, err := s.sceneRepo.GetByID(ctx, sceneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scene: %w", err)
	}

	variant, err := s.variantRepo.GetByID(ctx, scene.VariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video variant: %w", err)
	}

	// Deduct credits for scene regeneration (lower cost than full video)
	credits, err := s.creditService.GetUserCredits(ctx, variant.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %w", err)
	}

	creditsNeeded := scene.Duration * 1 // Lower multiplier for single scene
	if credits < creditsNeeded {
		return nil, errors.New("insufficient credits for scene regeneration")
	}

	// Create generation job for scene
	job := &model.GenerationJob{
		UserID:          variant.UserID,
		ProjectID:       variant.ProjectID,
		StoryboardID:    variant.StoryboardID,
		JobType:         "regenerate_scene",
		Status:          "queued",
		Priority:        3,
		SceneCount:      1,
		VideoDuration:   scene.Duration,
		Provider:        variant.Provider,
		Model:           variant.Model,
		Resolution:      variant.Resolution,
		CreditsRequired: creditsNeeded,
		MaxRetries:      3,
	}

	// Update scene prompt
	if newPrompt != "" {
		scene.Prompt = newPrompt
		if err := s.sceneRepo.Update(ctx, scene); err != nil {
			return nil, fmt.Errorf("failed to update scene prompt: %w", err)
		}
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create scene regeneration job: %w", err)
	}

	// Deduct credits
	if err := s.creditService.DeductCredits(ctx, variant.UserID, creditsNeeded, "scene_regeneration"); err != nil {
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}

	return job, nil
}

func (s *videoGenerationService) GetJobStatus(ctx context.Context, jobID uuid.UUID) (*model.GenerationJob, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

func (s *videoGenerationService) GetVideoVariants(ctx context.Context, storyboardID uuid.UUID) ([]model.VideoVariant, error) {
	return s.variantRepo.GetByStoryboardID(ctx, storyboardID)
}

func (s *videoGenerationService) GetVideoVariant(ctx context.Context, variantID uuid.UUID) (*model.VideoVariant, error) {
	return s.variantRepo.GetByID(ctx, variantID)
}

func (s *videoGenerationService) GetVideoVariantWithScenes(ctx context.Context, variantID uuid.UUID) (*model.VideoVariant, []model.SceneGeneration, error) {
	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return nil, nil, err
	}

	scenes, err := s.sceneRepo.GetByVariantID(ctx, variantID)
	if err != nil {
		return nil, nil, err
	}

	return variant, scenes, nil
}

func (s *videoGenerationService) CalculateCreditsForGeneration(duration int, sceneCount int, videoCount int) int {
	// Cost formula: duration * sceneCount * videoCount * base_multiplier
	// For standard tier: ~2 credits per second per video
	return duration * sceneCount * videoCount * 2
}

func (s *videoGenerationService) ProcessGenerationJob(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Update job status to processing
	if err := s.jobRepo.UpdateStatus(ctx, jobID, "processing", ""); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Get video variants for this job
	variants, err := s.variantRepo.GetByStoryboardID(ctx, job.StoryboardID)
	if err != nil {
		return fmt.Errorf("failed to get variants: %w", err)
	}

	// Process each variant's scenes
	for _, variant := range variants {
		scenes, err := s.sceneRepo.GetByVariantID(ctx, variant.ID)
		if err != nil {
			continue
		}

		// Update variant status
		s.variantRepo.UpdateStatus(ctx, variant.ID, "processing")

		// Generate each scene
		provider := s.providerFactory.GetProviderByModel(job.Model)

		for _, scene := range scenes {
			req := ai.VideoGenerationRequest{
				Prompt:     scene.Prompt,
				Duration:   scene.Duration,
				Resolution: job.Resolution,
				FPS:        30,
				Model:      job.Model,
			}

			resp, err := provider.GenerateScene(ctx, req)
			if err != nil {
				s.sceneRepo.UpdateStatus(ctx, scene.ID, "failed")
				s.jobRepo.UpdateStatus(ctx, jobID, "failed", fmt.Sprintf("Scene generation failed: %v", err))
				continue
			}

			// Store external job ID
			scene.ExternalJobID = resp.JobID
			scene.Status = "processing"
			s.sceneRepo.Update(ctx, &scene)
		}
	}

	return nil
}

func (s *videoGenerationService) PollJobStatus(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	if job.Status == "completed" || job.Status == "failed" {
		return nil
	}

	// Get video variants
	variants, err := s.variantRepo.GetByStoryboardID(ctx, job.StoryboardID)
	if err != nil {
		return fmt.Errorf("failed to get variants: %w", err)
	}

	provider := s.providerFactory.GetProviderByModel(job.Model)
	completedVariants := 0

	for _, variant := range variants {
		scenes, err := s.sceneRepo.GetByVariantID(ctx, variant.ID)
		if err != nil {
			continue
		}

		variantComplete := true
		for _, scene := range scenes {
			if scene.ExternalJobID == "" || scene.Status == "completed" {
				continue
			}

			// Poll provider for status
			resp, err := provider.GetJobStatus(ctx, scene.ExternalJobID)
			if err != nil {
				continue
			}

			if resp.Status == "completed" && resp.VideoURL != "" {
				s.sceneRepo.UpdateWithVideoURL(ctx, scene.ID, resp.VideoURL)
				s.sceneRepo.UpdateStatus(ctx, scene.ID, "completed")
			} else if resp.Status == "failed" {
				s.sceneRepo.UpdateStatus(ctx, scene.ID, "failed")
			} else {
				variantComplete = false
			}
		}

		// Check if all scenes are complete
		if variantComplete {
			s.variantRepo.UpdateStatus(ctx, variant.ID, "completed")
			completedVariants++

			// Mock: simulate video URL and thumbnail
			mockVideoURL := fmt.Sprintf("https://storage.example.com/videos/variants/%s.mp4", variant.ID)
			mockThumbnailURL := fmt.Sprintf("https://storage.example.com/thumbnails/%s.jpg", variant.ID)
			s.variantRepo.UpdateWithVideoURL(ctx, variant.ID, mockVideoURL, mockThumbnailURL, int64(1024*1024*50)) // 50MB mock
		}
	}

	// If all variants are complete, update job status
	if completedVariants == len(variants) {
		s.jobRepo.UpdateStatus(ctx, jobID, "completed", "")
	}

	return nil
}

// Helper functions

func (s *videoGenerationService) generatePromptForVariant(ctx context.Context, storyboardID uuid.UUID, variantNumber int, customPrompt string) string {
	variations := []string{
		"cinematic",
		"vibrant and dynamic",
		"professional and polished",
	}

	basePrompt := customPrompt
	if basePrompt == "" {
		basePrompt = "Generate a professional marketing video"
	}

	variation := variations[(variantNumber-1)%len(variations)]
	return fmt.Sprintf("%s with a %s style. This is variation %d.", basePrompt, variation, variantNumber)
}

func (s *videoGenerationService) generateScenePlan(sceneDuration int, sceneCount int) []map[string]interface{} {
	scenes := make([]map[string]interface{}, sceneCount)

	for i := 0; i < sceneCount; i++ {
		scenes[i] = map[string]interface{}{
			"scene_number": i + 1,
			"duration":     sceneDuration,
			"prompt":       fmt.Sprintf("Scene %d: Professional marketing content with smooth transitions", i+1),
		}
	}

	return scenes
}
