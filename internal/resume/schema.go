package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// SchemaURL é onde o JSON Schema publicado fica acessível para os editores
// (comentário "# yaml-language-server: $schema=..." dos YAMLs).
const SchemaURL = "https://raw.githubusercontent.com/fabiodrneles/cv-craft/main/schema/cv-craft.schema.json"

// SchemaComment é a primeira linha dos YAMLs de exemplo e do modelo do init:
// faz editores com a extensão YAML (VS Code, JetBrains) autocompletarem e
// validarem o arquivo enquanto se digita (spec 010).
const SchemaComment = "# yaml-language-server: $schema=" + SchemaURL

// Regras de validação que o schema reproduz (spec 010, FR-3). Os caminhos
// usam o nome YAML dos campos; "[]" indica os itens de uma lista. Um teste
// confere que o schema e Validate aceitam e rejeitam os mesmos documentos,
// então esta tabela não tem como divergir de validate.go sem quebrar o CI.
var (
	// requiredText: texto obrigatório, não pode ser vazio nem só espaços.
	requiredText = map[string]bool{
		"contact.name": true, "contact.email": true, "professional_title": true,
		"skills[].name":     true,
		"experience[].role": true, "experience[].company": true,
		"education[].degree": true, "education[].institution": true,
		"certificates[].name":  true,
		"languages[].language": true,
	}
	// requiredList: lista obrigatória com pelo menos um item.
	requiredList = map[string]bool{"skills": true, "experience": true, "education": true, "skills[].keywords": true}
	// requiredObject: bloco obrigatório.
	requiredObject = map[string]bool{"contact": true}
)

// fieldDocs são as descrições mostradas pelo editor ao passar o mouse ou
// autocompletar. Todo campo do modelo precisa ter uma (há teste).
var fieldDocs = map[string]string{
	"meta":        "Configurações do documento (não aparecem no currículo).",
	"meta.locale": "Idioma dos títulos das seções: pt-BR (padrão) ou en. A flag --lang tem prioridade.",

	"contact":           "Dados de contato. Obrigatório.",
	"contact.name":      "Nome completo. Obrigatório.",
	"contact.email":     "E-mail. Obrigatório; vira um link clicável no PDF.",
	"contact.phone":     "Telefone, como deve aparecer (ex.: +55 11 99999-0000).",
	"contact.location":  "Cidade e estado ou país (ex.: São Paulo, SP).",
	"contact.linkedin":  "Perfil do LinkedIn; https:// é acrescentado se faltar (ex.: linkedin.com/in/seu-perfil).",
	"contact.github":    "Perfil do GitHub; https:// é acrescentado se faltar (ex.: github.com/seu-usuario).",
	"contact.portfolio": "Site ou portfólio; https:// é acrescentado se faltar.",

	"professional_title": "Título profissional, logo abaixo do nome (ex.: Engenheira de Software Sênior). Obrigatório.",
	"summary":            "Resumo profissional em um parágrafo.",

	"skills":            "Habilidades agrupadas por categoria, na ordem em que devem aparecer. Obrigatório: pelo menos uma categoria.",
	"skills[].name":     "Nome da categoria (ex.: Backend). Obrigatório.",
	"skills[].level":    "Nível, opcional: expert, advanced, proficient, intermediate ou beginner (ou em português: especialista, avançado, proficiente, intermediário, básico). Aparece traduzido entre parênteses; um valor desconhecido gera aviso e é omitido.",
	"skills[].keywords": "Habilidades da categoria (ex.: [Go, PostgreSQL, Docker]). Obrigatório: pelo menos uma.",

	"experience":                    "Experiências profissionais, da mais recente para a mais antiga. Obrigatório: pelo menos uma.",
	"experience[].role":             "Cargo. Obrigatório.",
	"experience[].company":          "Empresa. Obrigatório.",
	"experience[].location":         "Cidade ou \"Remoto\".",
	"experience[].period":           "Período, em texto livre (ex.: Mar 2022 - Atual).",
	"experience[].work_type":        "Modalidade, exibida como está (ex.: Remoto, Híbrido, Presencial).",
	"experience[].description":      "Descrição curta da empresa ou da função.",
	"experience[].responsibilities": "Responsabilidades, uma por item.",
	"experience[].achievements":     "Conquistas, de preferência com números (ex.: reduziu a latência em 40%).",
	"experience[].technologies":     "Tecnologias usadas, mostradas numa linha.",

	"education":                    "Formação acadêmica. Obrigatório: pelo menos uma.",
	"education[].degree":           "Curso ou grau (ex.: Bacharelado em Ciência da Computação). Obrigatório.",
	"education[].institution":      "Instituição. Obrigatório.",
	"education[].location":         "Cidade e estado ou país.",
	"education[].period":           "Período, em texto livre (ex.: 2015 - 2019).",
	"education[].thesis":           "Título do trabalho de conclusão, se houver.",
	"education[].relevant_courses": "Disciplinas relevantes, em texto livre.",

	"certificates":               "Certificados e cursos. Opcional.",
	"certificates[].name":        "Nome do certificado. Obrigatório.",
	"certificates[].institution": "Emissor (ex.: AWS).",
	"certificates[].date":        "Data, em texto livre (ex.: 2024).",
	"certificates[].url":         "Link de verificação; vira um link clicável no PDF.",

	"languages":            "Idiomas. Opcional.",
	"languages[].language": "Idioma (ex.: Inglês). Obrigatório.",
	"languages[].level":    "Proficiência, em texto livre (ex.: Fluente, Intermediário).",
}

// O YAML aceita qualquer escalar num campo de texto (2020 vira "2020"), e um
// campo vazio vale como ausente; o schema aceita o mesmo.
var scalarTypes = []string{"string", "number", "boolean", "null"}

// nonBlank: pelo menos um caractere que não é espaço.
const nonBlank = `\S`

// JSONSchema devolve o JSON Schema (draft-07) do YAML do currículo, gerado a
// partir do modelo Go (spec 010). A saída é determinística.
func JSONSchema() []byte {
	root := objectSchema(reflect.TypeOf(Resume{}), "")
	root = append(obj{
		{"$schema", "http://json-schema.org/draft-07/schema#"},
		{"$id", SchemaURL},
		{"title", "Currículo do CV-Craft"},
		{"description", "Schema do YAML lido por cv-craft build e cv-craft validate. Gerado a partir do código: não edite à mão (make schema)."},
	}, root...)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(root); err != nil {
		panic(err) // só tipos fixos: não acontece
	}
	return buf.Bytes()
}

func objectSchema(t reflect.Type, path string) obj {
	props := obj{}
	var required []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("yaml"), ",")
		if name == "" || name == "-" {
			continue
		}
		child := joinPath(path, name)
		props = append(props, kv{name, fieldSchema(f.Type, child)})
		if requiredText[child] || requiredList[child] || requiredObject[child] {
			required = append(required, name)
		}
	}
	s := obj{{"type", "object"}}
	if path != "" && !requiredObject[path] && !strings.HasSuffix(path, "[]") {
		s = obj{{"type", []string{"object", "null"}}}
	}
	s = append(s, kv{"additionalProperties", false}, kv{"properties", props})
	if len(required) > 0 {
		s = append(s, kv{"required", required})
	}
	return s
}

func fieldSchema(t reflect.Type, path string) obj {
	doc, ok := fieldDocs[path]
	if !ok {
		panic(fmt.Sprintf("resume: campo %s sem descrição em fieldDocs", path))
	}
	var s obj
	switch {
	case t.Kind() == reflect.Struct:
		s = objectSchema(t, path)
	case t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct:
		s = obj{{"type", "array"}}
		if !requiredList[path] {
			s = obj{{"type", []string{"array", "null"}}}
		}
		if requiredList[path] {
			s = append(s, kv{"minItems", 1})
		}
		s = append(s, kv{"items", objectSchema(t.Elem(), path+"[]")})
	case t.Kind() == reflect.Slice:
		s = obj{{"type", []string{"array", "null"}}, {"items", obj{{"type", scalarTypes}}}}
		if path == "skills[].keywords" {
			// Pelo menos um item que não seja vazio (Validate usa NonBlank).
			s = obj{
				{"type", "array"},
				{"items", obj{{"type", scalarTypes}}},
				{"contains", obj{{"type", []string{"string", "number", "boolean"}}, {"pattern", nonBlank}}},
			}
		}
	case path == "meta.locale":
		s = obj{{"anyOf", []obj{
			{{"enum", []string{"pt-BR", "en"}}},
			// Mesmas variações que i18n.Get aceita: pt, PT_br, en-US...
			{{"type", []string{"string", "null"}}, {"pattern", `^\s*(([pP][tT]|[eE][nN])([-_].*)?)?\s*$`}},
		}}}
	case path == "contact.email":
		// Aproximação do net/mail usado por Validate, suficiente para o
		// editor apontar o erro; o cv-craft validate continua sendo a
		// referência.
		s = obj{{"type", "string"}, {"pattern", `^\s*[^@\s<>()",;:]+@[^@\s<>()",;:]+\s*$`}}
	case path == "skills[].level":
		s = obj{{"anyOf", []obj{
			{{"enum", levelNames()}},
			// Nível desconhecido é só aviso: o schema não o rejeita.
			{{"type", scalarTypes}},
		}}}
	case requiredText[path]:
		s = obj{{"type", []string{"string", "number", "boolean"}}, {"pattern", nonBlank}}
	default:
		s = obj{{"type", scalarTypes}}
	}
	return append(obj{{"description", doc}}, s...)
}

// levelNames são os níveis aceitos sem aviso, canônicos primeiro.
func levelNames() []string {
	names := append([]string(nil), Levels...)
	var aliases []string
	for a, c := range levelAliases {
		if a != c {
			aliases = append(aliases, a)
		}
	}
	sort.Strings(aliases)
	return append(names, aliases...)
}

// obj é um objeto JSON que preserva a ordem das chaves.
type obj []kv

type kv struct {
	Key   string
	Value any
}

func (o obj) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, p := range o {
		if i > 0 {
			buf.WriteByte(',')
		}
		k, _ := json.Marshal(p.Key)
		buf.Write(k)
		buf.WriteByte(':')
		var v bytes.Buffer
		enc := json.NewEncoder(&v)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(p.Value); err != nil {
			return nil, err
		}
		buf.Write(bytes.TrimRight(v.Bytes(), "\n"))
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
