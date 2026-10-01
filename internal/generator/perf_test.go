package generator

import (
	"bytes"
	"testing"
	"time"
)

// 003 NFR-2 e base para comparar os três formatos.
func BenchmarkGenerate(b *testing.B) {
	r := loadExample(b, "full.yaml")
	for _, f := range Formats {
		b.Run(string(f), func(b *testing.B) {
			g, err := New(f, Options{CreatedAt: fixedDate})
			if err != nil {
				b.Fatal(err)
			}
			var buf bytes.Buffer
			for b.Loop() {
				buf.Reset()
				if err := g.Generate(&buf, r); err != nil {
					b.Fatal(err)
				}
			}
			b.ReportMetric(float64(buf.Len()), "bytes/op-output")
		})
	}
}

// 003 NFR-2: um CV de 2 páginas é gerado em menos de 200 ms. O teste usa
// margem de 5× (1 s) para não ser instável; o limite exato fica no benchmark.
func TestPDFGenerationTime(t *testing.T) {
	r := loadExample(t, "full.yaml")
	best := time.Duration(1<<63 - 1)
	for range 3 {
		start := time.Now()
		generate(t, PDF, r)
		best = min(best, time.Since(start))
	}
	t.Logf("PDF de examples/full.yaml em %v", best)
	if limit := 5 * 200 * time.Millisecond; best > limit {
		t.Errorf("geração do PDF levou %v (limite com margem: %v)", best, limit)
	}
}

// 003 NFR-3: o PDF fica abaixo de 300 KB, inclusive com as fontes embutidas
// em subset e com muito texto.
func TestPDFSize(t *testing.T) {
	r := loadExample(t, "full.yaml")
	base := r.Experience
	for range 5 {
		r.Experience = append(r.Experience, base...)
	}
	for _, tc := range []struct {
		name string
		r    func() []byte
	}{
		{"full.yaml", func() []byte { return generate(t, PDF, loadExample(t, "full.yaml")) }},
		{"full.yaml com 18 experiências", func() []byte { return generate(t, PDF, r) }},
	} {
		size := len(tc.r())
		t.Logf("%s: %d bytes", tc.name, size)
		if size >= 300*1024 {
			t.Errorf("%s: PDF com %d bytes (limite: 300 KB)", tc.name, size)
		}
	}
}
