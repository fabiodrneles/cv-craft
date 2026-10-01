#!/usr/bin/env bash
# Smoke test de ponta a ponta: compila o binário real e o exercita como um
# usuário faria, verificando saídas e exit codes (spec 002).
#
# Uso: scripts/smoke.sh            (a partir da raiz do repositório)
# Requer: go, bash. Se pdftotext estiver disponível, também confere o texto do PDF.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

bin="$work/cv-craft"
[[ "$(go env GOOS)" == "windows" ]] && bin="$bin.exe"

echo "==> compilando"
(cd "$root" && go build -trimpath -ldflags "-X main.version=v0.0.0-smoke" -o "$bin" .)

fail=0
check() { # check <descrição> <exit esperado> <comando...>
  local desc="$1" want="$2"; shift 2
  set +e
  "$@" >"$work/out" 2>"$work/err" </dev/null
  local got=$?
  set -e
  if [[ "$got" == "$want" ]]; then
    echo "ok   $desc (exit $got)"
  else
    echo "FAIL $desc: exit $got, esperava $want"
    sed 's/^/     stdout: /' "$work/out"; sed 's/^/     stderr: /' "$work/err"
    fail=1
  fi
}
contains() { # contains <arquivo> <texto>
  if ! grep -qF -- "$2" "$1"; then
    echo "FAIL $1 não contém: $2"; fail=1
  fi
}

cd "$work"

check "version"                          0 "$bin" version
contains out "v0.0.0-smoke"
check "help"                             0 "$bin" --help
check "comando desconhecido"             2 "$bin" foo
check "init"                             0 "$bin" init
check "init não sobrescreve"             3 "$bin" init
check "validate"                         0 "$bin" validate curriculum.yaml
check "build pdf"                        0 "$bin" build curriculum.yaml
check "saída existe sem --force"         3 "$bin" build curriculum.yaml
check "build --force"                    0 "$bin" build curriculum.yaml --force
check "build all"                        0 "$bin" build -f all -o dist curriculum.yaml
check "build full.yaml en"               0 "$bin" build "$root/examples/full.yaml" --lang en -f txt -o full-en.txt
contains full-en.txt "PROFESSIONAL EXPERIENCE"
contains full-en.txt "São Paulo"

printf 'contact:\n  name: X\n' > invalid.yaml
check "YAML inválido"                    4 "$bin" validate invalid.yaml
contains err "contact.email"
check "arquivo inexistente"              1 "$bin" build nao-existe.yaml
check "ui sem terminal"                  2 "$bin" ui

for f in dist/curriculum.pdf dist/curriculum.md dist/curriculum.txt; do
  [[ -s "$f" ]] || { echo "FAIL $f ausente ou vazio"; fail=1; }
done
head -c 5 dist/curriculum.pdf | grep -q '%PDF-' || { echo "FAIL PDF sem cabeçalho %PDF-"; fail=1; }
contains dist/curriculum.md "## Experiência Profissional"
contains dist/curriculum.txt "Responsabilidades:"

if command -v pdftotext >/dev/null; then
  pdftotext -enc UTF-8 dist/curriculum.pdf pdf.txt
  contains pdf.txt "São Paulo, SP"
  contains pdf.txt "EXPERIÊNCIA PROFISSIONAL"
  contains pdf.txt "Comece com um verbo de ação"
  echo "ok   texto do PDF extraído com pdftotext"
else
  echo "skip pdftotext não instalado"
fi

if [[ "$fail" != 0 ]]; then
  echo "smoke test FALHOU"; exit 1
fi
echo "smoke test OK"
