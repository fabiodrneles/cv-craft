package generator

import (
	_ "embed"
	"io"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/fabiodrneles/cv-craft/internal/resume"
)

// Liberation Sans (SIL OFL 1.1, ver fonts/OFL.txt): fonte TTF embutida para
// renderizar UTF-8 corretamente (spec 003, FR-1; decisão D4).
var (
	//go:embed fonts/LiberationSans-Regular.ttf
	fontRegular []byte
	//go:embed fonts/LiberationSans-Bold.ttf
	fontBold []byte
	//go:embed fonts/LiberationSans-Italic.ttf
	fontItalic []byte
)

const (
	fontFamily = "Sans"
	marginX    = 18.0
	marginTop  = 15.0
	marginBot  = 15.0
	lineH      = 5.0
	bulletIndt = 5.0
)

// Tons de cinza (0 = preto). Todos ≤ 80 para manter o contraste (NFR-4).
const (
	colorTitle = 0
	colorBody  = 30
	colorMuted = 80
)

type pdfGenerator struct{ opts Options }

type pdfDoc struct {
	*fpdf.Fpdf
	pageH float64
}

func (g *pdfGenerator) Generate(w io.Writer, r *resume.Resume) error {
	l := g.opts.Labels
	f := fpdf.New("P", "mm", "A4", "")
	f.AddUTF8FontFromBytes(fontFamily, "", fontRegular)
	f.AddUTF8FontFromBytes(fontFamily, "B", fontBold)
	f.AddUTF8FontFromBytes(fontFamily, "I", fontItalic)
	f.SetMargins(marginX, marginTop, marginX)
	f.SetAutoPageBreak(true, marginBot)

	created := g.opts.CreatedAt
	if created.IsZero() {
		created = time.Now()
	}
	f.SetCreationDate(created)
	f.SetModificationDate(created)
	f.SetCatalogSort(true)

	name := trim(r.Contact.Name)
	f.SetTitle(titled(name, r.ProfessionalTitle), true)
	f.SetAuthor(name, true)
	f.SetSubject(l.ResumeOf+" — "+name, true)
	f.SetKeywords(strings.Join(allKeywords(r.Skills), ", "), true)
	creator := "cv-craft"
	if g.opts.Version != "" {
		creator += " " + g.opts.Version
	}
	f.SetCreator(creator, true)
	f.SetLang(l.Code)

	_, pageH := f.GetPageSize()
	d := &pdfDoc{Fpdf: f, pageH: pageH}
	f.AddPage()

	d.header(r, g)

	if s := trim(r.Summary); s != "" {
		d.section(l.Summary)
		d.paragraph(s)
	}

	if len(r.Skills) > 0 {
		d.section(l.Skills)
		for _, s := range r.Skills {
			d.labeled(skillTitle(s, l), strings.Join(resume.NonBlank(s.Keywords), ", "))
		}
	}

	if len(r.Experience) > 0 {
		d.section(l.Experience)
		for i, e := range r.Experience {
			if i > 0 {
				d.Ln(3)
			}
			// Cabeçalho da experiência nunca fica separado do primeiro item.
			d.ensure(4 * lineH)
			d.font("B", 11, colorTitle)
			d.MultiCell(0, 6, titled(e.Role, e.Company), "", "L", false)
			if m := experienceMeta(e); m != "" {
				d.font("I", 10, colorMuted)
				d.MultiCell(0, lineH, m, "", "L", false)
			}
			if desc := trim(e.Description); desc != "" {
				d.paragraph(desc)
			}
			d.list(l.Responsibilities, e.Responsibilities)
			d.list(l.Achievements, e.Achievements)
			if t := resume.NonBlank(e.Technologies); len(t) > 0 {
				d.labeled(l.Technologies, strings.Join(t, ", "))
			}
		}
	}

	if len(r.Education) > 0 {
		d.section(l.Education)
		for i, e := range r.Education {
			if i > 0 {
				d.Ln(2)
			}
			d.ensure(3 * lineH)
			d.font("B", 11, colorTitle)
			d.MultiCell(0, 6, titled(e.Degree, e.Institution), "", "L", false)
			if m := educationMeta(e); m != "" {
				d.font("I", 10, colorMuted)
				d.MultiCell(0, lineH, m, "", "L", false)
			}
			if t := trim(e.Thesis); t != "" {
				d.labeled(l.Thesis, t)
			}
			if c := trim(e.RelevantCourses); c != "" {
				d.labeled(l.RelevantCourses, c)
			}
		}
	}

	if len(r.Certificates) > 0 {
		d.section(l.Certificates)
		for _, c := range r.Certificates {
			d.ensure(2 * lineH)
			d.font("B", 10, colorBody)
			d.Write(lineH, trim(c.Name))
			if m := certificateMeta(c); m != "" {
				d.font("", 10, colorBody)
				d.Write(lineH, " — "+m)
			}
			if u := resume.NormalizeURL(c.URL); u != "" {
				d.font("", 10, colorBody)
				d.Write(lineH, " — ")
				d.WriteLinkString(lineH, displayURL(u), u)
			}
			d.Ln(lineH)
		}
	}

	if len(r.Languages) > 0 {
		d.section(l.Languages)
		for _, lang := range r.Languages {
			d.labeled(trim(lang.Language), trim(lang.Level))
		}
	}

	return f.Output(w)
}

func (d *pdfDoc) header(r *resume.Resume, g *pdfGenerator) {
	d.font("B", 18, colorTitle)
	d.MultiCell(0, 8, trim(r.Contact.Name), "", "L", false)
	if t := trim(r.ProfessionalTitle); t != "" {
		d.font("B", 12, colorMuted)
		d.MultiCell(0, 6, t, "", "L", false)
	}
	d.Ln(1)

	d.font("", 10, colorBody)
	first := true
	sep := func() {
		if !first {
			d.Write(lineH, "  |  ")
		}
		first = false
	}
	for _, p := range nonEmpty(r.Contact.Location, r.Contact.Phone) {
		sep()
		d.Write(lineH, p)
	}
	if e := trim(r.Contact.Email); e != "" {
		sep()
		d.WriteLinkString(lineH, e, "mailto:"+e)
	}
	if !first {
		d.Ln(lineH)
	}

	first = true
	left, _, right, _ := d.GetMargins()
	pageW, _ := d.GetPageSize()
	for _, lk := range profileLinks(r.Contact, g.opts.Labels) {
		shown := displayURL(lk.URL)
		// Rótulo e link nunca são separados por uma quebra de linha.
		if !first && d.GetX()+d.GetStringWidth("  |  "+lk.Label+": "+shown) > pageW-right {
			d.Ln(lineH)
			d.SetX(left)
			first = true
		}
		sep()
		d.Write(lineH, lk.Label+": ")
		d.WriteLinkString(lineH, shown, lk.URL)
	}
	if !first {
		d.Ln(lineH)
	}
}

// section escreve o título de uma seção com um filete abaixo. Garante espaço
// para o título e ao menos duas linhas de conteúdo (spec 003, FR-8).
func (d *pdfDoc) section(title string) {
	d.Ln(4)
	d.ensure(8 + 3*lineH)
	d.font("B", 12, colorTitle)
	d.CellFormat(0, 7, strings.ToUpper(title), "", 1, "L", false, 0, "")
	left, _, right, _ := d.GetMargins()
	pageW, _ := d.GetPageSize()
	d.SetDrawColor(160, 160, 160)
	d.SetLineWidth(0.3)
	y := d.GetY()
	d.Line(left, y, pageW-right, y)
	d.Ln(2)
}

func (d *pdfDoc) paragraph(text string) {
	d.font("", 10, colorBody)
	d.MultiCell(0, lineH, text, "", "L", false)
}

// labeled escreve "Rótulo: texto" com o rótulo em negrito, quebrando linhas.
func (d *pdfDoc) labeled(label, text string) {
	d.font("B", 10, colorBody)
	if text == "" {
		d.Write(lineH, label)
	} else {
		d.Write(lineH, label+": ")
		d.font("", 10, colorBody)
		d.Write(lineH, text)
	}
	d.Ln(lineH)
}

// list escreve um rótulo seguido de itens com marcador e recuo deslocado.
func (d *pdfDoc) list(label string, items []string) {
	items = resume.NonBlank(items)
	if len(items) == 0 {
		return
	}
	left, _, right, _ := d.GetMargins()
	pageW, _ := d.GetPageSize()
	d.font("", 10, colorBody)
	firstItem := len(d.SplitText(items[0], pageW-left-right-bulletIndt))
	// O rótulo nunca fica sozinho: reserva espaço para ele e o primeiro item.
	d.ensure(lineH + 1 + float64(firstItem)*lineH)
	d.font("B", 10, colorBody)
	d.CellFormat(0, lineH+1, label+":", "", 1, "L", false, 0, "")
	d.font("", 10, colorBody)
	for _, it := range items {
		d.SetX(left + 1)
		d.CellFormat(bulletIndt-1, lineH, "•", "", 0, "L", false, 0, "")
		d.SetLeftMargin(left + bulletIndt)
		d.MultiCell(0, lineH, it, "", "L", false)
		d.SetLeftMargin(left)
		d.SetX(left)
	}
}

// ensure inicia uma nova página se não couber h mm na página atual.
func (d *pdfDoc) ensure(h float64) {
	if d.GetY()+h > d.pageH-marginBot {
		d.AddPage()
	}
}

func (d *pdfDoc) font(style string, size float64, gray int) {
	d.SetFont(fontFamily, style, size)
	d.SetTextColor(gray, gray, gray)
}

// displayURL remove o esquema para exibição; o link continua completo.
func displayURL(u string) string {
	for _, p := range []string{"https://", "http://"} {
		u = strings.TrimPrefix(u, p)
	}
	return strings.TrimSuffix(u, "/")
}

// allKeywords une as palavras-chave de todas as categorias, sem repetição.
func allKeywords(skills []resume.Skill) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range skills {
		for _, k := range resume.NonBlank(s.Keywords) {
			if key := strings.ToLower(k); !seen[key] {
				seen[key] = true
				out = append(out, k)
			}
		}
	}
	return out
}
