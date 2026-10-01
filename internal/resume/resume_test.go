package resume

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const validYAML = `
contact:
  name: "Ana"
  email: "ana@example.com"
professional_title: "Dev"
skills:
  - name: "Backend"
    keywords: ["Go"]
experience:
  - role: "Dev"
    company: "ACME"
education:
  - degree: "BSc"
    institution: "USP"
`

func load(t *testing.T, src string) (Report, error) {
	t.Helper()
	r, rep, err := Parse([]byte(src))
	if err == nil {
		rep.Merge(Validate(r))
	}
	return rep, err
}

func paths(issues []Issue) []string {
	var out []string
	for _, i := range issues {
		out = append(out, i.Path)
	}
	return out
}

func hasPath(issues []Issue, path string) bool {
	for _, i := range issues {
		if i.Path == path {
			return true
		}
	}
	return false
}

func TestValidMinimal(t *testing.T) {
	rep, err := load(t, validYAML)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.OK() || len(rep.Warnings) > 0 {
		t.Fatalf("esperava válido sem avisos, veio erros=%v avisos=%v", rep.Errors, rep.Warnings)
	}
}

// 001/AC-1: todos os erros juntos, com caminho.
func TestValidateAggregatesErrors(t *testing.T) {
	src := strings.Replace(validYAML, `  name: "Ana"`, "", 1)
	src = strings.Replace(src, "skills:\n  - name: \"Backend\"\n    keywords: [\"Go\"]\n", "", 1)
	rep, err := load(t, src)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"contact.name", "skills"} {
		if !hasPath(rep.Errors, p) {
			t.Errorf("faltou erro em %q; erros: %v", p, paths(rep.Errors))
		}
	}
}

// 001/AC-2: chave desconhecida com caminho e linha.
func TestUnknownKey(t *testing.T) {
	src := strings.Replace(validYAML, `    company: "ACME"`, "    company: \"ACME\"\n    responsabilities: [\"x\"]", 1)
	rep, err := load(t, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Errors) != 1 {
		t.Fatalf("esperava 1 erro, veio %v", rep.Errors)
	}
	e := rep.Errors[0]
	if e.Path != "experience[0].responsabilities" || e.Line != 12 || !strings.Contains(e.Message, "campo desconhecido") {
		t.Errorf("erro inesperado: %+v", e)
	}
	if got := e.String(); !strings.Contains(got, "(linha 12)") {
		t.Errorf("String() sem linha: %q", got)
	}
}

// 001/AC-3 e AC-4: níveis normalizados; desconhecido vira aviso.
func TestSkillLevels(t *testing.T) {
	for _, tc := range []struct {
		level    string
		warnings int
	}{
		{"Advanced", 0}, {"  EXPERT ", 0}, {"avançado", 0}, {"", 0}, {"guru", 1},
	} {
		src := strings.Replace(validYAML, `  - name: "Backend"`, "  - name: \"Backend\"\n    level: \""+tc.level+"\"", 1)
		rep, err := load(t, src)
		if err != nil {
			t.Fatal(err)
		}
		if !rep.OK() || len(rep.Warnings) != tc.warnings {
			t.Errorf("level %q: erros=%v avisos=%v", tc.level, rep.Errors, rep.Warnings)
		}
	}
}

// 001/AC-5.
func TestInvalidEmail(t *testing.T) {
	for _, email := range []string{"não-é-email", "Ana <ana@example.com>", "ana@"} {
		src := strings.Replace(validYAML, "ana@example.com", email, 1)
		rep, _ := load(t, src)
		if !hasPath(rep.Errors, "contact.email") {
			t.Errorf("e-mail %q deveria ser inválido", email)
		}
	}
}

// 001/AC-6.
func TestEmptyFile(t *testing.T) {
	for _, src := range []string{"", "   \n", "# só comentário\n", "\xef\xbb\xbf"} {
		if _, _, err := Parse([]byte(src)); !errors.Is(err, ErrEmpty) {
			t.Errorf("Parse(%q) = %v, esperava ErrEmpty", src, err)
		}
	}
}

// 001/AC-7.
func TestBlankRequiredField(t *testing.T) {
	src := strings.Replace(validYAML, `"Ana"`, `"   "`, 1)
	rep, _ := load(t, src)
	if !hasPath(rep.Errors, "contact.name") {
		t.Errorf("nome só com espaços deveria ser erro; erros: %v", paths(rep.Errors))
	}
}

func TestNestedRequiredFields(t *testing.T) {
	src := validYAML + `
certificates:
  - institution: "AWS"
languages:
  - level: "Fluente"
`
	src = strings.Replace(src, `    company: "ACME"`, "", 1)
	src = strings.Replace(src, `keywords: ["Go"]`, `keywords: ["  "]`, 1)
	rep, _ := load(t, src)
	for _, p := range []string{"experience[0].company", "skills[0].keywords", "certificates[0].name", "languages[0].language"} {
		if !hasPath(rep.Errors, p) {
			t.Errorf("faltou erro em %q; erros: %v", p, paths(rep.Errors))
		}
	}
}

func TestInvalidLocale(t *testing.T) {
	rep, _ := load(t, "meta:\n  locale: xx\n"+validYAML)
	if !hasPath(rep.Errors, "meta.locale") {
		t.Errorf("locale inválido deveria ser erro; erros: %v", rep.Errors)
	}
}

func TestSyntaxAndTypeErrors(t *testing.T) {
	for _, src := range []string{"foo: [\n", "- a\n- b\n", "skills: \"texto\"\n"} {
		if _, _, err := Parse([]byte(src)); err == nil {
			t.Errorf("Parse(%q) deveria falhar", src)
		}
	}
}

func TestBOMAccepted(t *testing.T) {
	rep, err := load(t, "\xef\xbb\xbf"+validYAML)
	if err != nil || !rep.OK() {
		t.Fatalf("BOM deveria ser aceito: err=%v erros=%v", err, rep.Errors)
	}
}

func TestNormalizeURL(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ""},
		{"linkedin.com/in/x", "https://linkedin.com/in/x"},
		{" http://a.com ", "http://a.com"},
		{"https://a.com", "https://a.com"},
		{"mailto:a@b.com", "mailto:a@b.com"},
		{"github.com/fabiodrneles", "https://github.com/fabiodrneles"},
	} {
		if got := NormalizeURL(tc.in); got != tc.want {
			t.Errorf("NormalizeURL(%q) = %q, esperava %q", tc.in, got, tc.want)
		}
	}
}

// 001/AC-8: os exemplos versionados são válidos e sem avisos.
func TestExamplesAreValid(t *testing.T) {
	files, _ := filepath.Glob("../../examples/*.yaml")
	if len(files) < 3 {
		t.Fatalf("esperava ao menos 3 exemplos, achei %v", files)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		rep, err := load(t, string(data))
		if err != nil || !rep.OK() || len(rep.Warnings) > 0 {
			t.Errorf("%s: err=%v erros=%v avisos=%v", f, err, rep.Errors, rep.Warnings)
		}
	}
}

// 008/AC-2: o YAML de exemplo do README é válido.
func TestReadmeExampleIsValid(t *testing.T) {
	data, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	blocks := regexp.MustCompile("(?s)```yaml\n(.*?)```").FindAllStringSubmatch(string(data), -1)
	if len(blocks) == 0 {
		t.Fatal("README sem bloco ```yaml")
	}
	for i, b := range blocks {
		rep, err := load(t, b[1])
		if err != nil || !rep.OK() || len(rep.Warnings) > 0 {
			t.Errorf("bloco yaml %d do README: err=%v erros=%v avisos=%v", i, err, rep.Errors, rep.Warnings)
		}
	}
}
