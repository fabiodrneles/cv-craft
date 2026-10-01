package resume

import (
	"fmt"
	"strings"
)

// Issue é um problema encontrado no YAML, identificado pelo caminho do campo
// (ex.: "experience[1].role").
type Issue struct {
	Path    string
	Line    int // 0 quando desconhecida
	Message string
}

func (i Issue) String() string {
	var b strings.Builder
	if i.Path != "" {
		b.WriteString(i.Path)
		b.WriteString(": ")
	}
	b.WriteString(i.Message)
	if i.Line > 0 {
		fmt.Fprintf(&b, " (linha %d)", i.Line)
	}
	return b.String()
}

// Report agrega todos os erros e avisos de um currículo.
type Report struct {
	Errors   []Issue
	Warnings []Issue
}

// OK indica se não há erros (avisos são permitidos).
func (r Report) OK() bool { return len(r.Errors) == 0 }

func (r *Report) errorf(path, format string, args ...any) {
	r.Errors = append(r.Errors, Issue{Path: path, Message: fmt.Sprintf(format, args...)})
}

func (r *Report) warnf(path, format string, args ...any) {
	r.Warnings = append(r.Warnings, Issue{Path: path, Message: fmt.Sprintf(format, args...)})
}

// Merge acrescenta os problemas de outro relatório a este.
func (r *Report) Merge(o Report) {
	r.Errors = append(r.Errors, o.Errors...)
	r.Warnings = append(r.Warnings, o.Warnings...)
}
