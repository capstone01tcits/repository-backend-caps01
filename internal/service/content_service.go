package service

import (
	"errors"

	"go-auth/internal/model"
	"go-auth/internal/repository"

	"github.com/google/uuid"
)

type ContentService interface {
	// Content Pillar
	GenerateContentPillars(userID, projectID string) ([]model.ContentPillar, error)
	GetContentPillars(userID, projectID string) ([]model.ContentPillar, error)
	GetContentPillar(userID, pillarID string) (*model.ContentPillar, error)
	SelectContentPillar(userID, pillarID string) (*model.ContentPillar, error)

	// Content Theme
	GetContentThemes(userID, pillarID string) ([]model.ContentTheme, error)
	SelectContentTheme(userID, themeID string) (*model.ContentTheme, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
}

func NewContentService(contentRepo repository.ContentRepository, projectRepo repository.ProjectRepository, userRepo repository.UserRepository) ContentService {
	return &contentService{contentRepo, projectRepo, userRepo}
}

// ==================== Content Pillar ====================

func (s *contentService) GenerateContentPillars(userID, projectID string) ([]model.ContentPillar, error) {
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

	// Check user credits before generating
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Credits < 1 {
		return nil, errors.New("insufficient credits: need 1 credit to generate content pillars")
	}

	// Deduct 1 credit for pillar generation
	if err := s.userRepo.UpdateCredits(userID, user.Credits-1); err != nil {
		return nil, errors.New("failed to deduct credits")
	}

	pid, _ := uuid.Parse(projectID)

	// Stub: generate 3 sample content pillars
	// In production, this would call the AI service
	pillars := []model.ContentPillar{
		{UserID: uid, ProjectID: pid, Title: "Brand Awareness", Description: "Content focused on increasing brand visibility and recognition among target audience"},
		{UserID: uid, ProjectID: pid, Title: "Product Education", Description: "Content that educates the audience about product features and benefits"},
		{UserID: uid, ProjectID: pid, Title: "Social Proof", Description: "Content showcasing testimonials, case studies, and success stories"},
	}

	var created []model.ContentPillar
	for i := range pillars {
		if err := s.contentRepo.CreateContentPillar(&pillars[i]); err != nil {
			return nil, errors.New("failed to create content pillar")
		}
		// Generate themes for each pillar
		themes := []model.ContentTheme{
			{UserID: uid, ContentPillarID: pillars[i].ID, Title: pillars[i].Title + " - Theme A", Description: "First theme variation for " + pillars[i].Title},
			{UserID: uid, ContentPillarID: pillars[i].ID, Title: pillars[i].Title + " - Theme B", Description: "Second theme variation for " + pillars[i].Title},
		}
		for j := range themes {
			if err := s.contentRepo.CreateContentTheme(&themes[j]); err != nil {
				return nil, errors.New("failed to create content theme")
			}
		}
		pillars[i].ContentThemes = themes
		created = append(created, pillars[i])
	}

	return created, nil
}

func (s *contentService) GetContentPillars(userID, projectID string) ([]model.ContentPillar, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	return s.contentRepo.FindContentPillarsByProjectID(projectID)
}

func (s *contentService) GetContentPillar(userID, pillarID string) (*model.ContentPillar, error) {
	pillar, err := s.contentRepo.FindContentPillarByID(pillarID)
	if err != nil {
		return nil, errors.New("content pillar not found")
	}

	if pillar.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return pillar, nil
}

func (s *contentService) SelectContentPillar(userID, pillarID string) (*model.ContentPillar, error) {
	pillar, err := s.contentRepo.FindContentPillarByID(pillarID)
	if err != nil {
		return nil, errors.New("content pillar not found")
	}

	if pillar.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	// Deselect all pillars in this project, then select this one
	if err := s.contentRepo.DeselectAllPillarsByProjectID(pillar.ProjectID.String()); err != nil {
		return nil, errors.New("failed to update selection")
	}

	pillar.IsSelected = true
	if err := s.contentRepo.UpdateContentPillar(pillar); err != nil {
		return nil, errors.New("failed to select content pillar")
	}

	return pillar, nil
}

// ==================== Content Theme ====================

func (s *contentService) GetContentThemes(userID, pillarID string) ([]model.ContentTheme, error) {
	pillar, err := s.contentRepo.FindContentPillarByID(pillarID)
	if err != nil {
		return nil, errors.New("content pillar not found")
	}

	if pillar.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return s.contentRepo.FindContentThemesByPillarID(pillarID)
}

func (s *contentService) SelectContentTheme(userID, themeID string) (*model.ContentTheme, error) {
	theme, err := s.contentRepo.FindContentThemeByID(themeID)
	if err != nil {
		return nil, errors.New("content theme not found")
	}

	if theme.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	// Deselect all themes under this pillar, then select this one
	if err := s.contentRepo.DeselectAllThemesByPillarID(theme.ContentPillarID.String()); err != nil {
		return nil, errors.New("failed to update selection")
	}

	theme.IsSelected = true
	if err := s.contentRepo.UpdateContentTheme(theme); err != nil {
		return nil, errors.New("failed to select content theme")
	}

	return theme, nil
}
