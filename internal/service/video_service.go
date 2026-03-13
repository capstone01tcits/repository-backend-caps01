package service

import (
	"errors"

	"go-auth/internal/model"
	"go-auth/internal/repository"

	"github.com/google/uuid"
)

type VideoService interface {
	GenerateVideo(userID string, req *model.GenerateVideoRequest) (*model.Video, error)
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

func (s *videoService) GenerateVideo(userID string, req *model.GenerateVideoRequest) (*model.Video, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Verify project ownership
	project, err := s.projectRepo.FindByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	// Verify storyboard ownership
	storyboard, err := s.storyboardRepo.FindByID(req.StoryboardID)
	if err != nil {
		return nil, errors.New("storyboard not found")
	}
	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this storyboard")
	}

	// Check user credits
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Credits <= 0 {
		return nil, errors.New("insufficient credits, please contact admin to add credits")
	}

	// Deduct credit
	user.Credits--
	if err := s.userRepo.UpdateCredits(userID, user.Credits); err != nil {
		return nil, errors.New("failed to deduct credits")
	}

	pid, _ := uuid.Parse(req.ProjectID)
	sid, _ := uuid.Parse(req.StoryboardID)

	title := req.Title
	if title == "" {
		title = project.Name + " - Video"
	}
	format := req.Format
	if format == "" {
		format = "mp4"
	}
	resolution := req.Resolution
	if resolution == "" {
		resolution = "1080p"
	}

	// Calculate total duration from storyboard scenes
	totalDuration := 0
	for _, scene := range storyboard.Scenes {
		totalDuration += scene.Duration
	}

	video := &model.Video{
		ProjectID:    pid,
		UserID:       uid,
		StoryboardID: sid,
		Title:        title,
		Status:       "completed", // Stub: mark as completed immediately
		VideoURL:     "/videos/" + uuid.New().String() + "." + format,
		Duration:     totalDuration,
		Format:       format,
		Resolution:   resolution,
		CreditsUsed:  1,
	}

	if err := s.videoRepo.Create(video); err != nil {
		return nil, errors.New("failed to create video record")
	}

	// Update project status
	inProgress := "in_progress"
	s.projectRepo.Update(&model.Project{
		ID:     project.ID,
		UserID: project.UserID,
		Name:   project.Name,
		Status: inProgress,
	})

	return video, nil
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
