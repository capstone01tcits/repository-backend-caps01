package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type BriefHandler struct {
	briefService service.BriefService
}

func NewBriefHandler(briefService service.BriefService) *BriefHandler {
	return &BriefHandler{briefService}
}

// CreateProjectFromFE godoc
// POST /api/projects/initialize
// Creates project with business brief and creative brief atomically from FE form data
// This is the main endpoint for the unified wizard flow
func (h *BriefHandler) CreateProjectFromFE(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.CreateProjectFromFERequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate required fields from FE
	// Only institutionName is strictly required, others have defaults
	if req.InstitutionName == "" {
		return utils.BadRequest(c, "Institution name is required")
	}
	if req.EventContent == "" || req.ToneOfVoice == "" || req.SelectedKeyMessage == "" || req.SelectedTheme == "" {
		return utils.BadRequest(c, "Event content, tone, key message, and theme are required")
	}

	result, err := h.briefService.CreateProjectFromFE(userID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Project created successfully with briefs", result)
}
