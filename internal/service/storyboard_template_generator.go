package service

import (
	"fmt"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"
)

// TemplateGenerator generates storyboard templates based on project data
type TemplateGenerator struct {
	projectRepo repository.ProjectRepository
	briefRepo   repository.BriefRepository
}

func NewTemplateGenerator(projectRepo repository.ProjectRepository, briefRepo repository.BriefRepository) *TemplateGenerator {
	return &TemplateGenerator{
		projectRepo: projectRepo,
		briefRepo:   briefRepo,
	}
}

// GenerateTemplates creates multiple storyboard template variations
func (tg *TemplateGenerator) GenerateTemplates(projectID string, videoDuration int) ([]model.StoryboardTemplate, error) {
	// Fetch project and briefs to get data
	project, err := tg.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	// Get business brief for detailed information
	businessBrief, err := tg.briefRepo.FindBusinessBriefByProjectID(projectID)
	if err != nil {
		// Business brief is optional, continue without it
		businessBrief = nil
	}

	// Extract relevant data from project and briefs
	projectData := map[string]string{
		"theme":       project.Theme,
		"name":        project.Name,
		"objective":   "",
		"keyMessage":  "",
		"tone":        "",
		"copywriting": "",
	}

	// Add data from business brief if available
	if businessBrief != nil {
		if businessBrief.ProjectObjective != "" {
			projectData["objective"] = businessBrief.ProjectObjective
		}
		if businessBrief.InstituteName != "" {
			projectData["name"] = businessBrief.InstituteName
		} else if businessBrief.CompanyName != "" {
			projectData["name"] = businessBrief.CompanyName
		}
		if businessBrief.KeyMessage != "" {
			projectData["keyMessage"] = businessBrief.KeyMessage
		}
	}

	// Calculate section durations based on total video duration
	hookDuration := videoDuration / 3
	valueDuration := videoDuration / 3
	ctaDuration := videoDuration - hookDuration - valueDuration

	// Generate multiple template styles
	templates := []model.StoryboardTemplate{
		tg.generateDynamicTemplate(projectData, hookDuration, valueDuration, ctaDuration),
		tg.generateNarrativeTemplate(projectData, hookDuration, valueDuration, ctaDuration),
		tg.generateEnergeticTemplate(projectData, hookDuration, valueDuration, ctaDuration),
		tg.generateMinimalistTemplate(projectData, hookDuration, valueDuration, ctaDuration),
	}

	return templates, nil
}

// generateDynamicTemplate creates a fast-paced, attention-grabbing template for educational institutions
func (tg *TemplateGenerator) generateDynamicTemplate(data map[string]string, hookDur, valueDur, ctaDur int) model.StoryboardTemplate {
	institutionName := data["name"]
	if institutionName == "" {
		institutionName = "Your Institution"
	}
	
	tone := data["tone"]
	if tone == "" {
		tone = "Energetic, engaging, and fast-paced"
	}
	copywriting := data["copywriting"]
	if copywriting != "" {
		copywriting = "\n\nCAPTION / COPYWRITING FOR SOCIAL MEDIA:\n\"" + copywriting + "\"\nHashtags: " + data["hashtags"]
	}

	return model.StoryboardTemplate{
		TemplateID:  "dynamic_template",
		Style:       "Dynamic",
		Description: "Fast-paced, attention-grabbing with motion and energy. Ideal for attracting Gen Z students to undergraduate or postgraduate programs.",
		Duration:    hookDur + valueDur + ctaDur,
		Sections: []model.TemplateStoryboardSection{
			{
				SectionType:       "hook",
				Title:             "Hook - Grab Attention",
				SuggestedDuration: hookDur,
				Content: fmt.Sprintf(
					"START WITH IMPACT: Capture attention in the first %ds with striking campus visuals or compelling student testimonial.\n\nOpening ideas:\n- Student quote: \"I came to study, I found my future\"\n- Question: \"What if your college changed everything?\"\n- Statement: \"%s is reshaping education\"\n\nVisuals:\n- Campus landmarks with dynamic music\n- Students collaborating in modern facilities\n- Quick cuts of campus life and activities\n- Aerial campus shots with motion graphics\n\nMake it impossible to scroll past - use fast pacing, trending audio, and energy that appeals to Gen Z.",
					hookDur, institutionName,
				),
				Tips: fmt.Sprintf("Tone: %s. Use fast transitions, trending audio, and relatable student voices. Create urgency and curiosity.", tone),
			},
			{
				SectionType:       "value",
				Title:             "Value - Highlight Excellence",
				SuggestedDuration: valueDur,
				Content: fmt.Sprintf(
					"SHOWCASE WHAT MAKES YOU DIFFERENT:\n\nHighlight key institutional strengths:\n- Academic excellence: Rankings, accreditations, program highlights\n- Faculty expertise: Distinguished professors, research opportunities\n- Campus facilities: Modern labs, libraries, sports complexes, student centers\n- Student life: Clubs, events, cultural diversity, residential experience\n- Career outcomes: Graduate placement rates, alumni success stories, industry partnerships\n- Unique programs: Specializations in %s, scholarships, internship opportunities\n\nInstitution: %s\n\nUse graphic overlays for:\n- Ranking badges (\"#1 Engineering Program\")\n- Student testimonials (\"This changed my career\")\n- Statistics (\"95%% employment rate\")\n- Alumni company logos (\"Our graduates work at...\")\n\nShow campus diversity, modern infrastructure, and hands-on learning experiences.",
					data["theme"], institutionName,
				),
				Tips: "Use concrete data and visual proof. Show student diversity. Make excellence tangible and achievable.",
			},
			{
				SectionType:       "cta",
				Title:             "CTA - Drive Enrollment",
				SuggestedDuration: ctaDur,
				Content: fmt.Sprintf(
					"DRIVE ENROLLMENT (Next %ds):\n\nClear calls to action:\n- \"Apply now at [admissions URL]\"\n- \"Limited spots available - apply by [DATE]\"\n- \"Early application scholarship deadline [DATE]\"\n- \"Join [number]+ students at %s\"\n- \"Request information: [phone/email]\"\n\nCreate urgency:\n- Application deadline display\n- \"Full scholarships for top students\"\n- \"Early bird pricing ends [DATE]\"\n- \"Open house dates: [DATES]\"\n\nVisual elements:\n- Prominent website/application link\n- QR code to admission portal\n- Phone number and email displayed\n- Striking final campus image with logo\n\nEnding: Institution logo + inspirational tagline (e.g., \"Transform Your Future at %s\")",
					ctaDur, institutionName, institutionName,
				) + copywriting,
				Tips: "Make action obvious and easy. Include deadline to create urgency. Use bold overlays and clear CTAs.",
			},
		},
	}
}

// generateNarrativeTemplate creates a storytelling-driven template focused on student journey
func (tg *TemplateGenerator) generateNarrativeTemplate(data map[string]string, hookDur, valueDur, ctaDur int) model.StoryboardTemplate {
	institutionName := data["name"]
	if institutionName == "" {
		institutionName = "Our Institution"
	}
	
	tone := data["tone"]
	if tone == "" {
		tone = "Inspiring, emotional, and authentic"
	}
	copywriting := data["copywriting"]
	if copywriting != "" {
		copywriting = "\n\nCAPTION / COPYWRITING FOR SOCIAL MEDIA:\n\"" + copywriting + "\"\nHashtags: " + data["hashtags"]
	}

	return model.StoryboardTemplate{
		TemplateID:  "narrative_template",
		Style:       "Narrative",
		Description: "Story-driven approach following student journey from admission through graduation and career success. Emotionally engaging for prospective students.",
		Duration:    hookDur + valueDur + ctaDur,
		Sections: []model.TemplateStoryboardSection{
			{
				SectionType:       "hook",
				Title:             "Hook - Student Story Begins",
				SuggestedDuration: hookDur,
				Content: fmt.Sprintf(
					"ESTABLISH THE STUDENT JOURNEY:\n\nOpen with relatable student scenario:\n- \"Meet [student name] - she came with a dream\"\n- \"Three years ago, he stood where you are now\"\n- \"This is the story of how %s changed lives\"\n\nIntroduce diverse student characters:\n- First-generation student, career-changer, athlete, international student\n- Show their initial doubts or aspirations\n- Establish why they chose %s\n\nCinematic wide shots:\n- Campus entrance - student's first day\n- Classroom ambience - beginning of learning journey\n- Warm, inviting institutional atmosphere\n\nSet the emotional tone: hopeful, transformative, inspiring",
					institutionName, institutionName,
				),
				Tips: fmt.Sprintf("Tone: %s. Make it personal and relatable. Show diverse student backgrounds. Build emotional connection.", tone),
			},
			{
				SectionType:       "value",
				Title:             "Value - The Transformation",
				SuggestedDuration: valueDur,
				Content: fmt.Sprintf(
					"SHOW THE TRANSFORMATION:\n\nFollow student journey through key moments:\n- First day in class with inspiring professor\n- Collaborative project with peers\n- Research opportunity or internship\n- Campus life - clubs, events, friendships\n- Academic breakthrough or skill development\n- Real-world application of learning\n\nHighlight institutional offerings:\n- Faculty mentorship and guidance\n- Hands-on learning experiences at %s\n- Professional development and career services\n- Diverse campus community and belonging\n- Practical skills for post-graduation success\n\nShow visible growth and confidence:\n- Student gaining expertise in %s\n- Building meaningful relationships\n- Developing leadership and soft skills\n- Preparing for career opportunities\n\nUse visual storytelling:\n- Montage of learning moments\n- Student testimonial quotes\n- Before & after confidence levels\n- Success metrics (projects completed, internships landed)",
					institutionName, data["theme"],
				),
				Tips: "Show genuine transformation. Use real student testimonials. Build momentum through the journey.",
			},
			{
				SectionType:       "cta",
				Title:             "CTA - Your Turn Starts Now",
				SuggestedDuration: ctaDur,
				Content: fmt.Sprintf(
					"INSPIRE FUTURE STUDENTS (Next %ds):\n\nShow graduate success:\n- \"She graduated and was hired immediately\"\n- \"He now leads projects at [industry]\"\n- \"Alumni network spans [number]+ professionals\"\n- \"[number]%% placement rate within 6 months\"\n\nBring it full circle:\n- \"Your story at %s starts here\"\n- \"You could be next - join [number]+ successful graduates\"\n- \"Transform your future like thousands before you\"\n\nEmpower prospective students:\n- Display career outcomes\n- Alumni company logos and achievements\n- Salary/career growth statistics\n- Graduate testimonials\n\nClear call to action:\n- \"Apply now at [admissions URL]\"\n- \"Application deadline: [DATE]\"\n- \"Early bird scholarship: [DATE]\"\n- \"Request prospectus: [email/phone]\"\n- \"Campus tour booking: [link]\"\n\nClosing visuals:\n- Graduating class celebration\n- Alumni achievement montage\n- Institution logo with tagline\n- Strong inspirational message",
					ctaDur, institutionName,
				) + copywriting,
				Tips: "End on a high note showing success and achievement. Make enrollment feel like the next step in an exciting journey.",
			},
		},
	}
}

// generateEnergeticTemplate creates a high-energy, youth-focused template for campus promotion
func (tg *TemplateGenerator) generateEnergeticTemplate(data map[string]string, hookDur, valueDur, ctaDur int) model.StoryboardTemplate {
	institutionName := data["name"]
	if institutionName == "" {
		institutionName = "Our Campus"
	}
	
	tone := data["tone"]
	if tone == "" {
		tone = "Vibrant, youthful, and exciting"
	}
	copywriting := data["copywriting"]
	if copywriting != "" {
		copywriting = "\n\nCAPTION / COPYWRITING FOR SOCIAL MEDIA:\n\"" + copywriting + "\"\nHashtags: " + data["hashtags"]
	}

	return model.StoryboardTemplate{
		TemplateID:  "energetic_template",
		Style:       "Energetic",
		Description: "High-energy, youth-focused showcasing campus culture, student life, and social vibrancy. Perfect for social media and Gen Z audiences.",
		Duration:    hookDur + valueDur + ctaDur,
		Sections: []model.TemplateStoryboardSection{
			{
				SectionType:       "hook",
				Title:             "Hook - Campus Energy",
				SuggestedDuration: hookDur,
				Content: fmt.Sprintf(
					"GO BIG, GO FAST - CAPTURE CAMPUS ENERGY (%ds):\n\nTrending approaches for Gen Z:\n- Bold text overlay: \"THIS IS %s\", \"WAIT FOR IT\", \"POV: You're a student\"\n- Quick cut montage of exciting campus moments\n- Trending audio with fast-paced visuals\n- Question: \"What if your college was actually this fun?\"\n\nVisuals:\n- Students laughing and collaborating\n- Campus clubs and events in action\n- Sports, performances, late-night study sessions\n- Diverse student body having genuine fun\n- Fast cuts (2-3 frames per second)\n\nBuild curiosity and FOMO:\nShow the social experience, not just academics\nMake them want to be part of this community",
					hookDur, institutionName,
				),
				Tips: fmt.Sprintf("Tone: %s. Use trending audio and fast transitions. Make campus life look irresistible. Capture genuine student moments.", tone),
			},
			{
				SectionType:       "value",
				Title:             "Value - Why Students Love It",
				SuggestedDuration: valueDur,
				Content: fmt.Sprintf(
					"DELIVER THE ENERGY - SHOW CAMPUS LIFE (%ds):\n\nWhat makes %s amazing:\n- Vibrant student organizations and clubs\n- Social proof: \"Join [number]+ students\"\n- Campus events: concerts, festivals, competitions\n- Diverse student community and friendships\n- Modern facilities and amenities\n- Hands-on learning with real-world impact\n\nKey highlights related to %s:\n- Industry-leading programs\n- Career opportunities and internships\n- Student testimonials (short, punchy quotes)\n- Social media moments and viral content\n\nKeep cuts fast and visually interesting:\n- Quick montage of campus locations\n- Overlaid text with fun facts\n- Upbeat music and sound effects\n- Student reactions and celebrations\n- Real campus activity footage\n\nBuild momentum toward action",
					valueDur, institutionName, data["theme"],
				),
				Tips: "Show genuine campus culture. Use user-generated content and real student moments. Keep energy high.",
			},
			{
				SectionType:       "cta",
				Title:             "CTA - Join the Movement",
				SuggestedDuration: ctaDur,
				Content: fmt.Sprintf(
					"THE MOMENT - JOIN US (%ds):\n\nFinal call with urgency:\n- \"Your seat is waiting\"\n- \"Apply now and join the community\"\n- \"Don't miss your moment\"\n- \"Early application deadline: [DATE]\"\n\nDisplay in large, bold text:\n- Admissions website/link\n- Application URL\n- QR code for instant application\n- Social media handles (Instagram, TikTok, etc.)\n\nCREATE ENGAGEMENT:\n- \"Link in bio\"\n- Instagram/TikTok handles\n- Hashtag campaign (e.g., \"#MyFutureAt%s\")\n- Follow, tag, share to enter scholarship raffle\n\nFINALE FRAME:\n- %s logo with vibrant animation\n- \"APPLY TODAY\" in bold overlay\n- Admissions contact info\n- Campus photo or student celebration image\n\nEnd with energetic music and sense of belonging",
					ctaDur, institutionName, institutionName,
				) + copywriting,
				Tips: "Make action super easy and fun. Use multiple CTA methods (link, QR, social). Create urgency with deadline.",
			},
		},
	}
}

// generateMinimalistTemplate creates a clean, professional template emphasizing academic excellence
func (tg *TemplateGenerator) generateMinimalistTemplate(data map[string]string, hookDur, valueDur, ctaDur int) model.StoryboardTemplate {
	institutionName := data["name"]
	if institutionName == "" {
		institutionName = "Our Institution"
	}
	
	tone := data["tone"]
	if tone == "" {
		tone = "Elegant, professional, and prestigious"
	}
	copywriting := data["copywriting"]
	if copywriting != "" {
		copywriting = "\n\nCAPTION / COPYWRITING FOR SOCIAL MEDIA:\n\"" + copywriting + "\"\nHashtags: " + data["hashtags"]
	}

	return model.StoryboardTemplate{
		TemplateID:  "minimalist_template",
		Style:       "Minimalist",
		Description: "Clean, professional, and sophisticated. Emphasizes institutional prestige and academic excellence. Ideal for premium positioning.",
		Duration:    hookDur + valueDur + ctaDur,
		Sections: []model.TemplateStoryboardSection{
			{
				SectionType:       "hook",
				Title:             "Hook - Elegant Opening",
				SuggestedDuration: hookDur,
				Content: fmt.Sprintf(
					"SOPHISTICATED START - ESTABLISH PRESTIGE (%ds):\n\nMinimalist approach:\n- Single, striking campus image (architecture or research facility)\n- Elegant fade-in with subtle music\n- Clean typography displaying institution name: %s\n- Simple, powerful statement: \"Excellence in Education\"\n\nTechnique:\n- Soft, premium classical or modern music\n- Plenty of white space and negative space\n- High-contrast text on clean background\n- Premium color palette (institutional colors)\n- Slow, deliberate pacing\n\nTone: Elegant, trustworthy, premium, established\n\nSet expectation of quality and prestige",
					hookDur, institutionName,
				),
				Tips: fmt.Sprintf("Tone: %s. Less is more. Let silence speak. Quality over quantity. Build anticipation.", tone),
			},
			{
				SectionType:       "value",
				Title:             "Value - Academic Excellence",
				SuggestedDuration: valueDur,
				Content: fmt.Sprintf(
					"DELIVER WITH CLARITY - SHOWCASE EXCELLENCE (%ds):\n\nInstitutional strengths:\n- Academic prestige: Rankings, accreditations, awards\n- Faculty expertise: Nobel laureates, renowned researchers\n- Research impact: Major discoveries, publications\n- Graduate outcomes: Top employer placements\n- International recognition and partnerships\n- Specialized programs: Engineering, Medicine, Business, etc.\n\nKey Message: %s\n\nMinimalist presentation:\n- One quality point per scene\n- Clean data visualization (minimal charts/graphs)\n- Faculty or student testimonial (professional tone)\n- Elegant architecture or research facilities\n- Premium photography with plenty of breathing room\n\nVisuals:\n- Institution campus landmarks\n- Modern lab or research facility\n- Historic library or academic buildings\n- Award plaques and accreditation badges\n- Consistent institutional branding\n- Sophisticated color scheme\n\nBuild confidence through clarity and substance",
					valueDur, data["keyMessage"],
				),
				Tips: "Respect audience intelligence. Focus on achievements and substance. No hard sell.",
			},
			{
				SectionType:       "cta",
				Title:             "CTA - Premium Invitation",
				SuggestedDuration: ctaDur,
				Content: fmt.Sprintf(
					"PREMIUM CLOSE - REFINE FINAL MESSAGE (%ds):\n\nSophisticated call to action:\n- \"Discover Excellence\"\n- \"Join a Legacy of Achievement\"\n- \"Apply to %s\"\n- \"Begin Your Scholarly Journey\"\n\nDisplay with elegance:\n- Institution name prominently\n- Clean website URL in elegant typography\n- Minimalist QR code (tastefully positioned)\n- Institutional logo\n- Elegant tagline or mission statement\n\nFINAL FRAME:\n- %s Logo reveal (slow, elegant animation)\n- Single powerful statement\n- Website/admissions portal (clean display)\n- Contact information (if space permits)\n- Brief, powerful closing statement\n- Muted, sophisticated exit music\n\nLEGACY APPROACH:\n- Emphasize tradition and excellence\n- Show multi-generational impact\n- Suggest joining an elite community\n- No rush - quality institutions speak quietly\n\nRemember: Premium brands whisper, not shout. Leave them impressed and wanting to learn more.",
					ctaDur, institutionName, institutionName,
				) + copywriting,
				Tips: "Be concise and elegant. Emphasize prestige and tradition. Subtle but memorable.",
			},
		},
	}
}
