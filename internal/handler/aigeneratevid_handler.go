//temp handler ai
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

//
// =========================
// DOMAIN: REQUEST / RESPONSE
// =========================
//

type GenerateVideoRequest struct {
	ProjectID    string `json:"project_id"`
	StoryboardID string `json:"storyboard_id"`
	CustomPrompt string `json:"custom_prompt"`
}

type WaveSpeedRequest struct {
	Model        string `json:"model"`
	Prompt       string `json:"prompt"`
	StoryboardID string `json:"storyboard_id"`
	SceneNumber  int    `json:"scene_number"`
}

type WaveSpeedResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

type SubmitVideoResponse struct {
	GenerationJobID string `json:"generation_job_id"`
	Status          string `json:"status"`
}

//
// =========================
// GATEWAY (WaveSpeed AI)
// =========================
//

type WaveSpeedGateway struct {
	baseURL string
	client  *http.Client
}

func NewWaveSpeedGateway(baseURL string) *WaveSpeedGateway {
	return &WaveSpeedGateway{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (g *WaveSpeedGateway) GenerateVideo(req WaveSpeedRequest) (*WaveSpeedResponse, error) {

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		g.baseURL+"/generate-video",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("wavespeed error: %s", string(body))
	}

	var result WaveSpeedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

//
// =========================
// USECASE (Business Logic)
// =========================
//

type VideoGenerationUsecase struct {
	wavespeed *WaveSpeedGateway
}

func NewVideoGenerationUsecase(ws *WaveSpeedGateway) *VideoGenerationUsecase {
	return &VideoGenerationUsecase{
		wavespeed: ws,
	}
}

// simple dummy scene generator (bisa kamu ganti AI prompt service)
func buildScenePrompt(customPrompt string, scene int) string {
	return fmt.Sprintf("%s - cinematic scene %d with dramatic lighting", customPrompt, scene)
}

func (u *VideoGenerationUsecase) SubmitGeneration(
	ctx context.Context,
	req GenerateVideoRequest,
) (*SubmitVideoResponse, error) {

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, errors.New("invalid project_id")
	}

	storyboardID, err := uuid.Parse(req.StoryboardID)
	if err != nil {
		return nil, errors.New("invalid storyboard_id")
	}

	_ = projectID // placeholder kalau nanti mau DB job
	_ = storyboardID

	// simulate multi-scene generation
	sceneCount := 3

	var lastJobID string

	for i := 1; i <= sceneCount; i++ {

		waveReq := WaveSpeedRequest{
			Model:        "veo3",
			Prompt:       buildScenePrompt(req.CustomPrompt, i),
			StoryboardID: req.StoryboardID,
			SceneNumber:  i,
		}

		resp, err := u.wavespeed.GenerateVideo(waveReq)
		if err != nil {
			return nil, err
		}

		lastJobID = resp.JobID
	}

	return &SubmitVideoResponse{
		GenerationJobID: lastJobID,
		Status:          "processing",
	}, nil
}

//
// =========================
// HANDLER (Fiber HTTP)
// =========================
//

type Handler struct {
	usecase *VideoGenerationUsecase
}

func NewHandler(uc *VideoGenerationUsecase) *Handler {
	return &Handler{usecase: uc}
}

// POST /api/ai-video/generate
func (h *Handler) SubmitWaveSpeedGeneration(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	var req GenerateVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if req.ProjectID == "" || req.StoryboardID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "project_id and storyboard_id required",
		})
	}

	result, err := h.usecase.SubmitGeneration(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "video generation submitted",
		"data":    result,
	})
}