// Package i18n contém os rótulos fixos usados nas saídas do currículo
// (spec 006). O conteúdo escrito pelo usuário nunca é traduzido.
package i18n

import (
	"fmt"
	"strings"
)

// Default é o idioma usado quando nem --lang nem meta.locale são informados.
const Default = "pt-BR"

// Labels são os textos fixos de um idioma.
type Labels struct {
	Code string

	ResumeOf         string // assunto do PDF
	Summary          string
	Skills           string
	Experience       string
	Education        string
	Certificates     string
	Languages        string
	Responsibilities string
	Achievements     string
	Technologies     string
	Thesis           string
	RelevantCourses  string
	Portfolio        string

	// Níveis de habilidade, indexados pelo nome canônico (ver resume.Levels).
	Levels map[string]string
}

var catalog = map[string]Labels{
	"pt-BR": {
		Code:             "pt-BR",
		ResumeOf:         "Currículo",
		Summary:          "Resumo",
		Skills:           "Habilidades",
		Experience:       "Experiência Profissional",
		Education:        "Formação",
		Certificates:     "Certificados",
		Languages:        "Idiomas",
		Responsibilities: "Responsabilidades",
		Achievements:     "Conquistas",
		Technologies:     "Tecnologias",
		Thesis:           "Trabalho de conclusão",
		RelevantCourses:  "Disciplinas relevantes",
		Portfolio:        "Portfólio",
		Levels: map[string]string{
			"expert":       "Especialista",
			"advanced":     "Avançado",
			"proficient":   "Proficiente",
			"intermediate": "Intermediário",
			"beginner":     "Básico",
		},
	},
	"en": {
		Code:             "en",
		ResumeOf:         "Resume",
		Summary:          "Summary",
		Skills:           "Skills",
		Experience:       "Professional Experience",
		Education:        "Education",
		Certificates:     "Certificates",
		Languages:        "Languages",
		Responsibilities: "Responsibilities",
		Achievements:     "Achievements",
		Technologies:     "Technologies",
		Thesis:           "Thesis",
		RelevantCourses:  "Relevant courses",
		Portfolio:        "Portfolio",
		Levels: map[string]string{
			"expert":       "Expert",
			"advanced":     "Advanced",
			"proficient":   "Proficient",
			"intermediate": "Intermediate",
			"beginner":     "Beginner",
		},
	},
}

// Supported lista os códigos de idioma disponíveis.
func Supported() []string { return []string{"pt-BR", "en"} }

// Get devolve os rótulos de um idioma. Aceita variações como "pt", "pt_br",
// "PT-BR", "en-US". Código vazio resulta no idioma padrão.
func Get(code string) (Labels, error) {
	c := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(code), "_", "-"))
	switch {
	case c == "":
		return catalog[Default], nil
	case c == "pt" || strings.HasPrefix(c, "pt-"):
		return catalog["pt-BR"], nil
	case c == "en" || strings.HasPrefix(c, "en-"):
		return catalog["en"], nil
	}
	return Labels{}, fmt.Errorf("idioma não suportado %q (suportados: %s)", code, strings.Join(Supported(), ", "))
}
