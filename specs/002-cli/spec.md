# 002 — Interface de linha de comando

- **Prioridade:** P0
- **Status:** Draft
- **Código afetado:** `main.go`, novo `internal/app` (serviço de build), `cli/`
- **Resolve:** A2, A3, A4, M2, M3, M5, M9

## Contexto

A CLI é o produto. Hoje o parsing de argumentos é manual, há um "modo legado" que quase nunca é ativado, `build` está duplicado no modo interativo, `--help` sai com erro, e não dá para usar em scripts porque o prompt de sobrescrita bloqueia.

## Comandos

```
cv-craft build <arquivo.yaml> [flags]     Gera o currículo
cv-craft validate <arquivo.yaml>           Valida sem gerar
cv-craft init [arquivo.yaml]               Cria um YAML modelo (padrão: curriculum.yaml)
cv-craft version                           Mostra versão, commit e data do build
cv-craft help [comando]                    Ajuda
cv-craft ui | interactive                  Modo interativo (spec 005)
cv-craft                                   Sem argumentos: ajuda (ver decisão D2)
```

### Flags de `build`

| Flag | Curta | Padrão | Descrição |
|---|---|---|---|
| `--format` | `-f` | `pdf` | `pdf`, `md` (`markdown`), `txt` (`text`), `all` |
| `--output` | `-o` | derivado | Arquivo de saída; com `all`, é um diretório |
| `--force` | `-y` | `false` | Sobrescreve sem perguntar |
| `--lang` | | do YAML | Idioma das seções (spec 006) |
| `--quiet` | `-q` | `false` | Só erros |
| `--verbose` | `-v` | `false` | Logs de diagnóstico em stderr |

Flags de forma única com um hífen (`-format`) SHOULD continuar aceitas por compatibilidade.

## Requisitos funcionais

- **FR-1** Parsing de argumentos MUST usar `flag` da stdlib (um `FlagSet` por subcomando) ou biblioteca justificada; flags podem vir antes ou depois do arquivo.
- **FR-2** Existe **um único** serviço `app.Build(ctx, BuildOptions) (BuildResult, error)` usado pela CLI e pelo modo interativo. `main.go` só faz parsing e apresentação.
- **FR-3** O "modo legado" baseado em `os.Args[0]` MUST ser removido.
- **FR-4** Saída padrão: mesmo diretório e nome base do YAML, extensão do formato (`.pdf`, `.md`, `.txt`).
- **FR-5** Se o arquivo de saída existir: com `--force` sobrescreve; sem `--force` e stdin for TTY, pergunta (`s/N`); sem `--force` e stdin não-TTY, falha com exit 3 e mensagem sugerindo `--force`.
- **FR-6** A escrita MUST ser atômica (arquivo temporário + `os.Rename`) para nunca deixar uma saída corrompida.
- **FR-7** `--format all` gera os três formatos numa execução.
- **FR-8** `validate` imprime todos os erros e avisos (spec 001) e um resumo (nº de experiências, skills etc.).
- **FR-9** `init` MUST recusar sobrescrever arquivo existente sem `--force`. O modelo vem de `examples/` via `embed`.
- **FR-10** `version`/`--version` imprimem `cv-craft vX.Y.Z (commit, data)`, injetados por `-ldflags`; padrão `dev`.
- **FR-11** `help`, `-h`, `--help` saem com código 0. Uso incorreto sai com 2.
- **FR-12** stdout: só o resultado (caminho gerado). stderr: logs, avisos, erros. Geradores não chamam `log` diretamente.
- **FR-13** Emojis/cores MUST ser desativados quando stdout não for TTY ou `NO_COLOR` estiver definida.
- **FR-14** Todas as mensagens da CLI em um único idioma (ver spec 006).

## Exit codes

| Código | Significado |
|---|---|
| 0 | Sucesso |
| 1 | Erro inesperado / falha de I/O |
| 2 | Uso incorreto (comando ou flag inválidos) |
| 3 | Saída existente sem `--force` em modo não interativo, ou cancelado pelo usuário |
| 4 | Erro de validação do YAML |

## Critérios de aceite

- **AC-1** `cv-craft build examples/full.yaml -o out.pdf` cria `out.pdf` e sai com 0.
- **AC-2** `cv-craft build x.yaml --format all -o dist/` cria `dist/x.pdf`, `dist/x.md`, `dist/x.txt`.
- **AC-3** Com `out.pdf` existente e stdin `/dev/null`, sem `--force` ⇒ exit 3, arquivo intacto.
- **AC-4** Mesmo cenário com `--force` ⇒ exit 0, arquivo regravado.
- **AC-5** `cv-craft --help` ⇒ exit 0; `cv-craft foo` ⇒ exit 2.
- **AC-6** `cv-craft build -f pdf x.yaml` e `cv-craft build x.yaml -f pdf` são equivalentes.
- **AC-7** `cv-craft init` em diretório com `curriculum.yaml` existente ⇒ exit 3, arquivo intacto.
- **AC-8** `cv-craft validate invalido.yaml` ⇒ exit 4, todos os erros listados em stderr.
- **AC-9** `cv-craft build x.yaml | cat` não contém códigos ANSI.
- **AC-10** `cv-craft version` com binário de release mostra a tag.

## Fora de escopo

- Watch mode (`--watch`) — candidato a spec futura.
- Configuração global (`~/.cv-craft`).

## Decisões em aberto

- D2 — comportamento de `cv-craft` sem argumentos (hoje abre o modo interativo; recomendado: mostrar ajuda).
