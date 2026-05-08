package service

// EducationPromotionTemplates contains pre-prepared content templates for educational institution promotion
// These templates are combined with actual project data to create final storyboard content
type EducationPromotionTemplates struct {
	HookTemplates  []HookTemplate
	ValueTemplates []ValueTemplate
	CTATemplates   []CTATemplate
}

type HookTemplate struct {
	Name    string
	Pattern string
}

type ValueTemplate struct {
	Name    string
	Pattern string
}

type CTATemplate struct {
	Name    string
	Pattern string
}

// NewEducationPromotionTemplates returns all pre-prepared templates for education promotion
func NewEducationPromotionTemplates() *EducationPromotionTemplates {
	return &EducationPromotionTemplates{
		HookTemplates:  getHookTemplates(),
		ValueTemplates: getValueTemplates(),
		CTATemplates:   getCTATemplates(),
	}
}

// HOOK TEMPLATES - Designed to grab attention in first 10 seconds
func getHookTemplates() []HookTemplate {
	return []HookTemplate{
		{
			Name: "student_testimonial",
			Pattern: `START WITH IMPACT: Student testimonial in first 3 seconds.
			
Quote: "I came to study {{.FieldOfStudy}}, I found my future"
Alternative: "{{.InstitutionName}} changed everything for me"

Visuals:
- Quick campus landmark shot (3 seconds)
- Students collaborating in modern facility
- Fast-paced cuts with dynamic music
- Aerial campus shots

Make it impossible to miss - use trending audio and relatable student voice.`,
		},
		{
			Name: "question_opener",
			Pattern: `SPARK CURIOSITY: Open with compelling question.
			
Question: "What if your college transformed your future?"
Alternative: "Are you ready to join {{.StudentCount}}+ successful students?"

Show:
- Diverse students in learning environment
- Quick cuts of campus life
- Glimpses of state-of-the-art facilities
- Uplifting, inspiring tone

Goal: Make viewer want to learn more about {{.InstitutionName}}.`,
		},
		{
			Name: "achievement_highlight",
			Pattern: `LEAD WITH SUCCESS: Open with biggest achievement.
			
Statement: "{{.InstitutionName}} - Ranked #{{.Ranking}} {{.FieldOfStudy}} Program in {{.Region}}"
Alternative: "Where {{.EmploymentRate}}% of graduates find jobs within 6 months"

Show:
- Prestigious campus architecture or modern facilities
- Evidence of excellence (labs, research equipment)
- Student success stories in action
- Achievement badges/rankings

Immediately establish credibility and prestige.`,
		},
		{
			Name: "story_beginning",
			Pattern: `START WITH CHARACTER: Introduce a relatable student character.
			
Open: "Meet {{.StudentCharacter}} - she came with a dream of studying {{.FieldOfStudy}}"
Alternative: "3 years ago, {{.StudentCharacter}} stood where you are now - uncertain but hopeful"

Show:
- Diverse student character (first-generation, career-changer, athlete, international)
- Candid moment on campus
- Warm, inviting institution atmosphere
- Emotional connection

Set tone: transformative, inspiring, achievable.`,
		},
		{
			Name: "problem_solution",
			Pattern: `IDENTIFY & SOLVE: Start with relatable student problem.
			
Problem: "Worried about finding the right {{.FieldOfStudy}} program?"
Solution: "{{.InstitutionName}} offers hands-on learning with world-class faculty"

Show:
- Student looking uncertain
- Quick transformation through learning
- Access to resources and opportunities
- Clear path to success

Immediately present {{.InstitutionName}} as the solution.`,
		},
	}
}

// VALUE TEMPLATES - Designed to showcase institutional strengths and unique offerings
func getValueTemplates() []ValueTemplate {
	return []ValueTemplate{
		{
			Name: "rankings_and_achievements",
			Pattern: `SHOWCASE INSTITUTIONAL EXCELLENCE:

Academic Rankings & Accreditation:
- Ranked #{{.Ranking}} {{.FieldOfStudy}} Program in {{.Region}}
- Accredited by {{.Accreditation}}
- {{.AwardCount}}+ international awards and recognitions

Faculty Expertise:
- {{.PhDCount}}+ PhD holders
- Active researchers at top conferences
- Industry practitioners teaching courses

Campus Facilities:
- State-of-the-art {{.FacilityType}} labs
- Innovation hubs and maker spaces
- Modern library with digital resources

Show graphics with rankings, certificates, and facility photos.`,
		},
		{
			Name: "career_outcomes",
			Pattern: `HIGHLIGHT SUCCESS AFTER GRADUATION:

Employment & Career:
- {{.EmploymentRate}}% employment rate within 6 months
- {{.AlumniCount}}+ alumni working at Fortune 500 companies
- Average starting salary: {{.AverageSalary}}

Alumni Working At:
- Tech giants: {{.PartnerCompanies}}
- Global companies across {{.CountryCount}}+ countries
- Leadership positions in industry

Success Stories:
- {{.NotableAlumni}} [specific achievement]
- Graduates launching startups
- Alumni network spanning {{.NetworkSize}}+ professionals worldwide

Use company logos, statistics graphics, and success story quotes.`,
		},
		{
			Name: "student_life_community",
			Pattern: `SHOWCASE VIBRANT CAMPUS CULTURE:

Student Engagement:
- {{.ClubCount}}+ student clubs and organizations
- {{.EventCount}}+ campus events annually
- Active student government and leadership programs

Campus Experience:
- Diverse student body: {{.InternationalStudentPercentage}}% international students
- Modern residential facilities
- Student wellness and support services
- Sports, arts, and cultural activities

Community:
- Mentorship programs connecting students with professionals
- Peer learning groups and study communities
- Career development workshops and seminars
- Strong alumni network for networking

Show vibrant montage of student activities, clubs, campus life moments.`,
		},
		{
			Name: "unique_programs_offerings",
			Pattern: `HIGHLIGHT DISTINCTIVE PROGRAMS:

Specialized Programs:
- {{.UniqueProgram1}}: [specific details]
- {{.UniqueProgram2}}: [specific details]
- {{.UniqueProgram3}}: [specific details]

Hands-On Learning:
- Internship opportunities with {{.PartnerCompanies}}
- Research projects with real-world impact
- Capstone projects solving actual problems
- Industry partnerships and collaborations

Innovation & Entrepreneurship:
- Startup incubator program
- Innovation lab with cutting-edge technology
- Business mentorship and funding opportunities
- {{.StartupCount}}+ startups founded by alumni

Show students working on projects, using facilities, collaborating.`,
		},
		{
			Name: "diversity_and_inclusion",
			Pattern: `CELEBRATE DIVERSE COMMUNITY:

Global Community:
- {{.InternationalStudentPercentage}}% international students from {{.CountryCount}}+ countries
- Multicultural campus environment
- International exchange programs
- {{.LanguageCount}}+ languages spoken on campus

Inclusion & Support:
- Scholarships for underrepresented groups
- First-generation student support programs
- LGBTQ+ inclusive community
- Accessibility and disability support
- Cultural centers and affinity groups

Real Impact:
- Diverse student perspectives enhance learning
- Safe, welcoming community for all
- Mentorship from diverse faculty and staff
- Alumni network supporting equity and inclusion

Show diverse students collaborating, celebrating together, being themselves.`,
		},
		{
			Name: "scholarship_and_financial_aid",
			Pattern: `MAKE EDUCATION ACCESSIBLE:

Financial Support:
- {{.ScholarshipPercentage}}% of students receive some form of aid
- {{.FullScholarshipCount}} full scholarships available
- Merit-based scholarships: {{.MeritScholarshipAmount}}+
- Need-based financial aid available

Scholarship Categories:
- Academic excellence scholarships
- Sports and arts scholarships
- First-generation student scholarships
- International student scholarships
- Career-specific scholarships for {{.FieldOfStudy}}

Application Support:
- Simple, transparent application process
- Financial aid advisors available for guidance
- Payment plan options
- Work-study opportunities on campus

Show diverse students benefiting from scholarships, text about amounts and availability.`,
		},
	}
}

// CTA TEMPLATES - Designed to drive enrollment and action
func getCTATemplates() []CTATemplate {
	return []CTATemplate{
		{
			Name: "direct_application_cta",
			Pattern: `DRIVE IMMEDIATE ACTION:

Primary CTA:
- "Apply Now at {{.AdmissionURL}}"
- "Limited Spots Available"
- "Application Deadline: {{.DeadlineDate}}"

Create Urgency:
- Early application closes: {{.EarlyDeadlineDate}}
- Early scholarship deadline: {{.ScholarshipDeadlineDate}}
- Open house dates: {{.OpenHouseDates}}

Contact Information:
- Website: {{.AdmissionURL}}
- Phone: {{.PhoneNumber}}
- Email: {{.EmailAddress}}
- QR code to application portal

Show bold text overlay and logo with contact details prominently.`,
		},
		{
			Name: "enrollment_motivation_cta",
			Pattern: `INSPIRE & MOTIVATE:

Key Message:
- "Transform Your Future at {{.InstitutionName}}"
- "Join {{.StudentCount}}+ Students Building Their Success Stories"
- "Your {{.FieldOfStudy}} Journey Starts Here"

Take Action:
- "Apply now - your future awaits"
- "Early applicants get priority scholarship consideration"
- "Full scholarships for top {{.ScholarshipCount}} students"

Next Steps:
- Visit admissions website: {{.AdmissionURL}}
- Schedule a campus tour
- Attend webinar: {{.WebinarInfo}}
- Chat with admissions counselor

End with inspiring campus image and {{.InstitutionName}} logo.`,
		},
		{
			Name: "scholarship_focused_cta",
			Pattern: `EMPHASIZE FINANCIAL SUPPORT:

Scholarship Opportunity:
- "Full Scholarships Available for Top {{.ScholarshipCount}} Students"
- "{{.ScholarshipPercentage}}% of Our Students Receive Financial Aid"
- "Make {{.InstitutionName}} Affordable"

Scholarship Types:
- Merit-based (based on academics): {{.MeritScholarshipAmount}}+
- Need-based: Available for eligible students
- {{.FieldOfStudy}}-specific scholarships
- First-generation student scholarships

Action:
- "Learn about scholarship opportunities"
- "Apply by {{.ScholarshipDeadlineDate}} for full scholarship consideration"
- "Talk to financial aid advisor"
- "Calculate your aid package: {{.FinancialAidURL}}"

Show students with scholarship awards, text about amounts.`,
		},
		{
			Name: "event_invitation_cta",
			Pattern: `INVITE TO CAMPUS:

Event Invitation:
- "Join Our Open House"
- "Campus Tours Available {{.OpenHouseDates}}"
- "Experience {{.InstitutionName}} Firsthand"

What to Expect:
- Meet current students and faculty
- Tour state-of-the-art facilities
- Learn about academic programs
- Chat with admissions counselors
- Complimentary lunch and campus merchandise

How to Register:
- Visit: {{.EventRegistrationURL}}
- Call: {{.PhoneNumber}}
- Register online: {{.EventRegistrationURL}}
- Spots filling up fast!

Show excited students on campus, group tours, welcome atmosphere.`,
		},
		{
			Name: "information_request_cta",
			Pattern: `START THE CONVERSATION:

Get More Information:
- "Request {{.InstitutionName}} Prospectus"
- "Download Program Guide"
- "Subscribe to Student Stories Newsletter"

Why Connect:
- Personalized program recommendations
- Scholarship eligibility information
- Admission timeline and requirements
- Student success stories
- Campus life insights

Easy Ways to Learn More:
- Chat with admissions advisor
- Fill out quick information form
- Schedule 1-on-1 virtual meeting
- Join virtual webinar: {{.WebinarInfo}}

Contact:
- {{.AdmissionEmail}}
- {{.PhoneNumber}}
- {{.AdmissionURL}}

Show friendly admissions staff, welcoming environment.`,
		},
		{
			Name: "alumni_connection_cta",
			Pattern: `CONNECT WITH ALUMNI SUCCESS:

Connect with Our Alumni:
- "See Where {{.InstitutionName}} Graduates Are Now"
- "200,000+ Alumni Network Worldwide"
- "Join a Community of Achievers"

Alumni Success:
- Working at {{.PartnerCompanies}}
- Leaders in {{.IndustryCount}}+ industries
- Making impact globally
- Mentoring current students

What Alumni Say:
- [Quote 1]: "{{.AlumniQuote1}}"
- [Quote 2]: "{{.AlumniQuote2}}"
- Testimonials from successful graduates

Alumni Mentorship:
- Current students paired with alumni mentors
- Career guidance and networking
- Job placement support
- Lifelong connection

Show alumni in professional settings, alumni testimonials.`,
		},
	}
}
