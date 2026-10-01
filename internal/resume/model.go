// Package resume define o modelo do currículo, o parsing do YAML e a validação
// (spec 001).
package resume

// Meta contém configurações do documento que não são conteúdo do currículo.
type Meta struct {
	Locale string `yaml:"locale"`
}

// Contact representa as informações de contato.
type Contact struct {
	Name      string `yaml:"name"`
	Email     string `yaml:"email"`
	Phone     string `yaml:"phone"`
	Location  string `yaml:"location"`
	LinkedIn  string `yaml:"linkedin"`
	GitHub    string `yaml:"github"`
	Portfolio string `yaml:"portfolio"`
}

// Skill é uma categoria de habilidades com nível opcional.
type Skill struct {
	Name     string   `yaml:"name"`
	Level    string   `yaml:"level"`
	Keywords []string `yaml:"keywords"`
}

// Experience representa uma experiência profissional.
type Experience struct {
	Role             string   `yaml:"role"`
	Company          string   `yaml:"company"`
	Location         string   `yaml:"location"`
	Period           string   `yaml:"period"`
	WorkType         string   `yaml:"work_type"`
	Description      string   `yaml:"description"`
	Responsibilities []string `yaml:"responsibilities"`
	Achievements     []string `yaml:"achievements"`
	Technologies     []string `yaml:"technologies"`
}

// Education representa uma formação acadêmica.
type Education struct {
	Degree          string `yaml:"degree"`
	Institution     string `yaml:"institution"`
	Location        string `yaml:"location"`
	Period          string `yaml:"period"`
	Thesis          string `yaml:"thesis"`
	RelevantCourses string `yaml:"relevant_courses"`
}

// Certificate representa um certificado ou curso.
type Certificate struct {
	Name        string `yaml:"name"`
	Institution string `yaml:"institution"`
	Date        string `yaml:"date"`
	URL         string `yaml:"url"`
}

// Language representa a proficiência em um idioma.
type Language struct {
	Language string `yaml:"language"`
	Level    string `yaml:"level"`
}

// Resume é o currículo completo, fonte única da verdade para todos os formatos.
type Resume struct {
	Meta              Meta          `yaml:"meta"`
	Contact           Contact       `yaml:"contact"`
	ProfessionalTitle string        `yaml:"professional_title"`
	Summary           string        `yaml:"summary"`
	Skills            []Skill       `yaml:"skills"`
	Experience        []Experience  `yaml:"experience"`
	Education         []Education   `yaml:"education"`
	Certificates      []Certificate `yaml:"certificates"`
	Languages         []Language    `yaml:"languages"`
}
