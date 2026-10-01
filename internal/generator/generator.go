// Package generator transforma um currículo validado em PDF, Markdown ou
// texto puro. Os três formatos usam os mesmos helpers deste arquivo para
// garantir paridade de conteúdo (spec 004, FR-1).
package generator

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fabiodrneles/cv-craft/internal/i18n"
	"github.com/fabiodrneles/cv-craft/internal/resume"
)

// Format é um formato de saída suportado.
type Format string

const (
	PDF      Format = "pdf"
	Markdown Format = "md"
	Text     Format = "txt"
)

// Formats lista todos os formatos, na ordem usada por --format all.
var Formats = []Format{PDF, Markdown, Text}

// ParseFormat aceita os nomes e apelidos de formato da CLI.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "pdf":
		return PDF, nil
	case "md", "markdown":
		return Markdown, nil
	case "txt", "text":
		return Text, nil
	}
	return "", fmt.Errorf("formato inválido %q (use pdf, md, txt ou all)", s)
}

// Ext é a extensão de arquivo do formato, com ponto.
func (f Format) Ext() string { return "." + string(f) }

// Options configura a geração.
type Options struct {
	Labels    i18n.Labels
	Version   string    // gravado nos metadados do PDF
	CreatedAt time.Time // data de criação do PDF; zero usa time.Now()
}

// Generator escreve um currículo em um formato.
type Generator interface {
	Generate(w io.Writer, r *resume.Resume) error
}

// New devolve o gerador do formato.
func New(f Format, opts Options) (Generator, error) {
	if opts.Labels.Code == "" {
		opts.Labels, _ = i18n.Get(i18n.Default)
	}
	switch f {
	case PDF:
		return &pdfGenerator{opts: opts}, nil
	case Markdown:
		return &markdownGenerator{opts: opts}, nil
	case Text:
		return &textGenerator{opts: opts}, nil
	}
	return nil, fmt.Errorf("formato não suportado: %q", f)
}

// ---- helpers de conteúdo compartilhados ----

type link struct {
	Label string
	URL   string
}

// contactParts devolve localização, telefone e e-mail, nessa ordem, sem vazios.
func contactParts(c resume.Contact) []string {
	return nonEmpty(c.Location, c.Phone, c.Email)
}

// profileLinks devolve LinkedIn, GitHub e portfólio com URLs normalizadas.
func profileLinks(c resume.Contact, l i18n.Labels) []link {
	var out []link
	for _, p := range []link{
		{"LinkedIn", c.LinkedIn},
		{"GitHub", c.GitHub},
		{l.Portfolio, c.Portfolio},
	} {
		if u := resume.NormalizeURL(p.URL); u != "" {
			out = append(out, link{p.Label, u})
		}
	}
	return out
}

// skillTitle devolve "Categoria (Nível)" ou só "Categoria".
func skillTitle(s resume.Skill, l i18n.Labels) string {
	name := strings.TrimSpace(s.Name)
	if lvl, ok := resume.NormalizeLevel(s.Level); ok && lvl != "" {
		return fmt.Sprintf("%s (%s)", name, l.Levels[lvl])
	}
	return name
}

// experienceMeta devolve "Local | Modalidade | Período", sem repetir a
// modalidade quando ela é igual ao local (ex.: "Remoto | Remoto").
func experienceMeta(e resume.Experience) string {
	workType := e.WorkType
	if strings.EqualFold(trim(workType), trim(e.Location)) {
		workType = ""
	}
	return strings.Join(nonEmpty(e.Location, workType, e.Period), " | ")
}

// educationMeta devolve "Local | Período".
func educationMeta(e resume.Education) string {
	return strings.Join(nonEmpty(e.Location, e.Period), " | ")
}

// certificateMeta devolve "Emissor | Data".
func certificateMeta(c resume.Certificate) string {
	return strings.Join(nonEmpty(c.Institution, c.Date), " | ")
}

// titled junta duas partes com travessão, omitindo a segunda se vazia.
func titled(a, b string) string {
	return strings.Join(nonEmpty(a, b), " — ")
}

func nonEmpty(items ...string) []string { return resume.NonBlank(items) }

func trim(s string) string { return strings.TrimSpace(s) }
