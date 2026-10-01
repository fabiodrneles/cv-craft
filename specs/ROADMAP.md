# Roadmap e tarefas

Cada tarefa referencia a spec e os critérios de aceite que ela fecha. Ordem sugerida dentro de cada fase = ordem da lista.

## Fase 0 — Decisões (antes de codar)

- [ ] Responder D1–D6 em [ANALYSIS.md §7](ANALYSIS.md#7-decisões-em-aberto-precisam-do-dono-do-projeto) e mover specs para `Approved`.

## Fase 1 — Funcionar de verdade (P0) → `v0.1.0`

- [ ] **T1** Renomear módulo, `.gitignore`, remover binários/saídas, criar `examples/` — 007 FR-1..5
- [ ] **T2** Esqueleto de testes + golden files com o comportamento atual (rede de segurança) — 007 FR-6..8
- [ ] **T3** Separar `Parse` e `Validate`, erros agregados com caminho, `KnownFields`, normalização de níveis/URLs — 001 todos
- [ ] **T4** Extrair `internal/app` (Build/Validate/Init) e fazer CLI e UI usarem — 002 FR-2, 005 FR-3
- [ ] **T5** PDF: UTF-8 — 003 FR-1, AC-1/2
- [ ] **T6** PDF: conteúdo completo (achievements, description, courses, cert URL), skills sem descarte — 003 FR-2..5
- [ ] **T7** Camada `render.Document` + MD/TXT com paridade — 004 todos
- [ ] **T8** CLI com `flag`, exit codes, `--force`, `--help`, `version`, escrita atômica — 002 todos
- [ ] **T9** Modo interativo: EOF, vazio, não-TTY, remover recursos fictícios — 005 todos
- [ ] **T10** CI (lint + test + matriz de SO) — 007 FR-11..12

## Fase 2 — Confiável (P1) → `v0.2.0`

- [ ] **T11** Decisão D3 aplicada à flag `-ats` — 003 FR-9
- [ ] **T12** Metadados do PDF e links clicáveis — 003 FR-6..7
- [ ] **T13** Quebra de página inteligente — 003 FR-8
- [ ] **T14** `--format all` — 002 FR-7
- [ ] **T15** README novo + `docs/schema.md` + `CONTRIBUTING.md` — 008
- [ ] **T16** GoReleaser + CHANGELOG — 007 FR-13..15

## Fase 3 — Profissional (P2) → `v1.0.0`

- [ ] **T17** i18n `pt-BR`/`en` — 006
- [ ] **T18** README em inglês, prévia em PNG, badges
- [ ] **T19** (futuro) JSON Schema para autocompletar no VS Code; `--watch`; templates adicionais
