package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"Sevima-AI-Content-Creator/internal/ai"
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type VideoGenerationService interface {
	// Generate video from a storyboard — creates one job per scene (3 total)
	GenerateVideo(ctx context.Context, userID, projectID, storyboardID uuid.UUID, req model.GenerateVideoRequest) ([]*model.GenerationJob, error)

	// Regenerate a single scene video in-place
	RegenerateScene(ctx context.Context, userID uuid.UUID, videoID uuid.UUID, customPrompt string) (*model.GenerationJob, error)

	// Get generation job status
	GetJobStatus(ctx context.Context, jobID uuid.UUID) (*model.GenerationJob, error)

	// Get video by storyboard ID
	GetVideoByStoryboard(ctx context.Context, storyboardID uuid.UUID) (*model.Video, error)

	// Get all videos by storyboard ID
	GetVideosByStoryboard(ctx context.Context, storyboardID uuid.UUID) ([]model.Video, error)

	// Get video by ID
	GetVideoByID(ctx context.Context, videoID uuid.UUID) (*model.Video, error)

	// Get all videos by user ID
	GetVideosByUserID(ctx context.Context, userID uuid.UUID) ([]model.Video, error)

	// Calculate credits for generation
	CalculateCreditsForGeneration(duration int) int

	// Process generation job (called by worker)
	ProcessGenerationJob(ctx context.Context, jobID uuid.UUID) error

	// Poll and update job status (called by worker)
	PollJobStatus(ctx context.Context, jobID uuid.UUID) error
}

type videoGenerationService struct {
	jobRepo         repository.GenerationJobRepository
	videoRepo       repository.VideoRepository
	projectRepo     repository.ProjectRepository
	creditService   CreditService
	storageService  StorageService
	briefRepo       repository.BriefRepository
	storyboardRepo  repository.StoryboardRepository
	providerFactory *ai.ProviderFactory
}

func NewVideoGenerationService(
	jobRepo repository.GenerationJobRepository,
	videoRepo repository.VideoRepository,
	projectRepo repository.ProjectRepository,
	briefRepo repository.BriefRepository,
	storyboardRepo repository.StoryboardRepository,
	creditService CreditService,
	storageService StorageService,
) VideoGenerationService {
	return &videoGenerationService{
		jobRepo:         jobRepo,
		videoRepo:       videoRepo,
		projectRepo:     projectRepo,
		briefRepo:       briefRepo,
		storyboardRepo:  storyboardRepo,
		creditService:   creditService,
		storageService:  storageService,
		providerFactory: &ai.ProviderFactory{},
	}
}

func creditMultiplierForMode(mode string) int {
	switch mode {
	case "image-to-video":
		return 2
	case "start-end-to-video":
		return 3
	default:
		return 1
	}
}

func (s *videoGenerationService) GenerateVideo(ctx context.Context, userID, projectID, storyboardID uuid.UUID, req model.GenerateVideoRequest) ([]*model.GenerationJob, error) {
	// Normalize video mode
	videoMode := req.VideoMode
	if videoMode == "" {
		videoMode = "text-to-video"
	}
	resolution := req.Resolution
	if resolution == "" {
		resolution = "1080p"
	}

	// 1. Check for any active videos on this storyboard
	existingVideos, _ := s.videoRepo.FindByProjectID(projectID.String())
	for _, v := range existingVideos {
		if v.StoryboardID == storyboardID {
			if v.Status == "pending" || v.Status == "processing" || v.Status == "stitching_video" {
				return nil, errors.New("sedang ada proses generate video yang berjalan untuk storyboard ini")
			}
			if v.Status == "failed" {
				s.videoRepo.Delete(v.ID.String())
			}
		}
	}

	// 2. Get storyboard sections
	storyboardSections, err := s.storyboardRepo.FindSectionsByStoryboardID(storyboardID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get storyboard sections: %w", err)
	}

	// 3. Calculate total credits needed (multiplied by mode)
	totalDuration := 0
	for _, sec := range storyboardSections {
		totalDuration += sec.Duration
	}
	if totalDuration == 0 {
		totalDuration = 18 // fallback: 3 scenes × 6 sec
	}

	multiplier := creditMultiplierForMode(videoMode)
	requiredCredits := totalDuration * multiplier

	credits, err := s.creditService.GetUserCredits(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %w", err)
	}
	if credits < requiredCredits {
		return nil, fmt.Errorf("insufficient credits: need %d, have %d", requiredCredits, credits)
	}

	// 4. Create one Video + GenerationJob per scene (hook, value, cta)
	sectionOrder := []string{"hook", "value", "cta"}
	jobs := make([]*model.GenerationJob, 0, 3)

	for idx, sectionType := range sectionOrder {
		var section *model.StoryboardSection
		for i := range storyboardSections {
			if strings.ToLower(storyboardSections[i].SectionType) == sectionType {
				section = &storyboardSections[i]
				break
			}
		}

		sceneDuration := 6
		sectionContent := ""
		if section != nil {
			if section.Duration > 0 {
				sceneDuration = section.Duration
			}
			sectionContent = section.Content
		}

		video := &model.Video{
			UserID:       userID,
			ProjectID:    projectID,
			StoryboardID: storyboardID,
			Title:        fmt.Sprintf("Scene %d - %s", idx+1, sectionType),
			Status:       "pending",
			Duration:     sceneDuration,
			Resolution:   resolution,
			Format:       "mp4",
			CreditsUsed:  sceneDuration * multiplier,
			SectionType:  sectionType,
			SceneIndex:   idx + 1,
			VideoMode:    videoMode,
		}

		if err := s.videoRepo.Create(video); err != nil {
			return nil, fmt.Errorf("failed to create video record for scene %d: %w", idx+1, err)
		}

		promptData := map[string]interface{}{
			"section_type":    sectionType,
			"scene_index":     idx + 1,
			"section_content": sectionContent,
			"video_mode":      videoMode,
			"start_image":     req.StartImage,
			"end_image":       req.EndImage,
			"negative_prompt": req.NegativePrompt,
			"generate_audio":  req.GenerateAudio,
			"seed":            req.Seed,
			"resolution":      resolution,
		}
		if req.CustomPrompt != "" {
			promptData["custom_prompt"] = req.CustomPrompt
		}
		promptJSON, _ := json.Marshal(promptData)

		job := &model.GenerationJob{
			UserID:          userID,
			ProjectID:       projectID,
			StoryboardID:    storyboardID,
			VideoID:         &video.ID,
			JobType:         "generate",
			Status:          "queued",
			Priority:        1,
			VideoDuration:   sceneDuration,
			Provider:        "wavespeed",
			Model:           "veo-3.1-lite",
			Resolution:      resolution,
			CreditsRequired: sceneDuration * multiplier,
			MaxRetries:      3,
			Prompt:          promptJSON,
		}

		if err := s.jobRepo.Create(ctx, job); err != nil {
			return nil, fmt.Errorf("failed to create generation job for scene %d: %w", idx+1, err)
		}

		jobs = append(jobs, job)
	}

	// 5. Deduct Credits
	modeLabel := map[string]string{
		"text-to-video":      "Text-to-Video",
		"image-to-video":     "Image-to-Video",
		"start-end-to-video": "Start-End-to-Video",
	}[videoMode]
	deductReason := fmt.Sprintf("Generate Video 3 Scenes (%s)", modeLabel)
	if err := s.creditService.DeductCredits(ctx, userID, requiredCredits, deductReason); err != nil {
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}

	return jobs, nil
}

func (s *videoGenerationService) RegenerateScene(ctx context.Context, userID uuid.UUID, videoID uuid.UUID, customPrompt string) (*model.GenerationJob, error) {
	video, err := s.videoRepo.FindByID(videoID.String())
	if err != nil {
		return nil, fmt.Errorf("video not found: %w", err)
	}

	if video.UserID != userID {
		return nil, errors.New("unauthorized: video does not belong to user")
	}

	switch video.Status {
	case "pending", "processing", "stitching_video", "generating_assets":
		return nil, errors.New("scene is already being processed")
	}

	multiplier := creditMultiplierForMode(video.VideoMode)
	creditsRequired := video.Duration * multiplier
	if creditsRequired == 0 {
		creditsRequired = 6 * multiplier
	}

	credits, err := s.creditService.GetUserCredits(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credits: %w", err)
	}
	if credits < creditsRequired {
		return nil, fmt.Errorf("insufficient credits: need %d, have %d", creditsRequired, credits)
	}

	// Get latest section content from storyboard
	storyboardSections, _ := s.storyboardRepo.FindSectionsByStoryboardID(video.StoryboardID.String())
	var sectionContent string
	for _, sec := range storyboardSections {
		if strings.EqualFold(sec.SectionType, video.SectionType) {
			sectionContent = sec.Content
			break
		}
	}

	// Get images from business brief for image-based modes
	var startImage, endImage string
	if video.VideoMode == "image-to-video" || video.VideoMode == "start-end-to-video" {
		if bb, err := s.briefRepo.FindBusinessBriefByProjectID(video.ProjectID.String()); err == nil && bb != nil {
			startImage = bb.LogoPath
			if video.VideoMode == "start-end-to-video" {
				endImage = bb.EnvironmentPath
			}
		}
	}

	// Reset video in-place
	video.Status = "pending"
	video.VideoURL = ""
	video.ThumbnailURL = ""
	video.ErrorMessage = ""
	video.ExternalJobID = ""
	video.RegenerateCount++
	if err := s.videoRepo.Update(video); err != nil {
		return nil, fmt.Errorf("failed to reset video status: %w", err)
	}

	promptData := map[string]interface{}{
		"section_type":    video.SectionType,
		"scene_index":     video.SceneIndex,
		"section_content": sectionContent,
		"video_mode":      video.VideoMode,
		"start_image":     startImage,
		"end_image":       endImage,
		"resolution":      video.Resolution,
	}
	if customPrompt != "" {
		promptData["custom_prompt"] = customPrompt
	}
	promptJSON, _ := json.Marshal(promptData)

	job := &model.GenerationJob{
		UserID:          userID,
		ProjectID:       video.ProjectID,
		StoryboardID:    video.StoryboardID,
		VideoID:         &video.ID,
		JobType:         "regenerate_scene",
		Status:          "queued",
		Priority:        2,
		VideoDuration:   video.Duration,
		Provider:        "wavespeed",
		Model:           "veo-3.1-lite",
		Resolution:      video.Resolution,
		CreditsRequired: creditsRequired,
		MaxRetries:      3,
		Prompt:          promptJSON,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create regeneration job: %w", err)
	}

	deductReason := fmt.Sprintf("Regenerate Scene %d (%s)", video.SceneIndex, video.SectionType)
	if err := s.creditService.DeductCredits(ctx, userID, creditsRequired, deductReason); err != nil {
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}

	return job, nil
}

func (s *videoGenerationService) GetJobStatus(ctx context.Context, jobID uuid.UUID) (*model.GenerationJob, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

func (s *videoGenerationService) GetVideoByStoryboard(ctx context.Context, storyboardID uuid.UUID) (*model.Video, error) {
	return s.videoRepo.FindByStoryboardID(storyboardID.String())
}

func (s *videoGenerationService) GetVideosByStoryboard(ctx context.Context, storyboardID uuid.UUID) ([]model.Video, error) {
	return s.videoRepo.FindAllByStoryboardID(storyboardID.String())
}

func (s *videoGenerationService) GetVideoByID(ctx context.Context, videoID uuid.UUID) (*model.Video, error) {
	return s.videoRepo.FindByID(videoID.String())
}

func (s *videoGenerationService) GetVideosByUserID(ctx context.Context, userID uuid.UUID) ([]model.Video, error) {
	return s.videoRepo.FindByUserID(userID.String())
}

func (s *videoGenerationService) CalculateCreditsForGeneration(duration int) int {
	return duration
}

func (s *videoGenerationService) ProcessGenerationJob(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	if job.VideoID == nil {
		return fmt.Errorf("job has no video ID")
	}

	video, err := s.videoRepo.FindByID(job.VideoID.String())
	if err != nil {
		return fmt.Errorf("failed to find video: %w", err)
	}

	if err := s.jobRepo.UpdateStatus(ctx, jobID, "generating_assets", "Mengirim prompt ke AI Service"); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	video.Status = "processing"
	s.videoRepo.Update(video)

	// Persiapkan prompt
	var bb *model.BusinessBrief
	var cb *model.CreativeBrief
	var refImages []string

	bb, _ = s.briefRepo.FindBusinessBriefByProjectID(job.ProjectID.String())
	if bb != nil {
		cbs, _ := s.briefRepo.FindCreativeBriefsByBusinessBriefID(bb.ID.String())
		if len(cbs) > 0 {
			cb = &cbs[0]
		}
		if bb.LogoPath != "" {
			refImages = append(refImages, bb.LogoPath)
		}
		if bb.EnvironmentPath != "" {
			refImages = append(refImages, bb.EnvironmentPath)
		}
	}

	provider := s.providerFactory.GetProviderByModel(job.Model)
	storyboardSections, _ := s.storyboardRepo.FindSectionsByStoryboardID(job.StoryboardID.String())

	// Read per-scene info from job prompt
	var promptData struct {
		SectionType    string `json:"section_type"`
		SceneIndex     int    `json:"scene_index"`
		SectionContent string `json:"section_content"`
		VideoMode      string `json:"video_mode"`
		StartImage     string `json:"start_image"`
		EndImage       string `json:"end_image"`
		NegativePrompt string `json:"negative_prompt"`
		GenerateAudio  bool   `json:"generate_audio"`
		Seed           int    `json:"seed"`
		Resolution     string `json:"resolution"`
		CustomPrompt   string `json:"custom_prompt"`
	}
	if len(job.Prompt) > 0 {
		json.Unmarshal(job.Prompt, &promptData)
	}
	if promptData.VideoMode == "" {
		promptData.VideoMode = "text-to-video"
	}
	if promptData.Resolution == "" {
		promptData.Resolution = "1080p"
	}

	fullPrompt := "Generate video based on storyboard"
	if bb != nil && cb != nil && promptData.SectionType != "" {
		// Find the specific section for this job
		var targetSection model.StoryboardSection
		for _, sec := range storyboardSections {
			if strings.ToLower(sec.SectionType) == promptData.SectionType {
				targetSection = sec
				break
			}
		}
		// Fallback: use content stored in job prompt
		if targetSection.SectionType == "" {
			targetSection = model.StoryboardSection{
				SectionType: promptData.SectionType,
				Content:     promptData.SectionContent,
				Duration:    job.VideoDuration,
			}
		}
		fullPrompt = ai.BuildScenePrompt(bb, cb, targetSection, promptData.SceneIndex)
	} else if bb != nil && cb != nil {
		fullPrompt = ai.BuildVeo3Prompt(bb, cb, storyboardSections)
	} else if promptData.SectionContent != "" {
		fullPrompt = promptData.SectionContent
	} else if len(storyboardSections) > 0 {
		fullPrompt = storyboardSections[0].Content
	}

	if promptData.CustomPrompt != "" {
		fullPrompt = fullPrompt + ". Additional instructions: " + promptData.CustomPrompt
	}

	// Images are only passed to the provider when the mode explicitly requires them.
	// text-to-video must never receive images even if business brief has logo/environment.
	var startImage, endImage string
	switch promptData.VideoMode {
	case "image-to-video":
		startImage = promptData.StartImage
	case "start-end-to-video":
		startImage = promptData.StartImage
		endImage = promptData.EndImage
	}

	req := ai.VideoGenerationRequest{
		Prompt:         fullPrompt,
		Duration:       job.VideoDuration,
		Resolution:     promptData.Resolution,
		FPS:            30,
		Model:          job.Model,
		VideoMode:      promptData.VideoMode,
		StartImage:     startImage,
		EndImage:       endImage,
		NegativePrompt: promptData.NegativePrompt,
		GenerateAudio:  promptData.GenerateAudio,
		Seed:           promptData.Seed,
	}

	resp, err := provider.GenerateScene(ctx, req)
	if err != nil {
		s.jobRepo.UpdateStatus(ctx, jobID, "failed", fmt.Sprintf("Generation failed: %v", err))
		video.Status = "failed"
		video.ErrorMessage = err.Error()
		s.videoRepo.Update(video)
		return err
	}

	video.ExternalJobID = resp.JobID
	s.videoRepo.Update(video)

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

	if job.VideoID == nil {
		return nil
	}

	video, err := s.videoRepo.FindByID(job.VideoID.String())
	if err != nil {
		return err
	}

	if video.ExternalJobID == "" || video.Status == "completed" || video.Status == "failed" {
		return nil
	}

	provider := s.providerFactory.GetProviderByModel(job.Model)
	resp, err := provider.GetJobStatus(ctx, video.ExternalJobID)
	if err != nil {
		return err
	}

	if resp.Status == "stitching_video" {
		s.jobRepo.UpdateStatus(ctx, jobID, "stitching_video", "Menggabungkan video...")
	}

	if resp.Status == "completed" && resp.VideoURL != "" {
		// Download & Upload to Storage
		videoBytes, err := s.downloadVideoFromProvider(ctx, resp.VideoURL)
		if err == nil {
			videoFilename := fmt.Sprintf("video_%s.mp4", video.ID.String())
			videoPath, uploadErr := s.storageService.UploadVideo(ctx, videoFilename, videoBytes)
			if uploadErr == nil {
				video.VideoURL = s.storageService.GetPublicURL("videos", videoPath)
				video.FileSize = int64(len(videoBytes))
			} else {
				video.VideoURL = resp.VideoURL
			}
		} else {
			video.VideoURL = resp.VideoURL
		}

		// Thumbnail logic
		if resp.ThumbnailURL != "" {
			thumbBytes, err := s.downloadVideoFromProvider(ctx, resp.ThumbnailURL)
			if err == nil {
				thumbFilename := fmt.Sprintf("thumb_%s.jpg", video.ID.String())
				thumbPath, uploadErr := s.storageService.UploadVideo(ctx, thumbFilename, thumbBytes)
				if uploadErr == nil {
					video.ThumbnailURL = s.storageService.GetPublicURL("videos", thumbPath)
				} else {
					video.ThumbnailURL = resp.ThumbnailURL
				}
			} else {
				video.ThumbnailURL = resp.ThumbnailURL
			}
		}

		video.Status = "completed"
		s.videoRepo.Update(video)

		// Update Project status to ready
		if project, err := s.projectRepo.FindByID(video.ProjectID.String()); err == nil {
			project.Status = "ready"
			s.projectRepo.Update(project)
		}

		s.jobRepo.UpdateStatus(ctx, jobID, "completed", "Video siap diunduh")
	} else if resp.Status == "failed" {
		video.Status = "failed"
		video.ErrorMessage = "Provider reported failure"
		s.videoRepo.Update(video)
		
		s.jobRepo.UpdateStatus(ctx, jobID, "failed", "Provider failed")
	}

	return nil
}

func (s *videoGenerationService) downloadVideoFromProvider(ctx context.Context, videoURL string) ([]byte, error) {
	if videoURL == "" {
		return nil, errors.New("video URL is empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	videoBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read video data: %w", err)
	}

	if len(videoBytes) == 0 {
		return nil, errors.New("downloaded video is empty")
	}

	return videoBytes, nil
}
