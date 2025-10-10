package generator

import (
	"log"
	"strings"

	"cv-craft/internal/parser"

	"github.com/go-pdf/fpdf"
)

type PDFGenerator struct {
	IsATS bool
}

func (g *PDFGenerator) Generate(cv *parser.Curriculum, outputPath string) error {
	pdf := fpdf.New("P", "mm", "A4", "")

	// Professional ATS-friendly settings
	pdf.SetHeaderFunc(nil)
	pdf.SetFooterFunc(nil)
	pdf.SetAutoPageBreak(true, 20)
	pdf.SetMargins(20, 15, 20)

	log.Println("Generating professional ATS-optimized resume...")
	return g.generateProfessionalResume(pdf, cv, outputPath)
}

func (g *PDFGenerator) generateProfessionalResume(pdf *fpdf.Fpdf, cv *parser.Curriculum, outputPath string) error {
	pdf.AddPage()

	// ===== HEADER SECTION =====
	g.addHeaderSection(pdf, cv)

	// ===== SUMMARY SECTION =====
	g.addSummarySection(pdf, cv)

	// ===== SKILLS SECTION =====
	g.addSkillsSectionProfessional(pdf, cv)

	// ===== PROFESSIONAL EXPERIENCE SECTION =====
	g.addExperienceSectionProfessional(pdf, cv)

	// ===== EDUCATION SECTION =====
	g.addEducationSectionProfessional(pdf, cv)

	// ===== CERTIFICATES SECTION (se existir) =====
	if len(cv.Certificates) > 0 {
		g.addCertificatesSection(pdf, cv)
	}

	// ===== LANGUAGES SECTION (se existir) =====
	if len(cv.Languages) > 0 {
		g.addLanguagesSection(pdf, cv)
	}

	return pdf.OutputFileAndClose(outputPath)
}

// ===== SECTION: HEADER =====
func (g *PDFGenerator) addHeaderSection(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	// MAIN NAME
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, cv.Contact.Name, "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// PROFESSIONAL TITLE
	if cv.ProfessionalTitle != "" {
		pdf.SetFont("Helvetica", "B", 12)
		pdf.SetTextColor(70, 70, 70)
		pdf.CellFormat(0, 6, cv.ProfessionalTitle, "", 1, "L", false, 0, "")
		pdf.Ln(2)
	}

	// CONTACT INFORMATION - Formato profissional
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(80, 80, 80)

	// Linha 1: Location, Phone, Email
	line1 := []string{}
	if cv.Contact.Location != "" {
		line1 = append(line1, "Location: "+cv.Contact.Location)
	}
	if cv.Contact.Phone != "" {
		line1 = append(line1, "Phone: "+cv.Contact.Phone)
	}
	if cv.Contact.Email != "" {
		line1 = append(line1, "Email: "+cv.Contact.Email)
	}

	if len(line1) > 0 {
		pdf.CellFormat(0, 5, strings.Join(line1, " | "), "", 1, "L", false, 0, "")
	}

	// Linha 2: LinkedIn, Portfolio
	line2 := []string{}
	if cv.Contact.LinkedIn != "" {
		// Garantir que o LinkedIn tenha a URL completa
		linkedin := cv.Contact.LinkedIn
		if !strings.HasPrefix(linkedin, "http") {
			linkedin = "https://" + linkedin
		}
		line2 = append(line2, "LinkedIn: "+linkedin)
	}
	if cv.Contact.Portfolio != "" {
		portfolio := cv.Contact.Portfolio
		if !strings.HasPrefix(portfolio, "http") {
			portfolio = "https://" + portfolio
		}
		line2 = append(line2, "Portfolio: "+portfolio)
	}

	if len(line2) > 0 {
		pdf.CellFormat(0, 5, strings.Join(line2, " | "), "", 1, "L", false, 0, "")
	}

	pdf.Ln(8)
}

// ===== SECTION: SUMMARY =====
func (g *PDFGenerator) addSummarySection(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	if cv.Summary == "" {
		return
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "SUMMARY", "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(60, 60, 60)
	pdf.MultiCell(0, 5, cv.Summary, "", "J", false)
	pdf.Ln(8)
}

// ===== SECTION: SKILLS =====
func (g *PDFGenerator) addSkillsSectionProfessional(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	if len(cv.Skills) == 0 {
		return
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "SKILLS", "", 1, "L", false, 0, "")

	// Agrupar skills por nível
	proficientSkills := g.getSkillsByLevel(cv.Skills, "proficient")
	intermediateSkills := g.getSkillsByLevel(cv.Skills, "intermediate")
	beginnerSkills := g.getSkillsByLevel(cv.Skills, "beginner")

	// Proficient Skills
	if len(proficientSkills) > 0 {
		g.addSkillLevel(pdf, "Proficient:", proficientSkills)
	}

	// Intermediate Skills
	if len(intermediateSkills) > 0 {
		g.addSkillLevel(pdf, "Intermediate:", intermediateSkills)
	}

	// Beginner Skills
	if len(beginnerSkills) > 0 {
		g.addSkillLevel(pdf, "Beginner:", beginnerSkills)
	}

	pdf.Ln(5)
}

// ===== SECTION: PROFESSIONAL EXPERIENCE =====
func (g *PDFGenerator) addExperienceSectionProfessional(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	if len(cv.Experience) == 0 {
		return
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "PROFESSIONAL EXPERIENCE", "", 1, "L", false, 0, "")

	for i, exp := range cv.Experience {
		// Verificar se precisa de nova página
		if pdf.GetY() > 250 {
			pdf.AddPage()
			pdf.SetY(20)
		}

		// ROLE
		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 6, exp.Role, "", 1, "L", false, 0, "")

		// COMPANY, LOCATION e PERÍODO
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(70, 70, 70)

		companyInfo := exp.Company
		if exp.Location != "" {
			companyInfo += ", " + exp.Location
		}
		if exp.WorkType != "" {
			companyInfo += " -- " + exp.WorkType
		}
		if exp.Period != "" {
			companyInfo += " | " + exp.Period
		}

		pdf.CellFormat(0, 5, companyInfo, "", 1, "L", false, 0, "")

		// Core Responsibilities
		if len(exp.Responsibilities) > 0 {
			pdf.SetFont("Helvetica", "B", 10)
			pdf.SetTextColor(50, 50, 50)
			pdf.CellFormat(0, 6, "Core Responsibilities:", "", 1, "L", false, 0, "")

			pdf.SetFont("Helvetica", "", 10)
			pdf.SetTextColor(60, 60, 60)

			for _, responsibility := range exp.Responsibilities {
				responsibility = strings.TrimSpace(responsibility)
				if responsibility != "" {
					pdf.CellFormat(5, 5, "", "", 0, "L", false, 0, "")
					pdf.MultiCell(0, 5, "- "+responsibility, "", "L", false)
				}
			}
		}

		// Key Technologies and Tools
		if len(exp.Technologies) > 0 {
			pdf.SetFont("Helvetica", "B", 10)
			pdf.SetTextColor(50, 50, 50)
			pdf.CellFormat(0, 6, "Key Technologies and Tools:", "", 1, "L", false, 0, "")

			pdf.SetFont("Helvetica", "", 10)
			pdf.SetTextColor(80, 80, 80)
			pdf.MultiCell(0, 5, strings.Join(exp.Technologies, ", "), "", "L", false)
		}

		// Espaço entre experiências
		if i < len(cv.Experience)-1 {
			pdf.Ln(8)
		}
	}

	pdf.Ln(5)
}

// ===== SECTION: EDUCATION =====
func (g *PDFGenerator) addEducationSectionProfessional(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	if len(cv.Education) == 0 {
		return
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "EDUCATION", "", 1, "L", false, 0, "")

	for _, edu := range cv.Education {
		// Degree
		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 6, edu.Degree, "", 1, "L", false, 0, "")

		// Institution e Localização
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(80, 80, 80)

		institutionInfo := edu.Institution
		if edu.Location != "" {
			institutionInfo += ", " + edu.Location
		}
		if edu.Period != "" {
			institutionInfo += " | " + edu.Period
		}

		pdf.CellFormat(0, 5, institutionInfo, "", 1, "L", false, 0, "")

		// Thesis (se existir)
		if edu.Thesis != "" {
			pdf.SetFont("Helvetica", "I", 9)
			pdf.SetTextColor(100, 100, 100)
			pdf.CellFormat(0, 5, "Thesis: "+edu.Thesis, "", 1, "L", false, 0, "")
		}

		pdf.Ln(3)
	}
}

// ===== SECTION: CERTIFICATES =====
func (g *PDFGenerator) addCertificatesSection(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "CERTIFICATES", "", 1, "L", false, 0, "")

	for _, cert := range cv.Certificates {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(40, 40, 40)
		pdf.CellFormat(0, 6, cert.Name, "", 1, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(80, 80, 80)

		certInfo := cert.Institution
		if cert.Date != "" {
			certInfo += " | " + cert.Date
		}

		pdf.CellFormat(0, 5, certInfo, "", 1, "L", false, 0, "")
		pdf.Ln(2)
	}

	pdf.Ln(5)
}

// ===== SECTION: LANGUAGES =====
func (g *PDFGenerator) addLanguagesSection(pdf *fpdf.Fpdf, cv *parser.Curriculum) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, "LANGUAGES", "", 1, "L", false, 0, "")

	for _, lang := range cv.Languages {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(40, 40, 40)
		pdf.CellFormat(40, 6, lang.Language+":", "", 0, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(0, 6, lang.Level, "", 1, "L", false, 0, "")
	}

	pdf.Ln(5)
}

// ===== HELPER FUNCTIONS =====

func (g *PDFGenerator) getSkillsByLevel(skills []parser.Skill, level string) []string {
	var result []string
	for _, skill := range skills {
		if strings.EqualFold(skill.Level, level) {
			result = append(result, skill.Keywords...)
		}
	}
	return result
}

func (g *PDFGenerator) addSkillLevel(pdf *fpdf.Fpdf, levelTitle string, skills []string) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(40, 40, 40)
	pdf.CellFormat(0, 6, levelTitle, "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.MultiCell(0, 5, strings.Join(skills, ", "), "", "L", false)
	pdf.Ln(2)
}
