package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Todo "make <alvo>" citado na documentação existe no Makefile (008 FR-6).
func TestDocumentedMakeTargetsExist(t *testing.T) {
	makefile, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	targets := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^([a-z][a-z0-9-]*):`).FindAllStringSubmatch(string(makefile), -1) {
		targets[m[1]] = true
	}

	for _, doc := range []string{"README.md", "README.en.md", "CONTRIBUTING.md"} {
		data, err := os.ReadFile(doc)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range regexp.MustCompile("(?m)`make ([a-z][a-z0-9-]*)`|^\\s*make ([a-z][a-z0-9-]*)").FindAllStringSubmatch(string(data), -1) {
			target := m[1] + m[2]
			if !targets[target] {
				t.Errorf("%s cita `make %s`, que não existe no Makefile", doc, target)
			}
		}
	}
}

// Os formulários de issue são YAML válido com os campos que o GitHub exige.
func TestIssueTemplates(t *testing.T) {
	files, _ := filepath.Glob(".github/ISSUE_TEMPLATE/*.yml")
	if len(files) < 2 {
		t.Fatalf("esperava ao menos 2 formulários, achei %v", files)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var form struct {
			Name        string           `yaml:"name"`
			Description string           `yaml:"description"`
			Body        []map[string]any `yaml:"body"`
			// config.yml
			BlankIssues *bool `yaml:"blank_issues_enabled"`
		}
		if err := yaml.Unmarshal(data, &form); err != nil {
			t.Errorf("%s: YAML inválido: %v", f, err)
			continue
		}
		if filepath.Base(f) == "config.yml" {
			if form.BlankIssues == nil {
				t.Errorf("%s: falta blank_issues_enabled", f)
			}
			continue
		}
		if form.Name == "" || form.Description == "" || len(form.Body) == 0 {
			t.Errorf("%s: formulário precisa de name, description e body", f)
		}
	}
}

// As duas versões do README têm a mesma estrutura de títulos (008 FR-3, #13):
// uma seção nova em uma delas sem a tradução na outra quebra o teste. Cada uma
// aponta para a outra no topo.
func TestReadmeTranslationsMatch(t *testing.T) {
	headings := func(path string) []string {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var levels []string
		inCode := false
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "```") {
				inCode = !inCode
				continue
			}
			if h, _, ok := strings.Cut(line, " "); !inCode && ok && h != "" && strings.Trim(h, "#") == "" {
				levels = append(levels, h)
			}
		}
		return levels
	}
	pt, en := headings("README.md"), headings("README.en.md")
	if strings.Join(pt, " ") != strings.Join(en, " ") {
		t.Errorf("os READMEs têm estruturas de títulos diferentes:\nREADME.md:    %v\nREADME.en.md: %v", pt, en)
	}
	for file, link := range map[string]string{"README.md": "(README.en.md)", "README.en.md": "(README.md)"} {
		data, _ := os.ReadFile(file)
		top, _, _ := strings.Cut(string(data), "\n## ")
		if !strings.Contains(top, link) {
			t.Errorf("%s: falta o link de idioma %s no topo", file, link)
		}
	}
}
