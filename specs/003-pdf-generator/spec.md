# 003 — Gerador de PDF

- **Prioridade:** P0
- **Status:** Done
- **Código afetado:** `internal/generator/pdf.go`
- **Resolve:** C1, C3, A1, A5, M5, M6, M8

## Contexto

O PDF é o formato principal: é o que vai para recrutadores e para ATS. O layout atual (coluna única, Helvetica, sem gráficos) é um bom ponto de partida, mas há corrupção de acentos, conteúdo descartado e uma flag sem efeito.

## Estado atual (verificado)

- `São Paulo` sai como `SÃ£o Paulo` (fontes core cp1252 recebendo bytes UTF-8).
- `achievements`, `description`, `relevant_courses`, `certificates[].url` nunca aparecem.
- Skills agrupadas por nível; níveis desconhecidos são descartados.
- `IsATS` nunca é lido.
- Quebra de página manual quando `GetY() > 250`.

## Layout (ordem das seções)

1. Nome (16pt bold) · Título profissional (12pt bold)
2. Linha de contato: localização · telefone · e-mail
3. Linha de links: LinkedIn · GitHub · Portfólio (URLs completas, clicáveis)
4. Resumo
5. Habilidades
6. Experiência profissional
7. Formação
8. Certificados (se houver)
9. Idiomas (se houver)

Seções vazias não são renderizadas.

## Requisitos funcionais

- **FR-1** Todo texto MUST ser renderizado corretamente em UTF-8 para, no mínimo, Latin-1 + Latin Extended-A (pt, es, fr, de, it). Implementação conforme decisão D4: fonte TTF embutida via `embed` (preferida) ou `UnicodeTranslatorFromDescriptor("")`.
- **FR-2** Experiência MUST renderizar, nesta ordem e quando presentes: cargo; empresa, local, modalidade, período; `description`; `responsibilities` (lista); `achievements` (lista, rótulo próprio); `technologies` (linha única).
- **FR-3** Formação MUST renderizar `thesis` e `relevant_courses` quando presentes.
- **FR-4** Certificados MUST renderizar `url` quando presente (como link).
- **FR-5** Habilidades: cada categoria (`skills[].name`) em uma linha `Categoria: kw1, kw2, …`, **na ordem do YAML**. Nível, se exibido, entre parênteses após a categoria (decisão D6). Nenhuma skill é descartada.
- **FR-6** Links (e-mail, LinkedIn, GitHub, portfólio, certificados) MUST ser clicáveis (`pdf.LinkString`).
- **FR-7** O PDF MUST conter metadados: `Title` = "Nome — Título profissional", `Author` = nome, `Subject`, `Keywords` = união das keywords das skills, `Creator` = "cv-craft vX.Y.Z".
- **FR-8** Quebra de página: título de seção nunca fica sozinho no fim da página; o cabeçalho de uma experiência (cargo + empresa) nunca é separado do primeiro item. Remover o limiar fixo de 250 mm.
- **FR-9** A flag `-ats` não existe (decisão D3): foi aceita sem efeito, com aviso de depreciação, na v0.2.0, e removida na v1.0.0. Passá-la resulta em "flag desconhecida" (exit 2, spec 002).
- **FR-10** O gerador MUST ser determinístico em testes: data de criação injetável (ou `SOURCE_DATE_EPOCH`).
- **FR-11** O gerador não escreve em stdout/stderr; retorna erro.
- **FR-12** Rótulos das seções vêm do pacote de i18n (spec 006), não são literais no gerador.

## Requisitos não-funcionais

- **NFR-1** Texto 100% extraível: `pdftotext` do PDF gerado contém todo o conteúdo do YAML, na ordem visual.
- **NFR-2** Geração de um CV de 2 páginas em < 200 ms.
- **NFR-3** Arquivo < 300 KB (com fonte embutida em subset, se aplicável).
- **NFR-4** Margens de 15–20 mm; corpo ≥ 10 pt; contraste do texto ≥ cinza 80.
- **NFR-5** Verificação: NFR-1 pelos golden files e pelo teste de paridade; NFR-2 por `TestPDFGenerationTime` (margem de 5×) e `BenchmarkGenerate/pdf`; NFR-3 por `TestPDFSize` (exemplo completo e um CV com 18 experiências).

## Critérios de aceite

- **AC-1** Dado `location: "São Paulo, Brazil"`, quando gerar PDF, então `pdftotext` contém exatamente `São Paulo`.
- **AC-2** Dado um YAML com `ç, ã, é, ñ, ü, ß`, todos aparecem corretamente no texto extraído.
- **AC-3** Dado `achievements` em uma experiência, então eles aparecem no texto extraído.
- **AC-4** Dado `level: advanced`, as keywords aparecem.
- **AC-5** `pdfinfo` mostra `Title` e `Author` preenchidos.
- **AC-6** Dado um CV que transborda para a página 2, nenhum título de seção é a última linha da página 1.
- **AC-7** Golden test: texto extraído de `examples/full.yaml` é igual a `testdata/full.pdf.txt`.
- **AC-8** Duas gerações com a mesma data injetada produzem bytes idênticos.

## Fora de escopo

- Múltiplos templates/temas visuais (decisão D5).
- Foto, ícones, colunas.

## Decisões

- D3 — `-ats` foi aceita e ignorada, com aviso de depreciação em stderr, na v0.2.0, e removida na v1.0.0 (#14): o layout padrão já é otimizado para ATS.
- D4 — Liberation Sans (Regular, Bold, Italic; SIL OFL 1.1) embutida via `embed`, com subset no PDF. Cobre Latin, Latin Extended, grego e cirílico. PDF típico: ~50 KB.
- D6 — nível exibido traduzido entre parênteses (`Backend (Avançado): Go, …`).
- Links do cabeçalho são exibidos sem `https://` (o destino clicável é a URL completa) e rótulo e link nunca são separados por quebra de linha.
- Lista de itens usa marcador `•` com recuo deslocado; rótulos de lista ("Responsabilidades:") também nunca ficam órfãos no fim da página.
