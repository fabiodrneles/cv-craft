# Exemplos

Arquivos YAML prontos para gerar currículos de exemplo. Todos são validados pelos testes do projeto e servem de base para os golden files das saídas (`internal/generator/testdata/`).

| Arquivo | O que demonstra |
|---|---|
| [`minimal.yaml`](minimal.yaml) | O modelo criado por `cv-craft init`: todos os campos, com comentários explicando cada um e marcando os obrigatórios. |
| [`full.yaml`](full.yaml) | Um currículo completo em português: três experiências com responsabilidades, conquistas e tecnologias; duas formações (uma com TCC e disciplinas); certificados com link; idiomas; níveis de habilidade variados e uma categoria sem nível. Ocupa duas páginas no PDF. |
| [`en.yaml`](en.yaml) | Um currículo em inglês (`meta.locale: en`): os títulos das seções saem em inglês. |

## Como gerar

A partir da raiz do repositório:

```bash
cv-craft build examples/full.yaml --format all -o dist/
cv-craft build examples/en.yaml -f pdf -o dist/en.pdf
cv-craft validate examples/minimal.yaml
```

Para gerar o mesmo conteúdo com os títulos em outro idioma, use `--lang`:

```bash
cv-craft build examples/full.yaml --lang en -o dist/full-en.pdf
```

A referência de todos os campos está em [`docs/schema.md`](../docs/schema.md).
