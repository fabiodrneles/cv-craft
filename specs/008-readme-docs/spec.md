# 008 — README e documentação

- **Prioridade:** P1
- **Status:** Draft
- **Resolve:** C5 (parte de docs) e a seção 5 de [ANALYSIS.md](../ANALYSIS.md)

## Contexto

O README atual tem ~25 linhas, está truncado (bloco de código não fechado), tem instruções de instalação que não funcionam e promete recursos que não existem (otimização via `-ats`, "múltiplos formatos com um único comando"). É a porta de entrada do projeto e precisa refletir exatamente o que o binário faz.

## Requisitos

- **FR-1** O README MUST ser escrito **depois** das specs 001–004 estarem implementadas, ou marcar claramente o que é "planejado".
- **FR-2** Todo comando de exemplo no README MUST ser executável e coberto por teste (ou por um job de CI que roda os exemplos).
- **FR-3** Idioma: PT-BR como principal; `README.en.md` opcional com link no topo.
- **FR-4** Nenhum placeholder (`seu-usuario`) nem link quebrado (verificar com `lychee` no CI — MAY).

## Estrutura proposta

```
# CV-Craft
<tagline de 1 linha>
[badges: CI · Go Report Card · Go version · Release · License]
<imagem: prévia do PDF gerado, lado a lado com o YAML>

## Por que CV-Craft?            (3–4 bullets: fonte única, ATS, versionável no Git, offline)
## Instalação
   ### Binário pré-compilado    (link para Releases)
   ### Via Go                   (go install github.com/fabiodrneles/cv-craft@latest)
   ### A partir do código
## Início rápido               (init → editar → build, 3 comandos)
## Uso
   ### Comandos                 (tabela: build, validate, init, version, ui)
   ### Flags                    (tabela)
   ### Exit codes
## Formato do YAML             (exemplo mínimo + link para docs/schema.md e examples/)
## Formatos de saída           (PDF / Markdown / Texto — quando usar cada um)
## Boas práticas para ATS      (curta: verbos de ação, métricas, keywords da vaga)
## Roadmap                     (link para specs/ROADMAP.md)
## Contribuindo                (SDD: link para specs/README.md; como rodar testes)
## Licença
```

## Documentos de apoio

- **FR-5** `docs/schema.md` — referência campo a campo, gerada a partir da spec 001.
- **FR-6** `CONTRIBUTING.md` — fluxo SDD, `make test`, convenção de commits.
- **FR-7** `examples/README.md` — o que cada exemplo demonstra.

## Critérios de aceite

- **AC-1** Uma pessoa sem contexto consegue, só com o README, instalar e gerar um PDF em < 5 minutos.
- **AC-2** Todos os comandos do README rodam com exit 0 em CI.
- **AC-3** Nenhuma funcionalidade descrita no README deixa de existir no binário.
