#!/bin/bash

REPO="Kamaliev/gogit"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$OS" == "darwin" ]]; then
  OS_NAME="macos"
else
  OS_NAME=$OS
fi

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

URL="https://github.com/$REPO/releases/$VERSION/download/gogit-${OS_NAME}-${ARCH}"
echo "⬇️  Скачивание $URL..."
curl -L -o gogit "$URL"

chmod +x gogit
sudo mv gogit "$INSTALL_DIR/gogit"

if [[ "$OS_NAME" == "macos" ]]; then sudo xattr -d com.apple.quarantine "$INSTALL_DIR/gogit"; fi

echo "✅ Установка завершена! Теперь можно использовать команду: gogit"
