package resume

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/fabiodrneles/cv-craft/internal/i18n"
)

// Validate verifica o currículo e devolve todos os erros e avisos de uma vez,
// cada um com o caminho do campo (spec 001, FR-3..FR-7).
func Validate(r *Resume) Report {
	var rep Report
	if r == nil {
		rep.errorf("", "currículo ausente")
		return rep
	}

	if r.Meta.Locale != "" {
		if _, err := i18n.Get(r.Meta.Locale); err != nil {
			rep.errorf("meta.locale", "%v", err)
		}
	}

	required(&rep, "contact.name", r.Contact.Name)
	if required(&rep, "contact.email", r.Contact.Email) {
		if a, err := mail.ParseAddress(strings.TrimSpace(r.Contact.Email)); err != nil || a.Address != strings.TrimSpace(r.Contact.Email) {
			rep.errorf("contact.email", "e-mail inválido %q", r.Contact.Email)
		}
	}
	required(&rep, "professional_title", r.ProfessionalTitle)

	if len(r.Skills) == 0 {
		rep.errorf("skills", "informe pelo menos uma categoria de habilidades")
	}
	for i, s := range r.Skills {
		p := fmt.Sprintf("skills[%d]", i)
		required(&rep, p+".name", s.Name)
		if len(NonBlank(s.Keywords)) == 0 {
			rep.errorf(p+".keywords", "informe pelo menos uma habilidade")
		}
		if _, ok := NormalizeLevel(s.Level); !ok {
			rep.warnf(p+".level", "nível desconhecido %q será omitido (aceitos: %s)", s.Level, strings.Join(Levels, ", "))
		}
	}

	if len(r.Experience) == 0 {
		rep.errorf("experience", "informe pelo menos uma experiência")
	}
	for i, e := range r.Experience {
		p := fmt.Sprintf("experience[%d]", i)
		required(&rep, p+".role", e.Role)
		required(&rep, p+".company", e.Company)
	}

	if len(r.Education) == 0 {
		rep.errorf("education", "informe pelo menos uma formação")
	}
	for i, e := range r.Education {
		p := fmt.Sprintf("education[%d]", i)
		required(&rep, p+".degree", e.Degree)
		required(&rep, p+".institution", e.Institution)
	}

	for i, c := range r.Certificates {
		required(&rep, fmt.Sprintf("certificates[%d].name", i), c.Name)
	}
	for i, l := range r.Languages {
		required(&rep, fmt.Sprintf("languages[%d].language", i), l.Language)
	}
	return rep
}

func required(rep *Report, path, value string) bool {
	if strings.TrimSpace(value) == "" {
		rep.errorf(path, "campo obrigatório")
		return false
	}
	return true
}

// NonBlank devolve os itens não vazios, sem espaços nas pontas.
func NonBlank(items []string) []string {
	var out []string
	for _, s := range items {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
