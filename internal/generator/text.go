// internal/generator/text.go
package generator

import (
	"fmt"
	"os"
	"strings"

	"cv-craft/internal/parser"
)

type TextGenerator struct{}

func (g *TextGenerator) Generate(cv *parser.Curriculum, outputPath string) error {
	var builder strings.Builder

	// Contato
	builder.WriteString(cv.Contact.Name + "\n")
	builder.WriteString(fmt.Sprintf("%s | %s | %s | %s\n", cv.Contact.Email, cv.Contact.Phone, cv.Contact.LinkedIn, cv.Contact.Portfolio))
	builder.WriteString("\n" + strings.Repeat("-", 40) + "\n\n")

	// Resumo
	builder.WriteString("RESUMO\n")
	builder.WriteString(cv.Summary + "\n\n")

	// Habilidades
	builder.WriteString("HABILIDADES\n")
	for _, skill := range cv.Skills {
		builder.WriteString(fmt.Sprintf("- %s: %s\n", skill.Name, strings.Join(skill.Keywords, ", ")))
	}
	builder.WriteString("\n")

	// Experiência
	builder.WriteString("EXPERIÊNCIA\n")
	for _, exp := range cv.Experience {
		builder.WriteString(fmt.Sprintf("%s | %s (%s)\n", exp.Role, exp.Company, exp.Period))
		builder.WriteString(exp.Description + "\n\n")
	}

	// Educação
	builder.WriteString("EDUCAÇÃO\n")
	for _, edu := range cv.Education {
		builder.WriteString(fmt.Sprintf("%s - %s\n", edu.Degree, edu.Institution))
	}

	return os.WriteFile(outputPath, []byte(builder.String()), 0644)
}
