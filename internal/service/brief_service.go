package service

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type BriefService interface {
	// Business Brief
	CreateBusinessBrief(userID string, req *model.CreateBusinessBriefRequest) (*model.BusinessBrief, error)
	GetBusinessBrief(userID, briefID string) (*model.BusinessBrief, error)
	GetBusinessBriefs(userID string) ([]model.BusinessBrief, error)
	UpdateBusinessBrief(userID, briefID string, req *model.UpdateBusinessBriefRequest) (*model.BusinessBrief, error)
	DeleteBusinessBrief(userID, briefID string) error

	// Creative Brief
	CreateCreativeBrief(userID string, req *model.CreateCreativeBriefRequest) (*model.CreativeBrief, error)
	GetCreativeBrief(userID, briefID string) (*model.CreativeBrief, error)
	GetCreativeBriefs(userID string) ([]model.CreativeBrief, error)
	GetCreativeBriefsByBusinessBrief(userID, businessBriefID string) ([]model.CreativeBrief, error)
	UpdateCreativeBrief(userID, briefID string, req *model.UpdateCreativeBriefRequest) (*model.CreativeBrief, error)
	DeleteCreativeBrief(userID, briefID string) error
}

type briefService struct {
	briefRepo repository.BriefRepository
}

func NewBriefService(briefRepo repository.BriefRepository) BriefService {
	return &briefService{briefRepo}
}

// ==================== Business Brief ====================

func (s *briefService) CreateBusinessBrief(userID string, req *model.CreateBusinessBriefRequest) (*model.BusinessBrief, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	brief := &model.BusinessBrief{
		UserID:           uid,
		ProjectID:        projectID,
		ProjectName:      req.ProjectName,
		CompanyName:      req.CompanyName,
		Industry:         req.Industry,
		TargetAudience:   req.TargetAudience,
		ProjectObjective: req.ProjectObjective,
		KeyMessage:       req.KeyMessage,
		Budget:           req.Budget,
		Timeline:         req.Timeline,
		Competitors:      req.Competitors,
		AdditionalNotes:  req.AdditionalNotes,
		Status:           "draft",
	}

	if err := s.briefRepo.CreateBusinessBrief(brief); err != nil {
		return nil, errors.New("failed to create business brief")
	}

	return brief, nil
}

func (s *briefService) GetBusinessBrief(userID, briefID string) (*model.BusinessBrief, error) {
	brief, err := s.briefRepo.FindBusinessBriefByID(briefID)
	if err != nil {
		return nil, errors.New("business brief not found")
	}

	if brief.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this brief")
	}

	return brief, nil
}

func (s *briefService) GetBusinessBriefs(userID string) ([]model.BusinessBrief, error) {
	return s.briefRepo.FindBusinessBriefsByUserID(userID)
}

func (s *briefService) UpdateBusinessBrief(userID, briefID string, req *model.UpdateBusinessBriefRequest) (*model.BusinessBrief, error) {
	brief, err := s.briefRepo.FindBusinessBriefByID(briefID)
	if err != nil {
		return nil, errors.New("business brief not found")
	}

	if brief.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this brief")
	}

	if req.ProjectName != nil {
		brief.ProjectName = *req.ProjectName
	}
	if req.CompanyName != nil {
		brief.CompanyName = *req.CompanyName
	}
	if req.Industry != nil {
		brief.Industry = *req.Industry
	}
	if req.TargetAudience != nil {
		brief.TargetAudience = *req.TargetAudience
	}
	if req.ProjectObjective != nil {
		brief.ProjectObjective = *req.ProjectObjective
	}
	if req.KeyMessage != nil {
		brief.KeyMessage = *req.KeyMessage
	}
	if req.Budget != nil {
		brief.Budget = *req.Budget
	}
	if req.Timeline != nil {
		brief.Timeline = *req.Timeline
	}
	if req.Competitors != nil {
		brief.Competitors = *req.Competitors
	}
	if req.AdditionalNotes != nil {
		brief.AdditionalNotes = *req.AdditionalNotes
	}
	if req.Status != nil {
		brief.Status = *req.Status
	}

	if err := s.briefRepo.UpdateBusinessBrief(brief); err != nil {
		return nil, errors.New("failed to update business brief")
	}

	return brief, nil
}

func (s *briefService) DeleteBusinessBrief(userID, briefID string) error {
	brief, err := s.briefRepo.FindBusinessBriefByID(briefID)
	if err != nil {
		return errors.New("business brief not found")
	}

	if brief.UserID.String() != userID {
		return errors.New("unauthorized access to this brief")
	}

	return s.briefRepo.DeleteBusinessBrief(briefID)
}

// ==================== Creative Brief ====================

func (s *briefService) CreateCreativeBrief(userID string, req *model.CreateCreativeBriefRequest) (*model.CreativeBrief, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	bbID, err := uuid.Parse(req.BusinessBriefID)
	if err != nil {
		return nil, errors.New("invalid business brief ID")
	}

	// Verify the business brief exists and belongs to the user
	bb, err := s.briefRepo.FindBusinessBriefByID(req.BusinessBriefID)
	if err != nil {
		return nil, errors.New("business brief not found")
	}
	if bb.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this business brief")
	}

	brief := &model.CreativeBrief{
		UserID:           uid,
		BusinessBriefID:  bbID,
		Title:            req.Title,
		VideoType:        req.VideoType,
		Duration:         req.Duration,
		Style:            req.Style,
		Tone:             req.Tone,
		Script:           req.Script,
		Storyboard:       req.Storyboard,
		VisualReferences: req.VisualReferences,
		MusicPreference:  req.MusicPreference,
		CallToAction:     req.CallToAction,
		OutputFormat:     req.OutputFormat,
		Resolution:       req.Resolution,
		AdditionalNotes:  req.AdditionalNotes,
		Status:           "draft",
	}

	if err := s.briefRepo.CreateCreativeBrief(brief); err != nil {
		return nil, errors.New("failed to create creative brief")
	}

	return brief, nil
}

func (s *briefService) GetCreativeBrief(userID, briefID string) (*model.CreativeBrief, error) {
	brief, err := s.briefRepo.FindCreativeBriefByID(briefID)
	if err != nil {
		return nil, errors.New("creative brief not found")
	}

	if brief.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this brief")
	}

	return brief, nil
}

func (s *briefService) GetCreativeBriefs(userID string) ([]model.CreativeBrief, error) {
	return s.briefRepo.FindCreativeBriefsByUserID(userID)
}

func (s *briefService) GetCreativeBriefsByBusinessBrief(userID, businessBriefID string) ([]model.CreativeBrief, error) {
	// Verify user owns the business brief
	bb, err := s.briefRepo.FindBusinessBriefByID(businessBriefID)
	if err != nil {
		return nil, errors.New("business brief not found")
	}
	if bb.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this business brief")
	}

	return s.briefRepo.FindCreativeBriefsByBusinessBriefID(businessBriefID)
}

func (s *briefService) UpdateCreativeBrief(userID, briefID string, req *model.UpdateCreativeBriefRequest) (*model.CreativeBrief, error) {
	brief, err := s.briefRepo.FindCreativeBriefByID(briefID)
	if err != nil {
		return nil, errors.New("creative brief not found")
	}

	if brief.UserID.String() != userID {
		return nil, errors.New("unauthorized access to this brief")
	}

	if req.Title != nil {
		brief.Title = *req.Title
	}
	if req.VideoType != nil {
		brief.VideoType = *req.VideoType
	}
	if req.Duration != nil {
		brief.Duration = *req.Duration
	}
	if req.Style != nil {
		brief.Style = *req.Style
	}
	if req.Tone != nil {
		brief.Tone = *req.Tone
	}
	if req.Script != nil {
		brief.Script = *req.Script
	}
	if req.Storyboard != nil {
		brief.Storyboard = *req.Storyboard
	}
	if req.VisualReferences != nil {
		brief.VisualReferences = *req.VisualReferences
	}
	if req.MusicPreference != nil {
		brief.MusicPreference = *req.MusicPreference
	}
	if req.CallToAction != nil {
		brief.CallToAction = *req.CallToAction
	}
	if req.OutputFormat != nil {
		brief.OutputFormat = *req.OutputFormat
	}
	if req.Resolution != nil {
		brief.Resolution = *req.Resolution
	}
	if req.AdditionalNotes != nil {
		brief.AdditionalNotes = *req.AdditionalNotes
	}
	if req.Status != nil {
		brief.Status = *req.Status
	}

	if err := s.briefRepo.UpdateCreativeBrief(brief); err != nil {
		return nil, errors.New("failed to update creative brief")
	}

	return brief, nil
}

func (s *briefService) DeleteCreativeBrief(userID, briefID string) error {
	brief, err := s.briefRepo.FindCreativeBriefByID(briefID)
	if err != nil {
		return errors.New("creative brief not found")
	}

	if brief.UserID.String() != userID {
		return errors.New("unauthorized access to this brief")
	}

	return s.briefRepo.DeleteCreativeBrief(briefID)
}
