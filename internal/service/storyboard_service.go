package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"

	"github.com/google/uuid"
)

type StoryboardService interface {
	CreateManualStoryboard(userID string, req *model.CreateManualStoryboardRequest) (*model.Storyboard, error)
	GetStoryboardByProject(userID, projectID string) (*model.Storyboard, error)
	GetStoryboard(userID, storyboardID string) (*model.Storyboard, error)
	UpdateStoryboard(userID, storyboardID string, req *model.UpdateManualStoryboardRequest) (*model.Storyboard, error)
	GetStoryboardSections(userID, storyboardID string) ([]model.StoryboardSection, error)
	DeleteStoryboard(userID, storyboardID string) error
	RestoreStoryboard(userID, storyboardID string) error
	AutoGenerateStoryboard(userID string, projectID uuid.UUID, bb *model.BusinessBrief, cb *model.CreativeBrief) (*model.Storyboard, error)
}

type storyboardService struct {
	storyboardRepo repository.StoryboardRepository
	projectRepo    repository.ProjectRepository
	briefRepo      repository.BriefRepository
}

func NewStoryboardService(
	storyboardRepo repository.StoryboardRepository,
	projectRepo repository.ProjectRepository,
	briefRepo repository.BriefRepository,
) StoryboardService {
	return &storyboardService{storyboardRepo, projectRepo, briefRepo}
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

func (s *storyboardService) GetStoryboardByProject(userID, projectID string) (*model.Storyboard, error) {
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

// helper to parse duration safely
func parseDurationStr(durStr string) int {
	cleaned := strings.TrimSpace(strings.ToLower(durStr))
	parts := strings.Fields(cleaned)
	if len(parts) > 0 {
		num, err := strconv.Atoi(parts[0])
		if err == nil && num > 0 {
			return num
		}
	}
	return 30
}

func (s *storyboardService) AutoGenerateStoryboard(userID string, projectID uuid.UUID, bb *model.BusinessBrief, cb *model.CreativeBrief) (*model.Storyboard, error) {
	durInt := parseDurationStr(cb.VideoDuration)
	
	// 0. Hapus storyboard lama jika ada (untuk fitur edit/update)
	existing, err := s.storyboardRepo.FindByProjectID(projectID.String())
	if err == nil && existing != nil {
		s.storyboardRepo.Delete(existing.ID.String())
	}

	// 1. Create Storyboard entity
	storyboard := &model.Storyboard{
		ProjectID:     projectID,
		UserID:        uuid.MustParse(userID),
		Title:         "Promotional Video: " + bb.InstitutionName,
		Description:   "Auto-generated 3-scene storyboard for " + cb.EventContent,
		Style:         cb.Theme,
		TotalDuration: durInt,
	}

	if err := s.storyboardRepo.Create(storyboard); err != nil {
		return nil, fmt.Errorf("gagal membuat storyboard otomatis: %w", err)
	}

	// 2. Allocate duration for 3 scenes to fit total Veo 3 Lite video constraints [4, 6, 8]
	hookDur, valueDur, ctaDur := 2, 3, 1 // default for 6s

	if durInt == 4 {
		hookDur, valueDur, ctaDur = 1, 2, 1
	} else if durInt == 6 {
		hookDur, valueDur, ctaDur = 2, 3, 1
	} else if durInt == 8 {
		hookDur, valueDur, ctaDur = 2, 4, 2
	}

	// 3. Construct 3 mandatory scenes (stored as JSON string to support frontend narration & visual fields)
	hookNarration := fmt.Sprintf(`"Halo generasi masa depan! Tahukah kamu bahwa %s"`, strings.ToLower(cb.KeyMessage))
	hookVisual := fmt.Sprintf("Visual bergaya %s. Menampilkan gerbang utama %s. Sesuai instruksi: %s.", cb.ToneOfVoice, bb.InstitutionName, cb.Prompt)

	valueNarration := fmt.Sprintf(`"Di %s, kami siap membantumu mewujudkan impian itu melalui program unggulan kami."`, bb.InstitutionName)
	valueVisual := fmt.Sprintf("Gaya visual: %s. Memperlihatkan mahasiswa sedang beraktivitas, fasilitas modern.", cb.Theme)

	ctaNarration := fmt.Sprintf(`"Jangan lewatkan momen %s tahun ini. Yuk, raih mimpimu bersama kami!"`, cb.EventContent)
	ctaVisual := fmt.Sprintf("Logo %s muncul di tengah layar dengan teks ajakan (Call to Action).", bb.InstitutionName)

	sections := []model.StoryboardSection{
		{
			StoryboardID: storyboard.ID,
			UserID:       uuid.MustParse(userID),
			SectionType:  "Intro & Hook",
			Content:      fmt.Sprintf(`{"narration": %q, "visual": %q}`, hookNarration, hookVisual),
			Duration:     hookDur,
		},
		{
			StoryboardID: storyboard.ID,
			UserID:       uuid.MustParse(userID),
			SectionType:  "Suasana & Keunggulan Kampus",
			Content:      fmt.Sprintf(`{"narration": %q, "visual": %q}`, valueNarration, valueVisual),
			Duration:     valueDur,
		},
		{
			StoryboardID: storyboard.ID,
			UserID:       uuid.MustParse(userID),
			SectionType:  "Promosi & Call to Action",
			Content:      fmt.Sprintf(`{"narration": %q, "visual": %q}`, ctaNarration, ctaVisual),
			Duration:     ctaDur,
		},
	}

	for _, sec := range sections {
		if err := s.storyboardRepo.CreateSection(&sec); err != nil {
			return nil, fmt.Errorf("gagal menyimpan scene %s: %w", sec.SectionType, err)
		}
	}

	storyboard.Sections = sections
	return storyboard, nil
}
