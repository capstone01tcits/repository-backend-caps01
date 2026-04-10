package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SafeGetUserID safely retrieves the user ID from the context with proper error handling
func SafeGetUserID(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: invalid or missing user ID")
	}
	return userID, nil
}

// SafeGetUserIDWithResponse retrieves user ID and returns error response if invalid
func SafeGetUserIDWithResponse(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	return userID, nil
}
