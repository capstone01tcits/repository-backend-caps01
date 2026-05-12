package service

import (
	"errors"
	"fmt"

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
	GetVeo3TestPayload(userID, storyboardID string) (*model.Veo3TestPayload, error)
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

func (s *storyboardService) AutoGenerateStoryboard(userID string, projectID uuid.UUID, bb *model.BusinessBrief, cb *model.CreativeBrief) (*model.Storyboard, error) {
	// 1. Create Storyboard entity
	storyboard := &model.Storyboard{
		ProjectID:     projectID,
		UserID:        uuid.MustParse(userID),
		Title:         "Promotional Video: " + bb.InstituteName,
		Description:   "Auto-generated 3-scene storyboard for " + cb.VideoType,
		Style:         cb.Style,
		TotalDuration: cb.Duration,
	}

	if err := s.storyboardRepo.Create(storyboard); err != nil {
		return nil, fmt.Errorf("gagal membuat storyboard otomatis: %w", err)
	}

	// 2. Allocate duration for Hook (0-5s), Value (5-12s), CTA (12-15s)
	hookDur, valueDur, ctaDur := 5, 7, 3
	if cb.Duration != 15 && cb.Duration > 0 {
		hookDur = int(float64(cb.Duration) * 0.33)
		ctaDur = int(float64(cb.Duration) * 0.20)
		valueDur = cb.Duration - hookDur - ctaDur
	}

	// 3. Construct 3 mandatory scenes
	sections := []model.StoryboardSection{
		{
			StoryboardID: storyboard.ID,
			UserID:       uuid.MustParse(userID),
			SectionType:  "hook",
			Content:      fmt.Sprintf("Visual memukau menampilkan landmark/suasana %s. Tone: %s.", bb.InstituteName, cb.Tone),
			Duration:     hookDur,
		},
		{
			StoryboardID: storyboard.ID,
			UserID:       uuid.MustParse(userID),
			SectionType:  "value",
			Content:      fmt.Sprintf("Sorot fasilitas dan keunggulan spesifik. Catatan: %s", cb.AdditionalNotes),
			Duration:     valueDur,
		},
		{
			StoryboardID: storyboard.ID,
			UserID:       uuid.MustParse(userID),
			SectionType:  "cta",
			Content:      fmt.Sprintf("Pesan penutup: '%s'. %s", cb.CallToAction, cb.Copywriting),
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

func (s *storyboardService) GetVeo3TestPayload(userID, storyboardID string) (*model.Veo3TestPayload, error) {
	storyboard, err := s.storyboardRepo.FindByID(storyboardID)
	if err != nil {
		return nil, errors.New("storyboard tidak ditemukan")
	}

	if storyboard.UserID.String() != userID {
		return nil, errors.New("unauthorized access")
	}

	// Fetch business brief and creative brief to get meta info
	bb, err := s.briefRepo.FindBusinessBriefByProjectID(storyboard.ProjectID.String())
	if err != nil {
		return nil, errors.New("business brief tidak ditemukan")
	}

	cb, err := s.briefRepo.FindCreativeBriefByProjectID(storyboard.ProjectID.String())
	if err != nil {
		return nil, errors.New("creative brief tidak ditemukan")
	}

	// Fetch sections
	sections, err := s.storyboardRepo.FindSectionsByStoryboardID(storyboardID)
	if err != nil {
		return nil, errors.New("gagal mengambil sections storyboard")
	}

	// Construct "Kata-kata Pakem" for Veo3 Prompt
	var hook, value, cta model.StoryboardSection
	for _, sec := range sections {
		switch sec.SectionType {
		case "hook":
			hook = sec
		case "value":
			value = sec
		case "cta":
			cta = sec
		}
	}

	// Standard phrasing (kata-kata pakem) menggunakan template yang lebih detail
	prompt := fmt.Sprintf(
		`Buatlah video promosi %s berkualitas tinggi untuk iklan institusi pendidikan.

		Detail Institusi:
		- Nama: %s
		- Tingkat Pendidikan: %s
		- Program Studi: %s
		- Latar Belakang/Sejarah: %s
		- Gunakan logo dan foto lingkungan kampus sebagai referensi visual.

		Tujuan Video:
		Membuat video promosi yang menarik, modern, profesional, dan membangun kepercayaan, dengan menonjolkan kualitas akademik, fasilitas, serta peluang masa depan.

		Gaya & Tone:
		%s
		Target Audiens: Calon siswa/mahasiswa dan orang tua.
		Gaya Visual: Sinematik, bersih, profesional, transisi halus, kualitas produksi tinggi.

		SCENE STRUCTURE:

		SCENE 1 (%ds–%ds): HOOK
		%s

		SCENE 2 (%ds–%ds): NILAI UNGGULAN
		%s

		SCENE 3 (%ds–%ds): CALL TO ACTION
		%s

		Panduan Teknis:
		- Pertahankan kesinambungan sinematik
		- Gunakan karakter yang konsisten di setiap scene
		- Transisi antar scene harus halus dan natural
		- Tampilkan nama institusi dengan jelas
		- Tonjolkan kepercayaan, prestasi, dan masa depan cerah
		- Tambahkan nuansa musik latar inspiratif
		- Hindari klaim berlebihan atau tidak realistis`,
		cb.Style,
		bb.InstituteName,
		bb.SchoolLevel,
		bb.OfferedDegrees,
		bb.InstitutionHistory,
		cb.Style,
		0, hook.Duration, hook.Content,
		hook.Duration, hook.Duration+value.Duration, value.Content,
		hook.Duration+value.Duration, hook.Duration+value.Duration+cta.Duration, cta.Content,
	)

	// Reference images from initial input (logo and environment)
	var refImages []string
	if bb.LogoPath != "" {
		refImages = append(refImages, bb.LogoPath)
	}
	if bb.EnvironmentPath != "" {
		refImages = append(refImages, bb.EnvironmentPath)
	}

	return &model.Veo3TestPayload{
		// Model:           "veo3",
		Prompt:          prompt,
		// ReferenceImages: refImages,
	}, nil
}

