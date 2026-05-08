package service

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type StoryboardService interface {
	GenerateTemplates(userID, projectID string, videoDuration int) ([]model.StoryboardTemplate, error)
	CreateManualStoryboard(userID string, req *model.CreateManualStoryboardRequest) (*model.Storyboard, error)
	GetStoryboards(userID, projectID string) ([]model.Storyboard, error)
	GetStoryboard(userID, storyboardID string) (*model.Storyboard, error)
	SelectStoryboard(userID, storyboardID string) (*model.Storyboard, error)
	UpdateStoryboard(userID, storyboardID string, req *model.UpdateManualStoryboardRequest) (*model.Storyboard, error)
	GetStoryboardSections(userID, storyboardID string) ([]model.StoryboardSection, error)
	DeleteStoryboard(userID, storyboardID string) error
	RestoreStoryboard(userID, storyboardID string) error
}

type storyboardService struct {
	storyboardRepo repository.StoryboardRepository
	projectRepo    repository.ProjectRepository
	briefRepo      repository.BriefRepository
	templateGen    *TemplateGenerator
}

func NewStoryboardService(
	storyboardRepo repository.StoryboardRepository,
	projectRepo repository.ProjectRepository,
	briefRepo repository.BriefRepository,
) StoryboardService {
	templateGen := NewTemplateGenerator(projectRepo, briefRepo)
	return &storyboardService{storyboardRepo, projectRepo, briefRepo, templateGen}
}

// GenerateTemplates creates multiple storyboard template options for user selection
func (s *storyboardService) GenerateTemplates(userID, projectID string, videoDuration int) ([]model.StoryboardTemplate, error) {
	// Verify project exists and user has access
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	// Validate video duration
	if videoDuration < 15 || videoDuration > 300 {
		return nil, errors.New("video duration must be between 15 and 300 seconds")
	}

	// Generate templates
	templates, err := s.templateGen.GenerateTemplates(projectID, videoDuration)
	if err != nil {
		return nil, errors.New("failed to generate templates")
	}

	return templates, nil
}

func (s *storyboardService) CreateManualStoryboard(userID string, req *model.CreateManualStoryboardRequest) (*model.Storyboard, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	project, err := s.projectRepo.FindByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this project")
	}

	pid, _ := uuid.Parse(req.ProjectID)

	// Create the storyboard
	storyboard := &model.Storyboard{
		ProjectID:     pid,
		UserID:        uid,
		Title:         req.Title,
		Description:   req.Description,
		Style:         req.Style,
		TotalDuration: req.Duration,
	}

	if err := s.storyboardRepo.Create(storyboard); err != nil {
		return nil, errors.New("failed to create storyboard")
	}

	// Create the 3 sections: hook, value, cta
	var sections []model.StoryboardSection
	for _, sectionInput := range req.Sections {
		section := &model.StoryboardSection{
			StoryboardID: storyboard.ID,
			UserID:       uid,
			SectionType:  sectionInput.SectionType,
			Content:      sectionInput.Content,
			Duration:     sectionInput.Duration,
		}
		if err := s.storyboardRepo.CreateSection(section); err != nil {
			return nil, errors.New("failed to create storyboard section")
		}
		sections = append(sections, *section)
	}

	storyboard.Sections = sections
	return storyboard, nil
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

// UpdateStoryboard godoc
// Updates storyboard and its sections (manual storyboard)
func (s *storyboardService) UpdateStoryboard(userID, storyboardID string, req *model.UpdateManualStoryboardRequest) (*model.Storyboard, error) {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return nil, errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	// Update basic storyboard info
	if req.Title != nil {
		storyboard.Title = *req.Title
	}
	if req.Description != nil {
		storyboard.Description = *req.Description
	}

	if err := s.storyboardRepo.Update(storyboard); err != nil {
		return nil, errors.New("failed to update storyboard")
	}

	// Update sections if provided
	if len(req.Sections) > 0 {
		// Get existing sections
		sections, err := s.storyboardRepo.FindSectionsByStoryboardID(storyboardID)
		if err != nil {
			return nil, errors.New("failed to retrieve sections")
		}

		// Update each section
		for i, sectionInput := range req.Sections {
			if i < len(sections) {
				sections[i].SectionType = sectionInput.SectionType
				sections[i].Content = sectionInput.Content
				if err := s.storyboardRepo.UpdateSection(&sections[i]); err != nil {
					return nil, errors.New("failed to update section")
				}
			}
		}

		storyboard.Sections = sections
	}

	return storyboard, nil
}

func (s *storyboardService) GetStoryboardSections(userID, storyboardID string) ([]model.StoryboardSection, error) {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return nil, errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	return s.storyboardRepo.FindSectionsByStoryboardID(storyboardID)
}

func (s *storyboardService) DeleteStoryboard(userID, storyboardID string) error {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return errors.New("unauthorized access")
	}

	return s.storyboardRepo.Delete(storyboardID)
}

func (s *storyboardService) RestoreStoryboard(userID, storyboardID string) error {
	storyboard, err := s.storyboardRepo.UnscopedFindByID(storyboardID)
	if err != nil {
		return errors.New("storyboard not found")
	}

	if storyboard.UserID.String() != userID {
		return errors.New("unauthorized access")
	}

	return s.storyboardRepo.Restore(storyboardID)
}

