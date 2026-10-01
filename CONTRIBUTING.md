# Contribuindo com o CV-Craft

Obrigado pelo interesse! Este guia descreve como o projeto é desenvolvido: o processo é o mesmo para quem mantém o projeto e para quem contribui pela primeira vez.

## Princípios

O projeto segue **Spec Driven Development (SDD)**: nenhuma mudança de comportamento entra no código sem uma spec que a descreva e critérios de aceite que a verifiquem.

- As specs ficam em [`specs/`](specs/README.md); os princípios inegociáveis, em [`specs/constitution.md`](specs/constitution.md).
- Cada critério de aceite (`AC-*`) vira ao menos um teste automatizado.
- O CI precisa estar verde para qualquer merge.

## Fluxo de trabalho

```text
spec → ticket → branch → testes + código → PR → revisão → merge → fechamento da fase
```

### 1. Fases e tickets

O trabalho é organizado em **fases** (ver [`specs/ROADMAP.md`](specs/ROADMAP.md)). Cada fase tem uma issue **épico** (label `épico`), e cada tarefa é uma **sub-issue** do épico, com:

- contexto (o problema);
- o que fazer;
- critérios de aceite verificáveis;
- spec(s) afetada(s).

Labels: `fase-N`, `tipo:feature|docs|ci|teste|chore` (defeitos usam o label padrão `bug`) e prioridade `P1` (alta) a `P3` (baixa).

Mudanças pequenas e óbvias (erro de digitação, link quebrado) podem ir direto para PR, sem ticket.

### 2. Branches

Uma branch por ticket, criada a partir da `main` (ou da branch da fase em andamento, quando a fase anterior ainda não foi mergeada):

```text
<tipo>/<nº-da-issue>-<descrição-curta>
```

Exemplos: `feat/12-watch-mode`, `fix/31-acentos-no-pdf`, `docs/7-schema`, `ci/6-goreleaser`.

### 3. Commits

[Conventional Commits](https://www.conventionalcommits.org/pt-br/), em inglês, no imperativo:

```text
feat: add --watch mode to build
fix: keep list label with its first item on page breaks
docs: document YAML schema
test: cover overwrite prompt on EOF
ci: validate goreleaser config on every PR
chore: raise minimum Go version to 1.26
```

Mudanças incompatíveis usam `!` (`feat!: ...`) e explicam o impacto no corpo do commit. O changelog das releases é gerado a partir desses prefixos.

### 4. Pull requests

- Um PR por ticket. A descrição começa com `Closes #nº · Épico #nº · Spec NNN`; o template já traz essa linha.
- Diga quais specs e critérios de aceite o PR implementa e **como foi testado**.
- Se o comportamento mudou, atualize a spec no mesmo PR.
- Mantenha o PR focado: o que não é do ticket vira outro ticket.
- PRs que dependem de outro PR ainda não mergeado são **empilhados** (base = branch do PR anterior) e dizem isso no topo da descrição.

### 5. Revisão e merge

- Todo PR precisa de CI verde e da aprovação de um mantenedor (ver [`CODEOWNERS`](.github/CODEOWNERS)).
- Os PRs de uma fase são revisados e mergeados na **ordem sugerida no épico**.
- PRs de fase inteira (que outros PRs usam como base) são mergeados com **merge commit**, para que os PRs empilhados sejam redirecionados à `main` sem rebase. Os demais podem usar **squash**.

### 6. Fechamento da fase

Depois que todos os PRs da fase forem mergeados, um **PR de fechamento** atualiza o status das specs, o [`ROADMAP.md`](specs/ROADMAP.md) e o `CHANGELOG.md`. Assim, os PRs da fase não entram em conflito por editarem os mesmos arquivos. Em seguida, cria-se a tag da versão.

O processo completo é normativo na [spec 009](specs/009-delivery-process/spec.md). Quem usa o Claude Code tem a skill [`sdd-delivery`](.claude/skills/sdd-delivery/SKILL.md), que aplica o mesmo fluxo neste e em outros repositórios.

## Ambiente de desenvolvimento

Requisitos:

- **Go** na versão declarada no [`go.mod`](go.mod) (ou mais nova);
- **make** e **bash** (no Windows, o Git Bash serve);
- **[golangci-lint](https://golangci-lint.run/welcome/install/)** v2, para `make lint`;
- **poppler** (`poppler-utils` no Debian/Ubuntu, `poppler` no Homebrew), para os testes que leem o texto dos PDFs. Sem ele, esses testes são pulados localmente; no CI eles são obrigatórios.

Comandos:

| Comando | O que faz |
|---|---|
| `make build` | Compila `./bin/cv-craft` |
| `make test` | Testes unitários, golden files e CLI |
| `make race` | Testes com race detector |
| `make cover` | Cobertura, com mínimo de 80% em `internal/` |
| `make lint` | `go vet` + `golangci-lint` |
| `make smoke` | Smoke test de ponta a ponta com o binário real |
| `make golden` | Regrava os golden files após uma mudança **intencional** nas saídas |
| `make bench` | Benchmarks de parsing e de cada formato |
| `make docs` | Lint de Markdown, comandos da documentação e links (requer Node) |
| `make release-snapshot` | Gera os arquivos da release em `./dist`, sem publicar (requer GoReleaser) |
| `make ci` | Verificações de código do CI; rode antes de abrir o PR (e `make docs` se mexeu em documentação) |

### Golden files

As saídas de `examples/*.yaml` em todos os formatos ficam em `internal/generator/testdata/`. Se você mudou a renderização de propósito, rode `make golden` e **revise o diff dos golden files** no PR: ele mostra exatamente o que mudou para o usuário.

## Reportando bugs

Abra uma issue com o formulário **Bug**. Inclua a versão (`cv-craft version`), o comando executado e um YAML mínimo que reproduza o problema. **Não publique dados pessoais reais**: troque nome, e-mail e telefone por valores fictícios.
