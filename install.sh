#!/bin/bash
set -e

REPO="indigo-net/Brf.it"
INSTALL_DIR="${BRFIT_INSTALL_DIR:-/usr/local/bin}"

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}==>${NC} $1"; }
warn() { echo -e "${YELLOW}==>${NC} $1"; }
error() { echo -e "${RED}Error:${NC} $1" >&2; exit 1; }

# Detect OS/Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) error "Unsupported architecture: $ARCH" ;;
esac

info "Detected: ${OS}/${ARCH}"

# Get latest version
info "Fetching latest version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
[ -z "$VERSION" ] && error "Failed to get latest version"
info "Latest version: $VERSION"

# Filename and URLs
FILENAME="brfit_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"
CHECKSUM_URL="https://github.com/$REPO/releases/download/$VERSION/checksums.txt"

# Download
info "Downloading $FILENAME..."
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT
curl -fsSL "$URL" -o "$TMPDIR/$FILENAME"
curl -fsSL "$CHECKSUM_URL" -o "$TMPDIR/checksums.txt"

# Verify checksum
info "Verifying checksum..."
cd "$TMPDIR"
EXPECTED=$(grep "$FILENAME" checksums.txt | awk '{print $1}')
if command -v sha256sum &> /dev/null; then
  ACTUAL=$(sha256sum "$FILENAME" | awk '{print $1}')
elif command -v shasum &> /dev/null; then
  ACTUAL=$(shasum -a 256 "$FILENAME" | awk '{print $1}')
else
  error "No SHA256 utility found"
fi

[ "$EXPECTED" != "$ACTUAL" ] && error "Checksum mismatch!"
info "Checksum verified"

# Install
info "Installing to $INSTALL_DIR..."
tar -xzf "$FILENAME"
if [ -w "$INSTALL_DIR" ]; then
  mv brfit "$INSTALL_DIR/"
else
  sudo mv brfit "$INSTALL_DIR/"
fi

info "brfit $VERSION installed successfully!"

# macOS user guidance
if [ "$OS" = "darwin" ]; then
  warn "macOS users: If 'brfit' is blocked, run:"
  echo "    xattr -d com.apple.quarantine $INSTALL_DIR/brfit"
fi

echo ""
echo "Run 'brfit --help' to get started."
