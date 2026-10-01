package generator

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/fabiodrneles/cv-craft/internal/i18n"
	"github.com/fabiodrneles/cv-craft/internal/resume"
)

var update = flag.Bool("update", false, "regrava os golden files em testdata/")

var fixedDate = time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

func loadExample(t *testing.T, name string) *resume.Resume {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "examples", name))
	if err != nil {
		t.Fatal(err)
	}
	r, _, err := resume.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func generate(t *testing.T, f Format, r *resume.Resume) []byte {
	t.Helper()
	labels, err := i18n.Get(r.Meta.Locale)
	if err != nil {
		t.Fatal(err)
	}
	g, err := New(f, Options{Labels: labels, Version: "test", CreatedAt: fixedDate})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := g.Generate(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// havePoppler informa se o pdftotext (poppler-utils) está instalado. No CI,
// CV_CRAFT_REQUIRE_POPPLER=1 transforma a ausência em falha.
func havePoppler(t *testing.T) bool {
	t.Helper()
	if _, err := exec.LookPath("pdftotext"); err == nil {
		return true
	}
	if os.Getenv("CV_CRAFT_REQUIRE_POPPLER") != "" {
		t.Fatal("pdftotext não encontrado e CV_CRAFT_REQUIRE_POPPLER está definido")
	}
	t.Log("pdftotext (poppler-utils) não encontrado; verificações do texto do PDF puladas")
	return false
}

// requirePoppler pula o teste se o pdftotext não estiver instalado.
func requirePoppler(t *testing.T) {
	t.Helper()
	if !havePoppler(t) {
		t.Skip("pdftotext (poppler-utils) não encontrado")
	}
}

// outputText gera o formato e devolve o texto visível; ok é false para PDF
// quando o poppler não está disponível.
func outputText(t *testing.T, f Format, r *resume.Resume) (text string, ok bool) {
	t.Helper()
	out := generate(t, f, r)
	if f != PDF {
		return string(out), true
	}
	if !havePoppler(t) {
		return "", false
	}
	return pdfText(t, out), true
}

func poppler(t *testing.T, tool string, data []byte, args ...string) string {
	t.Helper()
	requirePoppler(t)
	path := filepath.Join(t.TempDir(), "cv.pdf")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	args = append(append([]string{"-enc", "UTF-8"}, args...), path)
	if tool == "pdftotext" {
		args = append(args, "-") // escreve na saída padrão
	}
	out, err := exec.Command(tool, args...).Output()
	if err != nil {
		t.Fatalf("%s: %v", tool, err)
	}
	return string(out)
}

// pdfPages extrai o texto de cada página com o pdftotext, sem linhas em branco.
func pdfPages(t *testing.T, data []byte) [][]string {
	t.Helper()
	raw := poppler(t, "pdftotext", data)
	var pages [][]string
	for _, page := range strings.Split(strings.TrimRight(raw, "\f\n"), "\f") {
		var lines []string
		for _, l := range strings.Split(page, "\n") {
			if l = strings.TrimSpace(l); l != "" {
				lines = append(lines, l)
			}
		}
		pages = append(pages, lines)
	}
	return pages
}

func pdfText(t *testing.T, data []byte) string {
	t.Helper()
	var all []string
	for _, p := range pdfPages(t, data) {
		all = append(all, p...)
	}
	return strings.Join(all, "\n") + "\n"
}

func golden(t *testing.T, name string, got []byte, canon func(string) string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (rode: go test ./internal/generator -update)", err)
	}
	if canon == nil {
		canon = func(s string) string { return s }
	}
	if canon(string(got)) != canon(string(want)) {
		t.Errorf("%s difere do golden file; se a mudança é intencional, rode go test ./internal/generator -update\n--- obtido ---\n%s", name, got)
	}
}

// 003/AC-7, 004/AC-4.
func TestGolden(t *testing.T) {
	for _, ex := range []string{"full", "en", "minimal"} {
		r := loadExample(t, ex+".yaml")
		t.Run(ex, func(t *testing.T) {
			golden(t, ex+".md", generate(t, Markdown, r), nil)
			golden(t, ex+".txt", generate(t, Text, r), nil)
			// O texto extraído do PDF é comparado com espaços normalizados:
			// versões diferentes do poppler podem agrupar linhas de outro jeito.
			golden(t, ex+".pdf.txt", []byte(pdfText(t, generate(t, PDF, r))), normalize)
		})
	}
}

var spaces = regexp.MustCompile(`\s+`)

func normalize(s string) string { return strings.TrimSpace(spaces.ReplaceAllString(s, " ")) }

// contentStrings devolve todo texto preenchido no YAML que deve aparecer nas
// saídas, já no formato em que é exibido.
func contentStrings(r *resume.Resume) []string {
	var out []string
	var walk func(v reflect.Value, path string)
	walk = func(v reflect.Value, path string) {
		switch v.Kind() {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				tag, _, _ := strings.Cut(v.Type().Field(i).Tag.Get("yaml"), ",")
				walk(v.Field(i), path+"."+tag)
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i), path+"[]")
			}
		case reflect.String:
			s := normalize(v.String())
			switch path {
			case ".meta.locale", ".skills[].level": // traduzidos/omitidos
				return
			case ".contact.linkedin", ".contact.github", ".contact.portfolio", ".certificates[].url":
				s = resume.NormalizeURL(s)
			}
			if s != "" {
				out = append(out, s)
			}
		}
	}
	walk(reflect.ValueOf(*r), "")
	return out
}

// 004/AC-1, 004/AC-5, 003/NFR-1: nenhum campo preenchido some em nenhum formato.
func TestContentParity(t *testing.T) {
	for _, ex := range []string{"full.yaml", "en.yaml", "minimal.yaml"} {
		r := loadExample(t, ex)
		for _, f := range Formats {
			out, ok := outputText(t, f, r)
			if !ok {
				continue
			}
			if f == Markdown {
				out = strings.ReplaceAll(out, `\`, "")
			}
			out = normalize(out)
			for _, s := range contentStrings(r) {
				if f == PDF {
					s = displayURL(s) // o PDF exibe links sem "https://"
				}
				if !strings.Contains(out, s) {
					t.Errorf("%s/%s: conteúdo ausente: %q", ex, f, s)
				}
			}
		}
	}
}

// 003/AC-1, AC-2: UTF-8 correto no PDF.
func TestPDFUnicode(t *testing.T) {
	r := loadExample(t, "minimal.yaml")
	r.Contact.Location = "São Paulo"
	r.Summary = "Ação, coração, café, señor, über, Straße, Øresund, Łódź."
	text := pdfText(t, generate(t, PDF, r))
	for _, s := range []string{"São Paulo", "Ação, coração, café, señor, über, Straße, Øresund, Łódź."} {
		if !strings.Contains(text, s) {
			t.Errorf("PDF não contém %q:\n%s", s, text)
		}
	}
}

// 003/AC-4 e 001/AC-3/AC-4: skills com qualquer nível nunca são descartadas.
func TestSkillsNeverDropped(t *testing.T) {
	r := loadExample(t, "minimal.yaml")
	r.Skills = []resume.Skill{
		{Name: "A", Level: "Advanced", Keywords: []string{"Rust"}},
		{Name: "B", Level: "guru", Keywords: []string{"Zig"}},
		{Name: "C", Keywords: []string{"Elixir"}},
	}
	for _, f := range Formats {
		out, ok := outputText(t, f, r)
		if !ok {
			continue
		}
		for _, kw := range []string{"Rust", "Zig", "Elixir", "A (Avançado)"} {
			if !strings.Contains(out, kw) {
				t.Errorf("%s: faltou %q", f, kw)
			}
		}
		if strings.Contains(out, "guru") {
			t.Errorf("%s: nível desconhecido não deveria ser exibido", f)
		}
	}
}

// 004/AC-2: campos vazios não deixam separadores nem links vazios.
func TestEmptyContactFields(t *testing.T) {
	r := loadExample(t, "minimal.yaml")
	r.Contact = resume.Contact{Name: "Ana", Email: "ana@example.com"}
	for _, f := range []Format{Markdown, Text} {
		out := string(generate(t, f, r))
		for _, bad := range []string{"|  |", "| |", "](https://)", "LinkedIn", "GitHub", "Portfólio"} {
			if strings.Contains(out, bad) {
				t.Errorf("%s contém %q:\n%s", f, bad, out)
			}
		}
	}
}

// 004/AC-3: texto do usuário é escapado no Markdown.
func TestMarkdownEscaping(t *testing.T) {
	r := loadExample(t, "minimal.yaml")
	r.Summary = "Uso *Go* e [C] com_underscore"
	r.Experience[0].Responsibilities = []string{"- começa com hífen", "1. começa com número"}
	out := string(generate(t, Markdown, r))
	for _, want := range []string{`Uso \*Go\* e \[C\] com\_underscore`, `- \- começa com hífen`, `- 1\. começa com número`} {
		if !strings.Contains(out, want) {
			t.Errorf("Markdown não contém %q:\n%s", want, out)
		}
	}
}

// 006/AC-1/AC-2: títulos no idioma pedido em todos os formatos.
func TestLabelsFollowLocale(t *testing.T) {
	r := loadExample(t, "full.yaml")
	for code, want := range map[string][]string{
		"pt-BR": {"Resumo", "Experiência Profissional", "Formação"},
		"en":    {"Summary", "Professional Experience", "Education"},
	} {
		r.Meta.Locale = code
		for _, f := range Formats {
			out, ok := outputText(t, f, r)
			if !ok {
				continue
			}
			for _, w := range want {
				if !strings.Contains(strings.ToLower(out), strings.ToLower(w)) {
					t.Errorf("%s/%s: faltou %q", code, f, w)
				}
			}
		}
	}
}

// 003/AC-5: metadados do PDF.
func TestPDFMetadata(t *testing.T) {
	info := poppler(t, "pdfinfo", generate(t, PDF, loadExample(t, "full.yaml")))
	for _, want := range []string{
		"Title:           Carlos Silva — Engenheiro de Software Sênior | Go e Cloud",
		"Author:          Carlos Silva",
		"Creator:         cv-craft test",
		"Subject:         Currículo — Carlos Silva",
	} {
		if !strings.Contains(info, want+"\n") {
			t.Errorf("pdfinfo sem %q:\n%s", want, info)
		}
	}
	if !regexp.MustCompile(`Keywords:.*Kubernetes`).MatchString(info) {
		t.Errorf("Keywords sem as habilidades:\n%s", info)
	}
}

// 003/AC-8: mesma entrada e mesma data geram bytes idênticos.
func TestPDFDeterministic(t *testing.T) {
	r := loadExample(t, "full.yaml")
	a, b := generate(t, PDF, r), generate(t, PDF, r)
	if !bytes.Equal(a, b) {
		t.Error("PDF não é determinístico")
	}
}

// 003/AC-6: título de seção, rótulo de lista ou cabeçalho de experiência
// nunca é a última linha de uma página. O teste desloca o conteúdo linha a
// linha para que as quebras de página caiam em todas as posições possíveis.
func TestPDFPageBreaks(t *testing.T) {
	labels, _ := i18n.Get("pt-BR")
	base := loadExample(t, "full.yaml")
	for shift := 0; shift < 40; shift++ {
		r := *base
		r.Experience = nil
		for i := 0; i < 4; i++ {
			for _, e := range base.Experience {
				e.Role = fmt.Sprintf("%s %d", e.Role, i)
				r.Experience = append(r.Experience, e)
			}
		}
		r.Summary = strings.Repeat("Texto de deslocamento. ", 1+shift*5)

		orphans := map[string]bool{
			labels.Responsibilities + ":": true,
			labels.Achievements + ":":     true,
		}
		for _, h := range []string{labels.Summary, labels.Skills, labels.Experience, labels.Education, labels.Certificates, labels.Languages} {
			orphans[strings.ToUpper(h)] = true
		}
		for _, e := range r.Experience {
			orphans[titled(e.Role, e.Company)] = true
		}
		for _, e := range r.Education {
			orphans[titled(e.Degree, e.Institution)] = true
		}

		pages := pdfPages(t, generate(t, PDF, &r))
		for i, p := range pages[:len(pages)-1] {
			if last := p[len(p)-1]; orphans[last] {
				t.Errorf("shift %d: página %d termina com %q", shift, i+1, last)
			}
		}
	}
}

func TestParseFormat(t *testing.T) {
	for in, want := range map[string]Format{"pdf": PDF, "PDF": PDF, "md": Markdown, "markdown": Markdown, "txt": Text, "text": Text} {
		if got, err := ParseFormat(in); err != nil || got != want {
			t.Errorf("ParseFormat(%q) = %q, %v", in, got, err)
		}
	}
	if _, err := ParseFormat("docx"); err == nil {
		t.Error("docx deveria ser inválido")
	}
}

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, io.ErrShortWrite }

func TestWriteErrorsPropagate(t *testing.T) {
	r := loadExample(t, "minimal.yaml")
	for _, f := range Formats {
		g, _ := New(f, Options{CreatedAt: fixedDate})
		if err := g.Generate(failWriter{}, r); err == nil {
			t.Errorf("%s: erro de escrita deveria ser devolvido", f)
		}
	}
}
