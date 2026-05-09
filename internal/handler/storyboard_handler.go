package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StoryboardHandler struct {
	storyboardService service.StoryboardService
	videoGenService   service.VideoGenerationService
}

func NewStoryboardHandler(storyboardService service.StoryboardService, videoGenService service.VideoGenerationService) *StoryboardHandler {
	return &StoryboardHandler{storyboardService, videoGenService}
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

// GetStoryboardByProject godoc
// GET /api/storyboard/:project_id
// Get the storyboard for a project
func (h *StoryboardHandler) GetStoryboardByProject(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	projectID := c.Params("project_id")
	if projectID == "" {
		return utils.BadRequest(c, "project_id is required")
	}

	storyboard, err := h.storyboardService.GetStoryboardByProject(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Storyboard retrieved successfully", storyboard)
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

// GetVeo3TestPayload godoc
// GET /api/storyboard/:storyboard_id/veo3-test
// Returns a payload for testing Veo3 video generation based on storyboard data
func (h *StoryboardHandler) GetVeo3TestPayload(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	payload, err := h.storyboardService.GetVeo3TestPayload(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Veo3 test payload generated successfully", payload)
}

// SendVeo3TestPayload godoc
// POST /api/storyboard/:storyboard_id/veo3-test
// Generates and sends a Veo3 test payload to the AI service
func (h *StoryboardHandler) SendVeo3TestPayload(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	storyboardID := c.Params("storyboard_id")
	if storyboardID == "" {
		return utils.BadRequest(c, "storyboard_id is required")
	}

	// 1. Get the payload using the pakem logic
	payload, err := h.storyboardService.GetVeo3TestPayload(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	// 2. We can either trigger a standard job OR send it directly to the AI provider for "test"
	// To keep it simple and truly "testing" the connection, let's trigger the standard generation job
	// but using the generated prompt.

	storyboardUUID, _ := uuid.Parse(storyboardID)
	
	// We need the ProjectID from the storyboard
	storyboard, err := h.storyboardService.GetStoryboard(userID, storyboardID)
	if err != nil {
		return utils.BadRequest(c, "Gagal mengambil detail storyboard")
	}

	userUUID, _ := uuid.Parse(userID)
	
	// Trigger generation job
	job, err := h.videoGenService.GenerateVideoVariants(c.Context(), userUUID, storyboard.ProjectID, storyboardUUID, payload.Prompt)
	if err != nil {
		return utils.BadRequest(c, "Gagal mengirim request ke AI Service: "+err.Error())
	}

	return utils.Created(c, "Veo3 test request sent successfully", map[string]interface{}{
		"generation_job_id": job.ID,
		"status":            job.Status,
		"payload_sent":     payload,
	})
}

