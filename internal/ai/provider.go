package ai

import (
	"context"
)

// VideoGenerationRequest represents a request to generate a video scene
type VideoGenerationRequest struct {
	Prompt          string   // Scene description/prompt
	Duration        int      // Duration in seconds (4, 6, or 8)
	Resolution      string   // "720p" | "1080p"
	FPS             int      // Frames per second
	Model           string   // Model name/version
	ReferenceImages []string // Legacy: URLs for reference (unused by Wavespeed)

	// Mode selection
	VideoMode string // "text-to-video" | "image-to-video" | "start-end-to-video"

	// Image inputs (CDN URLs)
	StartImage string // Reference/start frame image URL
	EndImage   string // End frame image URL (start-end-to-video only)

	// Generation controls
	NegativePrompt string // What to avoid in the video
	GenerateAudio  bool   // Enable synchronized audio (text-to-video only)
	Seed           int    // -1 = random
}

// VideoGenerationResponse represents the response from video generation
type VideoGenerationResponse struct {
	JobID        string // External job ID for polling
	Status       string // pending, processing, completed, failed
	VideoURL     string // URL to download video (if completed)
	ThumbnailURL string // URL to download thumbnail (if completed)
	Message      string // Status message from provider
	Credits      int    // Credits consumed
	ErrorCode    string // Error code if failed
}

// VideoProvider interface defines methods for different video generation providers
type VideoProvider interface {
	// GenerateScene generates a video scene based on prompt
	GenerateScene(ctx context.Context, req VideoGenerationRequest) (*VideoGenerationResponse, error)

	// GetJobStatus polls the status of a generation job
	GetJobStatus(ctx context.Context, jobID string) (*VideoGenerationResponse, error)

	// CancelJob cancels an ongoing generation
	CancelJob(ctx context.Context, jobID string) error

	// DownloadVideo downloads the generated video
	DownloadVideo(ctx context.Context, videoURL string) ([]byte, error)

	// GetProviderName returns the name of the provider
	GetProviderName() string

	// GetModelName returns the model name
	GetModelName() string

	// CalculateCredits calculates credits needed for generation
	CalculateCredits(duration int) int
}

// ProviderFactory creates appropriate provider based on tier
type ProviderFactory struct {
	// Could inject different providers here
}

// GetProvider returns the appropriate provider based on tier
func (pf *ProviderFactory) GetProvider(tier string, model string) VideoProvider {
	// Default to Veo3 since other providers are mocks and being removed
	return NewVeo3Provider()
}

// GetProviderByModel returns provider based on specific model name
func (pf *ProviderFactory) GetProviderByModel(model string) VideoProvider {
	switch model {
	case "veo-3.1-lite":
		return NewVeo3Provider()
	default:
		// Fallback: all models go through Veo3Provider → Wavespeed
		return NewVeo3Provider()
	}
}
