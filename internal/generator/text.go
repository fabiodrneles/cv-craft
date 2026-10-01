package generator

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/fabiodrneles/cv-craft/internal/resume"
)

type textGenerator struct{ opts Options }

func (g *textGenerator) Generate(w io.Writer, r *resume.Resume) error {
	l := g.opts.Labels
	var b strings.Builder
	line := func(s string) { b.WriteString(s); b.WriteByte('\n') }
	section := func(title string) {
		title = strings.ToUpper(title)
		line("")
		line(title)
		line(strings.Repeat("=", utf8.RuneCountInString(title)))
	}
	list := func(label string, items []string) {
		if items = resume.NonBlank(items); len(items) == 0 {
			return
		}
		line(label + ":")
		for _, it := range items {
			line("- " + it)
		}
	}

	line(trim(r.Contact.Name))
	if t := trim(r.ProfessionalTitle); t != "" {
		line(t)
	}
	if c := contactParts(r.Contact); len(c) > 0 {
		line(strings.Join(c, " | "))
	}
	if links := profileLinks(r.Contact, l); len(links) > 0 {
		parts := make([]string, len(links))
		for i, lk := range links {
			parts[i] = lk.Label + ": " + lk.URL
		}
		line(strings.Join(parts, " | "))
	}

	if s := trim(r.Summary); s != "" {
		section(l.Summary)
		line(s)
	}

	if len(r.Skills) > 0 {
		section(l.Skills)
		for _, s := range r.Skills {
			line(fmt.Sprintf("- %s: %s", skillTitle(s, l), strings.Join(resume.NonBlank(s.Keywords), ", ")))
		}
	}

	if len(r.Experience) > 0 {
		section(l.Experience)
		for i, e := range r.Experience {
			if i > 0 {
				line("")
			}
			line(titled(e.Role, e.Company))
			if m := experienceMeta(e); m != "" {
				line(m)
			}
			if d := trim(e.Description); d != "" {
				line(d)
			}
			list(l.Responsibilities, e.Responsibilities)
			list(l.Achievements, e.Achievements)
			if t := resume.NonBlank(e.Technologies); len(t) > 0 {
				line(l.Technologies + ": " + strings.Join(t, ", "))
			}
		}
	}

	if len(r.Education) > 0 {
		section(l.Education)
		for i, e := range r.Education {
			if i > 0 {
				line("")
			}
			line(titled(e.Degree, e.Institution))
			if m := educationMeta(e); m != "" {
				line(m)
			}
			if t := trim(e.Thesis); t != "" {
				line(l.Thesis + ": " + t)
			}
			if c := trim(e.RelevantCourses); c != "" {
				line(l.RelevantCourses + ": " + c)
			}
		}
	}

	if len(r.Certificates) > 0 {
		section(l.Certificates)
		for _, c := range r.Certificates {
			parts := nonEmpty(c.Name, certificateMeta(c), resume.NormalizeURL(c.URL))
			line("- " + strings.Join(parts, " — "))
		}
	}

	if len(r.Languages) > 0 {
		section(l.Languages)
		for _, lang := range r.Languages {
			s := "- " + trim(lang.Language)
			if lv := trim(lang.Level); lv != "" {
				s += ": " + lv
			}
			line(s)
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}
