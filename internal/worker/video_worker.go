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
	// Identify the 3 scenes (hook, value, cta)
	var hook, value, cta model.StoryboardSection
	for _, sec := range sections {
		switch strings.ToLower(sec.SectionType) {
		case "hook":
			hook = sec
		case "value":
			value = sec
		case "cta":
			cta = sec
		}
	}

	// Fallback jika tidak ada section (untuk amannya)
	if hook.Duration == 0 { hook.Duration = 5 }
	if value.Duration == 0 { value.Duration = 5 }
	if cta.Duration == 0 { cta.Duration = 5 }

	// Standard phrasing (kata-kata pakem)
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
		cb.Theme,
		bb.InstitutionName,
		bb.SchoolLevel,
		bb.OfferedDegrees,      
		bb.InstitutionHistory,  
		cb.Theme, // Style / Tone diwakili Theme atau ToneOfVoice
		0, hook.Duration, hook.Content,
		hook.Duration, hook.Duration+value.Duration, value.Content,
		hook.Duration+value.Duration, hook.Duration+value.Duration+cta.Duration, cta.Content,
	)

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
		Prompt:          prompt,
		ReferenceImages: images,
	}
}
