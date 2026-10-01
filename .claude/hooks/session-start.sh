#!/bin/bash
# Prepara o ambiente do Claude Code na web para rodar `make ci` e `make docs`
# sem instalar nada no meio do trabalho (#37). Só roda em sessões remotas;
# na máquina local de quem contribui, não faz nada.
set -euo pipefail

if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

cd "${CLAUDE_PROJECT_DIR:-$(pwd)}"

# Módulos Go (o container é guardado em cache depois do hook).
go mod download

# golangci-lint na mesma versão do CI: o binário pré-instalado no ambiente pode
# ter sido compilado com um Go mais antigo que o do go.mod e se recusar a rodar.
want="$(sed -n 's/^ *version: *\(v2\.[0-9][0-9.]*\).*/\1/p' .github/workflows/ci.yml | head -n1)"
[[ "$want" =~ ^v[0-9]+\.[0-9]+$ ]] && want="$want.0"
gobin="$(go env GOPATH)/bin"
if ! "$gobin/golangci-lint" version 2>/dev/null | grep -q "version ${want#v} "; then
  GOBIN="$gobin" go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$want"
fi
if [ -n "${CLAUDE_ENV_FILE:-}" ]; then
  echo "export PATH=\"$gobin:\$PATH\"" >> "$CLAUDE_ENV_FILE"
fi

# poppler (pdftotext e pdfinfo): os testes de texto do PDF são obrigatórios no CI.
if ! command -v pdftotext >/dev/null || ! command -v pdfinfo >/dev/null; then
  if command -v apt-get >/dev/null; then
    (apt-get update -qq && apt-get install -y -qq poppler-utils) >/dev/null
  fi
fi
