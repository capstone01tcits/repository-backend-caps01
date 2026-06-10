package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"Sevima-AI-Content-Creator/config"
	"Sevima-AI-Content-Creator/internal/model"

	"github.com/bregydoc/gtranslate"
)

// themeVisualStyle maps frontend theme selection to Veo 3 cinematographic descriptors.
// Each theme produces a distinctly different visual style in the generated video.
var themeVisualStyle = map[string]struct {
	globalStyle string
	hookLabel   string
	valueLabel  string
	ctaLabel    string
}{
	"tur kampus sinematik": {
		globalStyle: "Cinematic 4K, aerial drone shots, golden hour lighting, smooth dolly movements, epic orchestral score",
		hookLabel:   "aerial drone reveal sweeping over campus at golden hour",
		valueLabel:  "smooth cinematic walkthrough showcasing facilities and campus spaces",
		ctaLabel:    "dramatic wide campus shot with logo reveal, orchestral crescendo",
	},
	"cerita kehidupan mahasiswa": {
		globalStyle: "Documentary-cinematic, warm natural lighting, handheld intimate shots, emotional acoustic score",
		hookLabel:   "intimate close-up of student face, warm natural light, personal story opening",
		valueLabel:  "candid documentary moments — students in labs, libraries, and collaborative spaces",
		ctaLabel:    "warm group shot with hopeful expressions, soft emotional resolution",
	},
	"keunggulan akademik": {
		globalStyle: "Corporate cinematic, clean bright uniform lighting, precise structured composition, professional uplifting score",
		hookLabel:   "clean precision establishing shot of modern academic building, bright and authoritative",
		valueLabel:  "structured showcase of laboratories, equipment, faculty, and academic achievements",
		ctaLabel:    "confident graduate celebration shot, achievement highlight, clear institution branding",
	},
	"tren & gaya hidup cepat": {
		globalStyle: "Dynamic fast-cut style, vivid saturated colors, bold text overlays, modern upbeat electronic score",
		hookLabel:   "fast-cut montage of vibrant campus life, high energy, bold framing",
		valueLabel:  "rapid visual showcase with dynamic transitions and bold on-screen text highlights",
		ctaLabel:    "high-energy finale with punchy call-to-action text, music peak",
	},
}

// toneAtmosphere maps frontend tone selection to visual atmosphere descriptors.
var toneAtmosphere = map[string]string{
	"santai & ramah":         "approachable, warm, friendly, welcoming atmosphere",
	"profesional & formal":   "authoritative, structured, formal, trustworthy atmosphere",
	"kreatif & inovatif":     "bold, creative, experimental, forward-thinking atmosphere",
	"berwibawa & meyakinkan": "powerful, prestigious, commanding, inspiring atmosphere",
}

// BuildVeo3Prompt merangkai prompt sesuai dengan standar Veo 3 untuk video promosi pendidikan.
// Format: bilingual — data institusi tetap dalam bahasa aslinya, arahan sinematik dalam English.
// Tema dan tone yang dipilih user menghasilkan visual style yang berbeda-beda.
func BuildVeo3Prompt(bb *model.BusinessBrief, cb *model.CreativeBrief, sections []model.StoryboardSection) string {
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

	if hook.Duration == 0 {
		hook.Duration = 5
	}
	if value.Duration == 0 {
		value.Duration = 5
	}
	if cta.Duration == 0 {
		cta.Duration = 5
	}

	// Resolve visual style from theme
	style, ok := themeVisualStyle[strings.ToLower(cb.Theme)]
	if !ok {
		style = themeVisualStyle["tur kampus sinematik"]
	}

	// Append tone atmosphere modifier to global style
	globalStyle := style.globalStyle
	if atm, found := toneAtmosphere[strings.ToLower(cb.ToneOfVoice)]; found {
		globalStyle = globalStyle + ", " + atm
	}

	// Optional metadata — only included when provided by user
	var optionalLines strings.Builder
	if cb.EventContent != "" {
		optionalLines.WriteString(fmt.Sprintf("Event context: %s\n", cb.EventContent))
	}
	if cb.KeyMessage != "" {
		optionalLines.WriteString(fmt.Sprintf("Key message: %s\n", cb.KeyMessage))
	}
	if cb.Copywriting != "" {
		optionalLines.WriteString(fmt.Sprintf("Tagline: %s\n", cb.Copywriting))
	}
	if cb.Hashtags != "" {
		optionalLines.WriteString(fmt.Sprintf("Brand keywords: %s\n", cb.Hashtags))
	}

	scene2Start := hook.Duration
	scene3Start := hook.Duration + value.Duration
	totalDuration := hook.Duration + value.Duration + cta.Duration

	prompt := fmt.Sprintf(
		`%s.

Institution: %s | Level: %s | Programs: %s
Background: %s
%s
SCENE 1 (%ds–%ds) — HOOK — %s:
Cinematic scene showing: %s

SCENE 2 (%ds–%ds) — KEY VALUES — %s:
Cinematic scene showing: %s

SCENE 3 (%ds–%ds) — CALL TO ACTION — %s:
Cinematic scene showing: %s. Display institution name "%s" prominently on screen.

Total duration: %ds. Consistent characters across all scenes, natural smooth transitions between scenes.`,
		globalStyle,
		bb.InstitutionName, bb.SchoolLevel, bb.OfferedDegrees,
		bb.InstitutionHistory,
		optionalLines.String(),
		0, hook.Duration, style.hookLabel, hook.Content,
		scene2Start, scene3Start, style.valueLabel, value.Content,
		scene3Start, totalDuration, style.ctaLabel, cta.Content, bb.InstitutionName,
		totalDuration,
	)

	return prompt
}

// BuildScenePrompt builds a rich, scene-specific video generation prompt for WaveSpeed/Veo.
// It parses the storyboard section content (stored as JSON {narration, visual}) and
// combines all brief fields into a structured cinematic prompt.
func BuildScenePrompt(bb *model.BusinessBrief, cb *model.CreativeBrief, section model.StoryboardSection, sceneIndex int) string {
	duration := section.Duration
	if duration <= 0 {
		duration = 6
	}

	// Parse narration & visual from JSON content (saved by FE as {narration, visual})
	var narration, visualDesc string
	var contentJSON struct {
		Narration string `json:"narration"`
		Visual    string `json:"visual"`
	}
	if err := json.Unmarshal([]byte(section.Content), &contentJSON); err == nil {
		narration = contentJSON.Narration
		visualDesc = contentJSON.Visual
	} else {
		narration = section.Content
	}

	sceneType := strings.ToLower(section.SectionType)

	// Per-scene cinematic direction
	sceneDirections := map[string]string{
		"hook":  fmt.Sprintf("OPENING HOOK — Establish %s's identity immediately. Captivate within the first second. Wide cinematic establishing shot transitioning to a revealing close-up.", bb.InstitutionName),
		"value": fmt.Sprintf("VALUE SHOWCASE — Highlight the excellence and unique strengths of %s. Tracking shots, close-ups of facilities and engaged students, dynamic movement.", bb.InstitutionName),
		"cta":   fmt.Sprintf("CALL TO ACTION — Clear brand moment for %s. Logo reveal, institution name displayed prominently. End on a motivating, action-driving visual.", bb.InstitutionName),
	}
	direction := sceneDirections[sceneType]
	if direction == "" {
		direction = "Dynamic scene with smooth transitions and professional cinematography."
	}

	// Theme-specific visual style
	themeStyles := map[string]string{
		"Tur Kampus Sinematik":       "Cinematic aerial drone and ground-level shots. Golden hour warm lighting, smooth slow-motion camera movements. Architectural focus on campus beauty and landmarks.",
		"Cerita Kehidupan Mahasiswa": "Authentic student life moments. Warm natural color grading, slightly handheld naturalistic footage. Social, communal, genuine emotions.",
		"Keunggulan Akademik":        "State-of-the-art laboratory and research environments. Professional studio lighting, technical precision. Achievement trophies, publications, rankings visible.",
		"Tren & Gaya Hidup Cepat":   "Fast-cut energetic editing, vibrant saturated colors. Social media vertical-style aesthetic. Gen Z energy, dynamic quick movements, trending visual style.",
	}
	visualStyle := themeStyles[cb.Theme]
	if visualStyle == "" {
		visualStyle = "Professional cinematic quality with smooth transitions and warm, inviting color palette."
	}

	// Tone guidance
	toneGuides := map[string]string{
		"Santai & Ramah":         "Warm, approachable, and friendly atmosphere. Bright, welcoming visuals. People smiling naturally. Inclusive and supportive environment.",
		"Profesional & Formal":   "Authoritative and polished. Clean, structured composition. Professional gravitas. Confident and trustworthy visual language.",
		"Kreatif & Inovatif":     "Creative energy and innovation. Unexpected angles, bold compositions, innovative visual storytelling. Inspiring and forward-thinking.",
		"Berwibawa & Meyakinkan": "Prestigious and commanding presence. Sophisticated color palette, measured confident pacing. Excellence and tradition radiating from every frame.",
	}
	toneGuide := toneGuides[cb.ToneOfVoice]
	if toneGuide == "" {
		toneGuide = "Professional, engaging, and inspiring tone."
	}

	// Build the prompt — structured for WaveSpeed Veo 3.1 video generation
	prompt := fmt.Sprintf(
		`Create a high-quality %d-second promotional video scene for %s, a %s institution in Indonesia.

SCENE %d OF 3 — %s
%s

INSTITUTION DETAILS:
- Name: %s
- Level: %s
- Programs: %s
- Background: %s

CAMPAIGN:
- Event: %s
- Key Message: "%s"
- Additional Direction: %s

VISUAL STYLE: %s
MOOD & TONE: %s

WHAT TO SHOW (Visual Prompt):
%s

NARRATION / VOICEOVER:
%s

TECHNICAL SPECS:
- Duration: exactly %d seconds
- Resolution: 1080p, cinematic widescreen
- Show institution name "%s" clearly at least once
- Smooth transitions, no abrupt cuts
- Professional color grading consistent throughout
- High production value — broadcast quality`,
		duration,
		bb.InstitutionName, bb.SchoolLevel,
		sceneIndex, strings.ToUpper(sceneType),
		direction,
		bb.InstitutionName, bb.SchoolLevel, bb.OfferedDegrees, bb.InstitutionHistory,
		cb.EventContent, cb.KeyMessage,
		func() string {
			if cb.Prompt != "" {
				return cb.Prompt
			}
			return "None"
		}(),
		visualStyle,
		toneGuide,
		func() string {
			if visualDesc != "" {
				return visualDesc
			}
			return fmt.Sprintf("Cinematic scene showcasing %s campus and student life.", bb.InstitutionName)
		}(),
		func() string {
			if narration != "" {
				return narration
			}
			return cb.KeyMessage
		}(),
		duration,
		bb.InstitutionName,
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
	translatedPrompt := req.Prompt

	// Menerjemahkan prompt dari Bahasa Indonesia ke Bahasa Inggris menggunakan gtranslate
	translated, err := gtranslate.TranslateWithParams(
		req.Prompt,
		gtranslate.TranslationParams{
			From: "id",
			To:   "en",
		},
	)
	if err == nil {
		translatedPrompt = translated
		log.Printf("[Veo3Provider] Successfully translated prompt to English: %.100s...", translatedPrompt)
	} else {
		log.Printf("[Veo3Provider] Warning: translation failed (%v), using original prompt", err)
	}

	// Resolve resolution (default 480p if not set)
	resolution := req.Resolution
	if resolution == "" {
		resolution = "480p"
	}

	// Seed: use 0 as default when caller passes 0 (means "random")
	seed := req.Seed
	if seed == 0 {
		seed = -1
	}

	videoMode := req.VideoMode
	if videoMode == "" {
		videoMode = "text-to-video"
	}

	// POST /generate → Python AI Service → Wavespeed
	payload := map[string]interface{}{
		"prompt":         translatedPrompt,
		"duration":       req.Duration,
		"ratio":          "16:9",
		"task_type":      "veo3",
		"video_mode":     videoMode,
		"resolution":     resolution,
		"generate_audio": req.GenerateAudio,
		"seed":           seed,
	}
	// Only include images when mode explicitly requires them
	if videoMode == "image-to-video" || videoMode == "start-end-to-video" {
		if req.StartImage != "" {
			payload["start_image"] = req.StartImage
		}
	}
	if videoMode == "start-end-to-video" {
		if req.EndImage != "" {
			payload["end_image"] = req.EndImage
		}
	}
	if req.NegativePrompt != "" {
		payload["negative_prompt"] = req.NegativePrompt
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
				VideoURL     string `json:"video_url"`
				ThumbnailURL string `json:"thumbnail_url"`
			}
			json.NewDecoder(resResp.Body).Decode(&resData)

			return &VideoGenerationResponse{
				JobID:        jobID,
				Status:       "completed",
				VideoURL:     resData.VideoURL,
				ThumbnailURL: resData.ThumbnailURL,
				Message:      "Video ready",
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
