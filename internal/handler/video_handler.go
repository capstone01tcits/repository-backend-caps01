package handler

import (
	"go-auth/internal/model"
	"go-auth/internal/service"
	"go-auth/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type VideoHandler struct {
	videoService service.VideoService
}

func NewVideoHandler(videoService service.VideoService) *VideoHandler {
	return &VideoHandler{videoService}
}

// GenerateVideo godoc
// POST /api/videos/generate
func (h *VideoHandler) GenerateVideo(c *fiber.Ctx) error {
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

	video, err := h.videoService.GenerateVideo(userID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Video generated successfully", video)
}

// GetVideo godoc
// GET /api/videos/:id
func (h *VideoHandler) GetVideo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	videoID := c.Params("id")

	video, err := h.videoService.GetVideo(userID, videoID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Video retrieved", video)
}

// GetVideosByProject godoc
// GET /api/projects/:id/videos
func (h *VideoHandler) GetVideosByProject(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	projectID := c.Params("id")

	videos, err := h.videoService.GetVideosByProject(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Videos retrieved", videos)
}

// GetMyVideos godoc
// GET /api/videos
func (h *VideoHandler) GetMyVideos(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	videos, err := h.videoService.GetVideosByUser(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Videos retrieved", videos)
}

// DownloadVideo godoc
// GET /api/videos/:id/download
func (h *VideoHandler) DownloadVideo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	videoID := c.Params("id")

	video, err := h.videoService.GetVideo(userID, videoID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	if video.Status != "completed" {
		return utils.BadRequest(c, "Video is not ready for download")
	}

	// Stub: return download info (in production, this would stream the file or return a signed URL)
	return utils.OK(c, "Video download ready", fiber.Map{
		"video_id":   video.ID,
		"title":      video.Title,
		"format":     video.Format,
		"resolution": video.Resolution,
		"file_size":  video.FileSize,
		"download_url": video.VideoURL,
	})
}
