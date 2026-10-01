# 001 — Schema YAML e validação

- **Prioridade:** P0
- **Status:** Done
- **Código afetado:** `internal/resume` (antes `internal/parser/parser.go`)
- **Resolve:** C3, A5, A6, A7 (ver [ANALYSIS.md](../ANALYSIS.md))

## Contexto

O YAML é a fonte única da verdade (Constituição §1). Hoje o parser aceita chaves desconhecidas em silêncio, valida em dois lugares diferentes, para no primeiro erro e não define o que acontece com valores fora do esperado (ex.: `level: advanced`).

## Estado atual

- `ParseFile` valida `contact.name`, `contact.email`, `professional_title`.
- `ValidateCurriculum` valida que `skills`, `experience`, `education` não são vazios.
- Chaves desconhecidas são ignoradas.
- `Skill.Level` é texto livre; o PDF só reconhece 3 valores.

## Schema (v1)

```yaml
meta:                     # opcional
  locale: pt-BR           # ver spec 006; padrão pt-BR
contact:                  # obrigatório
  name: string            # obrigatório
  email: string           # obrigatório, formato e-mail
  phone: string
  location: string
  linkedin: string        # URL ou "linkedin.com/in/..." (normalizado para https://)
  portfolio: string       # idem
  github: string          # NOVO, opcional
professional_title: string  # obrigatório
summary: string
skills:                   # obrigatório, ≥ 1
  - name: string          # obrigatório
    level: expert|advanced|proficient|intermediate|beginner   # opcional
    keywords: [string]    # obrigatório, ≥ 1
experience:               # obrigatório, ≥ 1
  - role: string          # obrigatório
    company: string       # obrigatório
    location: string
    period: string        # texto livre ("Mar 2022 - Atual")
    work_type: string     # Remote|Hybrid|On-site (texto livre, exibido como está)
    description: string
    responsibilities: [string]
    achievements: [string]
    technologies: [string]
education:                # obrigatório, ≥ 1
  - degree: string        # obrigatório
    institution: string   # obrigatório
    location: string
    period: string
    thesis: string
    relevant_courses: string
certificates:             # opcional
  - name: string          # obrigatório
    institution: string
    date: string
    url: string
languages:                # opcional
  - language: string      # obrigatório
    level: string
```

## Requisitos funcionais

- **FR-1** O parser MUST ler o arquivo como UTF-8 e aceitar BOM.
- **FR-2** O parser MUST rejeitar chaves desconhecidas (`yaml.Decoder.KnownFields(true)`), informando linha e chave. Uma flag `--lenient` MAY rebaixar isso para aviso.
- **FR-3** Toda a validação MUST ficar em uma única função `Validate(cv) []ValidationError`, separada do parsing (`Parse` só decodifica).
- **FR-4** A validação MUST acumular **todos** os erros, cada um com caminho do campo (`experience[2].company`) e mensagem.
- **FR-5** Campos obrigatórios: os marcados acima. Strings só com espaços contam como vazias.
- **FR-6** `contact.email` MUST ter formato válido (`net/mail.ParseAddress`).
- **FR-7** `skills[].level` MUST ser normalizado (case-insensitive, trim). Valor fora da lista ⇒ **aviso** (não erro) e a skill é tratada como "sem nível". Skills nunca são descartadas.
- **FR-8** URLs (`linkedin`, `portfolio`, `github`, `certificates[].url`) SHOULD ser normalizadas para incluir `https://` quando não houver esquema, em um único helper usado por todos os geradores.
- **FR-9** O modelo MUST expor os avisos (`[]ValidationWarning`) para que a CLI os imprima em stderr.

## Requisitos não-funcionais

- **NFR-1** Parsing + validação de um YAML de 50 KB em < 50 ms.
- **NFR-2** Sem dependências novas além de `yaml.v3`.
- **NFR-3** Verificação: `TestParseValidate50KBPerformance` (limite com margem de 5×, para não ser instável em runners compartilhados), `BenchmarkParseValidate50KB` (limite exato; `make bench` e o workflow `bench.yml`) e `TestResumeDependencies` (`go list -deps`).

## Critérios de aceite

- **AC-1** Dado um YAML com `contact.name` e `skills` ausentes, quando validar, então ambos os erros são retornados juntos, com caminhos `contact.name` e `skills`.
- **AC-2** Dado `responsabilities:` (typo) em uma experiência, quando validar, então erro `experience[0]: campo desconhecido "responsabilities"` com número de linha.
- **AC-3** Dado `level: Advanced`, quando gerar qualquer formato, então as keywords da skill aparecem na saída.
- **AC-4** Dado `level: guru`, quando validar, então há 1 aviso e 0 erros, e as keywords aparecem na saída.
- **AC-5** Dado `email: "não-é-email"`, quando validar, então erro em `contact.email`.
- **AC-6** Dado um arquivo vazio ou só com comentários, então erro "arquivo vazio".
- **AC-7** Dado `name: "  "`, então erro de campo obrigatório.
- **AC-8** Os três YAMLs de `examples/` validam sem erros nem avisos.

## Fora de escopo

- JSON Schema publicado para autocompletar em editores: virou a [spec 010](../010-json-schema/spec.md).
- Datas estruturadas em `period`.

## Decisões

- D6 — níveis `expert/advanced/proficient/intermediate/beginner`, aceitando também os equivalentes em português (`especialista`, `avançado`, `proficiente`, `intermediário`, `básico`/`iniciante`). O nível aparece traduzido entre parênteses após a categoria.
- FR-2: a flag `--lenient` (MAY) não foi implementada; chaves desconhecidas são sempre erro.
