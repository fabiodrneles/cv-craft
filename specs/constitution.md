# Constituição do CV-Craft

Princípios que toda spec e todo PR devem respeitar. Mudá-los exige uma spec própria.

1. **Fonte única da verdade.** O YAML é o único lugar onde o conteúdo do currículo vive. Todos os formatos de saída são derivados dele e **não podem omitir silenciosamente** um campo preenchido que o schema declara.
2. **ATS primeiro.** Toda saída deve ser legível por máquinas: coluna única, texto selecionável, sem tabelas, imagens ou ícones carregando informação, ordem de leitura igual à ordem visual.
3. **Texto correto em qualquer idioma latino.** UTF-8 de ponta a ponta. Um acento corrompido é bug crítico.
4. **Determinismo.** Mesma entrada + mesma versão ⇒ mesma saída (exceto metadados de data do PDF, que devem ser controláveis em testes).
5. **Offline e local.** Nenhuma chamada de rede, telemetria ou dependência de serviço externo. Os dados pessoais nunca saem da máquina.
6. **Falhar alto, falhar cedo.** Erros de validação são reportados todos de uma vez, com caminho do campo (`experience[1].role`) e exit code ≠ 0. Nada é descartado em silêncio.
7. **Scriptável.** Todo comando funciona sem TTY: sem prompts quando stdin não é terminal, flags para tudo, exit codes documentados.
8. **Uma implementação por comportamento.** CLI e modo interativo chamam o mesmo serviço; não há lógica de negócio duplicada.
9. **Testado.** Todo critério de aceite tem teste automatizado; geradores têm golden files; CI bloqueia merge vermelho.
10. **Dependências mínimas.** Só biblioteca padrão + dependências justificadas em spec.
