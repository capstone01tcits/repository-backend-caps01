package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"Sevima-AI-Content-Creator/internal/model"
)

// BuildVeo3Prompt merangkai prompt sesuai dengan standar Veo 3 untuk video promosi pendidikan.
// Logic ini dipindahkan dari video_worker.go untuk sentralisasi.
func BuildVeo3Prompt(bb *model.BusinessBrief, cb *model.CreativeBrief, sections []model.StoryboardSection) string {
	// 1. Identity Mapping & Context Awareness
	contextKeywords := ""
	schoolLevelStr := strings.ToLower(bb.SchoolLevel)

	if strings.Contains(schoolLevelStr, "university") || strings.Contains(schoolLevelStr, "perguruan tinggi") || strings.Contains(schoolLevelStr, "kampus") {
		contextKeywords = "Higher Education, Campus Life, Research, Independent Learning"
	} else {
		contextKeywords = "Vibrant Classroom, Practical Skills, Student Discipline, Nurturing Environment"
	}

	// 2. Visual Quality Standard
	visualDirection := "Realistic skin textures, natural cinematic lighting, 4K resolution, no plastic-look faces"

	// 3. Merangkai Header Prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString(fmt.Sprintf(
		"Create a cinematic promotional video for %s. Type: %s. Event: %s. Tone: %s. Theme: %s. Duration: %d seconds. Key message: %s. Core emotion: %s. Visual direction: %s.\n\n",
		bb.InstituteName, bb.SchoolLevel, cb.VideoType, cb.Tone, cb.Style, cb.Duration, cb.CallToAction, contextKeywords, visualDirection,
	))

	// 4. Memasukkan Scene (Hook, Value, CTA) beserta perhitungan detik
	timeAccumulator := 0
	for i, sec := range sections {
		startTime := timeAccumulator
		endTime := timeAccumulator + sec.Duration
		timeAccumulator = endTime

		promptBuilder.WriteString(fmt.Sprintf(
			"SCENE %d (%d–%ds): %s [%s]\n",
			i+1, startTime, endTime, strings.ToUpper(sec.SectionType), sec.Content,
		))
	}

	// 5. Aturan Transisi Wajib
	promptBuilder.WriteString("\nMaintain cinematic continuity, same characters, natural skin textures, and smooth transitions.")

	return promptBuilder.String()
}

// Veo3Provider implements VideoProvider for the new Veo 3 AI Service
type Veo3Provider struct {
	baseURL string
}

// NewVeo3Provider creates a new instance of Veo3Provider
func NewVeo3Provider() VideoProvider {
	return &Veo3Provider{
		baseURL: "http://localhost:8000/api/veo3/generate",
	}
}

// GenerateScene sends the prompt and reference images to the AI Service
func (p *Veo3Provider) GenerateScene(ctx context.Context, req VideoGenerationRequest) (*VideoGenerationResponse, error) {
	// Create payload for the new endpoint
	payload := map[string]interface{}{
		"model":            req.Model,
		"prompt":           req.Prompt,
		"reference_images": req.ReferenceImages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal veo3 payload: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ai service request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai service error (%d): %s", resp.StatusCode, string(body))
	}

	// Decode response
	var genResp struct {
		JobID   string `json:"job_id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return nil, fmt.Errorf("failed to decode ai service response: %w", err)
	}

	return &VideoGenerationResponse{
		JobID:   genResp.JobID,
		Status:  genResp.Status,
		Message: genResp.Message,
		Credits: p.CalculateCredits(req.Duration),
	}, nil
}

// GetJobStatus polls the AI Service for job status
func (p *Veo3Provider) GetJobStatus(ctx context.Context, jobID string) (*VideoGenerationResponse, error) {
	// 1. Check status
	statusURL := fmt.Sprintf("http://localhost:8000/status/%s", jobID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", statusURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var statusData struct {
		JobID  string `json:"job_id"`
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&statusData); err != nil {
		return nil, err
	}

	// 2. If done, get the result URL
	if statusData.Status == "done" {
		resultURL := fmt.Sprintf("http://localhost:8000/result/%s", jobID)
		resReq, _ := http.NewRequestWithContext(ctx, "GET", resultURL, nil)
		resResp, err := client.Do(resReq)
		if err == nil {
			defer resResp.Body.Close()
			var resData struct {
				VideoURL string `json:"video_url"`
			}
			json.NewDecoder(resResp.Body).Decode(&resData)
			
			return &VideoGenerationResponse{
				JobID:    jobID,
				Status:   "completed",
				VideoURL: resData.VideoURL,
				Message:  "Video ready",
			}, nil
		}
	}

	// Map Python status to Go standard
	goStatus := statusData.Status
	if goStatus == "done" {
		goStatus = "completed"
	}

	return &VideoGenerationResponse{
		JobID:   jobID,
		Status:  goStatus,
		Message: statusData.Error,
	}, nil
}

// CancelJob cancels an ongoing generation
func (p *Veo3Provider) CancelJob(ctx context.Context, jobID string) error {
	return nil
}

// DownloadVideo downloads the final video file
func (p *Veo3Provider) DownloadVideo(ctx context.Context, videoURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetProviderName returns "Veo3"
func (p *Veo3Provider) GetProviderName() string {
	return "Veo3"
}

// GetModelName returns "veo3"
func (p *Veo3Provider) GetModelName() string {
	return "veo3"
}

// CalculateCredits returns fixed cost for Veo 3
func (p *Veo3Provider) CalculateCredits(duration int) int {
	return duration * 10
}
