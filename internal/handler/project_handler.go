package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService}
}

// CreateProject godoc
// POST /api/projects
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" {
		return utils.BadRequest(c, "Project name is required")
	}

	project, err := h.projectService.CreateProject(userID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Project created successfully", project)
}

// GetProjects godoc
// GET /api/projects
func (h *ProjectHandler) GetProjects(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	projects, err := h.projectService.GetProjects(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Projects retrieved", projects)
}

// GetProject godoc
// GET /api/projects/:id
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	project, err := h.projectService.GetProject(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Project retrieved", project)
}

// UpdateProject godoc
// PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	var req model.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	project, err := h.projectService.UpdateProject(userID, projectID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Project updated successfully", project)
}

// DeleteProject godoc
// DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	err := h.projectService.DeleteProject(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Project deleted successfully", nil)
}
