package handler

import (
	"Sevima-AI-Content-Creator/internal/service"
	"Sevima-AI-Content-Creator/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type CreditHandler struct {
	creditService service.CreditService
}

func NewCreditHandler(creditService service.CreditService) *CreditHandler {
	return &CreditHandler{creditService}
}

// GetMyCredits godoc
// GET /api/credits
func (h *CreditHandler) GetMyCredits(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	credits, err := h.creditService.GetCredits(userID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Credits retrieved", fiber.Map{
		"credits": credits,
	})
}

// AddCredits godoc
// POST /api/admin/credits
func (h *CreditHandler) AddCredits(c *fiber.Ctx) error {
	adminUserID, ok := c.Locals("userID").(string)
	if !ok || adminUserID == "" {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var body struct {
		UserID string `json:"user_id"`
		Amount int    `json:"amount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if body.UserID == "" {
		return utils.BadRequest(c, "user_id is required")
	}
	if body.Amount <= 0 {
		return utils.BadRequest(c, "amount must be positive")
	}

	newCredits, err := h.creditService.AddCredits(adminUserID, body.UserID, body.Amount)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "Credits added successfully", fiber.Map{
		"user_id":       body.UserID,
		"credits_added": body.Amount,
		"total_credits": newCredits,
	})
}
