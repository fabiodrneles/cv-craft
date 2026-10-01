package main

import (
	"os"
	"path/filepath"
	"regexp"
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

	for _, doc := range []string{"README.md", "CONTRIBUTING.md"} {
		data, err := os.ReadFile(doc)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range regexp.MustCompile("`make ([a-z][a-z0-9-]*)`|^make ([a-z][a-z0-9-]*)").FindAllStringSubmatch(string(data), -1) {
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
