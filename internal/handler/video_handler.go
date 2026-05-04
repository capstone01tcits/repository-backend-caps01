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

// generateVideoVariants is internal method used by GenerateVideo
// Generates 3 video variants from a storyboard
func (h *VideoHandler) generateVideoVariants(c *fiber.Ctx) error {
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

	job, err := h.videoGenService.GenerateVideoVariants(c.Context(), userUUID, projectUUID, storyboardUUID, req.CustomPrompt)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Video generation job created", map[string]interface{}{
		"generation_job_id": job.ID,
		"status":            job.Status,
		"created_at":        job.CreatedAt,
	})
}

// DownloadVideo godoc
// GET /api/videos/download/:id
// Downloads generated video variant
func (h *VideoHandler) DownloadVideo(c *fiber.Ctx) error {
	variantIDStr := c.Params("id")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid variant ID format")
	}

	variant, err := h.videoGenService.GetVideoVariant(c.Context(), variantID)
	if err != nil {
		return utils.BadRequest(c, "Video variant not found")
	}

	if variant.Status != "completed" {
		return utils.BadRequest(c, "Video is not ready for download")
	}

	// Validate video URL exists (file existence check)
	if variant.VideoURL == "" {
		return utils.InternalError(c, "Video file not available - URL missing")
	}

	// Return download info with signed URL
	return utils.OK(c, "Video download ready", map[string]interface{}{
		"variant_id":   variant.ID.String(),
		"download_url": variant.VideoURL,
		"file_size":    variant.FileSize,
		"format":       "mp4",
		"resolution":   variant.Resolution,
	})
}

// GenerateVideo godoc
// POST /api/videos/generate
// Generates video from storyboard
func (h *VideoHandler) GenerateVideo(c *fiber.Ctx) error {
	return h.generateVideoVariants(c)
}

// GetVideo godoc
// GET /api/videos/:id
// Gets single video/variant by ID
func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	variantIDStr := c.Params("id")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid video ID format")
	}

	variant, scenes, err := h.videoGenService.GetVideoVariantWithScenes(c.Context(), variantID)
	if err != nil {
		return utils.NotFound(c, "Video not found")
	}

	sceneResponses := make([]model.SceneStatusResponse, len(scenes))
	for i, scene := range scenes {
		sceneResponses[i] = model.SceneStatusResponse{
			ID:           scene.ID.String(),
			SceneNumber:  scene.SceneNumber,
			Status:       scene.Status,
			VideoURL:     scene.VideoURL,
			Duration:     scene.Duration,
			ErrorMessage: scene.ErrorMessage,
			UpdatedAt:    scene.UpdatedAt,
		}
	}

	return utils.OK(c, "Video retrieved", map[string]interface{}{
		"id":             variant.ID.String(),
		"variant_number": variant.VariantNumber,
		"status":         variant.Status,
		"video_url":      variant.VideoURL,
		"thumbnail_url":  variant.ThumbnailURL,
		"prompt_used":    variant.PromptUsed,
		"duration":       variant.Duration,
		"provider":       variant.Provider,
		"model":          variant.Model,
		"scenes":         sceneResponses,
		"created_at":     variant.CreatedAt,
		"updated_at":     variant.UpdatedAt,
	})
}

// RegenerateVideoVariant godoc
// POST /api/videos/:variantId/regenerate
// Regenerates a specific video variant with optional new prompt
func (h *VideoHandler) RegenerateVideoVariant(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	variantIDStr := c.Params("variantId")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid variant ID format")
	}

	var req model.RegenerateVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	job, err := h.videoGenService.RegenerateVideoVariant(c.Context(), variantID, req.NewPrompt)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Video regeneration job created", map[string]interface{}{
		"generation_job_id": job.ID,
		"status":            job.Status,
		"created_at":        job.CreatedAt,
	})
}

// RegenerateScene godoc
// POST /api/videos/scene/:sceneId/regenerate
// Regenerates a specific scene within a video variant
func (h *VideoHandler) RegenerateScene(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	sceneIDStr := c.Params("sceneId")
	sceneID, err := uuid.Parse(sceneIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid scene ID format")
	}

	var req model.RegenerateSceneRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	job, err := h.videoGenService.RegenerateScene(c.Context(), sceneID, req.NewPrompt)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Scene regeneration job created", map[string]interface{}{
		"generation_job_id": job.ID,
		"status":            job.Status,
		"created_at":        job.CreatedAt,
	})
}

// ListVideos godoc
// GET /api/videos
// Lists all videos for authenticated user
func (h *VideoHandler) ListVideos(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	// TODO: Implement video listing from database
	// For now return empty list - this would query all VideoVariants owned by user
	videos := make([]interface{}, 0)

	return utils.OK(c, "Videos retrieved", videos)
}
