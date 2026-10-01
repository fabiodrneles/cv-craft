# Changelog

Todas as mudanças relevantes do CV-Craft são registradas aqui.

O formato segue o [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/), e o projeto usa [Versionamento Semântico](https://semver.org/lang/pt-BR/). Enquanto a versão for `0.x`, mudanças incompatíveis podem ocorrer em versões minor e são sempre sinalizadas aqui.

## [Unreleased]

Primeira versão pública, a ser publicada como **v0.1.0** ao fim da Fase 1 ([#1](https://github.com/fabiodrneles/cv-craft/issues/1)).

### Adicionado

- Comandos `build`, `validate`, `init`, `version`, `help` e `ui` (modo interativo).
- Saídas em PDF, Markdown e texto puro com o mesmo conteúdo; `--format all` gera as três de uma vez.
- Títulos das seções em português (padrão) ou inglês, via `meta.locale` ou `--lang`.
- PDF com metadados (título, autor, palavras-chave), links clicáveis e quebra de página sem títulos órfãos.
- Validação que lista todos os problemas de uma vez, com o caminho do campo e a linha; chaves desconhecidas no YAML são rejeitadas.
- Níveis de habilidade em inglês ou português (`advanced`, `avançado`…).
- Exit codes documentados e flags `--force`, `--quiet` e `--verbose`, para uso em scripts.
- Exemplos em `examples/` e modelo comentado gerado por `cv-craft init`.

### Alterado

- **Incompatível:** a CLI antiga (`go run main.go -input ...`) foi substituída pelos subcomandos acima; `cv-craft` sem argumentos mostra a ajuda.
- Módulo renomeado para `github.com/fabiodrneles/cv-craft`, o que permite `go install`.

### Corrigido

- Acentos corrompidos no PDF (`São Paulo` saía como `SÃ£o Paulo`).
- Markdown e texto saíam sem responsabilidades, conquistas e tecnologias das experiências.
- Habilidades com nível fora de `proficient/intermediate/beginner` sumiam do PDF.
- O modo interativo entrava em loop infinito quando a entrada terminava.
- `init` sobrescrevia um arquivo existente sem perguntar.

### Depreciado

- `-ats`: aceita, mas sem efeito (o layout padrão já é otimizado para ATS). Será removida em uma versão futura.

[Unreleased]: https://github.com/fabiodrneles/cv-craft/commits/main
