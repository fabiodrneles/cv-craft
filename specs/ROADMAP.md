# Roadmap e tarefas

Cada tarefa referencia a spec e os critérios de aceite que ela fecha. Ordem sugerida dentro de cada fase = ordem da lista.

## Fase 0 — Decisões (antes de codar)

- [x] Responder D1–D6 em [ANALYSIS.md §7](ANALYSIS.md#7-decisões) e mover specs para `Approved`.

## Fase 1 — Funcionar de verdade (P0) → publicada na `v0.2.0` (sem tag própria)

- [x] **T1** Renomear módulo, `.gitignore`, remover binários/saídas, criar `examples/` — 007 FR-1..5
- [x] **T2** Testes + golden files — 007 FR-6..8 (goldens gravados já com o comportamento corrigido, pois o comportamento antigo era o bug)
- [x] **T3** Separar `Parse` e `Validate`, erros agregados com caminho, `KnownFields`, normalização de níveis/URLs — 001 todos
- [x] **T4** Extrair `internal/app` (Build/Validate/Init) e fazer CLI e UI usarem — 002 FR-2, 005 FR-3
- [x] **T5** PDF: UTF-8 — 003 FR-1, AC-1/2
- [x] **T6** PDF: conteúdo completo (achievements, description, courses, cert URL), skills sem descarte — 003 FR-2..5
- [x] **T7** MD/TXT com paridade (helpers compartilhados + teste de paridade) — 004 todos
- [x] **T8** CLI com `flag`, exit codes, `--force`, `--help`, `version`, escrita atômica — 002 todos
- [x] **T9** Modo interativo: EOF, vazio, não-TTY, remover recursos fictícios — 005 todos
- [x] **T10** CI completo: lint, testes em 3 SOs, race, cobertura ≥ 80%, smoke test, cross-build, govulncheck — 007 FR-11..12

## Fase 2 — Confiável (P1) → `v0.2.0`

- [x] **T11** Decisão D3 aplicada à flag `-ats` — 003 FR-9 *(antecipada)*
- [x] **T12** Metadados do PDF e links clicáveis — 003 FR-6..7 *(antecipada)*
- [x] **T13** Quebra de página inteligente — 003 FR-8 *(antecipada)*
- [x] **T14** `--format all` — 002 FR-7 *(antecipada)*
- [x] **T15** ~~README novo~~ (feito na Fase 1) + `docs/schema.md` + `CONTRIBUTING.md` — 008 (#5, #7)
- [x] **T16** GoReleaser + CHANGELOG — 007 FR-13..16 (#6)
- [x] Go 1.26 como versão mínima (#4); CI de documentação (#8); benchmarks dos NFRs (#9); spec 009 do processo (#19); skill `sdd-delivery` versionada (#25)

## Fase 3 — Profissional (P2) → `v1.0.0`

- [x] **T17** i18n `pt-BR`/`en` — 006 *(antecipada)*
- [ ] **T18** README em inglês (~~prévia em PNG, badges~~ feitos na Fase 1)
- [ ] **T19** (futuro) JSON Schema para autocompletar no VS Code; `--watch`; templates adicionais
