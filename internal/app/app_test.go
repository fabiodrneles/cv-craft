package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fabiodrneles/cv-craft/examples"
	"github.com/fabiodrneles/cv-craft/internal/generator"
)

func writeExample(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "cv.yaml")
	if err := os.WriteFile(path, examples.Minimal, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestOutputPaths(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join("docs", "meu cv.yaml")
	all := generator.Formats
	pdf := []generator.Format{generator.PDF}

	for _, tc := range []struct {
		name    string
		output  string
		formats []generator.Format
		want    []string
	}{
		{"padrão", "", pdf, []string{filepath.Join("docs", "meu cv.pdf")}},
		{"arquivo explícito", "out.pdf", pdf, []string{"out.pdf"}},
		{"diretório existente", dir, pdf, []string{filepath.Join(dir, "meu cv.pdf")}},
		{"all sem output", "", all, []string{
			filepath.Join("docs", "meu cv.pdf"), filepath.Join("docs", "meu cv.md"), filepath.Join("docs", "meu cv.txt"),
		}},
		{"all com diretório", "dist", all, []string{
			filepath.Join("dist", "meu cv.pdf"), filepath.Join("dist", "meu cv.md"), filepath.Join("dist", "meu cv.txt"),
		}},
	} {
		got, err := OutputPaths(in, tc.output, tc.formats)
		if err != nil || strings.Join(got, ",") != strings.Join(tc.want, ",") {
			t.Errorf("%s: %v, %v; esperava %v", tc.name, got, err, tc.want)
		}
	}

	file := filepath.Join(dir, "x.pdf")
	mustWrite(t, file, "")
	if _, err := OutputPaths(in, file, all); err == nil {
		t.Error("vários formatos com --output apontando para arquivo deveria falhar")
	}
}

func TestBuildAllFormats(t *testing.T) {
	dir := t.TempDir()
	in := writeExample(t, dir)
	res, err := Build(BuildOptions{Input: in, Output: filepath.Join(dir, "dist"), Formats: generator.Formats})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Files) != 3 || res.Lang != "pt-BR" {
		t.Fatalf("resultado inesperado: %+v", res)
	}
	for _, f := range res.Files {
		fi, err := os.Stat(f.Path)
		if err != nil || fi.Size() != f.Size || fi.Size() == 0 {
			t.Errorf("%s: stat=%v err=%v", f.Path, fi, err)
		}
	}
	entries, _ := os.ReadDir(filepath.Join(dir, "dist"))
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".cv-craft-") {
			t.Errorf("arquivo temporário esquecido: %s", e.Name())
		}
	}
}

func TestBuildOverwrite(t *testing.T) {
	dir := t.TempDir()
	in := writeExample(t, dir)
	out := filepath.Join(dir, "cv.pdf")
	mustWrite(t, out, "antigo")
	unchanged := func() bool { b, _ := os.ReadFile(out); return string(b) == "antigo" }

	// Não interativo e sem --force: falha sem tocar no arquivo.
	_, err := Build(BuildOptions{Input: in})
	if !errors.Is(err, ErrOutputExists) || !unchanged() {
		t.Fatalf("esperava ErrOutputExists e arquivo intacto, veio %v", err)
	}

	// Usuário recusa.
	asked := ""
	_, err = Build(BuildOptions{Input: in, Confirm: func(p string) (bool, error) { asked = p; return false, nil }})
	if !errors.Is(err, ErrCanceled) || asked != out || !unchanged() {
		t.Fatalf("esperava ErrCanceled perguntando por %s, veio %v (perguntou %q)", out, err, asked)
	}

	// Usuário aceita.
	if _, err = Build(BuildOptions{Input: in, Confirm: func(string) (bool, error) { return true, nil }}); err != nil || unchanged() {
		t.Fatalf("esperava sobrescrever, veio %v", err)
	}

	// --force.
	mustWrite(t, out, "antigo")
	if _, err = Build(BuildOptions{Input: in, Force: true}); err != nil || unchanged() {
		t.Fatalf("esperava sobrescrever com force, veio %v", err)
	}
}

func TestBuildInvalidInput(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.yaml")
	mustWrite(t, bad, "contact:\n  nome: x\n")
	_, err := Build(BuildOptions{Input: bad})
	var inv *InvalidError
	if !errors.As(err, &inv) || len(inv.Report.Errors) < 4 {
		t.Fatalf("esperava InvalidError com vários erros, veio %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "bad.pdf")); err == nil {
		t.Error("não deveria gerar saída para YAML inválido")
	}

	_, err = Build(BuildOptions{Input: filepath.Join(dir, "nao-existe.yaml")})
	if err == nil || errors.As(err, &inv) || !strings.Contains(err.Error(), "não encontrado") {
		t.Errorf("arquivo inexistente: %v", err)
	}
}

func TestBuildLangOverride(t *testing.T) {
	dir := t.TempDir()
	in := writeExample(t, dir)
	res, err := Build(BuildOptions{Input: in, Lang: "en", Formats: []generator.Format{generator.Text}})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(res.Files[0].Path)
	if res.Lang != "en" || !strings.Contains(string(b), "PROFESSIONAL EXPERIENCE") {
		t.Errorf("esperava títulos em inglês, veio lang=%s:\n%s", res.Lang, b)
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "novo.yaml")
	if err := Init(path, false); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Load(path); err != nil {
		t.Fatalf("modelo gerado pelo init deveria ser válido: %v", err)
	}
	mustWrite(t, path, "meu conteúdo")
	if err := Init(path, false); !errors.Is(err, ErrOutputExists) {
		t.Fatalf("esperava ErrOutputExists, veio %v", err)
	}
	if b, _ := os.ReadFile(path); string(b) != "meu conteúdo" {
		t.Fatal("init sobrescreveu o arquivo do usuário")
	}
	if err := Init(path, true); err != nil {
		t.Fatal(err)
	}
}
