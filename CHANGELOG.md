# Changelog

Todas as mudanças relevantes do CV-Craft são registradas aqui.

O formato segue o [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/), e o projeto usa [Versionamento Semântico](https://semver.org/lang/pt-BR/). A partir da `1.0.0`, mudanças incompatíveis na CLI ou no formato do YAML só acontecem numa nova versão major.

## [Unreleased]

## [1.0.0] - 2026-10-01

Fase 3, "profissional" ([#10](https://github.com/fabiodrneles/cv-craft/issues/10)). A partir desta versão, a CLI e o formato do YAML são estáveis.

### Adicionado

- `cv-craft build --watch` (`-w`): gera de novo a cada vez que o YAML é salvo, em menos de um segundo, até `Ctrl+C`. Erros de validação aparecem sem encerrar a observação ([#12](https://github.com/fabiodrneles/cv-craft/issues/12)).
- JSON Schema do YAML em [`schema/cv-craft.schema.json`](schema/cv-craft.schema.json): no VS Code e nas IDEs JetBrains, autocompletar, descrição dos campos e erros marcados enquanto se digita. Os exemplos e o modelo do `init` já apontam para ele, e `cv-craft schema` imprime uma cópia para uso offline ([#11](https://github.com/fabiodrneles/cv-craft/issues/11)).
- README em inglês ([`README.en.md`](README.en.md)), com prévia do PDF em inglês ([#13](https://github.com/fabiodrneles/cv-craft/issues/13)).

### Alterado

- Dependências atualizadas: `fatih/color` 1.19.0 e `go-isatty` 0.0.24.

### Removido

- **Incompatível:** a flag `-ats`, depreciada desde a 0.2.0. Agora ela é uma flag desconhecida (exit code `2`); o layout padrão já é otimizado para ATS ([#14](https://github.com/fabiodrneles/cv-craft/issues/14)).

## [0.2.0] - 2026-10-01

Fase 2, "confiável e publicável" ([#3](https://github.com/fabiodrneles/cv-craft/issues/3)). Primeira versão publicada e a primeira com binários pré-compilados; inclui também a Fase 1, descrita no fim desta seção.

### Adicionado

- Binários pré-compilados para Linux, macOS e Windows (amd64 e arm64), com `checksums.txt`, publicados automaticamente a cada tag.
- Referência completa do formato YAML em [`docs/schema.md`](docs/schema.md) e guia dos exemplos em [`examples/README.md`](examples/README.md).
- Guia de contribuição ([`CONTRIBUTING.md`](CONTRIBUTING.md)), formulários de issue e template de PR.

### Alterado

- Versão mínima do Go para compilar ou instalar com `go install`: **1.26**.

### Fase 1 — "funcionar de verdade" ([#1](https://github.com/fabiodrneles/cv-craft/issues/1))

Incluída na 0.2.0: a Fase 1 não recebeu tag própria.

#### Adicionado

- Comandos `build`, `validate`, `init`, `version`, `help` e `ui` (modo interativo).
- Saídas em PDF, Markdown e texto puro com o mesmo conteúdo; `--format all` gera as três de uma vez.
- Títulos das seções em português (padrão) ou inglês, via `meta.locale` ou `--lang`.
- PDF com metadados (título, autor, palavras-chave), links clicáveis e quebra de página sem títulos órfãos.
- Validação que lista todos os problemas de uma vez, com o caminho do campo e a linha; chaves desconhecidas no YAML são rejeitadas.
- Níveis de habilidade em inglês ou português (`advanced`, `avançado`…).
- Exit codes documentados e flags `--force`, `--quiet` e `--verbose`, para uso em scripts.
- Exemplos em `examples/` e modelo comentado gerado por `cv-craft init`.

#### Alterado

- **Incompatível:** a CLI antiga (`go run main.go -input ...`) foi substituída pelos subcomandos acima; `cv-craft` sem argumentos mostra a ajuda.
- Módulo renomeado para `github.com/fabiodrneles/cv-craft`, o que permite `go install`.

#### Corrigido

- Acentos corrompidos no PDF (`São Paulo` saía como `SÃ£o Paulo`).
- Markdown e texto saíam sem responsabilidades, conquistas e tecnologias das experiências.
- Habilidades com nível fora de `proficient/intermediate/beginner` sumiam do PDF.
- O modo interativo entrava em loop infinito quando a entrada terminava.
- `init` sobrescrevia um arquivo existente sem perguntar.

#### Depreciado

- `-ats`: aceita, mas sem efeito (o layout padrão já é otimizado para ATS). Será removida em uma versão futura.

[Unreleased]: https://github.com/fabiodrneles/cv-craft/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/fabiodrneles/cv-craft/compare/v0.2.0...v1.0.0
[0.2.0]: https://github.com/fabiodrneles/cv-craft/releases/tag/v0.2.0
