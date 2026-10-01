package resume

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// fieldPaths devolve o caminho YAML de todo campo do modelo, no formato
// usado em docs/schema.md (ex.: "experience[].role").
func fieldPaths(t reflect.Type, prefix string) []string {
	var paths []string
	for i := 0; i < t.NumField(); i++ {
		name, _, _ := strings.Cut(t.Field(i).Tag.Get("yaml"), ",")
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		switch ft := t.Field(i).Type; {
		case ft.Kind() == reflect.Struct:
			paths = append(paths, fieldPaths(ft, path)...)
		case ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Struct:
			paths = append(paths, fieldPaths(ft.Elem(), path+"[]")...)
		default:
			paths = append(paths, path)
		}
	}
	return paths
}

// 008 FR-5: todo campo do modelo está documentado em docs/schema.md, para a
// referência não ficar desatualizada quando o schema mudar.
func TestSchemaDocCoversEveryField(t *testing.T) {
	data, err := os.ReadFile("../../docs/schema.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(data)
	paths := fieldPaths(reflect.TypeOf(Resume{}), "")
	if len(paths) < 30 {
		t.Fatalf("esperava ao menos 30 campos no modelo, achei %d: %v", len(paths), paths)
	}
	for _, p := range paths {
		if !strings.Contains(doc, "| `"+p+"` |") {
			t.Errorf("docs/schema.md não documenta o campo %q (esperava uma linha de tabela começando com `%s`)", p, p)
		}
	}
	for _, lvl := range Levels {
		if !strings.Contains(doc, "| `"+lvl+"` |") {
			t.Errorf("docs/schema.md não documenta o nível %q", lvl)
		}
	}
}

// 008 FR-5: os exemplos YAML de docs/schema.md são válidos e sem avisos (os do
// README são verificados em TestReadmeExampleIsValid).
func TestDocYAMLExamplesAreValid(t *testing.T) {
	for _, doc := range []string{"../../docs/schema.md", "../../README.md", "../../README.en.md"} {
		data, err := os.ReadFile(doc)
		if err != nil {
			t.Fatal(err)
		}
		blocks := regexp.MustCompile("(?s)```yaml\\n(.*?)```").FindAllStringSubmatch(string(data), -1)
		if len(blocks) == 0 {
			t.Errorf("%s: nenhum bloco ```yaml", doc)
		}
		for i, b := range blocks {
			rep, err := load(t, b[1])
			if err != nil || !rep.OK() || len(rep.Warnings) > 0 {
				t.Errorf("%s, bloco yaml %d: err=%v erros=%v avisos=%v", doc, i, err, rep.Errors, rep.Warnings)
			}
		}
	}
}
