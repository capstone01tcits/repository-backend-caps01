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
	// Generate video from a storyboard
	GenerateVideo(ctx context.Context, userID, projectID, storyboardID uuid.UUID, customPrompt string) (*model.GenerationJob, error)

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

func (s *videoGenerationService) GenerateVideo(ctx context.Context, userID, projectID, storyboardID uuid.UUID, customPrompt string) (*model.GenerationJob, error) {
	// 1. Check if there's already an active or completed video for this storyboard
	existingVideos, _ := s.videoRepo.FindByProjectID(projectID.String())
	for _, v := range existingVideos {
		if v.StoryboardID == storyboardID {
			if v.Status == "pending" || v.Status == "processing" || v.Status == "stitching_video" {
				return nil, errors.New("sedang ada proses generate video yang berjalan untuk storyboard ini")
			}
			if v.Status == "failed" {
				// Delete failed videos to allow retry
				s.videoRepo.Delete(v.ID.String())
			}
		}
	}

	// 2. Validate credits
	credits, err := s.creditService.GetUserCredits(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credits: %w", err)
	}

	// 2.5 Get Storyboard duration
	storyboardSections, err := s.storyboardRepo.FindSectionsByStoryboardID(storyboardID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get storyboard sections: %w", err)
	}
	
	videoDuration := 0
	for _, sec := range storyboardSections {
		videoDuration += sec.Duration
	}
	if videoDuration == 0 {
		videoDuration = 15 // fallback
	}
	
	creditsNeeded := videoDuration

	if credits < creditsNeeded {
		return nil, errors.New("insufficient credits for video generation")
	}

	// 3. Get Project Title for Video Name
	projectName := "Generated Video"

	// 4. Create Video record
	video := &model.Video{
		UserID:       userID,
		ProjectID:    projectID,
		StoryboardID: storyboardID,
		Title:        projectName,
		Status:       "pending",
		Duration:     videoDuration,
		Resolution:   "1080p",
		Format:       "mp4",
		CreditsUsed:  creditsNeeded,
	}

	if err := s.videoRepo.Create(video); err != nil {
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	// 5. Create Job
	job := &model.GenerationJob{
		UserID:          userID,
		ProjectID:       projectID,
		StoryboardID:    storyboardID,
		VideoID:         &video.ID,
		JobType:         "generate",
		Status:          "queued",
		Priority:        1,
		VideoDuration:   videoDuration,
		Provider:        "wavespeed",
		Model:           "veo-3.1-lite",
		Resolution:      "1080p",
		CreditsRequired: creditsNeeded,
		MaxRetries:      3,
	}

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

	// 6. Deduct Credits
	if err := s.creditService.DeductCredits(ctx, userID, creditsNeeded, "Generate Video"); err != nil {
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
	
	fullPrompt := "Generate video based on storyboard"
	if bb != nil && cb != nil {
		fullPrompt = ai.BuildVeo3Prompt(bb, cb, storyboardSections)
	} else if len(storyboardSections) > 0 {
		fullPrompt = storyboardSections[0].Content
	}

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
