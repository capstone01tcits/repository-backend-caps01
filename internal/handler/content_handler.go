package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ContentHandler struct {
	contentService service.ContentService
}

func NewContentHandler(contentService service.ContentService) *ContentHandler {
	return &ContentHandler{contentService}
}

// GenerateContentPillars godoc
// POST /api/projects/:id/content-pillars/generate
func (h *ContentHandler) GenerateContentPillars(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	pillars, err := h.contentService.GenerateContentPillars(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Content pillars generated successfully", pillars)
}

// GetContentPillars godoc
// GET /api/projects/:id/content-pillars
func (h *ContentHandler) GetContentPillars(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	pillars, err := h.contentService.GetContentPillars(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Content pillars retrieved", pillars)
}

// GetContentPillar godoc
// GET /api/content-pillars/:id
func (h *ContentHandler) GetContentPillar(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	pillarID := c.Params("id")

	pillar, err := h.contentService.GetContentPillar(userID, pillarID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Content pillar retrieved", pillar)
}

// SelectContentPillar godoc
// POST /api/content-pillars/:id/select
func (h *ContentHandler) SelectContentPillar(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	pillarID := c.Params("id")

	pillar, err := h.contentService.SelectContentPillar(userID, pillarID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Content pillar selected", pillar)
}

// GetContentThemes godoc
// GET /api/content-pillars/:id/themes
func (h *ContentHandler) GetContentThemes(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	pillarID := c.Params("id")

	themes, err := h.contentService.GetContentThemes(userID, pillarID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Content themes retrieved", themes)
}

// SelectContentTheme godoc
// POST /api/content-themes/:id/select
func (h *ContentHandler) SelectContentTheme(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	themeID := c.Params("id")

	theme, err := h.contentService.SelectContentTheme(userID, themeID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Content theme selected", theme)
}

// UpdateContentPillar godoc
// PUT /api/content-pillars/:id
func (h *ContentHandler) UpdateContentPillar(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	pillarID := c.Params("id")

	var req model.UpdateContentPillarRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate prompt length (max 1000 characters)
	if req.Prompt != nil && len(*req.Prompt) > 1000 {
		return utils.BadRequest(c, "Prompt must be less than 1000 characters")
	}

	pillar, err := h.contentService.UpdateContentPillar(userID, pillarID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Content pillar updated", pillar)
}
