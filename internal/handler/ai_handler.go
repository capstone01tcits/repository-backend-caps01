package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-auth/config"
	"go-auth/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	client *http.Client
}

func NewAIHandler() *AIHandler {
	return &AIHandler{
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 min timeout for AI processing
		},
	}
}

// Proxy forwards any request under /api/ai/* to the Python AI service.
// The sub-path after /api/ai is appended to AI_SERVICE_URL.
// User context from JWT (X-User-ID, X-User-Email) is injected as headers.
func (h *AIHandler) Proxy(c *fiber.Ctx) error {
	// Build target URL: strip the /api/ai prefix and forward the rest
	subPath := c.Params("*")
	targetURL := fmt.Sprintf("%s/%s", strings.TrimRight(config.Cfg.AIServiceURL, "/"), subPath)

	// Preserve query string
	if qs := string(c.Request().URI().QueryString()); qs != "" {
		targetURL += "?" + qs
	}

	// Create outgoing request
	req, err := http.NewRequestWithContext(c.Context(), c.Method(), targetURL, strings.NewReader(string(c.Body())))
	if err != nil {
		return utils.InternalError(c, "Failed to create proxy request")
	}

	// Forward original headers
	c.Request().Header.VisitAll(func(key, value []byte) {
		k := string(key)
		// Skip hop-by-hop headers
		if strings.EqualFold(k, "Connection") || strings.EqualFold(k, "Host") {
			return
		}
		req.Header.Set(k, string(value))
	})

	// Inject authenticated user context for the AI service
	if userID, ok := c.Locals("userID").(string); ok {
		req.Header.Set("X-User-ID", userID)
	}
	if email, ok := c.Locals("email").(string); ok {
		req.Header.Set("X-User-Email", email)
	}

	// Execute request to AI service
	resp, err := h.client.Do(req)
	if err != nil {
		return utils.InternalError(c, "AI service unavailable: "+err.Error())
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.InternalError(c, "Failed to read AI service response")
	}

	// Forward response headers from AI service
	for key, values := range resp.Header {
		for _, v := range values {
			c.Set(key, v)
		}
	}

	// Return with original status code and body
	return c.Status(resp.StatusCode).Send(body)
}

// HealthCheck pings the Python AI service health endpoint
func (h *AIHandler) HealthCheck(c *fiber.Ctx) error {
	targetURL := fmt.Sprintf("%s/health", strings.TrimRight(config.Cfg.AIServiceURL, "/"))

	resp, err := h.client.Get(targetURL)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"message": "AI service is unreachable",
			"error":   err.Error(),
		})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	return c.Status(resp.StatusCode).Send(body)
}
