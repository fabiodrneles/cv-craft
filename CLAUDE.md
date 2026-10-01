# CLAUDE.md

Guia rápido para agentes (Claude Code) trabalharem neste repositório sem redescobrir o projeto a cada sessão. O processo completo está na [spec 009](specs/009-delivery-process/spec.md) e na skill [`sdd-delivery`](https://github.com/fabiodrneles/sdd-kit) (plugin do sdd-kit, habilitado em `.claude/settings.json`).

## Retomar o trabalho (sessão nova ou contexto perdido)

1. Leia o **comentário "Estado da fase"** mais recente no épico aberto (issues com o label `épico`): ele lista os PRs, o estado do CI, as decisões e o próximo passo.
2. Liste os **PRs abertos** e o CI de cada um, e as **issues abertas** da fase.
3. Continue do próximo passo registrado. Não refaça análise que já está em specs, issues ou PRs.

O estado do trabalho vive no GitHub, e não na conversa. Abra o ticket e o PR assim que a tarefa começar e terminar, e atualize o comentário de estado do épico a cada marco.

## O projeto

CLI em Go (`go 1.26`) que gera currículos em PDF, Markdown e texto a partir de um YAML.

| Caminho | O que tem |
|---|---|
| `main.go` | Só monta o `cli.Env` (stdin/stdout, TTY, versão via `-ldflags`) |
| `internal/cli` | Subcomandos, flags, exit codes, modo interativo (`shell.go`), `--watch` (`watch.go`) |
| `internal/app` | Serviço único de build, validate e init; escrita atômica |
| `internal/resume` | Modelo, parsing estrito, validação, JSON Schema gerado (`schema.go`) |
| `internal/generator` | PDF (fonte embutida), Markdown, texto; golden files em `testdata/` |
| `internal/i18n` | Títulos em `pt-BR` e `en` |
| `examples/` | YAMLs de exemplo; `minimal.yaml` é o modelo do `init` |
| `specs/` | Specs SDD (`NNN-nome/spec.md`), `ROADMAP.md`, `ANALYSIS.md` |
| `schema/cv-craft.schema.json` | JSON Schema publicado (gerado; não edite à mão) |

## Comandos

```bash
make ci        # lint + testes com race + cobertura ≥ 80% + smoke test (rode antes de todo push)
make docs      # markdownlint + comandos dos READMEs + links (se mexeu em .md)
make golden    # regrava os golden files após mudança intencional nas saídas
make schema    # regrava o JSON Schema após mudar o modelo
go test -run NomeDoTeste ./internal/pacote   # um teste só
```

Numa sessão na web, o hook `.claude/hooks/session-start.sh` já instala o golangci-lint na versão do CI e o poppler.

## Convenções

- **Idioma:** specs, issues, PRs e documentação em português; commits e código (identificadores) em inglês.
- **Commits:** Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`, `ci:`, `chore:`; `!` para mudança incompatível).
- **Branch:** uma por ticket, `<tipo>/<nº-da-issue>-<descrição>`, a partir da `main`.
- **PR:** começa com `Closes #N · Épico #M · Spec NNN` e segue o template.
- **Arquivos de status** (status das specs, checkboxes do ROADMAP, CHANGELOG) só mudam no PR de fechamento da fase.
- **Merge, tag e release** são do dono, salvo delegação explícita para uma rodada.

## Armadilhas já conhecidas

- **`.gitignore` ignora `/*.md` e `/*.yaml` na raiz** (protege currículos pessoais). Um arquivo novo na raiz precisa de exceção (`!/ARQUIVO.md`).
- **Blocos ` ```bash ` dos READMEs são executados** pelo `scripts/doc-commands.sh`. Um comando que não termina (ex.: `--watch`) vai num bloco ` ```text `.
- **Blocos ` ```yaml ` de `docs/schema.md` e dos READMEs precisam ser currículos válidos** (há teste). Trechos que não são currículo vão em ` ```text `.
- **Os exemplos começam com o comentário `# yaml-language-server: $schema=...`** (há teste).
- **Conflito em `go.mod`/`go.sum`:** resolva regenerando com `go mod tidy`, nunca à mão.
- **No Windows, o poppler termina as linhas com CRLF**; os testes normalizam a saída.

## Economia de uso

- Leia trechos (`sed -n 'a,bp'`, `grep -n`) em vez de arquivos inteiros, e não releia o que já leu nesta sessão.
- Para conferir CI, peça só o resumo das conclusões dos checks. Para investigar uma falha, leia o fim do log do job que falhou.
- Junte a validação num comando só (`make ci`) em vez de rodar etapas avulsas.
- Detalhes vão nos PRs e nas issues; no chat, só o resumo e o próximo passo.
