package handler

import (
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

// GenerateStoryboard godoc
// POST /api/storyboard/generate
// Simplified storyboard generation from project (uses auto-selected content theme from initialize step)
func (h *StoryboardHandler) GenerateStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req struct {
		ProjectID string `json:"project_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.ProjectID == "" {
		return utils.BadRequest(c, "project_id is required")
	}

	// For simplified workflow, use first available content theme (auto-selected during initialize)
	// In production, retrieve from CreativeBrief created during project initialization
	storyboards, err := h.storyboardService.GenerateStoryboards(userID, req.ProjectID, "")
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	var storyboardID string
	if len(storyboards) > 0 {
		storyboardID = storyboards[0].ID.String()
	}

	return utils.Created(c, "Storyboard generated successfully", map[string]interface{}{
		"storyboard_id": storyboardID,
		"project_id":    req.ProjectID,
		"status":        "ready",
		"scenes":        storyboards,
	})
}
