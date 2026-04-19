package utils

import (
	"strconv"
	"strings"
)

// ParseVideoDuration converts FE duration string (e.g., "15 detik", "30 detik") to integer seconds
func ParseVideoDuration(duration string) int {
	// Remove extra spaces and convert to lowercase
	cleaned := strings.TrimSpace(strings.ToLower(duration))

	// Split by space to extract number
	parts := strings.Fields(cleaned)
	if len(parts) > 0 {
		num, err := strconv.Atoi(parts[0])
		if err == nil && num > 0 {
			return num
		}
	}

	// Default to 30 seconds if parsing fails
	return 30
}

// MapEventContentToVideoType maps FE event content to video type
func MapEventContentToVideoType(eventContent string) string {
	switch strings.ToLower(eventContent) {
	case "penerimaan mahasiswa baru":
		return "recruitment"
	case "dies natalis / ulang tahun":
		return "anniversary"
	case "promosi beasiswa":
		return "scholarship_promo"
	case "pengenalan kehidupan kampus":
		return "campus_introduction"
	default:
		return "promotional"
	}
}

// MapThemeToStyle maps FE theme selection to visual style
func MapThemeToStyle(theme string) string {
	switch strings.ToLower(theme) {
	case "tur kampus sinematik":
		return "cinematic"
	case "cerita kehidupan mahasiswa":
		return "narrative"
	case "keunggulan akademik":
		return "educational"
	case "tren & gaya hidup cepat":
		return "trendy"
	default:
		return "cinematic"
	}
}

// MapToneToMusicPreference maps tone of voice to music preference
func MapToneToMusicPreference(tone string) string {
	switch strings.ToLower(tone) {
	case "santai & ramah":
		return "upbeat_casual"
	case "profesional & formal":
		return "professional"
	case "kreatif & inovatif":
		return "modern_creative"
	case "berwibawa & meyakinkan":
		return "powerful_inspirational"
	default:
		return "professional"
	}
}
