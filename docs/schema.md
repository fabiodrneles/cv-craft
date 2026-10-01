# Referência do YAML

Esta página descreve **todos** os campos aceitos no arquivo do currículo, como eles aparecem nos formatos de saída e como resolver os erros mais comuns. Para um ponto de partida, rode `cv-craft init` ou veja os [exemplos](../examples/README.md).

> A especificação formal (requisitos e critérios de aceite) está na [spec 001](../specs/001-yaml-schema/spec.md). Esta página é a versão para quem escreve o currículo.

## Regras gerais

- O arquivo é YAML em **UTF-8**; acentos e caracteres especiais funcionam em todos os formatos.
- Campos marcados como **obrigatórios** não podem ficar vazios nem conter só espaços.
- **Chaves desconhecidas são rejeitadas.** Um erro de digitação como `responsabilities` gera um erro com o caminho e a linha, em vez de o conteúdo sumir do currículo sem aviso.
- Campos opcionais vazios são simplesmente omitidos das saídas, junto com os separadores.
- Use aspas quando o texto tiver `:` ou começar com caracteres especiais do YAML (`-`, `*`, `&`, `#`).
- Listas podem ser escritas em bloco (um item por linha com `-`) ou em linha (`["Go", "Docker"]`).

## `meta`

Configurações do documento. A seção inteira é opcional.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `meta.locale` | texto | não | Idioma dos títulos das seções: `pt-BR` (padrão) ou `en`. Variações como `pt`, `pt_BR` e `en-US` também são aceitas. A flag `--lang` tem prioridade sobre este campo. O conteúdo escrito por você nunca é traduzido. |

## `contact`

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `contact.name` | texto | **sim** | Nome completo; é o título do currículo e o autor nos metadados do PDF. |
| `contact.email` | texto | **sim** | Um endereço de e-mail simples (`ana@exemplo.com`). Formatos como `Ana <ana@exemplo.com>` são rejeitados. No PDF e no Markdown vira link `mailto:`. |
| `contact.phone` | texto | não | Exibido como escrito. |
| `contact.location` | texto | não | Cidade e estado/país. |
| `contact.linkedin` | texto | não | URL do perfil. |
| `contact.github` | texto | não | URL do perfil. |
| `contact.portfolio` | texto | não | URL do site ou portfólio. |

**Links.** Em `linkedin`, `github`, `portfolio` e `certificates[].url`, o `https://` é opcional: `linkedin.com/in/ana` vira `https://linkedin.com/in/ana`. No PDF, o link é exibido sem o `https://`, mas continua clicável. No texto e no Markdown, aparece completo.

## Título e resumo

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `professional_title` | texto | **sim** | Cargo ou especialidade, logo abaixo do nome. Alinhe com o cargo da vaga. |
| `summary` | texto | não | Resumo profissional de 2 a 4 frases. Para textos longos, use `>-` e quebre as linhas à vontade: o YAML junta tudo num parágrafo. |

## `skills`

**Obrigatório:** ao menos uma categoria. Cada categoria vira uma linha `Categoria (Nível): palavra-chave, palavra-chave…`, na ordem do arquivo.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `skills[].name` | texto | **sim** | Nome da categoria (ex.: `Backend`). |
| `skills[].level` | texto | não | Nível da categoria, exibido traduzido entre parênteses. Veja a tabela abaixo. |
| `skills[].keywords` | lista de textos | **sim** (≥ 1) | As habilidades. Use os mesmos termos das vagas: o PDF também os grava nos metadados como palavras-chave. |

Níveis aceitos (maiúsculas e minúsculas tanto faz):

| Valor canônico | Também aceito | Exibido em `pt-BR` | Exibido em `en` |
|---|---|---|---|
| `expert` | `especialista` | Especialista | Expert |
| `advanced` | `avançado`, `avancado` | Avançado | Advanced |
| `proficient` | `proficiente` | Proficiente | Proficient |
| `intermediate` | `intermediário`, `intermediario` | Intermediário | Intermediate |
| `beginner` | `básico`, `basico`, `iniciante` | Básico | Beginner |

Um nível fora da lista gera um **aviso**, e não um erro: a categoria e as habilidades aparecem normalmente, só sem o nível.

## `experience`

**Obrigatório:** ao menos uma experiência. Elas aparecem na ordem do arquivo; coloque a mais recente primeiro.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `experience[].role` | texto | **sim** | Cargo. |
| `experience[].company` | texto | **sim** | Empresa. Exibida como `Cargo — Empresa`. |
| `experience[].location` | texto | não | Local. |
| `experience[].work_type` | texto | não | Modalidade (`Remoto`, `Híbrido`, `Presencial`…), exibida como escrita. Omitida se for igual ao local. |
| `experience[].period` | texto | não | Período em texto livre (`Mar 2022 - Atual`). |
| `experience[].description` | texto | não | Uma frase de contexto sobre a empresa ou o time. |
| `experience[].responsibilities` | lista de textos | não | O que você fazia. Comece cada item com um verbo de ação. |
| `experience[].achievements` | lista de textos | não | Resultados, de preferência com números. Aparecem em uma lista própria. |
| `experience[].technologies` | lista de textos | não | Tecnologias usadas, exibidas numa linha. |

## `education`

**Obrigatório:** ao menos uma formação.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `education[].degree` | texto | **sim** | Curso. |
| `education[].institution` | texto | **sim** | Instituição. Exibida como `Curso — Instituição`. |
| `education[].location` | texto | não | Local. |
| `education[].period` | texto | não | Período em texto livre. |
| `education[].thesis` | texto | não | Título do TCC, tese ou dissertação. |
| `education[].relevant_courses` | texto | não | Disciplinas relevantes, separadas por vírgula. |

## `certificates`

Opcional.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `certificates[].name` | texto | **sim** | Nome do certificado. |
| `certificates[].institution` | texto | não | Emissor. |
| `certificates[].date` | texto | não | Data ou ano. |
| `certificates[].url` | texto | não | Link de verificação (Credly, por exemplo). |

## `languages`

Opcional.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `languages[].language` | texto | **sim** | Idioma. |
| `languages[].level` | texto | não | Nível, em texto livre (`Nativo`, `Fluente`, `B2`…). |

## Exemplo completo

```yaml
meta:
  locale: pt-BR

contact:
  name: "Ana Souza"
  email: "ana@exemplo.com"
  phone: "+55 (11) 99999-9999"
  location: "São Paulo, SP"
  linkedin: "linkedin.com/in/ana-souza"
  github: "github.com/ana-souza"
  portfolio: "anasouza.dev"

professional_title: "Engenheira de Software | Go"

summary: >-
  Engenheira de software com 5 anos de experiência em sistemas
  distribuídos e APIs de alto volume.

skills:
  - name: "Backend"
    level: "avançado"
    keywords: ["Go", "PostgreSQL", "gRPC"]
  - name: "Infraestrutura"
    keywords: ["Docker", "Kubernetes"]

experience:
  - role: "Engenheira de Software"
    company: "ACME"
    location: "São Paulo, SP"
    work_type: "Híbrido"
    period: "Jan 2022 - Atual"
    description: "Plataforma de pagamentos com 1 milhão de transações por dia."
    responsibilities:
      - "Desenvolvi APIs gRPC em Go"
    achievements:
      - "Reduzi a latência p99 em 35%"
    technologies: ["Go", "PostgreSQL", "Kubernetes"]

education:
  - degree: "Ciência da Computação"
    institution: "USP"
    location: "São Paulo, SP"
    period: "2016 - 2020"
    thesis: "Consistência eventual em sistemas de pagamento"
    relevant_courses: "Sistemas Distribuídos, Banco de Dados"

certificates:
  - name: "Certified Kubernetes Application Developer (CKAD)"
    institution: "CNCF"
    date: "2024"
    url: "credly.com/badges/exemplo"

languages:
  - language: "Português"
    level: "Nativo"
  - language: "Inglês"
    level: "Fluente"
```

## Erros e avisos comuns

Rode `cv-craft validate curriculum.yaml` para ver **todos os problemas de validação** de uma vez (campos obrigatórios, chaves desconhecidas, e-mail, idioma, níveis). Cada mensagem traz o caminho do campo e, quando possível, a linha.

A exceção são os erros de **sintaxe ou de tipo** do YAML (as duas últimas linhas da tabela abaixo): enquanto o arquivo não puder ser lido, a validação não roda, e só esse erro aparece. Corrija-o e rode `validate` de novo para ver o resto.

| Mensagem | Causa | Como corrigir |
|---|---|---|
| `experience[0].responsabilities: campo desconhecido "responsabilities" (linha 12)` | Chave com erro de digitação ou que não existe no schema. | Corrija o nome conforme as tabelas acima. |
| `contact.email: campo obrigatório` | Campo obrigatório ausente, vazio ou só com espaços. | Preencha o campo. |
| `contact.email: e-mail inválido "..."` | E-mail sem `@`, sem domínio ou com nome junto. | Use só o endereço: `ana@exemplo.com`. |
| `skills: informe pelo menos uma categoria de habilidades` | Lista vazia ou ausente (o mesmo vale para `experience` e `education`). | Adicione ao menos um item. |
| `skills[0].keywords: informe pelo menos uma habilidade` | Categoria sem palavras-chave. | Adicione ao menos uma em `keywords`. |
| `meta.locale: idioma não suportado "xx" (suportados: pt-BR, en)` | Idioma desconhecido. | Use `pt-BR` ou `en`. |
| `aviso: skills[0].level: nível desconhecido "guru" será omitido` | Nível fora da lista. É um aviso: o currículo é gerado mesmo assim. | Use um dos níveis da tabela ou remova o campo. |
| `YAML inválido: yaml: line 3: ...` | Erro de sintaxe: indentação, aspas não fechadas, `:` sem aspas. | Confira a linha indicada; textos com `:` precisam de aspas. |
| `YAML inválido: ... cannot unmarshal !!str ... into []string` | Texto onde se espera lista (ex.: `keywords: Go`). | Use lista: `keywords: ["Go"]`. |
| `arquivo vazio` | Arquivo sem conteúdo ou só com comentários. | Comece com `cv-craft init`. |
