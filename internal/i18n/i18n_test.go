package i18n

import (
	"reflect"
	"testing"
)

// 006/AC-4: todos os rótulos existem em todos os idiomas.
func TestCatalogComplete(t *testing.T) {
	levels := []string{"expert", "advanced", "proficient", "intermediate", "beginner"}
	for _, code := range Supported() {
		l, err := Get(code)
		if err != nil {
			t.Fatal(err)
		}
		v := reflect.ValueOf(l)
		for i := 0; i < v.NumField(); i++ {
			if f := v.Field(i); f.Kind() == reflect.String && f.String() == "" {
				t.Errorf("%s: rótulo %s vazio", code, v.Type().Field(i).Name)
			}
		}
		for _, lv := range levels {
			if l.Levels[lv] == "" {
				t.Errorf("%s: nível %s sem tradução", code, lv)
			}
		}
	}
}

func TestGetAliases(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", "pt-BR"}, {"pt", "pt-BR"}, {"PT_br", "pt-BR"}, {"pt-PT", "pt-BR"},
		{"en", "en"}, {"en-US", "en"}, {" EN ", "en"},
	} {
		l, err := Get(tc.in)
		if err != nil || l.Code != tc.want {
			t.Errorf("Get(%q) = %q, %v; esperava %q", tc.in, l.Code, err, tc.want)
		}
	}
	if _, err := Get("xx"); err == nil {
		t.Error("Get(xx) deveria falhar")
	}
}
