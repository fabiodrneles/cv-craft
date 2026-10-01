#!/usr/bin/env bash
# Falha se a cobertura dos pacotes em internal/ ficar abaixo do mínimo
# (spec 007, FR-10).
#
# Uso: scripts/coverage.sh <coverage.out> <mínimo-em-%>
set -euo pipefail

profile="${1:-coverage.out}"
min="${2:-80}"

filtered="$(mktemp)"
trap 'rm -f "$filtered"' EXIT
{ head -n1 "$profile"; grep '/internal/' "$profile"; } > "$filtered"

total="$(go tool cover -func="$filtered" | awk '/^total:/ { sub(/%/, "", $3); print $3 }')"
echo "cobertura de internal/: ${total}% (mínimo ${min}%)"
awk -v t="$total" -v m="$min" 'BEGIN { exit (t + 0 >= m + 0) ? 0 : 1 }' || {
  echo "cobertura abaixo do mínimo"
  exit 1
}
