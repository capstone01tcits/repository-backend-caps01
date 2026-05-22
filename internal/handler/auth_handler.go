package handler

import (
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

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
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	profile, err := h.authService.GetProfile(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Profile retrieved", profile)
}

// GetUserProfile godoc
// GET /api/auth/users/:user_id (Protected)
func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error {
	requestingUserID, ok := c.Locals("userID").(string)
	if !ok || requestingUserID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}
	requestingRole, _ := c.Locals("role").(string)
	targetUserID := c.Params("user_id")

	if targetUserID == "" {
		return utils.BadRequest(c, "User ID is required")
	}

	// Authorization check: user can only view their own profile, or admin can view any profile
	if requestingUserID != targetUserID && requestingRole != "admin" {
		return utils.Unauthorized(c, "You cannot access another user's profile")
	}

	profile, err := h.authService.GetProfile(targetUserID)
	if err != nil {
		return utils.NotFound(c, "User not found")
	}

	return utils.OK(c, "User profile retrieved", profile)
}

// GetAllUsers godoc
// GET /api/admin/users (Protected, Admin)
func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Users retrieved", users)
}

// ChangePassword godoc
// POST /api/auth/change-password (Protected)
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req model.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return utils.BadRequest(c, "Old password and new password are required")
	}

	if len(req.NewPassword) < 6 {
		return utils.BadRequest(c, "New password must be at least 6 characters")
	}

	if req.OldPassword == req.NewPassword {
		return utils.BadRequest(c, "New password must be different from old password")
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.OK(c, "Password changed successfully", nil)
}

// DeleteAccount godoc
// DELETE /api/auth/account (Protected)
func (h *AuthHandler) DeleteAccount(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	err := h.authService.DeleteAccount(userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.OK(c, "Account deleted successfully", nil)
}

// RestoreAccount godoc
// POST /api/auth/restore
func (h *AuthHandler) RestoreAccount(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return utils.BadRequest(c, "Refresh token is required")
	}

	profile, err := h.authService.RestoreAccount(req.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.OK(c, "Account restored successfully", profile)
}
