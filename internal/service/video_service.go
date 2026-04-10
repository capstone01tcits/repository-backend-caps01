package service

// DEPRECATED: This service is not currently used in the application.
// All video generation functionality is handled by VideoGenerationService.
// This file is retained for potential future use or can be safely removed.
// Consider using VideoGenerationService instead for any video-related operations.

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"
)

type VideoService interface {
	GetVideo(userID, videoID string) (*model.Video, error)
	GetVideosByProject(userID, projectID string) ([]model.Video, error)
	GetVideosByUser(userID string) ([]model.Video, error)
}

type videoService struct {
	videoRepo      repository.VideoRepository
	storyboardRepo repository.StoryboardRepository
	projectRepo    repository.ProjectRepository
	userRepo       repository.UserRepository
}

func NewVideoService(
	videoRepo repository.VideoRepository,
	storyboardRepo repository.StoryboardRepository,
	projectRepo repository.ProjectRepository,
	userRepo repository.UserRepository,
) VideoService {
	return &videoService{videoRepo, storyboardRepo, projectRepo, userRepo}
}

func (s *videoService) GetVideo(userID, videoID string) (*model.Video, error) {
	video, err := s.videoRepo.FindByID(videoID)
	if err != nil {
		return nil, errors.New("video not found")
	}

	if video.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return video, nil
}

func (s *videoService) GetVideosByProject(userID, projectID string) ([]model.Video, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return s.videoRepo.FindByProjectID(projectID)
}

func (s *videoService) GetVideosByUser(userID string) ([]model.Video, error) {
	return s.videoRepo.FindByUserID(userID)
}
