package service

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type StoryboardService interface {
	GenerateStoryboards(userID, projectID, contentThemeID string) ([]model.Storyboard, error)
	GetStoryboards(userID, projectID string) ([]model.Storyboard, error)
	GetStoryboard(userID, storyboardID string) (*model.Storyboard, error)
	SelectStoryboard(userID, storyboardID string) (*model.Storyboard, error)
	GetScenes(userID, storyboardID string) ([]model.Scene, error)
}

type storyboardService struct {
	storyboardRepo repository.StoryboardRepository
	projectRepo    repository.ProjectRepository
	contentRepo    repository.ContentRepository
}

func NewStoryboardService(
	storyboardRepo repository.StoryboardRepository,
	projectRepo repository.ProjectRepository,
	contentRepo repository.ContentRepository,
) StoryboardService {
	return &storyboardService{storyboardRepo, projectRepo, contentRepo}
}

func (s *storyboardService) GenerateStoryboards(userID, projectID, contentThemeID string) ([]model.Storyboard, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	// Verify content theme exists and belongs to user
	theme, err := s.contentRepo.FindContentThemeByID(contentThemeID)
	if err != nil {
		return nil, errors.New("content theme not found")
	}
	if theme.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this content theme")
	}

	pid, _ := uuid.Parse(projectID)

	// Stub: generate 2 storyboard variations with scenes
	// In production, this would call the AI service
	storyboardData := []struct {
		title  string
		desc   string
		scenes []struct {
			num  int
			t    string
			d    string
			v    string
			dur  int
		}
	}{
		{
			title: "Storyboard A - Dynamic",
			desc:  "A dynamic, fast-paced storyboard based on theme: " + theme.Title,
			scenes: []struct {
				num int
				t   string
				d   string
				v   string
				dur int
			}{
				{1, "Opening Hook", "Attention-grabbing opening sequence", "Wide shot with dramatic lighting and motion graphics", 5},
				{2, "Problem Statement", "Present the challenge or pain point", "Close-up shots with text overlay highlighting the problem", 8},
				{3, "Solution Reveal", "Introduce the product/service as solution", "Product showcase with smooth transitions and feature highlights", 10},
				{4, "Call to Action", "End with clear CTA", "Logo animation with contact info and call-to-action text", 7},
			},
		},
		{
			title: "Storyboard B - Narrative",
			desc:  "A storytelling-driven storyboard based on theme: " + theme.Title,
			scenes: []struct {
				num int
				t   string
				d   string
				v   string
				dur int
			}{
				{1, "Setting the Scene", "Establish context and atmosphere", "Cinematic wide shot establishing the environment", 6},
				{2, "Character Introduction", "Introduce relatable character/persona", "Medium shot of character in their environment", 8},
				{3, "Journey & Transformation", "Show the transformation journey", "Montage sequence showing before and after", 12},
				{4, "Resolution", "Happy ending with brand integration", "Warm lighting, character smiling, brand logo reveal", 4},
			},
		},
	}

	var created []model.Storyboard
	for _, sb := range storyboardData {
		storyboard := &model.Storyboard{
			ProjectID:   pid,
			UserID:      uid,
			Title:       sb.title,
			Description: sb.desc,
		}

		if err := s.storyboardRepo.Create(storyboard); err != nil {
			return nil, errors.New("failed to create storyboard")
		}

		var scenes []model.Scene
		for _, sc := range sb.scenes {
			scene := &model.Scene{
				StoryboardID: storyboard.ID,
				UserID:       uid,
				SceneNumber:  sc.num,
				Title:        sc.t,
				Description:  sc.d,
				VisualDesc:   sc.v,
				Duration:     sc.dur,
			}
			if err := s.storyboardRepo.CreateScene(scene); err != nil {
				return nil, errors.New("failed to create scene")
			}
			scenes = append(scenes, *scene)
		}

		storyboard.Scenes = scenes
		created = append(created, *storyboard)
	}

	return created, nil
}

func (s *storyboardService) GetStoryboards(userID, projectID string) ([]model.Storyboard, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return s.storyboardRepo.FindByProjectID(projectID)
}

func (s *storyboardService) GetStoryboard(userID, storyboardID string) (*model.Storyboard, error) {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return nil, errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return storyboard, nil
}

func (s *storyboardService) SelectStoryboard(userID, storyboardID string) (*model.Storyboard, error) {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return nil, errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	// Deselect all storyboards in this project, then select this one
	if err := s.storyboardRepo.DeselectAllByProjectID(storyboard.ProjectID.String()); err != nil {
		return nil, errors.New("failed to update selection")
	}

	storyboard.IsSelected = true
	if err := s.storyboardRepo.Update(storyboard); err != nil {
		return nil, errors.New("failed to select storyboard")
	}

	return storyboard, nil
}

func (s *storyboardService) GetScenes(userID, storyboardID string) ([]model.Scene, error) {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return nil, errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return s.storyboardRepo.FindScenesByStoryboardID(storyboardID)
}
