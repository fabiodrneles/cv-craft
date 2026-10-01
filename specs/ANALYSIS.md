# Análise e Verificação do CV-Craft

> Data: 2026-10-01 · Base: commit `93d342d` (branch `main`)
> Método: leitura completa do código + build, `go vet` e execução real de todos os formatos e casos de borda.

## 1. Resumo executivo

O CV-Craft compila, passa no `go vet` e gera arquivos nos três formatos, mas **ainda não está pronto para uso real**:

- O **PDF corrompe caracteres acentuados** (`São Paulo` → `SÃ£o Paulo`) — fatal para um público brasileiro e para ATS.
- **Markdown e TXT perdem praticamente todo o conteúdo das experiências** (renderizam um campo que quase ninguém preenche).
- O PDF **descarta silenciosamente** skills com nível diferente de `proficient|intermediate|beginner`.
- O **modo interativo entra em loop infinito** se o stdin fechar (EOF / pipe / CI).
- A flag **`-ats` não tem efeito** algum.
- **Zero testes**, sem CI, sem `.gitignore`, binário `.exe` de 5,8 MB versionado, e o `go install` do README **não funciona**.

A base (parser → modelo → interface `Generator`) é boa e simples; o trabalho é de **correção, consolidação e cobertura de testes**, não de reescrita.

## 2. O que foi verificado

| Verificação | Resultado |
|---|---|
| `go build` / `go vet ./...` (Go 1.24.7) | ✅ OK, sem avisos |
| `build -format pdf` | ⚠️ Gera, mas com acentos corrompidos |
| `build -format md` / `markdown` | ⚠️ Gera, mas sem responsabilidades/conquistas/tecnologias/idiomas/certificados |
| `build -format text` | ⚠️ Mesmo problema do Markdown |
| `build -format docx` | ✅ Erro claro, exit 1 |
| YAML inválido | ✅ Erro claro, exit 1 |
| Skill com `level: advanced` | ❌ Some do PDF sem aviso |
| Arquivo de saída existente + stdin não interativo | ❌ Cancela com exit 1 (não há `--force`) |
| `-ats` vs. sem `-ats` | ❌ Saída idêntica |
| `cv-craft --help` / `cv-craft version` | ❌ Exit 1; `version` não existe fora do modo interativo |
| `cv-craft ui < /dev/null` | ❌ Loop infinito (precisou de `timeout`) |
| Testes automatizados | ❌ Nenhum arquivo `_test.go` |

## 3. Observações por severidade

### 🔴 Críticas (bloqueiam o "funciona de verdade")

| # | Observação | Evidência | Spec |
|---|---|---|---|
| C1 | PDF usa fontes core do fpdf (cp1252) sem tradutor UTF-8 → `ã`, `ç`, `é` corrompidos | `internal/generator/pdf.go` — nenhum `UnicodeTranslatorFromDescriptor`/fonte TTF | 003 |
| C2 | MD/TXT só renderizam `Description` da experiência; ignoram `responsibilities`, `achievements`, `technologies` | `markdown.go:32-36`, `text.go:35-38`; ver `curriculum.md` versionado com blocos vazios | 004 |
| C3 | Skills com nível fora de `proficient/intermediate/beginner` (ou vazio) somem do PDF | `pdf.go:148-165` (`getSkillsByLevel`) | 001, 003 |
| C4 | Modo interativo em loop infinito em EOF (`scanner.Scan()` ignorado; entrada vazia → `help`) | `cli/ui.go:71`, `ui.go:162` | 005 |
| C5 | `go install github.com/seu-usuario/cv-craft@latest` não funciona: módulo se chama `cv-craft`, placeholder `seu-usuario`, link markdown dentro de bloco de código | `go.mod:1`, `README.md` | 007, 008 |

### 🟠 Altas

| # | Observação | Evidência | Spec |
|---|---|---|---|
| A1 | Flag `-ats` é aceita e documentada, mas não altera nada | `PDFGenerator.IsATS` nunca é lido | 003 |
| A2 | Lógica de build duplicada entre `main.go` (`buildResume`) e `cli/ui.go` (`executeBuild`), já divergindo (UI sobrescreve sem perguntar; mensagens em idiomas diferentes) | `main.go:147`, `ui.go:280` | 002 |
| A3 | Prompt de sobrescrita bloqueia uso em scripts/CI; não há `--force`/`-y` | `main.go:191-202` | 002 |
| A4 | `init` sobrescreve arquivo existente sem aviso (perda de dados do usuário) | `ui.go:392` | 002 |
| A5 | Campos do YAML ignorados em todos os formatos: `achievements`, `relevant_courses`, `certificates[].url`; `location` ausente em MD/TXT; certificados e idiomas ausentes em MD/TXT | `parser.go` vs. geradores | 001, 003, 004 |
| A6 | Typos em chaves do YAML são ignorados silenciosamente (ex.: `responsabilities`) | `yaml.Unmarshal` sem `KnownFields(true)` | 001 |
| A7 | Validação dividida: `ParseFile` valida campos obrigatórios e `ValidateCurriculum` valida listas; retorna só o primeiro erro | `parser.go:100-134` | 001 |
| A8 | Sem testes, sem CI, sem `.gitignore`; binário `cv-craft.exe` (5,8 MB) e PDFs gerados versionados | `git ls-files` | 007 |

### 🟡 Médias

| # | Observação | Spec |
|---|---|---|
| M1 | Idiomas misturados: PDF com títulos em inglês (`SUMMARY`), MD/TXT em português (`Resumo`), CLI ora PT ora EN, help do modo interativo em EN | 006 |
| M2 | Modo "legado" (`os.Args[0]` contém `main.go`) é praticamente código morto: com `go run`, `Args[0]` é um binário temporário | 002 |
| M3 | `--help`/`-h` retorna exit 1; não existe `cv-craft version`/`--version` fora do modo interativo; versão hard-coded `v1.0` | 002 |
| M4 | Help interativo promete TAB-completion e histórico (↑↓) que não existem; `templates list` lista 5 templates que não existem | 005 |
| M5 | `log.Println("Generating professional ATS-optimized resume...")` dentro do gerador sempre polui stderr; `-v` só muda flags do logger | 002, 003 |
| M6 | Quebra de página no PDF por heurística fixa (`GetY() > 250`) — título de seção pode ficar órfão no fim da página | 003 |
| M7 | Markdown: link do portfólio sem `https://`; campos vazios geram `\|  \|` e `[LinkedIn](https://)`; sem escape de caracteres Markdown | 004 |
| M8 | PDF sem metadados (`Title`, `Author`, `Subject`, `Keywords`) — úteis para ATS e para quem recebe o arquivo | 003 |
| M9 | Saída padrão derivada do nome do YAML (`curriculum.yaml` → `curriculum.md`) pode colidir com arquivos do próprio repo | 002 |

### 🟢 Baixas / higiene

- Nome `GeminiUI` / comentários "Gemini style" → renomear para `cli.Shell` ou similar.
- `printUsage` diz **"GO-Craft"**; produto se chama **CV-Craft**.
- `color.New(color.FgCyan)` sem uso em `Prompt` (`ui.go:64`).
- `ShowLoading` faz `sleep` artificial de 1 s antes de cada operação.
- Três YAMLs de exemplo na raiz (`curriculum.yaml`, `basic_model.yaml`, `my-resume.yaml`) — mover para `examples/`.
- Template do `init` hard-coded como string no código → usar `embed`.
- `go 1.24.1` no `go.mod` — ok, mas fixar também no CI.

## 4. Pontos positivos (manter)

- Separação clara `parser` → modelo `Curriculum` → interface `Generator`.
- Dependências mínimas e maduras (`yaml.v3`, `go-pdf/fpdf`, `fatih/color`).
- Layout do PDF já é ATS-friendly: coluna única, sem tabelas/imagens, texto selecionável, hífens em vez de bullets gráficos.
- Mensagens de erro de parsing claras.
- Licença MIT presente.

## 5. Avaliação do README atual

O README tem **~25 linhas e está truncado**: o bloco de código de instalação nunca é fechado. Problemas:

1. **Instalação quebrada** — placeholder `seu-usuario`, sintaxe de link dentro de bloco de código, módulo incompatível com `go install`.
2. **Falta o essencial**: uso (`build`, `init`, `validate`), flags, exemplo de YAML, referência do schema, exemplo de saída, requisitos (versão do Go), build a partir do código.
3. **Promete o que não entrega**: "otimização ATS" (flag sem efeito); "um único comando" para múltiplos formatos (hoje é um formato por execução).
4. Sem badges (CI, Go Report Card, licença, release), sem screenshot/preview do PDF, sem seção de contribuição, sem roadmap.

A estrutura proposta está em [`008-readme-docs/spec.md`](008-readme-docs/spec.md).

## 6. Melhorias recomendadas (priorizadas)

### Fase 1 — Funcionar de verdade (P0)

1. Corrigir UTF-8 no PDF (C1).
2. Paridade de conteúdo MD/TXT com o PDF (C2, A5).
3. Skills: nunca descartar; normalizar níveis (C3).
4. EOF no modo interativo (C4).
5. Renomear módulo para `github.com/fabiodrneles/cv-craft` e corrigir instalação (C5).
6. Extrair um único serviço de build (`internal/app`) usado por CLI e UI (A2).

### Fase 2 — Confiável (P1)

1. Testes: unitários do parser + golden files dos geradores + teste de CLI (A8).
2. CI no GitHub Actions (`vet`, `test`, `golangci-lint`), `.gitignore`, remover binários (A8).
3. CLI: `--force`, `--version`, `--help` com exit 0, `validate` e `init` como subcomandos não interativos, `init` sem sobrescrever (A3, A4, M3).
4. Validação agregada + modo estrito para chaves desconhecidas (A6, A7).
5. Definir comportamento real de `-ats` ou remover a flag (A1).

### Fase 3 — Profissional (P2)

1. i18n das seções (`pt-BR`/`en`) (M1).
2. Metadados do PDF, quebra de página inteligente (M6, M8).
3. `--format all` / múltiplos formatos numa execução.
4. README novo, `examples/`, CHANGELOG, release com GoReleaser (binários Win/macOS/Linux).

## 7. Decisões

Todas as recomendações abaixo foram **aprovadas** pelo dono do projeto em 2026-10-01 e implementadas na Fase 1.

| # | Pergunta | Decisão |
|---|---|---|
| D1 | Idioma das saídas: fixo PT, fixo EN ou configurável? | Configurável via `meta.locale` no YAML + `--lang`, padrão `pt-BR` |
| D2 | Manter o modo interativo? | Manter, mas como camada fina sobre os mesmos comandos; fora do caminho crítico |
| D3 | O que `-ats` deve fazer? | Remover a flag: o layout padrão já é ATS. Se mantida, definir: sem cor, sem itálico, rótulos explícitos |
| D4 | Fonte do PDF: tradutor cp1252 (leve, só Latin-1) ou TTF embutida (UTF-8 completo, +~300-700 KB no binário)? | TTF embutida via `embed` (ex.: Liberation Sans / DejaVu Sans) |
| D5 | Suportar múltiplos templates? | Fora do escopo do v1; remover o `templates list` fictício |
| D6 | Nome dos níveis de skill e se aparecem no PDF | Aceitar `expert/advanced/proficient/intermediate/beginner` + fallback "Outras" |
