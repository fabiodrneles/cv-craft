#!/usr/bin/env bash
# Confere o resultado do GoReleaser: 6 arquivos (linux/darwin/windows ×
# amd64/arm64), checksums e a versão embutida no binário (spec 002, AC-10).
#
# Uso: scripts/check-release.sh <diretório-dist>
set -euo pipefail

dist="${1:-dist}"
fail=0

archives=$(find "$dist" -maxdepth 1 \( -name 'cv-craft_*.tar.gz' -o -name 'cv-craft_*.zip' \) | wc -l)
if [[ "$archives" -ne 6 ]]; then
  echo "FAIL esperava 6 arquivos de release, achei $archives"; ls "$dist"; fail=1
else
  echo "ok   6 arquivos de release"
fi

for target in linux_amd64 linux_arm64 darwin_amd64 darwin_arm64 windows_amd64 windows_arm64; do
  ext=tar.gz; [[ "$target" == windows_* ]] && ext=zip
  ls "$dist"/cv-craft_*_"$target".$ext >/dev/null 2>&1 || { echo "FAIL falta o arquivo de $target"; fail=1; }
done

if [[ -s "$dist/checksums.txt" ]] && [[ "$(wc -l < "$dist/checksums.txt")" -eq 6 ]]; then
  echo "ok   checksums.txt com 6 entradas"
else
  echo "FAIL checksums.txt ausente ou incompleto"; fail=1
fi

# O conteúdo do arquivo inclui licenças e documentação.
listing=$(tar -tzf "$dist"/cv-craft_*_linux_amd64.tar.gz)
for f in cv-craft LICENSE README.md CHANGELOG.md licenses/OFL.txt; do
  grep -qx "$f" <<<"$listing" || { echo "FAIL arquivo linux_amd64 sem $f"; fail=1; }
done

# A versão é injetada via -ldflags.
bin=$(find "$dist" -path '*linux_amd64*' -name cv-craft -type f | head -n1)
out=$("$bin" version)
if [[ "$out" =~ ^cv-craft\ v[0-9]+\.[0-9]+\.[0-9]+.*\(commit\ [0-9a-f]{7,},\ [0-9]{4}- ]]; then
  echo "ok   $out"
else
  echo "FAIL versão inesperada: $out"; fail=1
fi

[[ "$fail" == 0 ]] && echo "release OK" || { echo "release FALHOU"; exit 1; }
