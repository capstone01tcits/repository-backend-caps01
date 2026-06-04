package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VideoHandler struct {
	videoGenService service.VideoGenerationService
	storageService  service.StorageService
}

func NewVideoHandler(videoGenService service.VideoGenerationService, storageService service.StorageService) *VideoHandler {
	return &VideoHandler{
		videoGenService: videoGenService,
		storageService:  storageService,
	}
}

// GenerateVideo godoc
// POST /api/videos/generate
func (h *VideoHandler) GenerateVideo(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.GenerateVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.ProjectID == "" {
		return utils.BadRequest(c, "project_id is required")
	}
	if req.StoryboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID format")
	}
	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return utils.BadRequest(c, "Invalid project_id format")
	}
	storyboardUUID, err := uuid.Parse(req.StoryboardID)
	if err != nil {
		return utils.BadRequest(c, "Invalid storyboard_id format")
	}

	job, err := h.videoGenService.GenerateVideo(c.Context(), userUUID, projectUUID, storyboardUUID, req.CustomPrompt)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Video generation job created", map[string]interface{}{
		"generation_job_id": job.ID,
		"status":            job.Status,
		"video_id":          job.VideoID,
		"created_at":        job.CreatedAt,
	})
}

// DownloadVideo godoc
// GET /api/videos/download/:id
func (h *VideoHandler) DownloadVideo(c *fiber.Ctx) error {
	videoIDStr := c.Params("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid video ID format")
	}

	video, err := h.videoGenService.GetVideoByID(c.Context(), videoID)
	if err != nil {
		return utils.BadRequest(c, "Video not found")
	}

	if video.Status != "completed" {
		return utils.BadRequest(c, "Video is not ready for download")
	}

	if video.VideoURL == "" {
		return utils.InternalError(c, "Video file not available - URL missing")
	}

	return utils.OK(c, "Video download ready", map[string]interface{}{
		"video_id":     video.ID.String(),
		"download_url": video.VideoURL,
		"file_size":    video.FileSize,
		"format":       video.Format,
		"resolution":   video.Resolution,
	})
}

// PreviewVideo godoc
// GET /api/videos/preview/:id
func (h *VideoHandler) PreviewVideo(c *fiber.Ctx) error {
	videoIDStr := c.Params("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid video ID format")
	}

	video, err := h.videoGenService.GetVideoByID(c.Context(), videoID)
	if err != nil {
		return utils.NotFound(c, "Video not found")
	}

	if video.VideoURL == "" {
		return utils.NotFound(c, "Video file not available")
	}

	return utils.OK(c, "Preview URL retrieved", map[string]interface{}{
		"preview_url": video.VideoURL,
	})
}

// GetVideo godoc
// GET /api/videos/:id
func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	videoIDStr := c.Params("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid video ID format")
	}

	video, err := h.videoGenService.GetVideoByID(c.Context(), videoID)
	if err != nil {
		return utils.NotFound(c, "Video not found")
	}

	return utils.OK(c, "Video retrieved", map[string]interface{}{
		"id":             video.ID.String(),
		"title":          video.Title,
		"status":         video.Status,
		"video_url":      video.VideoURL,
		"thumbnail_url":  video.ThumbnailURL,
		"duration":       video.Duration,
		"created_at":     video.CreatedAt,
		"updated_at":     video.UpdatedAt,
	})
}

// ListVideos godoc
// GET /api/videos
func (h *VideoHandler) ListVideos(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID format")
	}

	videos, err := h.videoGenService.GetVideosByUserID(c.Context(), userUUID)
	if err != nil {
		return utils.InternalError(c, "Failed to retrieve videos")
	}

	type videoResp struct {
		ID           string      `json:"id"`
		Title        string      `json:"title"`
		Status       string      `json:"status"`
		VideoURL     string      `json:"video_url"`
		ThumbnailURL string      `json:"thumbnail_url"`
		Duration     int         `json:"duration"`
		CreatedAt    interface{} `json:"created_at"`
		UpdatedAt    interface{} `json:"updated_at"`
	}

	result := make([]videoResp, len(videos))
	for i, v := range videos {
		result[i] = videoResp{
			ID:           v.ID.String(),
			Title:        v.Title,
			Status:       v.Status,
			VideoURL:     v.VideoURL,
			ThumbnailURL: v.ThumbnailURL,
			Duration:     v.Duration,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		}
	}

	return utils.OK(c, "Videos retrieved", result)
}

// GetVideosByStoryboard godoc
// GET /api/videos/storyboard/:storyboard_id
func (h *VideoHandler) GetVideosByStoryboard(c *fiber.Ctx) error {
	storyboardIDStr := c.Params("storyboard_id")
	storyboardID, err := uuid.Parse(storyboardIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid storyboard_id format")
	}

	videos, err := h.videoGenService.GetVideosByStoryboard(c.Context(), storyboardID)
	if err != nil || len(videos) == 0 {
		// Return empty list to match frontend expectations if not found
		return utils.OK(c, "No video found for storyboard", []interface{}{})
	}

	// Map to response struct
	type videoResp struct {
		ID           string      `json:"id"`
		Status       string      `json:"status"`
		VideoURL     string      `json:"video_url"`
		ThumbnailURL string      `json:"thumbnail_url"`
		Duration     int         `json:"duration"`
		CreatedAt    interface{} `json:"created_at"`
		UpdatedAt    interface{} `json:"updated_at"`
	}

	var result []videoResp
	for _, video := range videos {
		result = append(result, videoResp{
			ID:           video.ID.String(),
			Status:       video.Status,
			VideoURL:     video.VideoURL,
			ThumbnailURL: video.ThumbnailURL,
			Duration:     video.Duration,
			CreatedAt:    video.CreatedAt,
			UpdatedAt:    video.UpdatedAt,
		})
	}

	return utils.OK(c, "Videos retrieved", result)
}

// Stubs for removed methods to avoid compilation errors if referenced elsewhere
func (h *VideoHandler) RegenerateVideoVariant(c *fiber.Ctx) error {
	return utils.BadRequest(c, "Endpoint removed")
}

func (h *VideoHandler) RegenerateScene(c *fiber.Ctx) error {
	return utils.BadRequest(c, "Endpoint removed")
}
