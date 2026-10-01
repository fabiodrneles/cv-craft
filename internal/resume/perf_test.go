package resume

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

// bigYAML gera um currículo válido com pelo menos minBytes, repetindo
// experiências com o conteúdo de examples/full.yaml.
func bigYAML(t testing.TB, minBytes int) []byte {
	t.Helper()
	data, err := os.ReadFile("../../examples/full.yaml")
	if err != nil {
		t.Fatal(err)
	}
	head, rest, ok := strings.Cut(string(data), "experience:\n")
	if !ok {
		t.Fatal("examples/full.yaml sem a seção experience")
	}
	exps, tail, ok := strings.Cut(rest, "\neducation:\n")
	if !ok {
		t.Fatal("examples/full.yaml sem a seção education")
	}
	var b strings.Builder
	b.WriteString(head + "experience:\n")
	for i := 0; b.Len() < minBytes; i++ {
		b.WriteString(strings.ReplaceAll(exps, `role: "`, fmt.Sprintf(`role: "%d `, i)))
		b.WriteString("\n")
	}
	b.WriteString("\neducation:\n" + tail)
	return []byte(b.String())
}

func parseAndValidate(tb testing.TB, data []byte) {
	r, rep, err := Parse(data)
	if err != nil {
		tb.Fatal(err)
	}
	rep.Merge(Validate(r))
	if !rep.OK() {
		tb.Fatal(rep.Errors)
	}
}

// 001 NFR-1: parsing + validação de um YAML de 50 KB em menos de 50 ms.
func BenchmarkParseValidate50KB(b *testing.B) {
	data := bigYAML(b, 50*1024)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for b.Loop() {
		parseAndValidate(b, data)
	}
}

// 001 NFR-1, com margem de 5× (250 ms) para não ser instável em runners
// compartilhados ou com o race detector; o limite exato fica no benchmark.
func TestParseValidate50KBPerformance(t *testing.T) {
	data := bigYAML(t, 50*1024)
	best := time.Duration(1<<63 - 1)
	for range 5 {
		start := time.Now()
		parseAndValidate(t, data)
		best = min(best, time.Since(start))
	}
	t.Logf("%d bytes em %v", len(data), best)
	if limit := 5 * 50 * time.Millisecond; best > limit {
		t.Errorf("parsing + validação de %d bytes levou %v (limite com margem: %v)", len(data), best, limit)
	}
}

// 001 NFR-2: o pacote resume não depende de nada além da biblioteca padrão,
// do yaml.v3 e dos pacotes internos.
func TestResumeDependencies(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("comando go não encontrado")
	}
	out, err := exec.Command(goBin, "list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}", ".").Output()
	if err != nil {
		t.Fatal(err)
	}
	allowed := []string{
		"github.com/fabiodrneles/cv-craft/internal/resume",
		"github.com/fabiodrneles/cv-craft/internal/i18n",
		"gopkg.in/yaml.v3",
	}
	for _, dep := range strings.Fields(string(out)) {
		if !slices.Contains(allowed, dep) {
			t.Errorf("dependência não permitida no pacote resume: %s (spec 001 NFR-2)", dep)
		}
	}
}
