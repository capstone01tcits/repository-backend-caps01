package handler

import (
	"go-auth/internal/model"
	"go-auth/internal/service"
	"go-auth/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VideoHandler struct {
	videoGenService service.VideoGenerationService
}

func NewVideoHandler(videoGenService service.VideoGenerationService) *VideoHandler {
	return &VideoHandler{
		videoGenService: videoGenService,
	}
}

// GenerateVideoVariants godoc
// POST /api/videos/generate
// Generates 3 video variants from a storyboard
func (h *VideoHandler) GenerateVideoVariants(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

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

	userUUID, _ := uuid.Parse(userID)
	projectUUID, _ := uuid.Parse(req.ProjectID)
	storyboardUUID, _ := uuid.Parse(req.StoryboardID)

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

// GetGenerationJobStatus godoc
// GET /api/videos/generation/:jobId
func (h *VideoHandler) GetGenerationJobStatus(c *fiber.Ctx) error {
	jobIDStr := c.Params("jobId")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid job ID format")
	}

	job, err := h.videoGenService.GetJobStatus(c.Context(), jobID)
	if err != nil {
		return utils.NotFound(c, "Generation job not found")
	}

	return utils.OK(c, "Job status retrieved", job)
}

// GetVideoVariants godoc
// GET /api/videos/storyboard/:storyboardId
func (h *VideoHandler) GetVideoVariants(c *fiber.Ctx) error {
	storyboardIDStr := c.Params("storyboardId")
	storyboardID, err := uuid.Parse(storyboardIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid storyboard ID format")
	}

	variants, err := h.videoGenService.GetVideoVariants(c.Context(), storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	// Build response with variants and their scenes
	response := make([]map[string]interface{}, len(variants))
	for i, variant := range variants {
		vari, scenes, _ := h.videoGenService.GetVideoVariantWithScenes(c.Context(), variant.ID)
		
		sceneResponses := make([]model.SceneStatusResponse, len(scenes))
		for j, scene := range scenes {
			sceneResponses[j] = model.SceneStatusResponse{
				ID:           scene.ID.String(),
				SceneNumber:  scene.SceneNumber,
				Status:       scene.Status,
				VideoURL:     scene.VideoURL,
				Duration:     scene.Duration,
				ErrorMessage: scene.ErrorMessage,
				UpdatedAt:    scene.UpdatedAt,
			}
		}

		response[i] = map[string]interface{}{
			"id":              vari.ID.String(),
			"variant_number":  vari.VariantNumber,
			"status":          vari.Status,
			"video_url":       vari.VideoURL,
			"thumbnail_url":   vari.ThumbnailURL,
			"prompt_used":     vari.PromptUsed,
			"duration":        vari.Duration,
			"provider":        vari.Provider,
			"model":           vari.Model,
			"scenes":          sceneResponses,
			"created_at":      vari.CreatedAt,
			"updated_at":      vari.UpdatedAt,
		}
	}

	return utils.OK(c, "Video variants retrieved", response)
}

// GetVideoVariant godoc
// GET /api/videos/:variantId
func (h *VideoHandler) GetVideoVariant(c *fiber.Ctx) error {
	variantIDStr := c.Params("variantId")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid variant ID format")
	}

	variant, scenes, err := h.videoGenService.GetVideoVariantWithScenes(c.Context(), variantID)
	if err != nil {
		return utils.NotFound(c, "Video variant not found")
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

	return utils.OK(c, "Video variant retrieved", map[string]interface{}{
		"id":              variant.ID.String(),
		"variant_number":  variant.VariantNumber,
		"status":          variant.Status,
		"video_url":       variant.VideoURL,
		"thumbnail_url":   variant.ThumbnailURL,
		"prompt_used":     variant.PromptUsed,
		"duration":        variant.Duration,
		"provider":        variant.Provider,
		"model":           variant.Model,
		"scenes":          sceneResponses,
		"created_at":      variant.CreatedAt,
		"updated_at":      variant.UpdatedAt,
	})
}

// RegenerateVideoVariant godoc
// POST /api/videos/:variantId/regenerate
func (h *VideoHandler) RegenerateVideoVariant(c *fiber.Ctx) error {
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
	})
}

// RegenerateScene godoc
// POST /api/videos/scene/:sceneId/regenerate
func (h *VideoHandler) RegenerateScene(c *fiber.Ctx) error {
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

	variant, err := h.videoGenService.GetVideoVariant(c.Context(), videoID)
	if err != nil {
		return utils.NotFound(c, "Video variant not found")
	}

	_, scenes, _ := h.videoGenService.GetVideoVariantWithScenes(c.Context(), videoID)
	
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
		"id":              variant.ID.String(),
		"variant_number":  variant.VariantNumber,
		"status":          variant.Status,
		"video_url":       variant.VideoURL,
		"thumbnail_url":   variant.ThumbnailURL,
		"prompt_used":     variant.PromptUsed,
		"duration":        variant.Duration,
		"provider":        variant.Provider,
		"model":           variant.Model,
		"scenes":          sceneResponses,
		"created_at":      variant.CreatedAt,
		"updated_at":      variant.UpdatedAt,
	})
}

// DownloadVideo godoc
// GET /api/videos/:variantId/download
func (h *VideoHandler) DownloadVideo(c *fiber.Ctx) error {
	variantIDStr := c.Params("variantId")
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

	// Return download info with signed URL
	return utils.OK(c, "Video download ready", map[string]interface{}{
		"variant_id":   variant.ID.String(),
		"download_url": variant.VideoURL,
		"file_size":    variant.FileSize,
		"format":       "mp4",
		"resolution":   variant.Resolution,
	})
}
