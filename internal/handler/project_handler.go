package handler

import (
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

// GetProjects godoc
// GET /api/projects
// Lists all projects for authenticated user
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

// DeleteProject godoc
// DELETE /api/projects/:id
// Soft deletes a project
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

// RestoreProject godoc
// POST /api/projects/:id/restore
// Restores a soft-deleted project
func (h *ProjectHandler) RestoreProject(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	projectID := c.Params("id")

	err := h.projectService.RestoreProject(userID, projectID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Project restored successfully", nil)
}
