package resume

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/fabiodrneles/cv-craft/examples"
)

const schemaFile = "../../schema/cv-craft.schema.json"

// 010/FR-2: o schema publicado é exatamente o gerado a partir do modelo.
func TestSchemaFileUpToDate(t *testing.T) {
	got, err := os.ReadFile(schemaFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, JSONSchema()) {
		t.Fatalf("%s está desatualizado em relação ao modelo; rode 'make schema'", filepath.Base(schemaFile))
	}
}

// Toda descrição de fieldDocs corresponde a um campo que existe (o inverso,
// campo sem descrição, já faz JSONSchema entrar em pânico).
func TestSchemaDocsHaveNoStaleEntries(t *testing.T) {
	paths := map[string]bool{}
	collectPaths(JSONSchemaValue(t), "", paths)
	for p := range fieldDocs {
		if !paths[p] {
			t.Errorf("fieldDocs tem %q, que não é um campo do modelo", p)
		}
	}
	for _, m := range []map[string]bool{requiredText, requiredList, requiredObject} {
		for p := range m {
			if !paths[p] {
				t.Errorf("regra de obrigatoriedade para %q, que não é um campo do modelo", p)
			}
		}
	}
}

func JSONSchemaValue(t *testing.T) map[string]any {
	t.Helper()
	var v map[string]any
	if err := json.Unmarshal(JSONSchema(), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func collectPaths(s map[string]any, path string, out map[string]bool) {
	if props, ok := s["properties"].(map[string]any); ok {
		for name, child := range props {
			p := joinPath(path, name)
			out[p] = true
			collectPaths(child.(map[string]any), p, out)
		}
	}
	if items, ok := s["items"].(map[string]any); ok {
		collectPaths(items, path+"[]", out)
	}
}

// goAccepts informa se o CV-Craft aceita o documento: parsing sem erro e
// nenhum erro de validação (avisos não contam).
func goAccepts(doc string) bool {
	r, rep, err := Parse([]byte(doc))
	if err != nil || len(rep.Errors) > 0 {
		return false
	}
	return len(Validate(r).Errors) == 0
}

func compileSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(JSONSchema()))
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(SchemaURL, doc); err != nil {
		t.Fatal(err)
	}
	s, err := c.Compile(SchemaURL)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// schemaAccepts valida o YAML contra o schema, como faria um editor.
func schemaAccepts(t *testing.T, s *jsonschema.Schema, doc string) (bool, error) {
	t.Helper()
	var v any
	if err := yaml.Unmarshal([]byte(doc), &v); err != nil {
		return false, err // YAML com erro de sintaxe: o editor também aponta
	}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	err = s.Validate(inst)
	return err == nil, err
}

const validDoc = `meta:
  locale: pt-BR
contact:
  name: Ana Lima
  email: ana@example.com
  phone: "+55 11 99999-0000"
  location: São Paulo, SP
  linkedin: linkedin.com/in/ana
  github: github.com/ana
  portfolio: ana.dev
professional_title: Engenheira de Software
summary: Resumo.
skills:
  - name: Backend
    level: advanced
    keywords: [Go, SQL]
experience:
  - role: Dev
    company: ACME
    location: Remoto
    period: 2020 - Atual
    work_type: Remoto
    description: Descrição.
    responsibilities: [Uma]
    achievements: [Outra]
    technologies: [Go]
education:
  - degree: Bacharelado
    institution: USP
    location: São Paulo
    period: 2015 - 2019
    thesis: TCC
    relevant_courses: Algoritmos
certificates:
  - name: AWS
    institution: Amazon
    date: "2024"
    url: aws.amazon.com
languages:
  - language: Inglês
    level: Fluente
`

// edit troca a primeira ocorrência de old em validDoc (que precisa existir).
func edit(t *testing.T, old, repl string) string {
	t.Helper()
	if !strings.Contains(validDoc, old) {
		t.Fatalf("validDoc não contém %q", old)
	}
	return strings.Replace(validDoc, old, repl, 1)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// block troca o trecho de validDoc que vai de start (inclusive) até end
// (exclusive) por repl.
func block(t *testing.T, start, end, repl string) string {
	t.Helper()
	i, j := strings.Index(validDoc, start), strings.Index(validDoc, end)
	if i < 0 || j < i {
		t.Fatalf("validDoc não tem o trecho %q..%q", start, end)
	}
	return validDoc[:i] + repl + validDoc[j:]
}

// 010/FR-3 e AC-2: o schema e a validação do Go aceitam e rejeitam os mesmos
// documentos.
func TestSchemaMatchesValidate(t *testing.T) {
	s := compileSchema(t)
	cases := map[string]string{
		"exemplo minimal": string(examples.Minimal),
		"exemplo full":    readFile(t, "../../examples/full.yaml"),
		"exemplo en":      readFile(t, "../../examples/en.yaml"),
		"completo":        validDoc,

		// Campos obrigatórios ausentes, vazios, só com espaços ou nulos.
		"sem contact":               edit(t, "contact:\n  name: Ana Lima\n  email: ana@example.com\n  phone: \"+55 11 99999-0000\"\n  location: São Paulo, SP\n  linkedin: linkedin.com/in/ana\n  github: github.com/ana\n  portfolio: ana.dev\n", ""),
		"contact nulo":              edit(t, "contact:\n  name: Ana Lima\n  email: ana@example.com\n  phone: \"+55 11 99999-0000\"\n  location: São Paulo, SP\n  linkedin: linkedin.com/in/ana\n  github: github.com/ana\n  portfolio: ana.dev\n", "contact:\n"),
		"sem name":                  edit(t, "  name: Ana Lima\n", ""),
		"name vazio":                edit(t, "name: Ana Lima", `name: ""`),
		"name só espaços":           edit(t, "name: Ana Lima", `name: "   "`),
		"name nulo":                 edit(t, "name: Ana Lima", "name:"),
		"name numérico":             edit(t, "name: Ana Lima", "name: 42"),
		"sem email":                 edit(t, "  email: ana@example.com\n", ""),
		"email inválido":            edit(t, "email: ana@example.com", "email: não-é-email"),
		"email com nome":            edit(t, "email: ana@example.com", `email: "Ana <ana@example.com>"`),
		"email com espaço":          edit(t, "email: ana@example.com", "email: ana @example.com"),
		"email sem domínio":         edit(t, "email: ana@example.com", "email: ana@"),
		"email numérico":            edit(t, "email: ana@example.com", "email: 5"),
		"sem título":                edit(t, "professional_title: Engenheira de Software\n", ""),
		"título vazio":              edit(t, "professional_title: Engenheira de Software", "professional_title:"),
		"sem skills":                edit(t, "skills:\n  - name: Backend\n    level: advanced\n    keywords: [Go, SQL]\n", ""),
		"skills vazia":              edit(t, "skills:\n  - name: Backend\n    level: advanced\n    keywords: [Go, SQL]\n", "skills: []\n"),
		"skills nula":               edit(t, "skills:\n  - name: Backend\n    level: advanced\n    keywords: [Go, SQL]\n", "skills:\n"),
		"skill sem name":            edit(t, "  - name: Backend\n", "  - "),
		"skill sem keywords":        edit(t, "    keywords: [Go, SQL]\n", ""),
		"keywords vazia":            edit(t, "keywords: [Go, SQL]", "keywords: []"),
		"keywords só em branco":     edit(t, "keywords: [Go, SQL]", `keywords: ["", "  "]`),
		"keywords com um válido":    edit(t, "keywords: [Go, SQL]", `keywords: ["", Go]`),
		"keywords numérica":         edit(t, "keywords: [Go, SQL]", "keywords: [2024]"),
		"keywords com lista":        edit(t, "keywords: [Go, SQL]", "keywords: [[Go]]"),
		"experience vazia":          block(t, "experience:\n", "education:\n", "experience: []\n"),
		"sem experience":            block(t, "experience:\n", "education:\n", ""),
		"experience sem role":       edit(t, "  - role: Dev\n", "  - "),
		"experience sem company":    edit(t, "    company: ACME\n", ""),
		"education sem degree":      edit(t, "  - degree: Bacharelado\n", "  - "),
		"education sem institution": edit(t, "    institution: USP\n", ""),
		"education nula":            edit(t, "education:\n  - degree: Bacharelado\n    institution: USP\n    location: São Paulo\n    period: 2015 - 2019\n    thesis: TCC\n    relevant_courses: Algoritmos\n", "education:\n"),
		"certificado sem name":      edit(t, "  - name: AWS\n", "  - "),
		"idioma sem language":       edit(t, "  - language: Inglês\n", "  - "),

		// Opcionais: ausentes, nulos, vazios ou numéricos são aceitos.
		"sem meta":               edit(t, "meta:\n  locale: pt-BR\n", ""),
		"meta nulo":              edit(t, "meta:\n  locale: pt-BR\n", "meta:\n"),
		"summary nulo":           edit(t, "summary: Resumo.", "summary:"),
		"period numérico":        edit(t, "period: 2020 - Atual", "period: 2020"),
		"phone numérico":         edit(t, `phone: "+55 11 99999-0000"`, "phone: 11999990000"),
		"sem certificados":       edit(t, "certificates:\n  - name: AWS\n    institution: Amazon\n    date: \"2024\"\n    url: aws.amazon.com\n", ""),
		"certificados nulos":     edit(t, "certificates:\n  - name: AWS\n    institution: Amazon\n    date: \"2024\"\n    url: aws.amazon.com\n", "certificates:\n"),
		"technologies nula":      edit(t, "technologies: [Go]", "technologies:"),
		"achievements vazia":     edit(t, "achievements: [Outra]", "achievements: []"),
		"description lista":      edit(t, "description: Descrição.", "description: [a]"),
		"description mapa":       edit(t, "description: Descrição.", "description: {a: b}"),
		"responsibilities texto": edit(t, "responsibilities: [Uma]", "responsibilities: Uma"),

		// Níveis: em inglês ou português, qualquer caixa; desconhecido é só aviso.
		"level português":    edit(t, "level: advanced", "level: Avançado"),
		"level maiúsculo":    edit(t, "level: advanced", "level: EXPERT"),
		"level desconhecido": edit(t, "level: advanced", "level: guru"),
		"level nulo":         edit(t, "level: advanced", "level:"),

		// Idioma do documento.
		"locale en":          edit(t, "locale: pt-BR", "locale: en"),
		"locale en-US":       edit(t, "locale: pt-BR", "locale: en-US"),
		"locale PT_br":       edit(t, "locale: pt-BR", "locale: PT_br"),
		"locale pt-PT":       edit(t, "locale: pt-BR", "locale: pt-PT"),
		"locale vazio":       edit(t, "locale: pt-BR", `locale: ""`),
		"locale nulo":        edit(t, "locale: pt-BR", "locale:"),
		"locale com espaços": edit(t, "locale: pt-BR", `locale: " en "`),
		"locale fr":          edit(t, "locale: pt-BR", "locale: fr"),
		"locale english":     edit(t, "locale: pt-BR", "locale: english"),
		"locale numérico":    edit(t, "locale: pt-BR", "locale: 1"),

		// Chaves desconhecidas em todos os níveis (spec 001, FR-2).
		"chave desconhecida na raiz":        validDoc + "foo: bar\n",
		"chave desconhecida em meta":        edit(t, "  locale: pt-BR\n", "  locale: pt-BR\n  theme: dark\n"),
		"chave desconhecida em contact":     edit(t, "  name: Ana Lima\n", "  name: Ana Lima\n  emial: a@b.c\n"),
		"chave desconhecida em skills":      edit(t, "    level: advanced\n", "    level: advanced\n    nivel: 3\n"),
		"typo em experience":                edit(t, "    responsibilities: [Uma]\n", "    responsabilities: [Uma]\n"),
		"chave desconhecida em education":   edit(t, "    thesis: TCC\n", "    thesis: TCC\n    grade: 10\n"),
		"chave desconhecida em certificate": edit(t, "    url: aws.amazon.com\n", "    url: aws.amazon.com\n    id: 1\n"),
		"chave desconhecida em language":    edit(t, "    level: Fluente\n", "    level: Fluente\n    nota: 9\n"),

		// Estrutura errada.
		"documento lista":    "- a\n- b\n",
		"skills como mapa":   edit(t, "skills:\n  - name: Backend\n    level: advanced\n    keywords: [Go, SQL]\n", "skills:\n  name: Backend\n"),
		"contact como lista": edit(t, "contact:\n  name: Ana Lima\n", "contact:\n  - name: Ana Lima\n"),
	}

	accepted := 0
	for name, doc := range cases {
		goOK := goAccepts(doc)
		schemaOK, err := schemaAccepts(t, s, doc)
		if goOK != schemaOK {
			t.Errorf("%s: Validate aceita = %v, schema aceita = %v (%v)\n%s", name, goOK, schemaOK, err, doc)
		}
		if goOK {
			accepted++
		}
	}
	// Os casos precisam exercitar os dois lados, ou a equivalência não prova nada.
	if rejected := len(cases) - accepted; accepted < 20 || rejected < 20 {
		t.Errorf("casos desequilibrados: %d aceitos e %d rejeitados", accepted, rejected)
	}
}

// 010/FR-4: os exemplos (e o modelo do init, que é o minimal.yaml) apontam
// para o schema publicado na primeira linha.
func TestExamplesReferenceSchema(t *testing.T) {
	for _, name := range []string{"minimal.yaml", "full.yaml", "en.yaml"} {
		first, _, _ := strings.Cut(readFile(t, "../../examples/"+name), "\n")
		if first != SchemaComment {
			t.Errorf("%s: a primeira linha deveria ser %q, é %q", name, SchemaComment, first)
		}
	}
}
