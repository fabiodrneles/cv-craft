# 009 — Processo de entrega

- **Prioridade:** P1
- **Status:** Done — em uso desde a Fase 2 (CONTRIBUTING, templates e skill `sdd-delivery`)
- **Código afetado:** processo (`CONTRIBUTING.md`, `.github/`, `specs/`)
- **Resolve:** #19

## Contexto

O CV-Craft saiu de um protótipo sem testes para um projeto com specs, CI em três sistemas operacionais e release automatizada seguindo um fluxo que, até aqui, existia só na prática: análise → specs → tickets → branches → PRs validados → revisão → fechamento de fase → tag. Esse fluxo foi conduzido por um agente (Claude Code) com o dono do repositório decidindo os pontos em aberto e fazendo os merges.

Esta spec torna o fluxo **normativo e reaproveitável**. As seções 1–13 são genéricas e valem para qualquer repositório; a seção [Aplicação no cv-craft](#aplicação-no-cv-craft) fixa os valores deste projeto (labels, comandos, fases, números de issue).

Relação com outros documentos:

- Esta spec é a **fonte normativa** do processo. O [`CONTRIBUTING.md`](../../CONTRIBUTING.md) é o **guia prático** para quem contribui; os dois MUST concordar. Se divergirem, corrige-se o que estiver errado no mesmo PR que notar a divergência.
- A [constituição](../constitution.md) continua acima de tudo; em particular o princípio 9 (CI bloqueia merge vermelho).
- O fluxo SDD de [`specs/README.md`](../README.md) (especificar → decidir → testar → implementar → fechar) é o núcleo técnico; esta spec cobre o que vem em volta dele.

## Termos

| Termo | Significado |
|---|---|
| **Dono** | Quem mantém o repositório e tem a palavra final (decisões, merges, tags, configurações). |
| **Agente** | Quem executa o trabalho delegado: análise, specs, tickets, código, PRs. Pode ser uma pessoa ou um assistente de IA. |
| **Fase** | Conjunto de tarefas que leva a uma versão (ex.: Fase 1 → `v0.1.0`). |
| **Épico** | Issue que representa uma fase e agrega as tarefas como sub-issues. |
| **Ticket** | Issue de uma tarefa, sub-issue de um épico. |
| **PR empilhado** | PR cuja base é a branch de outro PR ainda não mergeado. |
| **PR de fechamento** | PR que encerra a fase: status das specs, ROADMAP e CHANGELOG. |

## Visão geral

```text
descoberta → ANALYSIS.md → decisões do dono → constituição + specs + ROADMAP
   → épico por fase → tickets (sub-issues) → branch por ticket
   → código + testes → validação local → push → PR → CI verde
   → revisão do dono → merge → PR de fechamento → tag → release
```

## Requisitos funcionais

### 1. Descoberta e análise

- **FR-1** Antes de alterar código num repositório novo para o agente, ele MUST auditá-lo: ler todo o código, compilar, rodar todos os comandos documentados e os casos de borda relevantes (entrada inválida, arquivo existente, stdin fechado, etc.).
- **FR-2** O resultado MUST ser registrado em `specs/ANALYSIS.md`, com: resumo executivo; o que foi verificado e como; achados classificados em **crítico / alto / médio / baixo**, cada um com evidência (`arquivo:linha` ou saída de comando) e a spec que o resolve; pontos positivos a preservar; avaliação do README; melhorias priorizadas por fase; e uma lista explícita de **decisões em aberto** `D1..Dn`, cada uma com as opções e uma recomendação.
- **FR-3** Nenhuma implementação começa antes de o dono responder as decisões em aberto (Fase 0). As respostas MUST ser registradas na seção "Decisões" das specs afetadas.
- **FR-4** A análise MUST distinguir o que foi **verificado** (executado, com evidência) do que é **inferido** (leitura de código).

### 2. Artefatos SDD

- **FR-5** O repositório MUST ter: `specs/constitution.md` (princípios inegociáveis); uma spec por área em `specs/NNN-nome/spec.md`; `specs/README.md` com o índice e o status de cada spec; `specs/ROADMAP.md` com as fases e as tarefas `T1..Tn`.
- **FR-6** Cada spec MUST seguir o formato: cabeçalho (Prioridade, Status, Código afetado e/ou Resolve), Contexto, requisitos `FR-*`/`NFR-*` com MUST/SHOULD/MAY (RFC 2119), critérios de aceite `AC-*` verificáveis (de preferência Dado/Quando/Então), Fora de escopo e Decisões.
- **FR-7** Cada tarefa do ROADMAP MUST citar os IDs que fecha (ex.: `003 FR-2..5`). As fases padrão são: **Fase 0** decisões; **Fase 1** "funcionar de verdade" → `v0.1.0`; **Fase 2** "confiável" → `v0.2.0`; **Fase 3** "profissional" → `v1.0.0`. Um repositório MAY usar outras fases, desde que cada uma termine numa versão.
- **FR-8** Status de spec: `Draft` → `Approved` (decisões respondidas) → `In Progress` → `Done`. Um status intermediário MAY ter uma nota curta (ex.: `In Progress (falta release)`).
- **FR-9** Quando a implementação divergir da spec, a spec MUST ser atualizada **no mesmo PR**, com uma entrada em "Decisões" (ou uma nota "Revisado na implementação" no próprio requisito) explicando o porquê.

### 3. Tickets (issues)

- **FR-10** Cada fase MUST ter um **épico** (labels `épico` e `fase-N`), e cada tarefa da fase MUST ser uma **sub-issue nativa** do épico, para que ele mostre o progresso.
- **FR-11** O corpo do ticket MUST ter as seções **Contexto**, **O que fazer**, **Critérios de aceite** (verificáveis) e as referências **Spec(s)** e **Épico**. Uma seção "Decisão para a revisão" MAY ser incluída quando o ticket embute uma escolha que o dono precisa ver.
- **FR-12** Labels do ticket: `fase-N`; `tipo:feature|docs|ci|teste|chore` ou `bug` (label padrão do GitHub); prioridade `P1` (alta), `P2` (média) ou `P3` (baixa).
- **FR-13** Trabalho descoberto no meio de uma fase MUST virar um novo ticket ligado ao épico da fase corrente (ou de uma fase futura, se não for necessário para a versão); não entra "de carona" num PR existente.
- **FR-14** Todo corpo de issue e todo comentário escrito pelo agente MUST terminar com uma linha em branco, `---` e `_Generated by [Claude Code](https://claude.ai/code)_` (ou o equivalente da ferramenta usada).
- **FR-15** Ao criar tickets por API, o agente SHOULD criar a issue primeiro (o que cria labels inexistentes) e só depois anexá-la como sub-issue: a criação já com o pai falha se alguma label ainda não existir.
- **FR-16** Mudanças triviais (erro de digitação, link quebrado) MAY ir direto para PR, sem ticket.

### 4. Branches

- **FR-17** Uma branch por ticket, com nome `<tipo>/<nº-da-issue>-<descrição-curta>`, onde `<tipo>` é o prefixo do Conventional Commit (`feat`, `fix`, `docs`, `ci`, `test`, `chore`, …).
- **FR-18** A branch parte da `main`; ou da branch da fase anterior, quando essa fase ainda não foi mergeada (o PR fica empilhado, FR-25).
- **FR-19** Uma fase inteira MAY ser entregue num único PR quando for uma reestruturação coesa que não pode ser dividida sem estados intermediários quebrados. Nesse caso o PR referencia o épico e lista as tarefas que fecha.

### 5. Commits

- **FR-20** Mensagens em [Conventional Commits](https://www.conventionalcommits.org/pt-br/), em inglês, no imperativo (`feat: add --watch mode`). Mudanças incompatíveis usam `!` e explicam o impacto no corpo.
- **FR-21** O changelog da release é agrupado a partir desses prefixos; um prefixo errado é um item no lugar errado do changelog.
- **FR-22** Commits do agente MUST trazer as linhas de atribuição exigidas pela ferramenta (ex.: `Co-Authored-By:`), e nunca reescrever histórico já publicado.

### 6. Pull requests

- **FR-23** Um PR por ticket. A primeira linha da descrição MUST ser `Closes #N · Épico #M · Spec NNN` (com `Spec —` quando não houver spec afetada); a nota de PR empilhado (FR-25), quando houver, vem logo em seguida. O título segue Conventional Commits.
- **FR-24** A descrição MUST ter as seções **O que muda** e **Como foi testado**, e SHOULD ter **Notas ou decisões para a revisão** quando houver algo que o dono precise decidir ou saber. Seções adicionais do template do repositório (specs e AC, checklist) MAY ser usadas.
- **FR-25** Um PR empilhado MUST dizer isso no topo, logo após a linha `Closes` (`> PR empilhado sobre #NN`). Quando o PR base é mergeado e sua branch apagada, o GitHub redireciona o empilhado para a `main`; se a branch base não for apagada, o agente MUST redirecionar os empilhados para a `main` (editar a base do PR). Em seguida, MUST conferir se o diff continua só com o ticket e se o CI segue verde.
- **FR-26** O PR MUST ficar no escopo do ticket. O que estiver fora vira ticket novo (FR-13).
- **FR-27** Para evitar conflitos entre PRs paralelos da mesma fase, PRs de ticket MUST NOT editar os arquivos de status compartilhados: status das specs em `specs/README.md` e no cabeçalho das specs, checkboxes do ROADMAP e entradas do CHANGELOG. Isso é feito no PR de fechamento (FR-42). Exceções: o ticket que **cria** um desses arquivos, e o conteúdo normativo das specs (FR-9), que muda junto com o código.
- **FR-28** PRs do agente MUST terminar com o rodapé de atribuição da ferramenta usada (com Claude Code: `🤖 Generated with [Claude Code](https://claude.com/claude-code)` e o link da sessão).

### 7. Validação antes do push

- **FR-29** Antes de cada push, o agente MUST rodar a verificação local completa do repositório (ex.: `make ci`) e só enviar com ela verde. Um push validado vale mais que vários especulativos.
- **FR-30** Para corrigir uma falha de CI, o agente MUST primeiro reproduzi-la (localmente ou pela leitura exata do log) e só então corrigir.
- **FR-31** Todo teste novo MUST passar por uma **checagem de mutação** manual: quebrar temporariamente o código coberto e confirmar que o teste falha; depois desfazer.
- **FR-32** O agente MUST reler o próprio diff de forma adversarial antes do push (escopo, arquivos esquecidos, segredos, saídas geradas, consistência entre README, specs e código).

### 8. Gates de qualidade (CI)

- **FR-33** O CI MUST rodar em todo PR e push para `main` e cobrir, no mínimo: formatação e lint, testes em todos os sistemas operacionais suportados (com race detector onde disponível), gate de cobertura, smoke test do artefato real, build para todos os alvos de release, verificação de vulnerabilidades e, quando houver release automatizada, um ensaio da release sem publicar.
- **FR-34** Um PR só está **pronto para revisão** com o CI verde. PRs com CI vermelho ou pendente SHOULD ficar como rascunho ou ser sinalizados como tal na descrição.
- **FR-35** O alvo local (`make ci` ou equivalente) MUST rodar o subconjunto do CI que é viável localmente, com os mesmos comandos e limites.

### 9. CI vermelho

- **FR-36** O agente MUST acompanhar os eventos dos PRs que abriu (inscrição em atividade do PR) e é responsável por eles até ficarem verdes.
- **FR-37** Diante de uma falha, o agente MUST: diagnosticar pelo log; achar a causa raiz; corrigir; validar (FR-29); enviar. MUST NOT: chamar a falha de instabilidade ("flake") sem evidência; pular, desativar ou enfraquecer testes ou gates; enviar commits vazios para disparar o CI de novo.
- **FR-38** Se a correção exigir mudança fora do escopo do PR mas for condição para ele ficar verde (ex.: atualizar uma ferramenta de lint incompatível), ela MAY entrar no mesmo PR, explicada em "O que muda".

### 10. Revisão e merge

- **FR-39** Com todos os PRs da fase abertos e verdes, o dono revisa na **ordem sugerida no épico**. O agente responde os comentários; corrige pedidos pequenos; para mudanças grandes ou de design, **propõe** (no comentário) e só implementa após concordância.
- **FR-40** Política de merge: PRs de fase sobre os quais outros estão empilhados MUST ser mergeados com **merge commit** (para que os empilhados sejam redirecionados à `main` sem rebase); PRs de ticket MAY usar **squash**.
- **FR-41** O dono SHOULD proteger a `main` (CI obrigatório e revisão de Code Owners). É configuração do dono, não do agente.
- **FR-41a** Antes da rodada de merges, o agente SHOULD simular a integração completa: mergear localmente todas as branches da fase, na ordem do épico, e rodar as verificações locais (`make ci` e `make docs`). Isso acha falhas que nenhum CI individual vê, como um lint novo de um PR reprovando um arquivo criado por outro. Durante a rodada, cada branch que ainda não tem a `main` atual SHOULD recebê-la (merge, não rebase) e ter o CI verde antes do seu merge.

### 11. Fechamento de fase

- **FR-42** Depois de mergeados todos os PRs da fase, um **PR de fechamento** MUST atualizar: o status das specs (`specs/README.md` e cabeçalhos) e a seção "Estado atual" das specs que a tiverem, os checkboxes do ROADMAP e o `CHANGELOG.md` ([Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/): `[Unreleased]` vira `[X.Y.Z] - AAAA-MM-DD`).
- **FR-43** Após o merge do PR de fechamento, o dono cria e envia a tag `vX.Y.Z` (SemVer). O workflow de release MUST rodar o CI completo antes de publicar os binários e os checksums.
- **FR-44** O épico é fechado quando a release da fase estiver publicada.

### 12. Papéis

| Responsabilidade | Dono | Agente |
|---|---|---|
| Responder decisões em aberto, mudar escopo ou prioridades | ✔ | propõe |
| Analisar, escrever specs, criar épicos e tickets | revisa | ✔ |
| Implementar, testar, validar, abrir PRs, manter CI verde | | ✔ |
| Revisar e aprovar PRs | ✔ | responde |
| Merge | ✔ | só com delegação explícita (FR-46) |
| Tags, releases | ✔ | |
| Configurações do repositório (proteção de branch, labels de sistema, segredos) | ✔ | sugere |

- **FR-45** Qualquer ação irreversível ou visível para fora além do fluxo combinado (merge, tag, release, apagar branch alheia, mudar configuração, publicar em outro lugar) é do dono.
- **FR-46** O agente MUST NOT: fazer merge, salvo quando o dono delega explicitamente uma rodada de merges (a delegação vale só para aquela rodada, segue a ordem do épico e exige CI verde no head de cada PR); dar force-push em branch de outra pessoa; reescrever histórico publicado; trabalhar fora das branches do seu ticket.

### 13. Comunicação

- **FR-47** Atualizações de status são curtas: o que foi feito, o que falta, o que bloqueia.
- **FR-48** Achados vão para o repositório (specs, issues, descrições de PR), não só para o chat; o chat aponta para eles.
- **FR-49** O agente MUST dizer claramente o que foi verificado e o que não foi (ex.: "testado no Linux; o Windows só no CI").

### 14. Retomada e economia de uso

O agente pode perder o contexto a qualquer momento: a conversa é compactada, a sessão é reiniciada ou o limite de uso se esgota. O trabalho MUST poder continuar a partir do repositório (#37).

- **FR-50** O ticket MUST ser criado quando a tarefa começa, e o PR MUST ser aberto assim que a tarefa termina e a verificação local passa. Trabalho pronto só na máquina do agente se perde com a sessão.
- **FR-51** Um pedido novo do dono que não cabe na tarefa em curso vira ticket **na hora**, antes de qualquer outra ação.
- **FR-52** O épico da fase MUST ter um comentário **"Estado da fase"**, criado ao abrir os primeiros PRs e atualizado a cada marco: PRs abertos com o estado do CI, decisões tomadas, conflitos previstos e próximo passo.
- **FR-53** Retomada: uma sessão nova lê o comentário de estado mais recente do épico aberto, os PRs e as issues abertas e o `CLAUDE.md`, e continua do próximo passo registrado, sem refazer análises que já estão no repositório.
- **FR-54** O repositório SHOULD ter um `CLAUDE.md` com o mapa do código, os comandos, as convenções e as armadilhas conhecidas, e um hook de início de sessão (`.claude/hooks/session-start.sh`) que instale as ferramentas do CI no ambiente na web.
- **FR-55** Práticas de economia (SHOULD):
  - ler trechos de arquivo (`sed -n`, `grep -n`) e não reler o que já foi lido;
  - para o CI, pedir só o resumo das conclusões e, numa falha, o fim do log do job;
  - validar num comando só (`make ci`);
  - usar subagentes só para buscas amplas;
  - no chat, só o resumo e o próximo passo; detalhes nos PRs.

## Requisitos não funcionais

- **NFR-1** **Rastreabilidade:** de qualquer linha mergeada é possível chegar, por links, ao PR, ao ticket, ao épico e à spec (`FR`/`AC`) que a justificam.
- **NFR-2** **Reaplicável:** o processo não depende de linguagem nem de ferramenta de build; só de GitHub (issues, sub-issues, Actions) e de um alvo local equivalente ao CI.
- **NFR-3** **Leve:** um ticket pequeno não exige mais que o corpo do ticket, uma branch e um PR com as seções mínimas.

## Critérios de aceite

- **AC-1** Dado um repositório novo para o agente, quando ele termina a descoberta, então existe `specs/ANALYSIS.md` com achados classificados por severidade, cada um com evidência, e uma lista `D1..Dn` com recomendação — e nenhum commit de código foi feito antes das respostas do dono.
- **AC-2** Dado uma fase do ROADMAP, então existe um épico com labels `épico` e `fase-N`, e cada tarefa da fase é sub-issue dele.
- **AC-3** Dado um ticket, então seu corpo tem Contexto, O que fazer, Critérios de aceite, Spec(s) e Épico, e ele tem uma label `fase-N`, uma `tipo:*` e uma de prioridade.
- **AC-4** Dado um PR de ticket, então o nome da branch segue `<tipo>/<nº>-<descrição>`, a descrição começa com `Closes #N · Épico #M · Spec NNN` e o CI está verde antes de o PR entrar na revisão.
- **AC-5** Dado um PR empilhado, então a primeira linha da descrição é a linha `Closes` (FR-23) e a seguinte diz sobre qual PR ele está empilhado.
- **AC-6** Dado uma divergência entre implementação e spec, então a spec é atualizada no mesmo PR com uma entrada em Decisões (ou nota "Revisado na implementação").
- **AC-7** Dado dois PRs de ticket da mesma fase, então nenhum deles altera status em `specs/README.md`, checkboxes do ROADMAP ou entradas do CHANGELOG (salvo a criação do arquivo, FR-27).
- **AC-8** Dado um teste novo num PR, quando o código coberto é quebrado de propósito, então o teste falha.
- **AC-9** Dado um CI vermelho num PR do agente, então o commit seguinte corrige a causa raiz descrita no PR ou no commit, sem nenhum teste (novo ou existente) pulado, desativado ou enfraquecido, sem gate rebaixado e sem commit vazio (FR-37).
- **AC-10** Dado o fim de uma fase, então existe um PR de fechamento que atualiza ROADMAP, status das specs e CHANGELOG, e a tag `vX.Y.Z` só é criada depois do merge dele.
- **AC-11** Dado um PR de fase com PRs empilhados sobre ele, quando é mergeado, então o merge é por merge commit e os empilhados passam a apontar para a `main`.
- **AC-12** Dado qualquer issue, PR ou comentário escrito pelo agente, então ele termina com o rodapé de atribuição (FR-14, FR-28).
- **AC-13** Dado o `CONTRIBUTING.md` e esta spec, então não há regra em um que contradiga o outro.
- **AC-14** Dada uma sessão nova sem a conversa anterior, quando ela lê o comentário "Estado da fase" do épico, os PRs e as issues abertas e o `CLAUDE.md`, então consegue dizer o próximo passo da fase sem refazer análises (FR-53).

## Reaplicação em outros repositórios

Checklist para iniciar o processo num repositório novo:

1. [ ] Descoberta e `specs/ANALYSIS.md` com decisões `D1..Dn` (FR-1..4); aguardar as respostas.
2. [ ] `specs/constitution.md` com 5–10 princípios verificáveis.
3. [ ] `specs/README.md` (fluxo SDD, convenções, índice de specs) e uma spec por área.
4. [ ] `specs/ROADMAP.md` com fases → versões e tarefas ligadas a `FR`/`AC`.
5. [ ] Labels: `épico`, `fase-0..N`, `tipo:feature|docs|ci|teste|chore` ou `bug` (label padrão do GitHub), `P1..P3`.
6. [ ] Um épico por fase; tickets como sub-issues (FR-15).
7. [ ] CI com os gates de FR-33 e um alvo local equivalente (FR-35).
8. [ ] Templates de issue (tarefa, bug) e de PR alinhados com FR-11 e FR-23..24.
9. [ ] `CONTRIBUTING.md` que resume esta spec na prática; `CODEOWNERS`.
10. [ ] Release automatizada por tag, precedida do CI completo; `CHANGELOG.md` em Keep a Changelog.
11. [ ] Pedir ao dono: proteção da `main` (CI obrigatório + Code Owners) e política de merge (FR-40).

## Aplicação no cv-craft

| Item | Valor |
|---|---|
| Verificação local | `make ci` = `lint` (`go vet` + `golangci-lint`) + `race` + `cover` (`scripts/coverage.sh`, ≥ 80% em `internal/`) + `smoke` (`scripts/smoke.sh`, binário real) |
| Golden files | `make golden` após mudança intencional nas saídas; o diff dos goldens é revisado no PR |
| CI (`.github/workflows/ci.yml`) | `go mod tidy` sem diff, `go vet`, `golangci-lint`; testes em Linux/macOS/Windows com `-race` (exceto Windows); cobertura; smoke; cross-build linux/darwin/windows × amd64/arm64; `govulncheck` |
| Documentação | job `docs`: `markdownlint-cli2` (inclui os golden files de Markdown), `lychee` (links) e `scripts/doc-commands.sh` (comandos do README); `make docs` localmente — entram com o ticket #8 |
| Ensaio de release | job `release-check` (`goreleaser release --snapshot` + `scripts/check-release.sh`) e `make release-snapshot` — entram com o ticket #6 |
| Release | tag `v*` → `.github/workflows/release.yml` chama o `ci.yml` e publica via GoReleaser com checksums — entra com o ticket #6 |
| Labels | `épico`, `fase-1..3`, `tipo:feature`, `bug`, `tipo:docs`, `tipo:ci`, `tipo:teste`, `tipo:chore`, `P1..P3` |
| Épicos | #1 Fase 1 → `v0.1.0`; #3 Fase 2 → `v0.2.0`; #10 Fase 3 → `v1.0.0` |
| Fase 1 | PR único #2 (branch `claude/determined-darwin-rmv0xj`), por ser uma reestruturação coesa (FR-19) |
| Fase 2 | #4 Go 1.26, #5 CONTRIBUTING, #6 GoReleaser, #7 `docs/schema.md`, #8 verificação de Markdown/links, #9 benchmarks, #19 esta spec — PRs empilhados sobre o #2 |
| Exemplos de branch | `chore/4-go-1.26`, `docs/5-contributing`, `ci/6-goreleaser`, `docs/19-delivery-process-spec` |
| Guia prático | `CONTRIBUTING.md`, `.github/pull_request_template.md`, `.github/ISSUE_TEMPLATE/`, `.github/CODEOWNERS` (entram com o ticket #5) |

Exemplos reais que motivaram requisitos:

- **FR-9:** a spec 007 FR-8 trocou a extração de texto do PDF de uma biblioteca Go (`ledongthuc/pdf`) pelo `pdftotext` do poppler, porque a biblioteca truncava caracteres acima de U+00FF — escondendo justamente os bugs de Unicode.
- **FR-30, FR-37:** o runner Windows falhou porque o Git for Windows traz `pdftotext` mas não `pdfinfo`; os testes passaram a detectar cada ferramenta separadamente. E subir o Go para 1.26 (#4) quebrou o `golangci-lint` v2.5, compilado com Go 1.25; a ferramenta foi atualizada no mesmo PR (FR-38).
- **FR-31:** a checagem de mutação revelou um teste de quebra de página que não testava nada e um teste de alvo do `Makefile` que passava com o alvo quebrado.
- **FR-13:** esta spec nasceu no meio da Fase 2 e virou o ticket #19, ligado ao épico #3.

## Fora de escopo

- Configurações do GitHub em si (proteção de branch, permissões de Actions, segredos): são do dono (FR-41, FR-45).
- Automação que verifique esta spec (lint de títulos de PR, de nomes de branch ou de rodapés). Candidata a ticket futuro.
- Processo de segurança (divulgação de vulnerabilidades) — `SECURITY.md` futuro.

## Decisões

- **Retomada e economia (#37).** Adotados: o `CLAUDE.md` na raiz, o hook de início de sessão (síncrono; só em sessões na web; instala o golangci-lint na versão do CI e o poppler) e o comentário de estado no épico. Não adotado: lista de permissões em `.claude/settings.json`, porque as sessões na web rodam em modo automático e quem contribui localmente configura as próprias permissões.
- **Duas escalas de prioridade.** Specs usam `P0..P2` (definida em `specs/README.md`: P0 bloqueia uso real); tickets usam labels `P1..P3` (alta → baixa). São escalas distintas: a da spec diz quão essencial é a área; a do ticket, a ordem de trabalho dentro da fase.
- **Status no fechamento.** O passo 5 de `specs/README.md` ("Fechar — atualizar o status") acontece no PR de fechamento da fase (FR-42), e não em cada PR de ticket, para evitar conflitos (FR-27).
- **Primeira linha do PR.** `Closes #N · Épico #M · Spec NNN` é o formato normativo; os IDs finos (`003 FR-8, AC-6`) vão na seção de specs do template ou em "O que muda".
- **Delegação de merge.** Na Fase 2 o dono delegou ao agente a rodada de merges dos 9 PRs. A regra padrão continua sendo "o dono mergeia" (FR-46); a delegação é pontual e não se estende a tags nem a releases.
- **Merge commit para PRs de fase.** Squash num PR que tem empilhados reescreveria a base deles e forçaria rebase de todos; com merge commit o GitHub só redireciona a base.
