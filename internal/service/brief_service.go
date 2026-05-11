package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"
	"Sevima-AI-Content-Creator/pkg/utils"

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

	// Unified FE Flow (matches frontend exactly)
	CreateProjectFromFE(userID string, req *model.CreateProjectFromFERequest) (map[string]interface{}, error)
}

type briefService struct {
	briefRepo      repository.BriefRepository
	projectRepo    repository.ProjectRepository
	storyboardSvc  StoryboardService
	storageSvc     StorageService
}

func NewBriefService(briefRepo repository.BriefRepository, projectRepo repository.ProjectRepository, storyboardSvc StoryboardService, storageSvc StorageService) BriefService {
	return &briefService{briefRepo, projectRepo, storyboardSvc, storageSvc}
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
		InstituteName:    req.InstituteName,
		Education:          req.Education,
		OfferedDegrees:     req.OfferedDegrees,
		InstitutionHistory: req.InstitutionHistory,
		Industry:           req.Industry,
		TargetAudience:   req.TargetAudience,
		ProjectObjective: req.ProjectObjective,
		KeyMessage:       req.KeyMessage,
		Budget:           req.Budget,
		Timeline:         req.Timeline,
		Deadline:         req.Deadline,
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
	if req.InstituteName != nil {
		brief.InstituteName = *req.InstituteName
	}
	if req.Education != nil {
		brief.Education = *req.Education
	}
	if req.OfferedDegrees != nil {
		brief.OfferedDegrees = *req.OfferedDegrees
	}
	if req.InstitutionHistory != nil {
		brief.InstitutionHistory = *req.InstitutionHistory
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
	if req.Deadline != nil {
		brief.Deadline = *req.Deadline
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

// ==================== Unified FE Flow (Matches Frontend Exactly) ====================

func (s *briefService) CreateProjectFromFE(userID string, req *model.CreateProjectFromFERequest) (map[string]interface{}, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Step 1: Create Project
	projectName := req.EventContent + " for " + req.InstitutionName
	description := req.InstitutionHistory
	if description == "" {
		description = "Video production project for " + req.InstitutionName
	}

	project := &model.Project{
		UserID:      uid,
		Name:        projectName,
		Description: description,
		Theme:       req.SelectedTheme,
		Status:      "draft",
	}
	if err := s.projectRepo.Create(project); err != nil {
		return nil, errors.New("failed to create project")
	}

	// Step 2: Create Business Brief (auto-fill missing fields)
	schoolLevel := req.SchoolLevel
	if schoolLevel == "" {
		schoolLevel = "Perguruan Tinggi" // default
	}

	// Handle Image Uploads to Bucket
	ctx := context.Background()
	bucketAssets := "assets" // Default bucket name

	logoURL := ""
	if req.LogoBase64 != "" {
		filename := fmt.Sprintf("logo_%s.png", project.ID.String())
		path, err := s.storageSvc.UploadBase64(ctx, bucketAssets, "logos/"+filename, req.LogoBase64)
		if err == nil {
			logoURL = s.storageSvc.GetPublicURL(bucketAssets, path)
		} else {
			fmt.Printf("Logo Upload Error: %v\n", err)
		}
	}

	envURL := ""
	if req.EnvBase64 != "" {
		filename := fmt.Sprintf("env_%s.jpg", project.ID.String())
		path, err := s.storageSvc.UploadBase64(ctx, bucketAssets, "environments/"+filename, req.EnvBase64)
		if err == nil {
			envURL = s.storageSvc.GetPublicURL(bucketAssets, path)
		} else {
			fmt.Printf("Env Upload Error: %v\n", err)
		}
	}

	docURL := ""
	if req.DocumentBase64 != "" {
		filename := fmt.Sprintf("doc_%s.pdf", project.ID.String())
		path, err := s.storageSvc.UploadBase64(ctx, bucketAssets, "documents/"+filename, req.DocumentBase64)
		if err == nil {
			docURL = s.storageSvc.GetPublicURL(bucketAssets, path)
		} else {
			fmt.Printf("Doc Upload Error: %v\n", err)
		}
	}

	businessBrief := &model.BusinessBrief{
		ID:               uuid.New(),
		UserID:           uid,
		ProjectID:        project.ID,
		ProjectName:      projectName,
		CompanyName:      req.InstitutionName,
		InstituteName:    req.InstitutionName,
		SchoolLevel:      schoolLevel,
		Education:        schoolLevel,
		Industry:           "Education", // auto-fill default
		TargetAudience:     "Students",  // auto-fill default
		ProjectObjective:   description,
		OfferedDegrees:     req.OfferedDegrees,
		InstitutionHistory: req.InstitutionHistory,
		KeyMessage:         req.SelectedKeyMessage,
		Budget:             "",      // optional
		Timeline:           "",      // optional
		Competitors:        "",      // optional
		AdditionalNotes:    "",      // previously req.OfferedDegrees was here
		LogoPath:           logoURL, // Store bucket URL instead of base64
		EnvironmentPath:  envURL,             // Store bucket URL instead of base64
		DocumentPath:     docURL,             // Store bucket URL instead of base64
		Status:           "draft",
	}
	if err := s.briefRepo.CreateBusinessBrief(businessBrief); err != nil {
		return nil, errors.New("failed to create business brief")
	}

	// Step 3: Create Creative Brief (auto-fill missing fields)
	duration := parseDurationToInt(req.VideoDuration)
	if duration == 0 {
		duration = 30 // default to 30 seconds if not provided
	}

	creativeBrief := &model.CreativeBrief{
		ID:               uuid.New(),
		UserID:           uid,
		BusinessBriefID:  businessBrief.ID,
		Title:            projectName,
		VideoType:        utils.MapEventContentToVideoType(req.EventContent),
		Duration:         duration,
		Style:            utils.MapThemeToStyle(req.SelectedTheme),
		Tone:             req.ToneOfVoice,
		Script:           req.Prompt, // use custom prompt as script
		VisualReferences: req.SelectedTheme,
		MusicPreference:  utils.MapToneToMusicPreference(req.ToneOfVoice),
		CallToAction:     req.SelectedKeyMessage,
		Copywriting:      req.EditableCopywriting, // social media caption
		Hashtags:         req.EditableHashtags,    // social media hashtags
		OutputFormat:     "mp4",                   // auto-fill default
		Resolution:       "1080p",                 // auto-fill default
		AdditionalNotes:  req.Prompt,
		Status:           "draft",
	}
	if err := s.briefRepo.CreateCreativeBrief(creativeBrief); err != nil {
		return nil, errors.New("failed to create creative brief")
	}

	// Step 4: Auto-Generate Storyboard
	storyboard, err := s.storyboardSvc.AutoGenerateStoryboard(userID, project.ID, businessBrief, creativeBrief)
	if err != nil {
		return nil, err
	}

	// Return response with IDs for frontend
	return map[string]interface{}{
		"project_id":        project.ID.String(),
		"business_brief_id": businessBrief.ID.String(),
		"creative_brief_id": creativeBrief.ID.String(),
		"storyboard_id":     storyboard.ID.String(),
		"project_name":      projectName,
		"theme":             req.SelectedTheme,
		"tone":              req.ToneOfVoice,
		"duration":          duration,
		"institution_name":  req.InstitutionName,
		"school_level":      schoolLevel,
		"event_content":     req.EventContent,
		"key_message":       req.SelectedKeyMessage,
		"copywriting":       req.EditableCopywriting,
		"hashtags":          req.EditableHashtags,
	}, nil
}

// Helper function to parse duration string to integer
func parseDurationToInt(duration string) int {
	cleaned := strings.TrimSpace(strings.ToLower(duration))
	parts := strings.Fields(cleaned)
	if len(parts) > 0 {
		num, err := strconv.Atoi(parts[0])
		if err == nil && num > 0 {
			return num
		}
	}
	return 30 // default
}
