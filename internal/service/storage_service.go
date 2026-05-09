package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StorageService defines file storage operations
type StorageService interface {
	// UploadFile uploads a file to specified bucket
	UploadFile(ctx context.Context, bucketName, filePath string, file []byte) (string, error)

	// UploadVideo uploads video file to videos bucket
	UploadVideo(ctx context.Context, filename string, file []byte) (string, error)

	// UploadThumbnail uploads thumbnail to thumbnails bucket
	UploadThumbnail(ctx context.Context, filename string, file []byte) (string, error)

	// UploadAsset uploads asset (logo, environment, document) to assets bucket
	UploadAsset(ctx context.Context, assetType, filename string, file []byte) (string, error)

	// UploadBase64 parses base64 string and uploads it
	UploadBase64(ctx context.Context, bucketName, filePath, base64Str string) (string, error)

	// DeleteFile deletes a file from bucket
	DeleteFile(ctx context.Context, bucketName, filePath string) error

	// GetPublicURL returns public URL for a file
	GetPublicURL(bucketName, filePath string) string

	// DownloadFile downloads a file from bucket
	DownloadFile(ctx context.Context, bucketName, filePath string) ([]byte, error)
}

type storageService struct {
	supabaseURL      string
	anonKey          string
	serviceRoleKey   string
	bucketVideos     string
	bucketThumbnails string
	bucketAssets     string
	httpClient       *http.Client
}

// NewStorageService creates a new storage service instance
func NewStorageService() StorageService {
	return &storageService{
		supabaseURL:      os.Getenv("SUPABASE_URL"),
		anonKey:          os.Getenv("SUPABASE_ANON_KEY"),
		serviceRoleKey:   os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		bucketVideos:     os.Getenv("STORAGE_BUCKET_VIDEOS"),
		bucketThumbnails: os.Getenv("STORAGE_BUCKET_THUMBNAILS"),
		bucketAssets:     os.Getenv("STORAGE_BUCKET_ASSETS"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadFile uploads a file to the specified bucket using raw bytes
func (s *storageService) UploadFile(ctx context.Context, bucketName, filePath string, file []byte) (string, error) {
	if s.supabaseURL == "" || s.serviceRoleKey == "" {
		return "", errors.New("supabase credentials not configured")
	}

	// Get directory and filename (ensure forward slashes for URL)
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "/" {
		dir = ""
	} else {
		dir = strings.ReplaceAll(dir, "\\", "/") + "/"
	}

	// Generate unique filename to prevent conflicts
	uniqueFilename := fmt.Sprintf("%s%s_%d_%s%s",
		dir,
		strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)),
		time.Now().UnixNano(),
		uuid.New().String()[:8],
		filepath.Ext(filePath),
	)

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.supabaseURL, bucketName, uniqueFilename)

	// Create request with raw bytes
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, bytes.NewReader(file))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.serviceRoleKey))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-upsert", "true") // Allow overwrite if file exists

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errStr := fmt.Sprintf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		fmt.Printf("Storage Error: %s\n", errStr)
		return "", errors.New(errStr)
	}

	return uniqueFilename, nil
}

// UploadVideo uploads a video file
func (s *storageService) UploadVideo(ctx context.Context, filename string, file []byte) (string, error) {
	path := fmt.Sprintf("videos/%s", filename)
	return s.UploadFile(ctx, s.bucketVideos, path, file)
}

// UploadThumbnail uploads a thumbnail file
func (s *storageService) UploadThumbnail(ctx context.Context, filename string, file []byte) (string, error) {
	path := fmt.Sprintf("thumbnails/%s", filename)
	return s.UploadFile(ctx, s.bucketThumbnails, path, file)
}

// UploadAsset uploads an asset file (logo, environment, document)
func (s *storageService) UploadAsset(ctx context.Context, assetType, filename string, file []byte) (string, error) {
	// assetType: "logo", "environment", "document"
	path := fmt.Sprintf("%s/%s", assetType, filename)
	return s.UploadFile(ctx, s.bucketAssets, path, file)
}

// UploadBase64 parses base64 string and uploads it
func (s *storageService) UploadBase64(ctx context.Context, bucketName, filePath, base64Str string) (string, error) {
	if base64Str == "" {
		return "", nil
	}

	// Handle data:image/png;base64,... format
	b64data := base64Str
	if strings.Contains(base64Str, ",") {
		parts := strings.Split(base64Str, ",")
		if len(parts) > 1 {
			b64data = parts[1]
		}
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	return s.UploadFile(ctx, bucketName, filePath, data)
}

// DeleteFile deletes a file from bucket
func (s *storageService) DeleteFile(ctx context.Context, bucketName, filePath string) error {
	if s.supabaseURL == "" || s.serviceRoleKey == "" {
		return errors.New("supabase credentials not configured")
	}

	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.supabaseURL, bucketName, filePath)

	req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.serviceRoleKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetPublicURL returns the public URL for a file
func (s *storageService) GetPublicURL(bucketName, filePath string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		s.supabaseURL, bucketName, filePath)
}

// DownloadFile downloads a file from bucket
func (s *storageService) DownloadFile(ctx context.Context, bucketName, filePath string) ([]byte, error) {
	if s.supabaseURL == "" || s.anonKey == "" {
		return nil, errors.New("supabase credentials not configured")
	}

	downloadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		s.supabaseURL, bucketName, filePath)

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.anonKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
