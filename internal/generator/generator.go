package generator

import "cv-craft/internal/parser"

// Generator define a interface para todos os tipos de geradores de currículo
type Generator interface {
	Generate(cv *parser.Curriculum, outputPath string) error
}
