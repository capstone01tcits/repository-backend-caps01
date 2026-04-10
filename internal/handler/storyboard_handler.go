package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type StoryboardHandler struct {
	storyboardService service.StoryboardService
}

func NewStoryboardHandler(storyboardService service.StoryboardService) *StoryboardHandler {
	return &StoryboardHandler{storyboardService}
}

// GenerateStoryboards godoc
// POST /api/projects/:id/storyboards/generate
func (h *StoryboardHandler) GenerateStoryboards(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	var body struct {
		ContentThemeID string `json:"content_theme_id"`
		Prompt         string `json:"prompt"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if body.ContentThemeID == "" {
		return utils.BadRequest(c, "content_theme_id is required")
	}

	// Validate prompt length (max 1000 characters)
	if body.Prompt != "" && len(body.Prompt) > 1000 {
		return utils.BadRequest(c, "Prompt must be less than 1000 characters")
	}

	storyboards, err := h.storyboardService.GenerateStoryboards(userID, projectID, body.ContentThemeID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Storyboards generated successfully", storyboards)
}

// GetStoryboards godoc
// GET /api/projects/:id/storyboards
func (h *StoryboardHandler) GetStoryboards(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	storyboards, err := h.storyboardService.GetStoryboards(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboards retrieved", storyboards)
}

// GetStoryboard godoc
// GET /api/storyboards/:id
func (h *StoryboardHandler) GetStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	storyboardID := c.Params("id")

	storyboard, err := h.storyboardService.GetStoryboard(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard retrieved", storyboard)
}

// SelectStoryboard godoc
// POST /api/storyboards/:id/select
func (h *StoryboardHandler) SelectStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	storyboardID := c.Params("id")

	storyboard, err := h.storyboardService.SelectStoryboard(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard selected", storyboard)
}

// GetScenes godoc
// GET /api/storyboards/:id/scenes
func (h *StoryboardHandler) GetScenes(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	storyboardID := c.Params("id")

	scenes, err := h.storyboardService.GetScenes(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Scenes retrieved", scenes)
}

// UpdateStoryboard godoc
// PUT /api/storyboards/:id
func (h *StoryboardHandler) UpdateStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	storyboardID := c.Params("id")

	var req model.UpdateStoryboardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate prompt length (max 1000 characters)
	if req.Prompt != nil && len(*req.Prompt) > 1000 {
		return utils.BadRequest(c, "Prompt must be less than 1000 characters")
	}

	storyboard, err := h.storyboardService.UpdateStoryboard(userID, storyboardID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard updated", storyboard)
}
