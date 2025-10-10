package generator

import (
	"fmt"
	"os"
	"strings"

	"cv-craft/internal/parser"
)

type MarkdownGenerator struct{}

func (g *MarkdownGenerator) Generate(cv *parser.Curriculum, outputPath string) error {
	var builder strings.Builder

	// Contato
	builder.WriteString(fmt.Sprintf("# %s\n", cv.Contact.Name))
	builder.WriteString(fmt.Sprintf("%s | %s | [%s](https://%s) | [Portfólio](%s)\n\n", cv.Contact.Email, cv.Contact.Phone, "LinkedIn", cv.Contact.LinkedIn, cv.Contact.Portfolio))

	// Resumo
	builder.WriteString("## Resumo\n")
	builder.WriteString(cv.Summary + "\n\n")

	// Habilidades
	builder.WriteString("## Habilidades\n")
	for _, skill := range cv.Skills {
		builder.WriteString(fmt.Sprintf("**%s:** %s\n\n", skill.Name, strings.Join(skill.Keywords, ", ")))
	}

	// Experiência
	builder.WriteString("## Experiência\n")
	for _, exp := range cv.Experience {
		builder.WriteString(fmt.Sprintf("### %s | %s\n", exp.Role, exp.Company))
		builder.WriteString(fmt.Sprintf("*%s*\n\n", exp.Period))
		builder.WriteString(exp.Description + "\n\n")
	}

	// Educação
	builder.WriteString("## Educação\n")
	for _, edu := range cv.Education {
		builder.WriteString(fmt.Sprintf("**%s** - %s\n", edu.Degree, edu.Institution))
	}

	return os.WriteFile(outputPath, []byte(builder.String()), 0644)
}
