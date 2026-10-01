package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fabiodrneles/cv-craft/examples"
)

type result struct {
	code           int
	stdout, stderr string
}

func run(t *testing.T, stdin string, tty bool, args ...string) result {
	t.Helper()
	var out, errb bytes.Buffer
	code := Run(args, Env{
		Stdin:      strings.NewReader(stdin),
		Stdout:     &out,
		Stderr:     &errb,
		IsTerminal: func() bool { return tty },
		Version:    "v9.9.9",
		Commit:     "abc123",
		Date:       "2025-01-01",
	})
	return result{code, out.String(), errb.String()}
}

// chdir muda para um diretório temporário com o modelo mínimo como cv.yaml.
func chdir(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	mustWrite(t, "cv.yaml", string(examples.Minimal))
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func expect(t *testing.T, r result, code int, contains ...string) {
	t.Helper()
	if r.code != code {
		t.Errorf("exit %d, esperava %d\nstdout: %s\nstderr: %s", r.code, code, r.stdout, r.stderr)
	}
	for _, c := range contains {
		if !strings.Contains(r.stdout+r.stderr, c) {
			t.Errorf("saída não contém %q\nstdout: %s\nstderr: %s", c, r.stdout, r.stderr)
		}
	}
}

// 002/AC-1.
func TestBuildPDF(t *testing.T) {
	chdir(t)
	r := run(t, "", false, "build", "cv.yaml", "-o", "out.pdf")
	expect(t, r, ExitOK, "Gerado: out.pdf")
	if b, err := os.ReadFile("out.pdf"); err != nil || !bytes.HasPrefix(b, []byte("%PDF-")) {
		t.Fatalf("out.pdf inválido: %v", err)
	}
	if r.stderr != "" {
		t.Errorf("stderr deveria estar vazio: %q", r.stderr)
	}
}

// 002/AC-2.
func TestBuildAll(t *testing.T) {
	chdir(t)
	expect(t, run(t, "", false, "build", "cv.yaml", "--format", "all", "-o", "dist"), ExitOK)
	for _, f := range []string{"cv.pdf", "cv.md", "cv.txt"} {
		if _, err := os.Stat(filepath.Join("dist", f)); err != nil {
			t.Error(err)
		}
	}
}

// 002/AC-3, AC-4.
func TestBuildOverwriteNonInteractive(t *testing.T) {
	chdir(t)
	mustWrite(t, "cv.pdf", "antigo")
	expect(t, run(t, "", false, "build", "cv.yaml"), ExitExists, "--force")
	if b, _ := os.ReadFile("cv.pdf"); string(b) != "antigo" {
		t.Fatal("arquivo foi alterado sem --force")
	}
	expect(t, run(t, "", false, "build", "cv.yaml", "--force"), ExitOK)
	expect(t, run(t, "", false, "build", "-y", "cv.yaml"), ExitOK)
}

func TestBuildOverwriteInteractive(t *testing.T) {
	chdir(t)
	mustWrite(t, "cv.pdf", "antigo")
	expect(t, run(t, "n\n", true, "build", "cv.yaml"), ExitExists, "Sobrescrever?", "Cancelado")
	expect(t, run(t, "", true, "build", "cv.yaml"), ExitExists) // EOF = não
	expect(t, run(t, "s\n", true, "build", "cv.yaml"), ExitOK, "Gerado: cv.pdf")
}

// 002/AC-5.
func TestHelpAndUsageErrors(t *testing.T) {
	expect(t, run(t, "", false), ExitOK, "Comandos:")
	expect(t, run(t, "", false, "--help"), ExitOK, "Comandos:")
	expect(t, run(t, "", false, "help", "build"), ExitOK, "--format")
	expect(t, run(t, "", false, "build", "-h"), ExitOK, "--format")
	expect(t, run(t, "", false, "foo"), ExitUsage, "comando desconhecido")
	expect(t, run(t, "", false, "help", "foo"), ExitUsage)
	expect(t, run(t, "", false, "build"), ExitUsage, "exatamente um arquivo")
	expect(t, run(t, "", false, "build", "a.yaml", "b.yaml"), ExitUsage)
	expect(t, run(t, "", false, "build", "cv.yaml", "--nope"), ExitUsage, "flag desconhecida: -nope")
	expect(t, run(t, "", false, "build", "cv.yaml", "-o"), ExitUsage, "precisa de um valor")
	expect(t, run(t, "", false, "build", "cv.yaml", "-f", "docx"), ExitUsage, "formato inválido")
	expect(t, run(t, "", false, "build", "cv.yaml", "--lang", "xx"), ExitUsage, "pt-BR, en")
}

// 002/AC-6.
func TestFlagsBeforeAndAfterFile(t *testing.T) {
	chdir(t)
	expect(t, run(t, "", false, "build", "-f", "md", "cv.yaml"), ExitOK, "cv.md")
	expect(t, run(t, "", false, "build", "cv.yaml", "-f", "txt"), ExitOK, "cv.txt")
	expect(t, run(t, "", false, "build", "-format", "text", "-output", "x.txt", "cv.yaml"), ExitOK, "x.txt")
}

// 002/AC-7.
func TestInit(t *testing.T) {
	t.Chdir(t.TempDir())
	expect(t, run(t, "", false, "init"), ExitOK, "Criado: curriculum.yaml")
	expect(t, run(t, "", false, "init"), ExitExists, "--force")
	expect(t, run(t, "", false, "init", "--force"), ExitOK)
	expect(t, run(t, "", false, "init", "meu cv.yaml"), ExitOK, `cv-craft build "meu cv.yaml"`)
	expect(t, run(t, "", false, "validate", "curriculum.yaml"), ExitOK, "Válido")
}

// 002/AC-8.
func TestValidate(t *testing.T) {
	chdir(t)
	mustWrite(t, "bad.yaml", "contact:\n  name: Ana\n  emial: x\n")
	r := run(t, "", false, "validate", "bad.yaml")
	expect(t, r, ExitInvalid, "contact.emial: campo desconhecido", "contact.email: campo obrigatório", "professional_title", "skills", "experience", "education")
	mustWrite(t, "syntax.yaml", "foo: [\n")
	expect(t, run(t, "", false, "validate", "syntax.yaml"), ExitInvalid, "YAML inválido")
	expect(t, run(t, "", false, "validate", "nao-existe.yaml"), ExitError, "não encontrado")
	expect(t, run(t, "", false, "build", "bad.yaml"), ExitInvalid)
}

func TestWarningsGoToStderr(t *testing.T) {
	chdir(t)
	data := strings.Replace(string(examples.Minimal), `level: "advanced"`, `level: "guru"`, 1)
	mustWrite(t, "warn.yaml", data)
	r := run(t, "", false, "build", "warn.yaml")
	expect(t, r, ExitOK, "aviso: skills[0].level")
	if strings.Contains(r.stdout, "aviso") {
		t.Error("avisos devem ir para stderr")
	}
	r = run(t, "", false, "build", "warn.yaml", "-q", "-y")
	if r.code != ExitOK || r.stdout != "" || r.stderr != "" {
		t.Errorf("--quiet deveria silenciar: %+v", r)
	}
}

func TestDeprecatedATSFlag(t *testing.T) {
	chdir(t)
	expect(t, run(t, "", false, "build", "cv.yaml", "-ats"), ExitOK, "aviso: --ats está obsoleta")
}

// 002/AC-9: sem ANSI quando a saída não é um terminal.
func TestNoANSIWhenNotTTY(t *testing.T) {
	chdir(t)
	r := run(t, "", false, "build", "cv.yaml")
	if strings.Contains(r.stdout+r.stderr, "\x1b[") {
		t.Errorf("saída contém códigos ANSI: %q", r.stdout)
	}
}

// 002/AC-10.
func TestVersion(t *testing.T) {
	expect(t, run(t, "", false, "version"), ExitOK, "cv-craft v9.9.9 (commit abc123, 2025-01-01)")
	expect(t, run(t, "", false, "--version"), ExitOK, "v9.9.9")
}

// ---- modo interativo (spec 005) ----

// 005/AC-1 (FR-8).
func TestInteractiveRequiresTTY(t *testing.T) {
	expect(t, run(t, "", false, "ui"), ExitUsage, "requer um terminal")
}

// 005/AC-3 (FR-1, FR-2): EOF encerra; linhas vazias não imprimem ajuda.
func TestInteractiveEOFAndEmptyLines(t *testing.T) {
	r := run(t, "\n\n   \n", true, "ui")
	expect(t, r, ExitOK)
	if strings.Contains(r.stdout, "Comandos:") {
		t.Error("linha vazia não deveria imprimir a ajuda")
	}
}

// 005/AC-2, AC-4, AC-5.
func TestInteractiveSession(t *testing.T) {
	chdir(t)
	input := strings.Join([]string{
		"help",
		"foo",
		`build "cv.yaml" -o "meu cv.pdf"`,
		"build cv.yaml -o 'meu cv.pdf'", // já existe: pergunta
		"n",
		`build "sem fim`,
		"ui",
		"exit",
		"version", // não deve executar após exit
	}, "\n") + "\n"
	r := run(t, input, true, "interactive")
	expect(t, r, ExitOK, "Comandos:", "comando desconhecido", "Gerado: meu cv.pdf", "Sobrescrever?", "Cancelado", "aspas não fechadas", "já está no modo interativo")
	if strings.Contains(r.stdout, "v9.9.9 (commit") {
		t.Error("comando após exit foi executado")
	}
	if _, err := os.Stat("meu cv.pdf"); err != nil {
		t.Error(err)
	}
}

func TestTokenize(t *testing.T) {
	for in, want := range map[string]string{
		`build cv.yaml`:                      "build|cv.yaml",
		`  build   "meu cv.yaml"  -o x.pdf `: "build|meu cv.yaml|-o|x.pdf",
		`build 'a b'"c d"`:                   "build|a bc d",
		`build C:\Users\ana\cv.yaml`:         `build|C:\Users\ana\cv.yaml`,
		`build ""`:                           "build|",
	} {
		got, err := tokenize(in)
		if err != nil || strings.Join(got, "|") != want {
			t.Errorf("tokenize(%q) = %q, %v; esperava %q", in, strings.Join(got, "|"), err, want)
		}
	}
	for _, bad := range []string{`"aberto`, "   "} {
		if _, err := tokenize(bad); err == nil {
			t.Errorf("tokenize(%q) deveria falhar", bad)
		}
	}
}
