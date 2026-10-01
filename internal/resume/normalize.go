package resume

import "strings"

// Levels são os níveis de habilidade canônicos, do maior para o menor.
var Levels = []string{"expert", "advanced", "proficient", "intermediate", "beginner"}

var levelAliases = map[string]string{
	"expert": "expert", "especialista": "expert",
	"advanced": "advanced", "avançado": "advanced", "avancado": "advanced",
	"proficient": "proficient", "proficiente": "proficient",
	"intermediate": "intermediate", "intermediário": "intermediate", "intermediario": "intermediate",
	"beginner": "beginner", "básico": "beginner", "basico": "beginner", "iniciante": "beginner",
}

// NormalizeLevel converte um nível escrito pelo usuário (case-insensitive,
// em inglês ou português) para o nome canônico. ok é false para valores
// desconhecidos; nível vazio devolve ("", true).
func NormalizeLevel(level string) (canonical string, ok bool) {
	l := strings.ToLower(strings.TrimSpace(level))
	if l == "" {
		return "", true
	}
	canonical, ok = levelAliases[l]
	return canonical, ok
}

// NormalizeURL garante que um link tenha esquema, para ser clicável:
// "linkedin.com/in/x" vira "https://linkedin.com/in/x".
func NormalizeURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" || strings.Contains(u, "://") || strings.HasPrefix(u, "mailto:") {
		return u
	}
	return "https://" + u
}
