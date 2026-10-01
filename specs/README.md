# Specs — CV-Craft

Este diretório organiza o desenvolvimento do CV-Craft em **Spec Driven Development (SDD)**: nenhuma mudança de comportamento entra no código sem uma spec que a descreva e critérios de aceite que a verifiquem.

## Fluxo

```text
spec.md (O QUÊ / POR QUÊ)  →  revisão  →  testes a partir dos critérios de aceite  →  implementação  →  status: Done
```

1. **Especificar** — requisitos (`FR-*`), não-funcionais (`NFR-*`) e critérios de aceite (`AC-*`) no formato Dado/Quando/Então.
2. **Resolver decisões** — itens em "Decisões em aberto" precisam ser respondidos antes de implementar.
3. **Testar primeiro** — cada `AC-*` vira ao menos um teste automatizado (unitário, golden file ou de CLI).
4. **Implementar** — PR referencia os IDs (ex.: `implements 003/FR-1, AC-2`).
5. **Fechar** — atualizar o status abaixo e a seção "Estado atual" da spec.

Convenções: `MUST`/`SHOULD`/`MAY` seguem a RFC 2119. Prioridades: **P0** (bloqueia uso real), **P1** (confiabilidade), **P2** (polimento).

## Documentos

| Documento | Conteúdo |
|---|---|
| [ANALYSIS.md](ANALYSIS.md) | Relatório da verificação: bugs, evidências, melhorias priorizadas |
| [constitution.md](constitution.md) | Princípios inegociáveis do projeto |
| [ROADMAP.md](ROADMAP.md) | Tarefas por fase, ligadas às specs |

## Specs

| ID | Spec | Prioridade | Status |
|---|---|---|---|
| 001 | [Schema YAML e validação](001-yaml-schema/spec.md) | P0 | Done |
| 002 | [Interface de linha de comando](002-cli/spec.md) | P0 | Done |
| 003 | [Gerador de PDF](003-pdf-generator/spec.md) | P0 | Done |
| 004 | [Saídas Markdown e Texto](004-text-outputs/spec.md) | P0 | Done |
| 005 | [Modo interativo](005-interactive-mode/spec.md) | P1 | Done |
| 006 | [Internacionalização das saídas](006-i18n/spec.md) | P2 | Done |
| 007 | [Qualidade, CI e release](007-quality-ci-release/spec.md) | P1 | In Progress (falta release) |
| 008 | [README e documentação](008-readme-docs/spec.md) | P1 | In Progress (falta `CONTRIBUTING.md`) |

Status possíveis: `Draft` → `Approved` → `In Progress` → `Done`.
