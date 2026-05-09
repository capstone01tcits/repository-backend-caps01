package worker

import (
	"fmt"
	"strings"

	"Sevima-AI-Content-Creator/internal/model"
)

// Veo3Payload adalah skema JSON standar yang akan dikirim ke AI Service untuk pemrosesan Veo 3
type Veo3Payload struct {
	Model           string   `json:"model"`
	Prompt          string   `json:"prompt"`
	ReferenceImages []string `json:"reference_images"`
}

// BuildVeo3Payload merangkai prompt sesuai dengan standar Veo 3 untuk video promosi pendidikan
func BuildVeo3Payload(bb *model.BusinessBrief, cb *model.CreativeBrief, sections []model.StoryboardSection) Veo3Payload {
	// 1. Identity Mapping & Context Awareness
	contextKeywords := ""
	schoolLevelStr := strings.ToLower(bb.SchoolLevel)
	
	if strings.Contains(schoolLevelStr, "university") || strings.Contains(schoolLevelStr, "perguruan tinggi") || strings.Contains(schoolLevelStr, "kampus") {
		contextKeywords = "Higher Education, Campus Life, Research, Independent Learning"
	} else {
		contextKeywords = "Vibrant Classroom, Practical Skills, Student Discipline, Nurturing Environment"
	}

	// 2. Visual Quality Standard (Sesuai Prompt dari User)
	visualDirection := "Realistic skin textures, natural cinematic lighting, 4K resolution, no plastic-look faces"

	// 3. Merangkai Header Prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString(fmt.Sprintf(
		"Create a cinematic promotional video for %s. Type: %s. Event: %s. Tone: %s. Theme: %s. Duration: %d seconds. Key message: %s. Core emotion: %s. Visual direction: %s.\n\n",
		bb.InstituteName, bb.SchoolLevel, cb.VideoType, cb.Tone, cb.Style, cb.Duration, cb.CallToAction, contextKeywords, visualDirection,
	))

	// 4. Memasukkan Scene (Hook, Value, CTA) beserta perhitungan detik
	timeAccumulator := 0
	for i, sec := range sections {
		startTime := timeAccumulator
		endTime := timeAccumulator + sec.Duration
		timeAccumulator = endTime
		
		// Format: SCENE 1 (0-5s): HOOK [Visual description]
		promptBuilder.WriteString(fmt.Sprintf(
			"SCENE %d (%d–%ds): %s [%s]\n",
			i+1, startTime, endTime, strings.ToUpper(sec.SectionType), sec.Content,
		))
	}

	// 5. Aturan Transisi Wajib
	promptBuilder.WriteString("\nMaintain cinematic continuity, same characters, natural skin textures, and smooth transitions.")

	// 6. Validasi Reference Images (Minimal 1)
	var images []string
	if bb.LogoPath != "" {
		images = append(images, bb.LogoPath) // Pastikan sudah berupa format base64
	}
	if bb.EnvironmentPath != "" {
		images = append(images, bb.EnvironmentPath)
	}
	// Jika tidak ada gambar, bisa ditambahkan fallback sesuai kebutuhan

	return Veo3Payload{
		Model:           "veo3",
		Prompt:          promptBuilder.String(),
		ReferenceImages: images,
	}
}
