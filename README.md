# CV-Craft

**Português** · [English](README.en.md)

**Currículos profissionais e otimizados para ATS, gerados a partir de um único arquivo YAML.**

[![CI](https://github.com/fabiodrneles/cv-craft/actions/workflows/ci.yml/badge.svg)](https://github.com/fabiodrneles/cv-craft/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/fabiodrneles/cv-craft.svg)](https://pkg.go.dev/github.com/fabiodrneles/cv-craft)
[![Release](https://img.shields.io/github/v/release/fabiodrneles/cv-craft)](https://github.com/fabiodrneles/cv-craft/releases/latest)
[![Licença: MIT](https://img.shields.io/badge/licen%C3%A7a-MIT-blue.svg)](LICENSE)

Você escreve o conteúdo uma vez em YAML; o CV-Craft gera o **PDF** para enviar a recrutadores, o **texto puro** para colar em formulários de candidatura e o **Markdown** para o GitHub ou portfólio, todos com o mesmo conteúdo.

<p align="center">
  <img src="docs/preview.png" alt="Primeira página do PDF gerado a partir de examples/full.yaml" width="520">
</p>

## Por que CV-Craft?

- **Fonte única da verdade.** Um arquivo YAML, três formatos sempre consistentes: nenhum campo preenchido fica de fora.
- **Feito para ATS.** Coluna única, texto selecionável, sem tabelas nem imagens, títulos de seção padronizados e metadados no PDF (título, autor, palavras-chave).
- **Português de verdade.** Acentos e caracteres especiais corretos em todos os formatos; títulos das seções em `pt-BR` ou `en`.
- **Versionável e offline.** Seu currículo vira texto no Git. Nada é enviado para a internet.
- **Feito para scripts.** Funciona sem terminal interativo, com flags para tudo e exit codes documentados.

## Instalação

Escolha o seu sistema e siga os passos. Cada bloco de comandos pode ser copiado e colado inteiro no terminal.

### Windows

1. Abra o **PowerShell**: aperte a tecla **Windows**, digite `PowerShell` e aperte **Enter**. Não precisa ser como administrador.
2. Copie e cole este comando e aperte **Enter**:

   ```powershell
   irm https://raw.githubusercontent.com/fabiodrneles/cv-craft/main/scripts/install.ps1 | iex
   ```

   Ele baixa a última versão, confere a integridade do arquivo (checksum), instala em `%LOCALAPPDATA%\Programs\cv-craft` e acrescenta essa pasta ao seu `PATH`, para que o comando `cv-craft` funcione em qualquer pasta.
3. Confira a instalação:

   ```powershell
   cv-craft version
   ```

   Deve aparecer algo como `cv-craft v1.0.0 (...)`. Nos terminais que já estavam abertos antes da instalação, feche e abra de novo para o `PATH` novo valer.

<details>
<summary>Instalação manual no Windows (sem script)</summary>

1. Em [Releases](https://github.com/fabiodrneles/cv-craft/releases/latest), baixe `cv-craft_<versão>_windows_amd64.zip` (ou `arm64`, em computadores com processador ARM).
2. Clique com o botão direito no arquivo → **Extrair tudo** → escolha a pasta `C:\Users\<seu usuário>\AppData\Local\Programs\cv-craft`.
3. Adicione essa pasta ao `PATH`: aperte **Windows**, digite **variáveis de ambiente**, abra **Editar as variáveis de ambiente para sua conta**, selecione **Path** → **Editar** → **Novo**, cole o caminho da pasta e confirme com **OK** nas janelas.
4. Abra um PowerShell **novo** e rode `cv-craft version`.

</details>

### macOS e Linux

1. Abra o **Terminal** (no macOS: **Cmd + Espaço**, digite `Terminal` e aperte **Enter**).
2. Copie e cole este comando:

   ```sh
   curl -fsSL https://raw.githubusercontent.com/fabiodrneles/cv-craft/main/scripts/install.sh | sh
   ```

   Ele detecta o sistema e o processador (Intel/AMD ou ARM, como os Macs M1/M2/M3), baixa a última versão, confere o checksum e instala em `~/.local/bin`.
3. Se o instalador avisar que a pasta não está no `PATH`, copie e rode o comando que ele mostra (ele acrescenta `~/.local/bin` ao `PATH` do seu shell).
4. Confira a instalação:

   ```sh
   cv-craft version
   ```

> Baixou o arquivo pelo navegador no macOS e apareceu "não é possível verificar o desenvolvedor"? Rode `xattr -d com.apple.quarantine ~/.local/bin/cv-craft` (troque pelo caminho onde você colocou o arquivo). O instalador acima não tem esse problema.

### Com Go

Se você já tem o Go 1.26 ou mais novo:

```bash
go install github.com/fabiodrneles/cv-craft@latest
```

### A partir do código

```bash
git clone https://github.com/fabiodrneles/cv-craft.git
cd cv-craft
make build        # gera ./bin/cv-craft (ou: go build .)
```

### Atualizar e desinstalar

- **Atualizar:** rode o mesmo comando de instalação de novo; ele substitui a versão instalada pela mais recente.
- **Desinstalar no Windows:** apague a pasta `%LOCALAPPDATA%\Programs\cv-craft` e remova-a do `Path` (mesma tela do passo 3 da instalação manual).
- **Desinstalar no macOS e no Linux:** `rm ~/.local/bin/cv-craft`.

## Início rápido

Do zero ao primeiro PDF:

1. Crie uma pasta para o seu currículo e entre nela:

   ```sh
   mkdir meu-curriculo
   cd meu-curriculo
   ```

2. Crie o modelo comentado:

   ```bash
   cv-craft init                      # cria curriculum.yaml com um modelo comentado
   ```

3. Abra o `curriculum.yaml` num editor e troque os exemplos pelos seus dados. Recomendação: o [VS Code](https://code.visualstudio.com/) com a extensão [YAML da Red Hat](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml), que autocompleta os campos e aponta erros enquanto você digita (veja [Autocompletar e validar no editor](#autocompletar-e-validar-no-editor)). No YAML, a **indentação com espaços** importa: mantenha o alinhamento do modelo.
4. Gere o PDF:

   ```bash
   cv-craft build curriculum.yaml     # gera curriculum.pdf
   ```

5. Abra o resultado: `start curriculum.pdf` (Windows), `open curriculum.pdf` (macOS) ou `xdg-open curriculum.pdf` (Linux).

Para gerar os três formatos de uma vez:

```bash
cv-craft build curriculum.yaml --format all -o dist/
```

Para ver o PDF se atualizar enquanto você edita, use o [`--watch`](#flags-do-build).

## Uso

### Comandos

| Comando | Descrição |
|---|---|
| `cv-craft build <arquivo.yaml>` | Gera o currículo (PDF, Markdown ou texto) |
| `cv-craft validate <arquivo.yaml>` | Valida o YAML e lista **todos** os problemas de validação de uma vez (um erro de sintaxe do YAML aparece sozinho, antes) |
| `cv-craft init [arquivo.yaml]` | Cria um YAML modelo (padrão: `curriculum.yaml`); não sobrescreve sem `--force` |
| `cv-craft schema` | Imprime o JSON Schema do YAML, para autocompletar no editor |
| `cv-craft version` | Mostra a versão |
| `cv-craft help [comando]` | Ajuda geral ou de um comando |
| `cv-craft ui` | Modo interativo, com os mesmos comandos |

### Flags do `build`

As flags podem vir antes ou depois do arquivo.

| Flag | Padrão | Descrição |
|---|---|---|
| `-f`, `--format` | `pdf` | `pdf`, `md` (`markdown`), `txt` (`text`) ou `all` |
| `-o`, `--output` | ao lado do YAML | Arquivo de saída; com `--format all`, um diretório |
| `--lang` | `meta.locale` ou `pt-BR` | Idioma dos títulos das seções: `pt-BR` ou `en` |
| `-y`, `--force` | — | Sobrescreve arquivos existentes sem perguntar |
| `-q`, `--quiet` | — | Mostra apenas erros |
| `-v`, `--verbose` | — | Mostra detalhes da execução |
| `-w`, `--watch` | — | Gera de novo a cada vez que o YAML é salvo, até `Ctrl+C` |

Se o arquivo de saída já existir, o CV-Craft pergunta antes de sobrescrever quando está em um terminal. Em scripts e no CI, ele falha com exit code `3`, a não ser que você use `--force`.

Para editar o currículo vendo o resultado, deixe o `--watch` rodando num terminal e o PDF aberto num visualizador que recarrega o arquivo:

```text
cv-craft build curriculum.yaml --watch
```

Cada vez que você salva o YAML, o currículo é gerado de novo em menos de um segundo. Se o YAML ficar inválido, os erros aparecem e a observação continua; ao corrigir e salvar, a geração volta a funcionar. `Ctrl+C` encerra com exit code `0`.

### Exit codes

| Código | Significado |
|---|---|
| `0` | Sucesso |
| `1` | Erro inesperado (ex.: arquivo não encontrado, falha de escrita) |
| `2` | Uso incorreto (comando, flag ou valor inválido) |
| `3` | Arquivo de saída já existe sem `--force`, ou operação cancelada |
| `4` | YAML inválido (sintaxe ou validação) |

## Formato do YAML

Exemplo mínimo válido:

```yaml
meta:
  locale: pt-BR            # pt-BR ou en (títulos das seções)

contact:
  name: "Ana Souza"        # obrigatório
  email: "ana@exemplo.com" # obrigatório
  phone: "+55 (11) 99999-9999"
  location: "São Paulo, SP"
  linkedin: "linkedin.com/in/ana-souza"
  github: "github.com/ana-souza"

professional_title: "Desenvolvedora Backend"   # obrigatório
summary: "Duas a quatro frases sobre você."

skills:                    # obrigatório: ao menos uma categoria
  - name: "Backend"
    level: "advanced"      # opcional: expert, advanced, proficient, intermediate, beginner
    keywords: ["Go", "PostgreSQL"]

experience:                # obrigatório: ao menos uma
  - role: "Desenvolvedora"
    company: "ACME"
    period: "2022 - Atual"
    responsibilities: ["Desenvolvi APIs REST em Go"]
    achievements: ["Reduzi a latência em 40%"]
    technologies: ["Go", "Docker"]

education:                 # obrigatório: ao menos uma
  - degree: "Ciência da Computação"
    institution: "USP"
```

Também são suportados `work_type` e `description` nas experiências; `location`, `period`, `thesis` e `relevant_courses` na formação; e as listas opcionais `certificates` (`name`, `institution`, `date`, `url`) e `languages` (`language`, `level`). A referência completa, com todos os campos, os níveis aceitos e os erros mais comuns, está em [`docs/schema.md`](docs/schema.md).

Os níveis de habilidade também podem ser escritos em português (`avançado`, `intermediário`, `básico`…). Chaves com erro de digitação são **rejeitadas**, com o caminho e a linha (`experience[0].responsabilities: campo desconhecido "responsabilities" (linha 12)`), para que nada suma do currículo sem aviso.

### Autocompletar e validar no editor

O YAML criado pelo `cv-craft init` (e os de `examples/`) começa com uma linha que aponta para o [JSON Schema](schema/cv-craft.schema.json) do CV-Craft. Com a extensão YAML no editor ([Red Hat YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) no VS Code, ou o suporte nativo das IDEs JetBrains), você ganha autocompletar dos campos, a descrição de cada um ao passar o mouse e os erros marcados enquanto digita: campo obrigatório vazio, chave com erro de digitação, e-mail inválido. Para trabalhar offline, gere uma cópia local com `cv-craft schema > cv-craft.schema.json` (detalhes em [`docs/schema.md`](docs/schema.md#autocompletar-e-validar-no-editor)).

Veja exemplos completos em [`examples/`](examples/README.md): [`full.yaml`](examples/full.yaml) (português, todos os campos), [`en.yaml`](examples/en.yaml) (inglês) e [`minimal.yaml`](examples/minimal.yaml) (o modelo do `init`).

## Problemas comuns

| Sintoma | Causa e solução |
|---|---|
| `cv-craft: command not found` ou "não é reconhecido como nome de cmdlet" | A pasta da instalação não está no `PATH`. No Windows, feche e abra o terminal; se continuar, refaça o passo 3 da instalação manual. No macOS/Linux, rode o comando de `PATH` que o instalador mostrou e abra um terminal novo. |
| `erro: ... já existe` e exit code `3` | O arquivo de saída já existe e o terminal não é interativo. Use `--force` para sobrescrever ou `-o` para outro nome. |
| `YAML inválido: ... line N` | Erro de sintaxe, quase sempre de indentação (use espaços, nunca Tab) ou de um texto com `:` sem aspas. Coloque o texto entre aspas. |
| `campo desconhecido "..."` | Nome de campo com erro de digitação. A mensagem mostra o caminho e a linha; confira o nome em [`docs/schema.md`](docs/schema.md). |
| `campo obrigatório` | Um campo obrigatório ficou vazio. Rode `cv-craft validate curriculum.yaml` para ver todos de uma vez. |
| O PowerShell bloqueia o script de instalação | Use exatamente o comando `irm ... \| iex` acima, que não depende da política de execução de scripts, ou siga a instalação manual. |

## Formatos de saída

| Formato | Quando usar |
|---|---|
| **PDF** | Enviar a recrutadores e anexar em candidaturas. Links clicáveis e metadados preenchidos. |
| **Texto** (`.txt`) | Colar em formulários de candidatura; máxima compatibilidade com ATS. |
| **Markdown** (`.md`) | README do GitHub, portfólio ou site pessoal. |

## Dicas para passar pelo ATS

- Comece cada responsabilidade com um **verbo de ação** e use **números** nas conquistas.
- Use nas `keywords` os mesmos termos da descrição da vaga (ex.: "Kubernetes", não só "K8s").
- Mantenha o título profissional alinhado ao cargo pretendido.

## Desenvolvimento

O projeto segue **Spec Driven Development**: toda mudança de comportamento começa por uma spec em [`specs/`](specs/README.md), e cada critério de aceite vira um teste. O fluxo completo (tickets, branches, commits, PRs e revisão) está no [guia de contribuição](CONTRIBUTING.md).

```bash
make test     # testes unitários, golden files e CLI
make lint     # go vet + golangci-lint
make cover    # cobertura (mínimo de 80% em internal/)
make smoke    # smoke test de ponta a ponta com o binário real
make docs     # lint de Markdown, links e comandos da documentação
make ci       # as verificações de código do CI
make golden   # regrava os golden files após uma mudança intencional nas saídas
make schema   # regrava schema/cv-craft.schema.json a partir do modelo Go
make release-snapshot  # gera os binários da release em ./dist, sem publicar
```

Para publicar uma versão, basta criar e enviar a tag (`git tag vX.Y.Z && git push origin vX.Y.Z`): o workflow de release roda o CI completo e publica os binários. As mudanças de cada versão ficam no [CHANGELOG](CHANGELOG.md).

Os testes que conferem o texto do PDF usam o `pdftotext` (pacote `poppler-utils` no Linux, `poppler` no Homebrew) e são pulados se ele não estiver instalado.

O roadmap está em [`specs/ROADMAP.md`](specs/ROADMAP.md).

## Licença

[MIT](LICENSE) © Fabio Dorneles. A fonte Liberation Sans, embutida no binário, é distribuída sob a [SIL Open Font License 1.1](internal/generator/fonts/OFL.txt).
