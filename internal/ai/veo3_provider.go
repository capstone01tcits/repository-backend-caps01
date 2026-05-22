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

	"Sevima-AI-Content-Creator/config"
	"Sevima-AI-Content-Creator/internal/model"
)

// BuildVeo3Prompt merangkai prompt sesuai dengan standar Veo 3 untuk video promosi pendidikan.
func BuildVeo3Prompt(bb *model.BusinessBrief, cb *model.CreativeBrief, sections []model.StoryboardSection) string {
	// Identify the 3 scenes (hook, value, cta)
	var hook, value, cta model.StoryboardSection
	for _, sec := range sections {
		switch strings.ToLower(sec.SectionType) {
		case "hook":
			hook = sec
		case "value":
			value = sec
		case "cta":
			cta = sec
		}
	}

	// Fallback jika tidak ada section (untuk amannya)
	if hook.Duration == 0 {
		hook.Duration = 5
	}
	if value.Duration == 0 {
		value.Duration = 5
	}
	if cta.Duration == 0 {
		cta.Duration = 5
	}

	// Standard phrasing (kata-kata pakem)
	prompt := fmt.Sprintf(
		`Buatlah video promosi %s berkualitas tinggi untuk iklan institusi pendidikan.

		Detail Institusi:
		- Nama: %s
		- Tingkat Pendidikan: %s
		- Program Studi: %s
		- Latar Belakang/Sejarah: %s
		- Gunakan logo dan foto lingkungan kampus sebagai referensi visual.

		Tujuan Video:
		Membuat video promosi yang menarik, modern, profesional, dan membangun kepercayaan, dengan menonjolkan kualitas akademik, fasilitas, serta peluang masa depan.

		Gaya & Tone:
		%s
		Target Audiens: Calon siswa/mahasiswa dan orang tua.
		Gaya Visual: Sinematik, bersih, profesional, transisi halus, kualitas produksi tinggi.

		SCENE STRUCTURE:

		SCENE 1 (%ds–%ds): HOOK
		%s

		SCENE 2 (%ds–%ds): NILAI UNGGULAN
		%s

		SCENE 3 (%ds–%ds): CALL TO ACTION
		%s

		Panduan Teknis:
		- Pertahankan kesinambungan sinematik
		- Gunakan karakter yang konsisten di setiap scene
		- Transisi antar scene harus halus dan natural
		- Tampilkan nama institusi dengan jelas
		- Tonjolkan kepercayaan, prestasi, dan masa depan cerah
		- Tambahkan nuansa musik latar inspiratif
		- Hindari klaim berlebihan atau tidak realistis
		- tolong hasilkan video dengan kualitas produksi tinggi, fokus pada detail visual dan storytelling yang kuat.`,
		cb.Theme,
		bb.InstitutionName,
		bb.SchoolLevel,
		bb.OfferedDegrees,
		bb.InstitutionHistory,
		cb.Theme,
		0, hook.Duration, hook.Content,
		hook.Duration, hook.Duration+value.Duration, value.Content,
		hook.Duration+value.Duration, hook.Duration+value.Duration+cta.Duration, cta.Content,
	)

	return prompt
}

// Veo3Provider implements VideoProvider — forwards to the Python AI Service
// which now routes all veo3/veo-3.1 requests to Wavespeed.
type Veo3Provider struct {
	aiServiceURL string
}

// NewVeo3Provider creates a new instance of Veo3Provider
func NewVeo3Provider() VideoProvider {
	return &Veo3Provider{
		aiServiceURL: strings.TrimRight(config.Cfg.AIServiceURL, "/"),
	}
}

// GenerateScene submits a video generation job to the Python AI Service.
// The AI Service routes the "veo3" task_type to Wavespeed internally.
func (p *Veo3Provider) GenerateScene(ctx context.Context, req VideoGenerationRequest) (*VideoGenerationResponse, error) {
	// POST /generate with task_type=veo3 → Python AI Service → Wavespeed
	payload := map[string]interface{}{
		"prompt":    req.Prompt,
		"duration":  req.Duration,
		"ratio":     "16:9",
		"task_type": "veo3",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/generate", p.aiServiceURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ai service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai service error (%d): %s", resp.StatusCode, string(body))
	}

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

// GetJobStatus polls the Python AI Service for job status, then fetches result if done.
func (p *Veo3Provider) GetJobStatus(ctx context.Context, jobID string) (*VideoGenerationResponse, error) {
	// GET /status/{job_id}
	statusURL := fmt.Sprintf("%s/status/%s", p.aiServiceURL, jobID)
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

	// If done, fetch result to get video URL
	if statusData.Status == "done" {
		resultURL := fmt.Sprintf("%s/result/%s", p.aiServiceURL, jobID)
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

	// Map Python statuses to our Go standard
	goStatus := statusData.Status
	switch goStatus {
	case "done":
		goStatus = "completed"
	case "pending", "processing":
		// keep as-is
	}

	return &VideoGenerationResponse{
		JobID:   jobID,
		Status:  goStatus,
		Message: statusData.Error,
	}, nil
}

// CancelJob is a no-op for Wavespeed-backed jobs.
func (p *Veo3Provider) CancelJob(ctx context.Context, jobID string) error {
	return nil
}

// DownloadVideo downloads the final video file from a URL.
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

// GetProviderName returns "Veo3/Wavespeed"
func (p *Veo3Provider) GetProviderName() string {
	return "Veo3/Wavespeed"
}

// GetModelName returns "wavespeed"
func (p *Veo3Provider) GetModelName() string {
	return "wavespeed"
}

// CalculateCredits returns cost estimate (1 credit per second).
func (p *Veo3Provider) CalculateCredits(duration int) int {
	return duration * 1
}
