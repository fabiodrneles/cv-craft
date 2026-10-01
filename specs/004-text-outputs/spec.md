# 004 — Saídas Markdown e Texto

- **Prioridade:** P0
- **Status:** Done
- **Código afetado:** `internal/generator/markdown.go`, `internal/generator/text.go`
- **Resolve:** C2, A5, M7

## Contexto

O `.txt` é o formato de máxima compatibilidade com ATS (colar em formulários); o `.md` serve para GitHub/portfólio. Hoje ambos renderizam apenas `experience[].description` — campo que nenhum exemplo preenche — então as experiências saem **vazias** (ver `curriculum.md` versionado).

## Requisitos funcionais (comuns)

- **FR-1 Paridade de conteúdo:** MD e TXT MUST conter todo campo preenchido que o PDF contém (spec 003, FR-2..FR-5), na mesma ordem de seções.
- **FR-2** Campos vazios MUST ser omitidos junto com seus separadores (sem `|  |`, sem `[LinkedIn](https://)`).
- **FR-3** Contato inclui `location`.
- **FR-4** Seções Certificados e Idiomas MUST existir quando houver dados.
- **FR-5** Rótulos das seções vêm do i18n (spec 006).
- **FR-6** Final de linha `\n`, arquivo termina com uma única quebra de linha, UTF-8 sem BOM.
- **FR-7** A paridade é garantida por helpers de conteúdo compartilhados (`internal/generator/generator.go`) e pelo teste de paridade (AC-5). Uma camada intermediária (ex.: `render.Document`) continua sendo uma opção futura, caso surjam novos formatos.

## Markdown

- **FR-MD-1** Nome como `#`, título profissional em negrito no mesmo parágrafo do contato, seções como `##`, experiências como `### Cargo — Empresa`. Linhas de metadados (local, período) sem ênfase, para não violar a regra MD036 do markdownlint.
- **FR-MD-2** Links em sintaxe Markdown com URL normalizada (`https://`).
- **FR-MD-3** Texto do usuário MUST ter escape de caracteres especiais de Markdown (`\`, `` ` ``, `*`, `_`, `[`, `]`, `<`; e `-`, `+`, `#`, `>`, `1.` quando iniciam uma linha). `|` não é escapado: fora de tabelas ele não altera a renderização.
- **FR-MD-4** Responsabilidades e conquistas como listas `-`.
- **FR-MD-5** Saída SHOULD passar no `markdownlint` com configuração padrão (exceto MD013 — tamanho de linha).

## Texto

- **FR-TXT-1** Sem marcação: títulos de seção em MAIÚSCULAS seguidos de linha de `=`/`-`.
- **FR-TXT-2** Listas com `-` seguido de espaço.
- **FR-TXT-3** URLs escritas por extenso.
- **FR-TXT-4** Quebra de linha em 100 colunas é **opcional** (MAY); padrão: não quebrar, para não atrapalhar colagem em formulários.

## Critérios de aceite

- **AC-1** Dado `examples/full.yaml`, então todas as `responsibilities`, `achievements` e `technologies` aparecem no `.md` e no `.txt`.
- **AC-2** Dado um contato sem `phone` e sem `portfolio`, então a linha de contato não contém separadores duplicados nem links vazios.
- **AC-3** Dado `summary: "Uso *Go* e [C]"`, então o `.md` renderiza o texto literal (asteriscos e colchetes escapados).
- **AC-4** Golden tests: `testdata/full.md` e `testdata/full.txt`.
- **AC-5** Teste de paridade: o conjunto de strings não vazias do YAML está contido na saída de cada formato.
