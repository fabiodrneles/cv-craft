package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// As skills versionadas em .claude/skills/ continuam válidas: frontmatter
// YAML com name igual ao nome da pasta, description dentro do limite e
// todos os links relativos apontando para arquivos que existem (#25).
func TestRepoSkills(t *testing.T) {
	skills, _ := filepath.Glob(".claude/skills/*/SKILL.md")
	if len(skills) == 0 {
		t.Fatal("nenhuma skill em .claude/skills/")
	}
	link := regexp.MustCompile(`\]\(([^)#\s]+)(?:#[^)]*)?\)`)

	for _, skill := range skills {
		dir := filepath.Dir(skill)
		data, err := os.ReadFile(skill)
		if err != nil {
			t.Fatal(err)
		}

		parts := strings.SplitN(string(data), "---\n", 3)
		if len(parts) != 3 || parts[0] != "" {
			t.Errorf("%s: o arquivo deve começar com um frontmatter entre linhas ---", skill)
			continue
		}
		var meta struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
		}
		if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
			t.Errorf("%s: frontmatter inválido: %v", skill, err)
			continue
		}
		if meta.Name != filepath.Base(dir) {
			t.Errorf("%s: name %q difere do nome da pasta %q", skill, meta.Name, filepath.Base(dir))
		}
		if n := len([]rune(meta.Description)); n == 0 || n > 1024 {
			t.Errorf("%s: description com %d caracteres (deve ter de 1 a 1024)", skill, n)
		}

		// Links relativos de todos os .md da skill (exemplos dentro de blocos de
		// código, como os modelos de templates.md, não são links reais).
		_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
				return err
			}
			md, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, m := range link.FindAllStringSubmatch(withoutCodeBlocks(string(md)), -1) {
				target := m[1]
				if strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
					continue
				}
				if _, err := os.Stat(filepath.Join(filepath.Dir(path), target)); err != nil {
					t.Errorf("%s: link quebrado para %s", path, target)
				}
			}
			return nil
		})
	}
}

// withoutCodeBlocks remove o conteúdo dos blocos de código cercados por ```.
func withoutCodeBlocks(md string) string {
	var b strings.Builder
	inside := false
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inside = !inside
			continue
		}
		if !inside {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return b.String()
}
