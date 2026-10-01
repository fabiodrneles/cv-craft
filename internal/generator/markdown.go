package generator

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/fabiodrneles/cv-craft/internal/resume"
)

type markdownGenerator struct{ opts Options }

func (g *markdownGenerator) Generate(w io.Writer, r *resume.Resume) error {
	l := g.opts.Labels
	var blocks []string
	add := func(format string, args ...any) { blocks = append(blocks, fmt.Sprintf(format, args...)) }

	add("# %s", md(r.Contact.Name))
	var contact []string
	for _, p := range nonEmpty(r.Contact.Location, r.Contact.Phone) {
		contact = append(contact, md(p))
	}
	if e := trim(r.Contact.Email); e != "" {
		contact = append(contact, fmt.Sprintf("[%s](mailto:%s)", md(e), e))
	}
	for _, lk := range profileLinks(r.Contact, l) {
		contact = append(contact, fmt.Sprintf("[%s](%s)", md(lk.Label), lk.URL))
	}
	// Título e contato no mesmo parágrafo, para não haver parágrafo só com
	// ênfase (markdownlint MD036).
	var intro []string
	if t := trim(r.ProfessionalTitle); t != "" {
		intro = append(intro, "**"+md(t)+"**")
	}
	if len(contact) > 0 {
		intro = append(intro, strings.Join(contact, " | "))
	}
	if len(intro) > 0 {
		add("%s", strings.Join(intro, "  \n"))
	}

	if s := trim(r.Summary); s != "" {
		add("## %s", l.Summary)
		add("%s", mdLines(s, ""))
	}

	if len(r.Skills) > 0 {
		add("## %s", l.Skills)
		var lines []string
		for _, s := range r.Skills {
			lines = append(lines, fmt.Sprintf("- **%s:** %s", md(skillTitle(s, l)), mdJoin(s.Keywords)))
		}
		add("%s", strings.Join(lines, "\n"))
	}

	if len(r.Experience) > 0 {
		add("## %s", l.Experience)
		for _, e := range r.Experience {
			add("### %s", md(titled(e.Role, e.Company)))
			if m := experienceMeta(e); m != "" {
				add("%s", mdLines(m, ""))
			}
			if d := trim(e.Description); d != "" {
				add("%s", mdLines(d, ""))
			}
			addList(&blocks, l.Responsibilities, e.Responsibilities)
			addList(&blocks, l.Achievements, e.Achievements)
			if t := resume.NonBlank(e.Technologies); len(t) > 0 {
				add("**%s:** %s", l.Technologies, mdJoin(t))
			}
		}
	}

	if len(r.Education) > 0 {
		add("## %s", l.Education)
		for _, e := range r.Education {
			add("### %s", md(titled(e.Degree, e.Institution)))
			if m := educationMeta(e); m != "" {
				add("%s", mdLines(m, ""))
			}
			var extra []string
			if t := trim(e.Thesis); t != "" {
				extra = append(extra, fmt.Sprintf("**%s:** %s", l.Thesis, md(t)))
			}
			if c := trim(e.RelevantCourses); c != "" {
				extra = append(extra, fmt.Sprintf("**%s:** %s", l.RelevantCourses, md(c)))
			}
			if len(extra) > 0 {
				add("%s", strings.Join(extra, "  \n"))
			}
		}
	}

	if len(r.Certificates) > 0 {
		add("## %s", l.Certificates)
		var lines []string
		for _, c := range r.Certificates {
			line := "- **" + md(trim(c.Name)) + "**"
			if m := certificateMeta(c); m != "" {
				line += " — " + md(m)
			}
			if u := resume.NormalizeURL(c.URL); u != "" {
				line += fmt.Sprintf(" — <%s>", u)
			}
			lines = append(lines, line)
		}
		add("%s", strings.Join(lines, "\n"))
	}

	if len(r.Languages) > 0 {
		add("## %s", l.Languages)
		var lines []string
		for _, lang := range r.Languages {
			line := "- **" + md(trim(lang.Language)) + "**"
			if lv := trim(lang.Level); lv != "" {
				line += ": " + md(lv)
			}
			lines = append(lines, line)
		}
		add("%s", strings.Join(lines, "\n"))
	}

	_, err := io.WriteString(w, strings.Join(blocks, "\n\n")+"\n")
	return err
}

func addList(blocks *[]string, label string, items []string) {
	items = resume.NonBlank(items)
	if len(items) == 0 {
		return
	}
	*blocks = append(*blocks, fmt.Sprintf("**%s:**", label))
	lines := make([]string, len(items))
	for i, it := range items {
		lines[i] = "- " + mdLines(it, "  ")
	}
	*blocks = append(*blocks, strings.Join(lines, "\n"))
}

func mdJoin(items []string) string {
	items = resume.NonBlank(items)
	for i := range items {
		items[i] = md(items[i])
	}
	return strings.Join(items, ", ")
}

var mdEscaper = strings.NewReplacer(
	`\`, `\\`, "`", "\\`", `*`, `\*`, `_`, `\_`, `[`, `\[`, `]`, `\]`, `<`, `\<`,
)

var lineBreaks = regexp.MustCompile(`\s*\n\s*`)

// md escapa texto do usuário para que seja exibido literalmente (FR-MD-3).
// É usado em contextos de uma linha só (títulos, rótulos, contato), então
// quebras de linha viram espaço.
func md(s string) string {
	return mdEscaper.Replace(lineBreaks.ReplaceAllString(strings.TrimSpace(s), " "))
}

// mdLines escapa texto que pode ter várias linhas (blocos "|" do YAML),
// aplicando mdLine a cada uma. indent é prefixado às linhas seguintes, para
// que a continuação de um item de lista permaneça dentro do item.
func mdLines(s, indent string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, l := range lines {
		lines[i] = mdLine(l)
	}
	return strings.Join(lines, "\n"+indent)
}

var mdBlockStart = regexp.MustCompile(`^(\d+)([.)])`)

// mdLine escapa texto que começa uma linha, onde "-", "+", "#", ">" ou "1."
// iniciariam uma lista, título ou citação.
func mdLine(s string) string {
	s = md(s)
	if s != "" && strings.ContainsRune("-+#>", rune(s[0])) {
		return `\` + s
	}
	return mdBlockStart.ReplaceAllString(s, `$1\$2`)
}
