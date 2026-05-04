package service

import (
	"context"
	"fmt"
	"log"
)

// Example usage of StorageService
// This shows how to integrate file uploads to Supabase in your services

/*

// In your video_generation_service.go or handler, inject StorageService:
type VideoGenerationService struct {
	videoRepo      repository.VideoRepository
	storageService service.StorageService  // Add this
	// ... other fields
}

// Then when saving video:
func (s *videoGenerationService) SaveGeneratedVideo(
	ctx context.Context,
	jobID string,
	videoBytes []byte,
	thumbnailBytes []byte,
) (videoURL, thumbnailURL string, err error) {

	// Upload video
	videoPath, err := s.storageService.UploadVideo(ctx, fmt.Sprintf("%s.mp4", jobID), videoBytes)
	if err != nil {
		log.Printf("Failed to upload video: %v", err)
		return "", "", err
	}

	// Upload thumbnail
	thumbnailPath, err := s.storageService.UploadThumbnail(ctx, fmt.Sprintf("%s.jpg", jobID), thumbnailBytes)
	if err != nil {
		log.Printf("Failed to upload thumbnail: %v", err)
		return "", "", err
	}

	// Get public URLs
	videoURL := s.storageService.GetPublicURL(os.Getenv("STORAGE_BUCKET_VIDEOS"), videoPath)
	thumbnailURL := s.storageService.GetPublicURL(os.Getenv("STORAGE_BUCKET_THUMBNAILS"), thumbnailPath)

	return videoURL, thumbnailURL, nil
}

// For project assets (logo, environment, document):
func (s *projectService) UploadProjectAsset(
	ctx context.Context,
	projectID string,
	assetType string, // "logo", "environment", "document"
	fileName string,
	fileBytes []byte,
) (assetURL string, err error) {

	assetPath, err := s.storageService.UploadAsset(ctx, assetType, fileName, fileBytes)
	if err != nil {
		return "", err
	}

	assetURL := s.storageService.GetPublicURL(os.Getenv("STORAGE_BUCKET_ASSETS"), assetPath)
	return assetURL, nil
}

// In your main.go or handler initialization:
func setupHandlers(router *fiber.App, db *gorm.DB) {
	// Initialize storage service
	storageService := service.NewStorageService()

	// Pass to other services
	videoGenService := service.NewVideoGenerationService(videoRepo, storageService)
	projectService := service.NewProjectService(projectRepo, storageService)

	// Setup handlers
	videoHandler := handler.NewVideoHandler(videoGenService)
	projectHandler := handler.NewProjectHandler(projectService)

	// ... rest of setup
}

*/

// QuickStart: Steps to integrate
//
// 1. Ensure .env has:
//    SUPABASE_URL=https://wkmvrwiesfpnfnaybpbx.supabase.co
//    SUPABASE_ANON_KEY=eyJ...
//    SUPABASE_SERVICE_ROLE_KEY=eyJ...
//    STORAGE_BUCKET_VIDEOS=videos
//    STORAGE_BUCKET_THUMBNAILS=thumbnails
//    STORAGE_BUCKET_ASSETS=assets
//
// 2. In Supabase Dashboard, create 3 buckets:
//    - videos (public)
//    - thumbnails (public)
//    - assets (public)
//
// 3. Import StorageService in your service or handler
//
// 4. Inject storageService into your service struct
//
// 5. Use methods:
//    - UploadVideo(ctx, filename, bytes) -> string (path)
//    - UploadThumbnail(ctx, filename, bytes) -> string (path)
//    - UploadAsset(ctx, type, filename, bytes) -> string (path)
//    - GetPublicURL(bucket, path) -> string (full URL)
//    - DeleteFile(ctx, bucket, path) -> error
//
// Example URLs returned:
// https://wkmvrwiesfpnfnaybpbx.supabase.co/storage/v1/object/public/videos/video_timestamp_uuid.mp4
// https://wkmvrwiesfpnfnaybpbx.supabase.co/storage/v1/object/public/thumbnails/thumb_timestamp_uuid.jpg
// https://wkmvrwiesfpnfnaybpbx.supabase.co/storage/v1/object/public/assets/logos/logo_timestamp_uuid.png

func ExampleStorageUsage(storageService StorageService) {
	ctx := context.Background()

	// Example 1: Upload video
	videoBytes := []byte{} // Your video bytes
	videoPath, err := storageService.UploadVideo(ctx, "sample-video.mp4", videoBytes)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return
	}
	videoURL := storageService.GetPublicURL("videos", videoPath)
	fmt.Printf("Video URL: %s\n", videoURL)

	// Example 2: Upload thumbnail
	thumbBytes := []byte{} // Your thumbnail bytes
	thumbPath, err := storageService.UploadThumbnail(ctx, "sample-thumb.jpg", thumbBytes)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return
	}
	thumbURL := storageService.GetPublicURL("thumbnails", thumbPath)
	fmt.Printf("Thumbnail URL: %s\n", thumbURL)

	// Example 3: Upload project asset
	logoBytes := []byte{} // Your logo bytes
	logoPath, err := storageService.UploadAsset(ctx, "logo", "company-logo.png", logoBytes)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return
	}
	logoURL := storageService.GetPublicURL("assets", logoPath)
	fmt.Printf("Logo URL: %s\n", logoURL)

	// Example 4: Delete file
	err = storageService.DeleteFile(ctx, "videos", videoPath)
	if err != nil {
		log.Printf("Delete failed: %v", err)
	}

	// Example 5: Download file
	downloadedBytes, err := storageService.DownloadFile(ctx, "videos", videoPath)
	if err != nil {
		log.Printf("Download failed: %v", err)
		return
	}
	fmt.Printf("Downloaded %d bytes\n", len(downloadedBytes))
}
