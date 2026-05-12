//ini adalah service untuk mengirim prompt ke Python API (FastAPI) dan menerima response berupa task_id dan result_url. Service ini digunakan untuk menghubungkan backend Go dengan AI service yang berjalan di Python.

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"Sevima-AI-Content-Creator/internal/model"
)

type PythonSubmitRequest struct {
	Prompt      string `json:"prompt"`
	StoryboardID string `json:"storyboard_id,omitempty"`
}

type PythonSubmitResponse struct {
	TaskID   string `json:"task_id"`
	ResultURL string `json:"result_url"`
}

func SendToPythonAPI(
	storyboardID *model.Storyboard,
	payload *model.Veo3TestPayload,
) (*PythonSubmitResponse, error) {

	// FastAPI endpoint
	url := "http://localhost:8000/submit"

	// Request body
	reqBody := PythonSubmitRequest{
		Prompt:      payload.Prompt,
		StoryboardID: storyboardID.ID.String(),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// HTTP client
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed send request to python api: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %w", err)
	}

	// Check status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"python api error: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	// Parse response
	var result PythonSubmitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed parse response: %w", err)
	}

	return &result, nil
}