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

// ==================== Business Brief Handlers ====================

// CreateBusinessBrief godoc
// POST /api/briefs/business
func (h *BriefHandler) CreateBusinessBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req model.CreateBusinessBriefRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.ProjectName == "" {
		return utils.BadRequest(c, "Project name is required")
	}

	brief, err := h.briefService.CreateBusinessBrief(userID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Business brief created successfully", brief)
}

// GetBusinessBriefs godoc
// GET /api/briefs/business
func (h *BriefHandler) GetBusinessBriefs(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	briefs, err := h.briefService.GetBusinessBriefs(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Business briefs retrieved", briefs)
}

// GetBusinessBrief godoc
// GET /api/briefs/business/:id
func (h *BriefHandler) GetBusinessBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	briefID := c.Params("id")

	brief, err := h.briefService.GetBusinessBrief(userID, briefID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Business brief retrieved", brief)
}

// UpdateBusinessBrief godoc
// PUT /api/briefs/business/:id
func (h *BriefHandler) UpdateBusinessBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	briefID := c.Params("id")

	var req model.UpdateBusinessBriefRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	brief, err := h.briefService.UpdateBusinessBrief(userID, briefID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Business brief updated successfully", brief)
}

// DeleteBusinessBrief godoc
// DELETE /api/briefs/business/:id
func (h *BriefHandler) DeleteBusinessBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	briefID := c.Params("id")

	err := h.briefService.DeleteBusinessBrief(userID, briefID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Business brief deleted successfully", nil)
}

// ==================== Creative Brief Handlers ====================

// CreateCreativeBrief godoc
// POST /api/briefs/creative
func (h *BriefHandler) CreateCreativeBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req model.CreateCreativeBriefRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Title == "" {
		return utils.BadRequest(c, "Title is required")
	}

	if req.BusinessBriefID == "" {
		return utils.BadRequest(c, "Business brief ID is required")
	}

	brief, err := h.briefService.CreateCreativeBrief(userID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Creative brief created successfully", brief)
}

// GetCreativeBriefs godoc
// GET /api/briefs/creative
func (h *BriefHandler) GetCreativeBriefs(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	briefs, err := h.briefService.GetCreativeBriefs(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Creative briefs retrieved", briefs)
}

// GetCreativeBriefsByBusinessBrief godoc
// GET /api/briefs/business/:id/creative
func (h *BriefHandler) GetCreativeBriefsByBusinessBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	businessBriefID := c.Params("id")

	briefs, err := h.briefService.GetCreativeBriefsByBusinessBrief(userID, businessBriefID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Creative briefs retrieved", briefs)
}

// GetCreativeBrief godoc
// GET /api/briefs/creative/:id
func (h *BriefHandler) GetCreativeBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	briefID := c.Params("id")

	brief, err := h.briefService.GetCreativeBrief(userID, briefID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Creative brief retrieved", brief)
}

// UpdateCreativeBrief godoc
// PUT /api/briefs/creative/:id
func (h *BriefHandler) UpdateCreativeBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	briefID := c.Params("id")

	var req model.UpdateCreativeBriefRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	brief, err := h.briefService.UpdateCreativeBrief(userID, briefID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Creative brief updated successfully", brief)
}

// DeleteCreativeBrief godoc
// DELETE /api/briefs/creative/:id
func (h *BriefHandler) DeleteCreativeBrief(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	briefID := c.Params("id")

	err := h.briefService.DeleteCreativeBrief(userID, briefID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Creative brief deleted successfully", nil)
}
