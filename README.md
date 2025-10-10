# CV-Craft

Uma ferramenta de linha de comando (CLI) para gerar currículos profissionais e otimizados para ATS (Applicant Tracking Systems) a partir de um único arquivo de dados YAML.

## O Problema que Resolve

Manter um currículo atualizado e otimizado para os robôs de recrutamento é um desafio. `CV-Craft` separa o **conteúdo** da **apresentação**, permitindo que você foque nas suas informações enquanto a ferramenta cuida da formatação, garantindo que seu CV seja legível tanto por humanos quanto por máquinas.

## Funcionalidades

- **Fonte Única da Verdade:** Todas as suas informações em um arquivo `curriculum.yaml` fácil de editar.
- **Saídas Múltiplas:** Gere seu currículo em formatos `.pdf`, `.txt` (para máxima compatibilidade com ATS) e `.md` com um único comando.
- **Otimização ATS:** O template padrão é projetado com as melhores práticas para garantir a correta interpretação por sistemas de recrutamento.

## Instalação
```bash
go install [github.com/seu-usuario/cv-craft@latest](https://github.com/seu-usuario/cv-craft@latest)