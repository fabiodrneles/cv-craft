package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Contact representa as informações de contato
type Contact struct {
	Name      string `yaml:"name"`
	Email     string `yaml:"email"`
	Phone     string `yaml:"phone"`
	LinkedIn  string `yaml:"linkedin"`
	Portfolio string `yaml:"portfolio"`
	Location  string `yaml:"location"`
}

// Skill representa uma habilidade com nível de proficiência
type Skill struct {
	Name     string   `yaml:"name"`
	Level    string   `yaml:"level"` // proficient, intermediate, beginner
	Keywords []string `yaml:"keywords"`
}

// Experience representa uma experiência profissional completa
type Experience struct {
	Role             string   `yaml:"role"`
	Company          string   `yaml:"company"`
	Location         string   `yaml:"location"`
	Period           string   `yaml:"period"`
	WorkType         string   `yaml:"work_type"` // Remote, Hybrid, On-site
	Description      string   `yaml:"description"`
	Responsibilities []string `yaml:"responsibilities"`
	Technologies     []string `yaml:"technologies"`
	Achievements     []string `yaml:"achievements"`
}

// Education representa uma formação acadêmica
type Education struct {
	Degree          string `yaml:"degree"`
	Institution     string `yaml:"institution"`
	Location        string `yaml:"location"`
	Period          string `yaml:"period"`
	Thesis          string `yaml:"thesis,omitempty"`
	RelevantCourses string `yaml:"relevant_courses,omitempty"`
}

// Certificate representa um certificado ou curso
type Certificate struct {
	Name        string `yaml:"name"`
	Institution string `yaml:"institution"`
	Date        string `yaml:"date"`
	URL         string `yaml:"url,omitempty"`
}

// Language representa proficiência em idiomas
type Language struct {
	Language string `yaml:"language"`
	Level    string `yaml:"level"` // Native, Fluent, Intermediate, Basic
}

// Curriculum representa a estrutura completa do currículo profissional
type Curriculum struct {
	Contact           Contact       `yaml:"contact"`
	ProfessionalTitle string        `yaml:"professional_title"`
	Summary           string        `yaml:"summary"`
	Skills            []Skill       `yaml:"skills"`
	Experience        []Experience  `yaml:"experience"`
	Education         []Education   `yaml:"education"`
	Certificates      []Certificate `yaml:"certificates,omitempty"`
	Languages         []Language    `yaml:"languages,omitempty"`
}

// ParseFile lê e decodifica o arquivo YAML do currículo
func ParseFile(filePath string) (*Curriculum, error) {
	// Verifica se o arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo não encontrado: %s", filePath)
	}

	// Lê o arquivo
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Verifica se o arquivo está vazio
	if len(data) == 0 {
		return nil, fmt.Errorf("arquivo vazio: %s", filePath)
	}

	// Tenta decodificar o YAML
	var cv Curriculum
	if err := yaml.Unmarshal(data, &cv); err != nil {
		return nil, fmt.Errorf("erro ao decodificar YAML: %w", err)
	}

	// Validações básicas da estrutura
	if cv.Contact.Name == "" {
		return nil, fmt.Errorf("campo 'name' é obrigatório no contato")
	}

	if cv.Contact.Email == "" {
		return nil, fmt.Errorf("campo 'email' é obrigatório no contato")
	}

	if cv.ProfessionalTitle == "" {
		return nil, fmt.Errorf("campo 'professional_title' é obrigatório")
	}

	return &cv, nil
}

// ValidateCurriculum realiza validações adicionais no currículo
func ValidateCurriculum(cv *Curriculum) error {
	if cv == nil {
		return fmt.Errorf("currículo não pode ser nulo")
	}

	if len(cv.Skills) == 0 {
		return fmt.Errorf("pelo menos uma habilidade deve ser fornecida")
	}

	if len(cv.Experience) == 0 {
		return fmt.Errorf("pelo menos uma experiência deve ser fornecida")
	}

	if len(cv.Education) == 0 {
		return fmt.Errorf("pelo menos uma educação deve ser fornecida")
	}

	return nil
}
