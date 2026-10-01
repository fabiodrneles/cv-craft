# 010 — JSON Schema do YAML

- **Prioridade:** P2
- **Status:** In Progress
- **Código afetado:** `internal/resume/schema.go`, `internal/cli` (comando `schema`), `schema/cv-craft.schema.json`, `examples/`
- **Resolve:** #11 (e o item "JSON Schema publicado" de "Fora de escopo" da [spec 001](../001-yaml-schema/spec.md))

## Contexto

Hoje quem escreve o currículo só descobre um campo obrigatório vazio ou uma chave com erro de digitação ao rodar `cv-craft validate`. Editores com a extensão YAML (VS Code com a extensão da Red Hat, IDEs JetBrains) leem um JSON Schema e passam a autocompletar os campos, mostrar a descrição de cada um e marcar os erros enquanto se digita.

## Requisitos funcionais

- **FR-1** O CV-Craft MUST publicar um JSON Schema (draft-07, o mais suportado pelos editores) do YAML do currículo em `schema/cv-craft.schema.json`, com `$id` igual à URL pública do arquivo na `main`.
- **FR-2** O schema MUST ser **gerado a partir do modelo Go** (`resume.JSONSchema`, decisão D9), e não escrito à mão. `make schema` regrava o arquivo, e um teste falha se o arquivo publicado diferir do gerado.
- **FR-3** O schema MUST aceitar e rejeitar os mesmos documentos que `Parse` + `Validate` (spec 001):
  - chaves desconhecidas são rejeitadas em todos os níveis (`additionalProperties: false`);
  - campos de texto obrigatórios não podem faltar nem ser vazios, nulos ou só espaços;
  - `skills`, `experience`, `education` e `skills[].keywords` precisam de pelo menos um item, e as `keywords` de pelo menos um item não vazio;
  - campos de texto aceitam qualquer escalar (o YAML converte `2020` em `"2020"`), e campos opcionais vazios valem como ausentes;
  - `meta.locale` aceita as mesmas variações que o pacote `i18n` (`pt`, `PT_br`, `en-US`…);
  - nível de habilidade desconhecido **não** é rejeitado, porque para o CV-Craft é só um aviso.
- **FR-4** Os exemplos de `examples/` e o modelo do `cv-craft init` MUST começar com `# yaml-language-server: $schema=<URL do FR-1>`.
- **FR-5** `cv-craft schema` MUST imprimir o schema em stdout (exit 0), para uso offline; argumentos extras são uso incorreto (exit 2).
- **FR-6** Cada campo MUST ter uma descrição em português no schema, mostrada pelo editor. Há teste para campo sem descrição e para descrição de um campo que não existe.
- **FR-7** O README e `docs/schema.md` MUST explicar como usar o schema no editor, inclusive com uma cópia local.

## Requisitos não funcionais

- **NFR-1** O schema gerado é determinístico: a mesma versão produz os mesmos bytes, com as chaves na ordem dos campos do modelo.
- **NFR-2** A biblioteca de validação de JSON Schema (`github.com/santhosh-tekuri/jsonschema/v6`) é usada **só nos testes** e não entra no binário (Constituição §10).

## Critérios de aceite

- **AC-1** Abrir `examples/full.yaml` no VS Code com a extensão YAML da Red Hat mostra o autocompletar dos campos e marca como erro um campo desconhecido (verificação manual, registrada no PR).
- **AC-2** `TestSchemaMatchesValidate`: para cada caso (os três exemplos, campos obrigatórios ausentes, vazios, nulos ou numéricos, opcionais nulos, listas vazias, níveis e idiomas válidos e inválidos, chaves desconhecidas em cada nível, estrutura errada), o schema e a validação do Go dão o mesmo resultado. O teste exige pelo menos 20 casos aceitos e 20 rejeitados.
- **AC-3** `TestSchemaFileUpToDate`: `schema/cv-craft.schema.json` é igual a `resume.JSONSchema()`.
- **AC-4** `TestExamplesReferenceSchema`: os três exemplos começam com o comentário do FR-4.
- **AC-5** `cv-craft schema` imprime o mesmo JSON (exit 0); `cv-craft schema x` sai com 2.

## Fora de escopo

- Publicar o schema no [SchemaStore](https://www.schemastore.org/), o que dispensaria o comentário na primeira linha. Fica para depois da v1.0.0, com um padrão de nome de arquivo (ex.: `*.cv.yaml`).
- Validar o formato do e-mail exatamente como o `net/mail`: o schema usa uma aproximação, e o `cv-craft validate` continua sendo a referência.

## Decisões

- **D9 — schema gerado a partir do código.** Um gerador lê as structs de `internal/resume` e as tabelas de obrigatoriedade e descrições de `schema.go`. Os testes garantem que essas tabelas não divergem de `validate.go`, então um campo novo exige atualizar o modelo, a validação e a descrição, e o CI aponta o que faltar. A alternativa, um schema escrito à mão, exigiria editar dois lugares a cada mudança.
- **Draft-07 em vez de 2020-12:** é a versão com suporte mais amplo na extensão YAML do VS Code e nas IDEs JetBrains.
- **URL na `main`:** o comentário aponta para o arquivo na `main`, que acompanha a versão mais recente. Para fixar a versão instalada, use a cópia local gerada por `cv-craft schema`.
