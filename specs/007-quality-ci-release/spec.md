# 007 — Qualidade, CI e release

- **Prioridade:** P1 (itens 1–3 são P0 por bloquearem a instalação)
- **Status:** Draft
- **Resolve:** C5, A8

## Contexto

Não há testes, CI, `.gitignore` nem processo de release. O repositório versiona um binário Windows de 5,8 MB e arquivos gerados. O nome do módulo (`cv-craft`) impede `go install`.

## Requisitos

### Repositório
- **FR-1** Renomear o módulo para `github.com/fabiodrneles/cv-craft` e atualizar imports.
- **FR-2** Mover `main.go` para `cmd/cv-craft/main.go` **ou** manter na raiz — decisão única documentada no README. (Raiz é aceitável e mantém `go install github.com/fabiodrneles/cv-craft@latest`.)
- **FR-3** Adicionar `.gitignore` (binários, `*.exe`, `dist/`, saídas geradas na raiz, `.DS_Store`, `.idea/`, `.vscode/`).
- **FR-4** Remover do controle de versão: `cv-craft.exe`, `meu_cv.pdf`, `my-resume.pdf`, `curriculum.md`, `my-resume.yaml`.
- **FR-5** Mover exemplos para `examples/` (`minimal.yaml`, `full.yaml`, e um `en.yaml`), e uma prévia do PDF para `docs/` (PNG para o README).

### Testes
- **FR-6** Unitários do parser/validação cobrindo os AC da spec 001.
- **FR-7** Golden files dos geradores em `internal/generator/testdata/`, com flag `-update` para regenerar.
- **FR-8** Teste do PDF via extração de texto (lib Go pura, ex.: `ledongthuc/pdf`, só em testes) — sem depender de `pdftotext` no CI.
- **FR-9** Teste de CLI de ponta a ponta (compila o binário em `TestMain` ou usa `testscript`) cobrindo os exit codes da spec 002.
- **FR-10** Cobertura mínima: 80% em `internal/`.

### CI (GitHub Actions)
- **FR-11** Workflow em PR e push para `main`: `go mod tidy` sem diff, `gofmt -l` vazio, `go vet`, `golangci-lint`, `go test -race -cover ./...`.
- **FR-12** Matriz: `ubuntu-latest`, `windows-latest`, `macos-latest` (o autor usa Windows — `.exe` versionado).

### Release
- **FR-13** GoReleaser disparado por tag `v*`: binários linux/darwin/windows × amd64/arm64, checksums, changelog automático.
- **FR-14** Versão injetada via `-ldflags "-X main.version=… -X main.commit=… -X main.date=…"` (spec 002 FR-10).
- **FR-15** `CHANGELOG.md` no formato Keep a Changelog; versionamento SemVer; primeira release `v0.1.0`.

## Critérios de aceite

- **AC-1** Em máquina limpa com Go ≥ 1.24: `go install github.com/fabiodrneles/cv-craft@latest && cv-craft version` funciona.
- **AC-2** `git ls-files` não contém `.exe` nem `.pdf` fora de `docs/`/`testdata/`.
- **AC-3** PR com `gofmt` quebrado ou teste falhando fica vermelho.
- **AC-4** `git tag v0.1.0 && git push --tags` publica release com 6 binários.
- **AC-5** `go test ./...` passa em Linux, macOS e Windows.
