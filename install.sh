#!/usr/bin/env bash
set -euo pipefail

REPO="Danneeb/recon"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY="recon"

# Detect OS
case "$(uname -s)" in
  Linux*)             OS="linux" ;;
  Darwin*)            OS="darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

# Detect architecture
case "$(uname -m)" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

# Get latest version from GitHub API
echo "Fetching latest version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

if [[ -z "$VERSION" ]]; then
  echo "Failed to determine latest version" >&2
  exit 1
fi

echo "Installing recon ${VERSION} (${OS}/${ARCH})..."

if [[ "$OS" == "windows" ]]; then
  ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.zip"
else
  ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
fi

URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$URL" -o "${TMP_DIR}/${ARCHIVE}"

cd "$TMP_DIR"
if [[ "$OS" == "windows" ]]; then
  unzip -q "$ARCHIVE"
  BINARY_NAME="${BINARY}.exe"
else
  tar -xzf "$ARCHIVE"
  BINARY_NAME="${BINARY}"
fi

chmod +x "$BINARY_NAME"

if [[ -w "$INSTALL_DIR" ]]; then
  mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
else
  sudo mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "recon ${VERSION} installed to ${INSTALL_DIR}/${BINARY_NAME}"
echo "Run: recon --help"
