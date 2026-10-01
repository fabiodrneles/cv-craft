#!/bin/sh
# Instala o CV-Craft no Linux ou no macOS a partir da última release do GitHub.
#
#   curl -fsSL https://raw.githubusercontent.com/fabiodrneles/cv-craft/main/scripts/install.sh | sh
#
# Variáveis opcionais:
#   CV_CRAFT_VERSION      versão a instalar (ex.: 1.0.0); padrão: a mais recente
#   CV_CRAFT_INSTALL_DIR  pasta de destino; padrão: ~/.local/bin
set -eu

repo="fabiodrneles/cv-craft"
dir="${CV_CRAFT_INSTALL_DIR:-$HOME/.local/bin}"

fail() { echo "erro: $*" >&2; exit 1; }
command -v curl >/dev/null || fail "o curl é necessário para baixar o CV-Craft"
command -v tar >/dev/null || fail "o tar é necessário para extrair o CV-Craft"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux | darwin) ;;
  *) fail "sistema não suportado: $os (no Windows, use o install.ps1)" ;;
esac
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) fail "arquitetura não suportada: $(uname -m)" ;;
esac

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# A versão vem do checksums.txt da release (e não da API do GitHub, que limita
# requisições anônimas por IP e falha em redes compartilhadas).
version="${CV_CRAFT_VERSION:-}"
if [ -n "$version" ]; then
  base="https://github.com/$repo/releases/download/v${version#v}"
else
  base="https://github.com/$repo/releases/latest/download"
fi
curl -fsSL -o "$tmp/checksums.txt" "$base/checksums.txt" || fail "não foi possível baixar $base/checksums.txt"
archive="$(grep -o "cv-craft_[^ ]*_${os}_${arch}\.tar\.gz" "$tmp/checksums.txt" | head -n 1)"
[ -n "$archive" ] || fail "a release não tem um arquivo para $os/$arch"
version="${archive#cv-craft_}"
version="${version%%_*}"

echo "Baixando o CV-Craft $version ($os/$arch)..."
curl -fsSL -o "$tmp/$archive" "$base/$archive" || fail "não foi possível baixar $base/$archive"

expected="$(grep " $archive\$" "$tmp/checksums.txt" | cut -d ' ' -f 1)"
if command -v sha256sum >/dev/null; then
  actual="$(sha256sum "$tmp/$archive" | cut -d ' ' -f 1)"
else
  actual="$(shasum -a 256 "$tmp/$archive" | cut -d ' ' -f 1)"
fi
[ -n "$expected" ] && [ "$expected" = "$actual" ] || fail "o checksum do arquivo baixado não confere"

tar -xzf "$tmp/$archive" -C "$tmp" cv-craft
mkdir -p "$dir"
mv "$tmp/cv-craft" "$dir/cv-craft"
chmod +x "$dir/cv-craft"
echo "Instalado em $dir/cv-craft"
"$dir/cv-craft" version

case ":$PATH:" in
  *":$dir:"*) ;;
  *)
    profile="$HOME/.profile"
    case "${SHELL:-}" in
      */zsh) profile="$HOME/.zshrc" ;;
      */bash) profile="$HOME/.bashrc" ;;
    esac
    echo
    echo "A pasta $dir não está no seu PATH. Para usar o comando cv-craft em qualquer pasta, rode:"
    echo
    echo "  echo 'export PATH=\"$dir:\$PATH\"' >> $profile && export PATH=\"$dir:\$PATH\""
    ;;
esac
