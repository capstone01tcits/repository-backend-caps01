package ai

import (
	"context"
)

// VideoGenerationRequest represents a request to generate a video scene
type VideoGenerationRequest struct {
	Prompt      string // Scene description/prompt
	Duration    int    // Duration in seconds (4-6 for scenes)
	Resolution  string // 1080p, 720p, etc
	FPS         int    // Frames per second
	Model           string   // Model name/version
	ReferenceImages []string // Base64 images or URLs for reference
}

// VideoGenerationResponse represents the response from video generation
type VideoGenerationResponse struct {
	JobID      string // External job ID for polling
	Status     string // pending, processing, completed, failed
	VideoURL   string // URL to download video (if completed)
	Message    string // Status message from provider
	Credits    int    // Credits consumed
	ErrorCode  string // Error code if failed
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
	case "veo3":
		return NewVeo3Provider()
	default:
		// Fallback to Veo3 or return nil/error if preferred. 
		// For now keeping it simple as per request to remove mocks.
		return NewVeo3Provider()
	}
}
