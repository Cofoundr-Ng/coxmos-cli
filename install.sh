#!/bin/sh
set -e

REPO="Cofoundr-Ng/coxmos-cli"
BIN="coxmos"
DEST="${DEST:-/usr/local/bin/$BIN}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Error: unsupported architecture $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "Error: unsupported OS $OS"; exit 1 ;;
esac

VERSION="${1:-latest}"
if [ "$VERSION" = "latest" ]; then
  echo "Fetching latest release..."
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BIN-${OS}-${ARCH}"

echo
echo "  Coxmos CLI $VERSION"
echo "  Platform:   $OS/$ARCH"
echo "  Destination: $DEST"
echo

curl -fSL# "$URL" -o "$DEST"

chmod +x "$DEST"

echo
echo "  Installed to $DEST"
echo "  Run 'coxmos --help' to get started"
echo
