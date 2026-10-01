package cli

const usage = `CV-Craft — currículos profissionais e otimizados para ATS a partir de um YAML.

Uso:
  cv-craft <comando> [argumentos] [flags]

Comandos:
  build <arquivo.yaml>    Gera o currículo (PDF, Markdown ou texto)
  validate <arquivo.yaml> Valida o YAML sem gerar arquivos
  init [arquivo.yaml]     Cria um YAML modelo (padrão: curriculum.yaml)
  schema                  Imprime o JSON Schema do YAML (autocompletar no editor)
  version                 Mostra a versão
  help [comando]          Mostra a ajuda de um comando
  ui                      Abre o modo interativo

Exemplos:
  cv-craft init
  cv-craft build curriculum.yaml
  cv-craft build curriculum.yaml --format all -o dist/
  cv-craft build curriculum.yaml -f md --lang en
  cv-craft build curriculum.yaml --watch

Exit codes: 0 sucesso · 1 erro inesperado · 2 uso incorreto ·
            3 saída já existe/cancelado · 4 YAML inválido
`

const usageBuild = `Uso: cv-craft build <arquivo.yaml> [flags]

Gera o currículo a partir do YAML. As flags podem vir antes ou depois do arquivo.

Flags:
  -f, --format string   pdf, md (markdown), txt (text) ou all (padrão "pdf")
  -o, --output string   arquivo de saída; com --format all, um diretório
                        (padrão: ao lado do YAML, com a extensão do formato)
      --lang string     idioma dos títulos: pt-BR ou en (padrão: meta.locale do YAML, senão pt-BR)
  -y, --force           sobrescreve arquivos existentes sem perguntar
  -q, --quiet           mostra apenas erros
  -v, --verbose         mostra detalhes da execução
  -w, --watch           gera de novo a cada vez que o YAML é salvo, até Ctrl+C
`

const usageValidate = `Uso: cv-craft validate <arquivo.yaml>

Valida o YAML e lista todos os erros e avisos de uma vez, sem gerar arquivos.
`

const usageInit = `Uso: cv-craft init [arquivo.yaml] [flags]

Cria um YAML modelo comentado (padrão: curriculum.yaml).

Flags:
  -y, --force   sobrescreve o arquivo se ele já existir
`

const usageSchema = `Uso: cv-craft schema

Imprime o JSON Schema do YAML do currículo. Editores com a extensão YAML
(VS Code, JetBrains) usam o schema para autocompletar e marcar erros enquanto
você digita. Os YAMLs criados pelo init já apontam para o schema publicado;
para usar uma cópia local, offline:

  cv-craft schema > cv-craft.schema.json

e troque a primeira linha do YAML por:

  # yaml-language-server: $schema=./cv-craft.schema.json
`

const usageUI = `Uso: cv-craft ui

Abre o modo interativo, que aceita os mesmos comandos da linha de comando
(build, validate, init, version, help). Requer um terminal.
`

var commandUsage = map[string]string{
	"build":       usageBuild,
	"validate":    usageValidate,
	"init":        usageInit,
	"version":     "Uso: cv-craft version\n\nMostra a versão, o commit e a data do build.\n",
	"help":        usage,
	"schema":      usageSchema,
	"ui":          usageUI,
	"interactive": usageUI,
}
