package resume

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrEmpty indica um arquivo sem conteúdo (vazio ou só com comentários).
var ErrEmpty = errors.New("arquivo vazio")

// Parse decodifica o YAML do currículo. Erros de sintaxe ou de tipo retornam
// error; chaves desconhecidas não interrompem o parsing e são devolvidas no
// Report para serem exibidas junto com os demais erros de validação.
func Parse(data []byte) (*Resume, Report, error) {
	var rep Report
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // BOM UTF-8

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, rep, fmt.Errorf("YAML inválido: %w", err)
	}
	if root.Kind == 0 || len(root.Content) == 0 || isNull(root.Content[0]) {
		return nil, rep, ErrEmpty
	}
	doc := root.Content[0]
	if doc.Kind != yaml.MappingNode {
		return nil, rep, fmt.Errorf("YAML inválido: o documento deve ser um mapa de campos (linha %d)", doc.Line)
	}

	checkUnknownKeys(doc, reflect.TypeOf(Resume{}), "", &rep)

	var r Resume
	if err := doc.Decode(&r); err != nil {
		return nil, rep, fmt.Errorf("YAML inválido: %w", err)
	}
	return &r, rep, nil
}

func isNull(n *yaml.Node) bool {
	return n.Kind == yaml.ScalarNode && n.Tag == "!!null"
}

// checkUnknownKeys percorre o YAML junto com o tipo Go correspondente e
// registra toda chave que não existe no schema (spec 001, FR-2).
func checkUnknownKeys(n *yaml.Node, t reflect.Type, path string, rep *Report) {
	switch {
	case t.Kind() == reflect.Struct && n.Kind == yaml.MappingNode:
		fields := yamlFields(t)
		for i := 0; i+1 < len(n.Content); i += 2 {
			key, val := n.Content[i], n.Content[i+1]
			child := joinPath(path, key.Value)
			ft, ok := fields[key.Value]
			if !ok {
				rep.Errors = append(rep.Errors, Issue{
					Path:    child,
					Line:    key.Line,
					Message: fmt.Sprintf("campo desconhecido %q", key.Value),
				})
				continue
			}
			checkUnknownKeys(val, ft, child, rep)
		}
	case t.Kind() == reflect.Slice && n.Kind == yaml.SequenceNode:
		for i, c := range n.Content {
			checkUnknownKeys(c, t.Elem(), fmt.Sprintf("%s[%d]", path, i), rep)
		}
	}
}

func yamlFields(t reflect.Type) map[string]reflect.Type {
	fields := make(map[string]reflect.Type, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("yaml"), ",")
		if name == "" || name == "-" {
			continue
		}
		fields[name] = f.Type
	}
	return fields
}

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}
