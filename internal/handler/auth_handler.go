package handler

import (
	"go-auth/internal/model"
	"go-auth/internal/service"
	"go-auth/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Register godoc
// POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Name, email, and password are required")
	}

	if len(req.Password) < 6 {
		return utils.BadRequest(c, "Password must be at least 6 characters")
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Registration successful", result)
}

// Login godoc
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Email and password are required")
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.OK(c, "Login successful", result)
}

// RefreshToken godoc
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return utils.BadRequest(c, "Refresh token is required")
	}

	result, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.OK(c, "Token refreshed", result)
}

// GetProfile godoc
// GET /api/auth/me (Protected)
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	profile, err := h.authService.GetProfile(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Profile retrieved", profile)
}
