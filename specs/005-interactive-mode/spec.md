# 005 — Modo interativo

- **Prioridade:** P1
- **Status:** Done
- **Código afetado:** `internal/cli/shell.go` (antes `cli/ui.go`)
- **Resolve:** C4, A2, M4

## Contexto

O modo interativo é um REPL que aceita `build`, `init`, `validate`, `templates`, `help`, `exit`. Tem bugs graves (loop infinito em EOF), duplica a lógica de build e anuncia funcionalidades que não existem.

## Requisitos funcionais

- **FR-1** EOF (Ctrl+D) ou erro de leitura MUST encerrar o REPL com exit 0. Ctrl+C encerra com 130.
- **FR-2** Entrada vazia MUST apenas reexibir o prompt (não imprimir a ajuda).
- **FR-3** Cada comando MUST delegar para o mesmo código da CLI (`app.Build`, `app.Validate`, `app.Init` — spec 002). Nenhuma lógica de negócio em `cli/`.
- **FR-4** Mesmos flags, mesmas mensagens, mesmo comportamento de sobrescrita da CLI (com TTY, pergunta).
- **FR-5** Argumentos com aspas (`build "meu cv.yaml"`) MUST funcionar (usar tokenizer com suporte a aspas, não `strings.Fields`).
- **FR-6** Remover: promessa de TAB/histórico na ajuda (a não ser que seja implementado), `templates list` fictício, `sleep` artificial do `ShowLoading`.
- **FR-7** Renomear `GeminiUI` para um nome do domínio (ex.: `cli.Shell`).
- **FR-8** Se stdin não for TTY, `cv-craft ui` MUST sair com mensagem e exit 2.

## Critérios de aceite

- **AC-1** `cv-craft ui < /dev/null` termina em < 1 s com exit 2 (não-TTY).
- **AC-2** Em PTY, enviar `build examples/full.yaml -o x.pdf\nexit\n` gera `x.pdf` e sai com 0.
- **AC-3** Em PTY, Ctrl+D no prompt encerra com 0.
- **AC-4** `help` lista apenas comandos que existem.
- **AC-5** Teste de unidade do REPL com `io.Reader`/`io.Writer` injetados cobre: vazio, comando desconhecido, EOF, `exit`.

## Decisões

- D2 — o modo interativo é mantido como camada fina sobre o mesmo `dispatch` da CLI e não é mais o padrão sem argumentos.
- O tokenizer respeita aspas simples e duplas, mas não trata `\` como escape, para aceitar caminhos do Windows.
- AC-2 e AC-3 são verificados com stdin simulado (`IsTerminal` injetável) em `internal/cli/cli_test.go`, sem PTY real.
