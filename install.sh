#!/bin/sh
set -e

# multi-gitter-agent installer
# This script installs the multi-gitter-agent binary from GitHub Releases.

OWNER="juftin"
REPO="multi-gitter-agent"
BINARY="multi-gitter-agent"

# Default installation directory
DEFAULT_BIN_DIR="$HOME/.local/bin"
BIN_DIR="${BIN_DIR:-$DEFAULT_BIN_DIR}"

# Parse CLI arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        -d|--bin-dir)
            BIN_DIR="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: install.sh [options]"
            echo ""
            echo "Options:"
            echo "  -d, --bin-dir DIR    Directory to install the binary (default: $DEFAULT_BIN_DIR)"
            echo "  -h, --help           Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="x86_64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux) OS="Linux" ;;
    darwin) OS="Darwin" ;;
    msys*|cygwin*|mingw*) OS="Windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest release tag
echo "🔍 Fetching latest version..."
LATEST_TAG=$(curl -s https://api.github.com/repos/$OWNER/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "❌ Could not find latest release tag."
    exit 1
fi

# Construct download URL
EXTENSION="tar.gz"
if [ "$OS" = "Windows" ]; then
    EXTENSION="zip"
fi

FILENAME="${BINARY}_${LATEST_TAG#v}_${OS}_${ARCH}.${EXTENSION}"
URL="https://github.com/$OWNER/$REPO/releases/download/$LATEST_TAG/$FILENAME"

echo "📥 Downloading $BINARY $LATEST_TAG for $OS ($ARCH)..."
TMP_DIR=$(mktemp -d)
curl -sL "$URL" -o "$TMP_DIR/$FILENAME"

# Extract and Install
echo "📦 Installing to $BIN_DIR..."
mkdir -p "$BIN_DIR"

if [ "$EXTENSION" = "tar.gz" ]; then
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"
else
    unzip -q "$TMP_DIR/$FILENAME" -d "$TMP_DIR"
fi

mv "$TMP_DIR/$BINARY" "$BIN_DIR/$BINARY"
chmod +x "$BIN_DIR/$BINARY"

rm -rf "$TMP_DIR"

echo "✨ Successfully installed $BINARY!"

# Check if BIN_DIR is in PATH
if ! echo "$PATH" | grep -q "$BIN_DIR"; then
    echo "⚠️  Warning: $BIN_DIR is not in your PATH."
    echo "You may need to add it to your shell configuration (e.g., ~/.bashrc or ~/.zshrc):"
    echo "  export PATH=\"\$PATH:$BIN_DIR\""
fi

"$BIN_DIR/$BINARY" --help | head -n 1
