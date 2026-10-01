# 006 — Internacionalização das saídas

- **Prioridade:** P2
- **Status:** Draft (depende da decisão D1)
- **Código afetado:** novo `internal/i18n`, geradores, CLI
- **Resolve:** M1

## Contexto

Hoje o PDF tem títulos em inglês (`SUMMARY`, `PROFESSIONAL EXPERIENCE`), o MD/TXT em português (`Resumo`, `Experiência`), a CLI mistura PT e EN e o modo interativo é todo em EN. O mesmo YAML gera currículos em idiomas diferentes conforme o formato.

## Requisitos funcionais

- **FR-1** Rótulos de seção e rótulos fixos (ex.: "Responsabilidades", "Conquistas", "Tecnologias", "Tese") MUST vir de um catálogo `i18n` com, no mínimo, `pt-BR` e `en`.
- **FR-2** Ordem de resolução do idioma: `--lang` > `meta.locale` no YAML > `pt-BR`.
- **FR-3** Idioma desconhecido ⇒ erro de uso (exit 2) listando os suportados.
- **FR-4** Todos os formatos usam o mesmo idioma numa execução.
- **FR-5** Mensagens da CLI: um único idioma em todo o binário (recomendado: PT-BR, público-alvo do projeto), independente do idioma do currículo.
- **FR-6** O conteúdo do usuário nunca é traduzido.

## Critérios de aceite

- **AC-1** Sem `meta.locale` e sem `--lang`, todos os formatos usam "Resumo", "Experiência Profissional", "Formação".
- **AC-2** `--lang en` ⇒ "Summary", "Professional Experience", "Education" em PDF, MD e TXT.
- **AC-3** `--lang xx` ⇒ exit 2 com lista `pt-BR, en`.
- **AC-4** Teste garante que todas as chaves existem em todos os idiomas do catálogo.

## Decisões em aberto

- D1 — confirmar idiomas suportados e idioma padrão.
