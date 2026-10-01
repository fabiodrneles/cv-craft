# 007 — Qualidade, CI e release

- **Prioridade:** P1 (itens 1–3 são P0 por bloquearem a instalação)
- **Status:** Done — release automatizada na Fase 2 (#6, #20)
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
- **FR-8** Teste do PDF via extração de texto com o `pdftotext` (poppler). **Revisado na implementação:** a lib Go pura avaliada (`ledongthuc/pdf`) trunca caracteres acima de U+00FF em fontes TTF embutidas (`—`, `•`, `Ł`), o que esconderia justamente os bugs de Unicode. O poppler é instalado no CI nos três sistemas (Linux, macOS e, desde #16, Windows via Chocolatey) e é obrigatório via `CV_CRAFT_REQUIRE_POPPLER=1`; localmente, os testes de texto do PDF são pulados se ele não existir.
- **FR-9** Teste de CLI de ponta a ponta cobrindo os exit codes da spec 002: `cli.Run` testado em processo (`internal/cli/cli_test.go`) e o binário real exercitado por `scripts/smoke.sh` em Linux, macOS e Windows.
- **FR-10** Cobertura mínima: 80% em `internal/`.

### CI (GitHub Actions)

- **FR-11** Workflow em PR e push para `main`: `go mod tidy` sem diff, `go vet`, `golangci-lint` (inclui gofmt/goimports), `go test -race` com cobertura mínima, smoke test do binário (`scripts/smoke.sh`), cross-compilação (linux/darwin/windows × amd64/arm64) e `govulncheck`.
- **FR-12** Matriz: `ubuntu-latest`, `windows-latest`, `macos-latest` com Go 1.26.x (mínimo do `go.mod`) + Ubuntu com Go `stable`.
- **FR-12b** Política de versão: o mínimo do Go é sempre a **penúltima versão suportada** pelo time do Go (que mantém só as duas últimas). Revisar a cada release do Go (fevereiro e agosto).
- **FR-12a** `Makefile` com os mesmos passos do CI (`make ci`) e Dependabot para módulos Go e GitHub Actions.

### Release

- **FR-13** GoReleaser disparado por tag `v*`: binários linux/darwin/windows × amd64/arm64, checksums, changelog automático.
- **FR-14** Versão injetada via `-ldflags "-X main.version=… -X main.commit=… -X main.date=…"` (spec 002 FR-10).
- **FR-15** `CHANGELOG.md` no formato Keep a Changelog; versionamento SemVer; primeira release publicada `v0.2.0` (ver Decisões).
- **FR-16** Implementação: `.goreleaser.yaml` (versão como `v{{ .Version }}`, `mod_timestamp` do commit para binários reproduzíveis, arquivos com `LICENSE`, `README.md`, `CHANGELOG.md` e `licenses/OFL.txt`); `.github/workflows/release.yml` chama o `ci.yml` completo (`workflow_call`) antes de publicar; o job `release-check` do CI roda `goreleaser release --snapshot` e `scripts/check-release.sh` em todo PR; `make release-snapshot` faz o mesmo localmente.

## Critérios de aceite

- **AC-1** Em máquina limpa com Go ≥ 1.26: `go install github.com/fabiodrneles/cv-craft@latest && cv-craft version` funciona.
- **AC-2** `git ls-files` não contém `.exe` nem `.pdf` fora de `docs/`/`testdata/`.
- **AC-3** PR com `gofmt` quebrado ou teste falhando fica vermelho.
- **AC-4** `git tag v0.2.0 && git push origin v0.2.0` publica release com 6 binários.
- **AC-5** `go test ./...` passa em Linux, macOS e Windows.

## Decisões

- **Primeira release publicada é a `v0.2.0`.** A automação de release só entrou na Fase 2, então o commit da Fase 1 (merge do #2) não tem o workflow de release. A Fase 1 não recebeu tag própria: a `v0.2.0` é a primeira versão publicada, com binários, e o `CHANGELOG.md` registra nela o conteúdo das duas fases.
