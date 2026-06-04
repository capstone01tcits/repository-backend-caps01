package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	// Get all video variants by user ID
	GetVideoVariantsByUserID(ctx context.Context, userID uuid.UUID) ([]model.VideoVariant, error)

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
	jobRepo         repository.GenerationJobRepository
	variantRepo     repository.VideoVariantRepository
	sceneRepo       repository.SceneGenerationRepository
	creditService   CreditService
	storageService  StorageService
	briefRepo       repository.BriefRepository
	storyboardRepo  repository.StoryboardRepository
	providerFactory *ai.ProviderFactory
}

func NewVideoGenerationService(
	jobRepo repository.GenerationJobRepository,
	variantRepo repository.VideoVariantRepository,
	sceneRepo repository.SceneGenerationRepository,
	briefRepo repository.BriefRepository,
	storyboardRepo repository.StoryboardRepository,
	creditService CreditService,
	storageService StorageService,
) VideoGenerationService {
	return &videoGenerationService{
		jobRepo:         jobRepo,
		variantRepo:     variantRepo,
		sceneRepo:       sceneRepo,
		briefRepo:       briefRepo,
		storyboardRepo:  storyboardRepo,
		creditService:   creditService,
		storageService:  storageService,
		providerFactory: &ai.ProviderFactory{},
	}
}

func (s *videoGenerationService) GenerateVideoVariants(ctx context.Context, userID, projectID, storyboardID uuid.UUID, customPrompt string) (*model.GenerationJob, error) {
	// Validate user has credits
	credits, err := s.creditService.GetUserCredits(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %w", err)
	}

	// Calculate credits needed for 1 video (Set to 1 for easy testing)
	videoDuration := 15 // Standard duration
	creditsNeeded := 1 

	// Disable strict check for testing
	if credits < creditsNeeded && false {
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
		SceneCount:      3,
		VideoDuration:   videoDuration,
		Provider:        "wavespeed",
		Model:           "veo-3.1",
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

	// Create 1 video variant (1 Project = 1 Storyboard = 1 Video)
	variant := &model.VideoVariant{
		UserID:        userID,
		ProjectID:     projectID,
		StoryboardID:  storyboardID,
		VariantNumber: 1,
		Status:        "pending",
		Duration:      videoDuration,
		Resolution:    "1080p",
		Provider:      "wavespeed",
		Model:         "veo-3.1",
	}

	// Ambil data sections untuk tracking (dan merangkai prompt jika diperlukan)
	storyboardSections, _ := s.storyboardRepo.FindSectionsByStoryboardID(storyboardID.String())

	// Gunakan customPrompt jika ada (test flow), jika tidak gunakan helper BuildVeo3Prompt
	fullPrompt := customPrompt
	if fullPrompt == "" {
		bb, _ := s.briefRepo.FindBusinessBriefByProjectID(projectID.String())
		var creativeBrief *model.CreativeBrief
		
		if bb != nil {
			cbList, _ := s.briefRepo.FindCreativeBriefsByBusinessBriefID(bb.ID.String())
			if len(cbList) > 0 {
				creativeBrief = &cbList[0]
			}
		}
		
		if bb != nil && creativeBrief != nil {
			fullPrompt = ai.BuildVeo3Prompt(bb, creativeBrief, storyboardSections)
		} else {
			fullPrompt = "Generate a professional marketing video based on the storyboard."
		}
	}
	variant.PromptUsed = fullPrompt

	if err := s.variantRepo.Create(ctx, variant); err != nil {
		return nil, fmt.Errorf("failed to create video variant: %w", err)
	}

	// Create scene generation tasks for tracking
	for i, sec := range storyboardSections {
		sceneGen := &model.SceneGeneration{
			VariantID:   variant.ID,
			SceneNumber: i + 1,
			SceneIndex:  i,
			Prompt:      sec.Content,
			Duration:    sec.Duration,
			Status:      "pending",
		}

		if err := s.sceneRepo.Create(ctx, sceneGen); err != nil {
			return nil, fmt.Errorf("failed to create scene generation task: %w", err)
		}
	}

	// Deduct credits (Skipped for testing)
	// if err := s.creditService.DeductCredits(ctx, userID, creditsNeeded, "video_generation"); err != nil {
	// 	return nil, fmt.Errorf("failed to deduct credits: %w", err)
	// }

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

func (s *videoGenerationService) GetVideoVariantsByUserID(ctx context.Context, userID uuid.UUID) ([]model.VideoVariant, error) {
	return s.variantRepo.GetByUserID(ctx, userID)
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

	// Update job status to generating_assets
	if err := s.jobRepo.UpdateStatus(ctx, jobID, "generating_assets", "Mengirim prompt ke Veo 3 AI Service"); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// 1. Persiapkan data Brief jika menggunakan model veo3/veo-3.1
	var bb *model.BusinessBrief
	var cb *model.CreativeBrief
	var refImages []string

	if job.Model == "veo-3.1" || job.Model == "veo3" {
		bb, _ = s.briefRepo.FindBusinessBriefByProjectID(job.ProjectID.String())
		if bb != nil {
			cbs, _ := s.briefRepo.FindCreativeBriefsByBusinessBriefID(bb.ID.String())
			if len(cbs) > 0 {
				cb = &cbs[0]
			}

			// Kumpulkan gambar referensi
			if bb.LogoPath != "" {
				refImages = append(refImages, bb.LogoPath)
			}
			if bb.EnvironmentPath != "" {
				refImages = append(refImages, bb.EnvironmentPath)
			}
		}
	}

	// 2. Get video variants for this job
	variants, err := s.variantRepo.GetByStoryboardID(ctx, job.StoryboardID)
	if err != nil {
		return fmt.Errorf("failed to get variants: %w", err)
	}

	// 3. Process each variant
	provider := s.providerFactory.GetProviderByModel(job.Model)

	for _, variant := range variants {
		scenes, err := s.sceneRepo.GetByVariantID(ctx, variant.ID)
		if err != nil || len(scenes) == 0 {
			continue
		}

		s.variantRepo.UpdateStatus(ctx, variant.ID, "processing")

		if job.Model == "veo-3.1" && bb != nil && cb != nil {
			// FLOW VEO 3: Kirim satu request gabungan untuk seluruh variant
			storyboardSections, _ := s.storyboardRepo.FindSectionsByStoryboardID(job.StoryboardID.String())
			fullPrompt := ai.BuildVeo3Prompt(bb, cb, storyboardSections)

			req := ai.VideoGenerationRequest{
				Prompt:          fullPrompt,
				Duration:        job.VideoDuration,
				Resolution:      job.Resolution,
				FPS:             30,
				Model:           job.Model,
				ReferenceImages: refImages,
			}

			resp, err := provider.GenerateScene(ctx, req)
			if err != nil {
				s.jobRepo.UpdateStatus(ctx, jobID, "failed", fmt.Sprintf("Veo 3 generation failed: %v", err))
				continue
			}

			// Tandai semua scene dengan JobID yang sama (karena diproses sekaligus oleh Veo 3)
			for _, scene := range scenes {
				scene.ExternalJobID = resp.JobID
				scene.Status = "processing"
				s.sceneRepo.Update(ctx, &scene)
			}
		} else {
			// FLOW STANDAR: Generate per scene (Legacy/Non-Veo3)
			for _, scene := range scenes {
				req := ai.VideoGenerationRequest{
					Prompt:          scene.Prompt,
					Duration:        scene.Duration,
					Resolution:      job.Resolution,
					FPS:             30,
					Model:           job.Model,
					ReferenceImages: refImages,
				}

				resp, err := provider.GenerateScene(ctx, req)
				if err != nil {
					s.sceneRepo.UpdateStatus(ctx, scene.ID, "failed")
					continue
				}

				scene.ExternalJobID = resp.JobID
				scene.Status = "processing"
				s.sceneRepo.Update(ctx, &scene)
			}
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

			if resp.Status == "stitching_video" {
				s.jobRepo.UpdateStatus(ctx, jobID, "stitching_video", "Menggabungkan adegan video menggunakan FFmpeg")
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
			s.jobRepo.UpdateStatus(ctx, jobID, "stitching_video", "Menggabungkan 3 scene dengan FFmpeg")
			s.variantRepo.UpdateStatus(ctx, variant.ID, "completed")
			completedVariants++

			// Get all completed scenes for this variant
			scenes, _ := s.sceneRepo.GetByVariantID(ctx, variant.ID)

			if len(scenes) > 0 && scenes[0].VideoURL != "" {
				// Download video from provider
				videoBytes, err := s.downloadVideoFromProvider(ctx, scenes[0].VideoURL)
				if err != nil {
					// Fallback: use provider URL directly
					videoURL := scenes[0].VideoURL
					thumbnailURL := fmt.Sprintf("https://provider-generated.com/thumbnails/%s.jpg", variant.ID)
					s.variantRepo.UpdateWithVideoURL(ctx, variant.ID, videoURL, thumbnailURL, int64(len(videoBytes)))
					continue
				}

				// Upload to Supabase
				videoFilename := fmt.Sprintf("video_%s.mp4", variant.ID.String())
				videoPath, err := s.storageService.UploadVideo(ctx, videoFilename, videoBytes)
				if err != nil {
					// Fallback: use provider URL if Supabase upload fails
					fmt.Printf("Supabase upload failed, using provider URL: %v\n", err)
					videoURL := scenes[0].VideoURL
					thumbnailURL := fmt.Sprintf("https://provider-generated.com/thumbnails/%s.jpg", variant.ID)
					s.variantRepo.UpdateWithVideoURL(ctx, variant.ID, videoURL, thumbnailURL, int64(len(videoBytes)))
					continue
				}

				// Get public URLs
				videoURL := s.storageService.GetPublicURL("videos", videoPath)
				thumbnailURL := fmt.Sprintf("https://provider-generated.com/thumbnails/%s.jpg", variant.ID)

				// Update database with Supabase URLs
				s.variantRepo.UpdateWithVideoURL(ctx, variant.ID, videoURL, thumbnailURL, int64(len(videoBytes)))
			}
		}
	}

	// If all variants are complete, update job status
	if completedVariants == len(variants) {
		s.jobRepo.UpdateStatus(ctx, jobID, "completed", "Video siap diunduh")
	}

	return nil
}

// Helper functions

// downloadVideoFromProvider downloads video bytes from provider URL
func (s *videoGenerationService) downloadVideoFromProvider(ctx context.Context, videoURL string) ([]byte, error) {
	if videoURL == "" {
		return nil, errors.New("video URL is empty")
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Read all bytes
	videoBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read video data: %w", err)
	}

	if len(videoBytes) == 0 {
		return nil, errors.New("downloaded video is empty")
	}

	return videoBytes, nil
}

func (s *videoGenerationService) generatePromptForVariant(variantNumber int, customPrompt string) string {
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
