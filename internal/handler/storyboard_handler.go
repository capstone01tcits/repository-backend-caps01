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

// GenerateStoryboardTemplates godoc
// POST /api/storyboard/templates/generate
// Generate multiple storyboard template options based on project data and desired duration
func (h *StoryboardHandler) GenerateStoryboardTemplates(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.GenerateStoryboardTemplatesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.ProjectID == "" {
		return utils.BadRequest(c, "project_id is required")
	}

	if req.VideoDuration < 15 || req.VideoDuration > 300 {
		return utils.BadRequest(c, "video_duration must be between 15 and 300 seconds")
	}

	templates, err := h.storyboardService.GenerateTemplates(userID, req.ProjectID, req.VideoDuration)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Templates generated successfully", map[string]interface{}{
		"project_id":     req.ProjectID,
		"video_duration": req.VideoDuration,
		"templates":      templates,
		"count":          len(templates),
	})
}

// CreateManualStoryboard godoc
// POST /api/storyboard/create
// Create manual storyboard with 3 sections: Hook/Intro, Value/Highlight, CTA
func (h *StoryboardHandler) CreateManualStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.CreateManualStoryboardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.ProjectID == "" {
		return utils.BadRequest(c, "project_id is required")
	}

	if req.Title == "" {
		return utils.BadRequest(c, "title is required")
	}

	if len(req.Sections) != 3 {
		return utils.BadRequest(c, "exactly 3 sections (hook, value, cta) are required")
	}

	// Validate section types
	sectionTypes := make(map[string]bool)
	for _, section := range req.Sections {
		if section.SectionType == "" || section.Content == "" {
			return utils.BadRequest(c, "section_type and content are required for all sections")
		}
		if section.SectionType != "hook" && section.SectionType != "value" && section.SectionType != "cta" {
			return utils.BadRequest(c, "section_type must be one of: hook, value, cta")
		}
		sectionTypes[section.SectionType] = true
	}

	if len(sectionTypes) != 3 {
		return utils.BadRequest(c, "all three section types (hook, value, cta) must be present")
	}

	storyboard, err := h.storyboardService.CreateManualStoryboard(userID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Manual storyboard created successfully", storyboard)
}

// GetStoryboards godoc
// GET /api/storyboard/:project_id
// Get all storyboards for a project
func (h *StoryboardHandler) GetStoryboards(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	projectID := c.Params("project_id")
	if projectID == "" {
		return utils.BadRequest(c, "project_id is required")
	}

	storyboards, err := h.storyboardService.GetStoryboards(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboards retrieved successfully", storyboards)
}

// GetStoryboard godoc
// GET /api/storyboard/detail/:storyboard_id
// Get a single storyboard with its sections
func (h *StoryboardHandler) GetStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	storyboard, err := h.storyboardService.GetStoryboard(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard retrieved successfully", storyboard)
}

// SelectStoryboard godoc
// POST /api/storyboard/select
// Select a storyboard for a project
func (h *StoryboardHandler) SelectStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.SelectStoryboardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.StoryboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	storyboard, err := h.storyboardService.SelectStoryboard(userID, req.StoryboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard selected successfully", storyboard)
}

// UpdateStoryboard godoc
// PUT /api/storyboard/:storyboard_id
// Update storyboard and its sections
func (h *StoryboardHandler) UpdateStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	var req model.UpdateManualStoryboardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	storyboard, err := h.storyboardService.UpdateStoryboard(userID, storyboardID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard updated successfully", storyboard)
}

// GetStoryboardSections godoc
// GET /api/storyboard/:storyboard_id/sections
// Get all sections for a storyboard
func (h *StoryboardHandler) GetStoryboardSections(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	sections, err := h.storyboardService.GetStoryboardSections(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Sections retrieved successfully", sections)
}

// DeleteStoryboard godoc
// DELETE /api/storyboard/:storyboard_id
// Soft deletes a storyboard
func (h *StoryboardHandler) DeleteStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	err := h.storyboardService.DeleteStoryboard(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard deleted successfully", nil)
}

// RestoreStoryboard godoc
// POST /api/storyboard/:storyboard_id/restore
// Restores a soft-deleted storyboard
func (h *StoryboardHandler) RestoreStoryboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	err := h.storyboardService.RestoreStoryboard(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard restored successfully", nil)
}

