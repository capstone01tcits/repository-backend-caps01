package service

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type ContentService interface {
	// Content Pillar
	GenerateContentPillars(userID, projectID string) ([]model.ContentPillar, error)
	GetContentPillars(userID, projectID string) ([]model.ContentPillar, error)
	GetContentPillar(userID, pillarID string) (*model.ContentPillar, error)
	SelectContentPillar(userID, pillarID string) (*model.ContentPillar, error)
	UpdateContentPillar(userID, pillarID string, req *model.UpdateContentPillarRequest) (*model.ContentPillar, error)

	// Content Pillar Adjustment (new workflow)
	AdjustContentPillar(userID, projectID string) (*model.ContentPillarAdjustmentResponse, error)
	SelectContentPillarAndGenerateCreativeBrief(userID, pillarID string) (*model.CreativeBrief, error)

	// Content Theme
	GetContentThemes(userID, pillarID string) ([]model.ContentTheme, error)
	SelectContentTheme(userID, themeID string) (*model.ContentTheme, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	projectRepo repository.ProjectRepository
	briefRepo   repository.BriefRepository
}

func NewContentService(contentRepo repository.ContentRepository, projectRepo repository.ProjectRepository, briefRepo repository.BriefRepository) ContentService {
	return &contentService{contentRepo, projectRepo, briefRepo}
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

// UpdateContentPillar godoc
// Updates prompt and/or video_url for a content pillar (Sprint 3)
func (s *contentService) UpdateContentPillar(userID, pillarID string, req *model.UpdateContentPillarRequest) (*model.ContentPillar, error) {
	pillar, err := s.contentRepo.FindContentPillarByID(pillarID)
	if err != nil {
		return nil, errors.New("content pillar not found")
	}

	if pillar.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	if req.Prompt != nil {
		pillar.Prompt = *req.Prompt
	}
	if req.VideoURL != nil {
		pillar.VideoURL = *req.VideoURL
	}

	if err := s.contentRepo.UpdateContentPillar(pillar); err != nil {
		return nil, errors.New("failed to update content pillar")
	}

	return pillar, nil
}

// ==================== Content Pillar Adjustment (New Workflow) ====================

// AdjustContentPillar godoc
// GET /api/projects/:id/content-pillars/adjustment
// Retrieves business brief, content brief and generates/returns content pillars for selection
func (s *contentService) AdjustContentPillar(userID, projectID string) (*model.ContentPillarAdjustmentResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Verify project exists and belongs to user
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	pid, _ := uuid.Parse(projectID)

	// Get existing content pillars for this project
	pillars, err := s.contentRepo.FindContentPillarsByProjectID(projectID)
	if err != nil {
		pillars = []model.ContentPillar{}
	}

	// If no pillars exist, generate them
	if len(pillars) == 0 {
		pillars = []model.ContentPillar{
			{UserID: uid, ProjectID: pid, Title: "Brand Awareness", Description: "Content focused on increasing brand visibility and recognition among target audience", IsSelected: false},
			{UserID: uid, ProjectID: pid, Title: "Product Education", Description: "Content that educates the audience about product features and benefits", IsSelected: false},
			{UserID: uid, ProjectID: pid, Title: "Social Proof", Description: "Content showcasing testimonials, case studies, and success stories", IsSelected: false},
		}

		for i := range pillars {
			if err := s.contentRepo.CreateContentPillar(&pillars[i]); err != nil {
				return nil, errors.New("failed to create content pillar")
			}
		}
	}

	// Load themes for each pillar
	for i := range pillars {
		themes, err := s.contentRepo.FindContentThemesByPillarID(pillars[i].ID.String())
		if err == nil {
			pillars[i].ContentThemes = themes
		}
	}

	response := &model.ContentPillarAdjustmentResponse{
		ContentPillars: pillars,
		BusinessBrief:  nil,
		CreativeBrief:  nil,
	}

	return response, nil
}

// SelectContentPillarAndGenerateCreativeBrief godoc
// POST /api/projects/:id/content-pillars/adjustment/select
// Selects a content pillar and generates a creative brief
func (s *contentService) SelectContentPillarAndGenerateCreativeBrief(userID, pillarID string) (*model.CreativeBrief, error) {
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

	// Generate creative brief based on the selected content pillar
	uid, _ := uuid.Parse(userID)

	// Find the business brief for this project
	businessBrief, err := s.briefRepo.FindBusinessBriefByProjectID(pillar.ProjectID.String())
	if err != nil {
		return nil, errors.New("business brief not found for this project")
	}

	// Generate creative brief from content pillar and business brief
	creativeBrief := &model.CreativeBrief{
		UserID:           uid,
		BusinessBriefID:  businessBrief.ID,
		Title:            pillar.Title + " - " + businessBrief.ProjectName,
		VideoType:        "promotional", // default, can be adjusted
		Duration:         60,            // default, can be adjusted
		Style:            "cinematic",   // default, can be adjusted
		Tone:             businessBrief.KeyMessage,
		Script:           pillar.Description, // use pillar description as initial script
		Storyboard:       "",
		VisualReferences: pillar.Title,
		MusicPreference:  "professional",
		CallToAction:     businessBrief.KeyMessage,
		OutputFormat:     "mp4",
		Resolution:       "1080p",
		AdditionalNotes:  "Generated from content pillar: " + pillar.Title,
		Status:           "draft",
	}

	// Create the creative brief
	if err := s.briefRepo.CreateCreativeBrief(creativeBrief); err != nil {
		return nil, errors.New("failed to create creative brief")
	}

	return creativeBrief, nil
}
