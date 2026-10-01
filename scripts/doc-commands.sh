#!/usr/bin/env bash
# Executa os comandos "cv-craft ..." dos blocos ```bash da documentação e
# exige exit 0 em todos (spec 008, FR-2 e AC-2): um exemplo do README que deixa
# de funcionar quebra o CI.
#
# Uso: scripts/doc-commands.sh [arquivo.md ...]   (padrão: README.md examples/README.md)
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
docs=("$@")
if [[ ${#docs[@]} -eq 0 ]]; then
  for d in README.md examples/README.md README.en.md; do
    [[ -f "$root/$d" ]] && docs+=("$d")
  done
fi
for d in "${docs[@]}"; do
  [[ -f "$root/$d" ]] || { echo "FAIL arquivo não encontrado: $d"; exit 1; }
done

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
mkdir -p "$work/bin"
bin="$work/bin/cv-craft"
[[ "$(go env GOOS)" == "windows" ]] && bin="$bin.exe"
(cd "$root" && go build -o "$bin" .)

fail=0
total=0
for doc in "${docs[@]}"; do
  # Cada documento roda num diretório limpo, com uma cópia de examples/.
  dir="$work/run-$(basename "$doc" .md)-$RANDOM"
  mkdir -p "$dir"
  cp -r "$root/examples" "$dir/examples"

  # Linhas que começam com "cv-craft " dentro de blocos ```bash.
  # (while read em vez de mapfile: funciona também no bash 3.2 do macOS)
  cmds=()
  while IFS= read -r line; do
    cmds+=("$line")
  done < <(awk '
    /^```bash/ { inside = 1; next }
    /^```/     { inside = 0; next }
    inside && /^cv-craft / { print }
  ' "$root/$doc")

  for cmd in ${cmds[@]+"${cmds[@]}"}; do
    total=$((total + 1))
    if (cd "$dir" && PATH="$work/bin:$PATH" bash -c "$cmd" >"$work/out" 2>&1 </dev/null); then
      echo "ok   [$doc] $cmd"
    else
      echo "FAIL [$doc] $cmd"
      sed 's/^/     /' "$work/out"
      fail=1
    fi
  done
done

if [[ "$total" -eq 0 ]]; then
  echo "FAIL nenhum comando cv-craft encontrado em ${docs[*]}"
  exit 1
fi
[[ "$fail" == 0 ]] && echo "$total comandos da documentação OK" || { echo "comandos da documentação FALHARAM"; exit 1; }
