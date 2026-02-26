#!/bin/sh
# Brf.it installer for Linux and macOS
# Usage: ./install.sh [version]
# Example: ./install.sh v0.5.0

set -e

REPO="indigo-net/Brf.it"
INSTALL_DIR="${BRFIT_INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="brfit"
ARCHIVE_PREFIX="Brf.it"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}==>${NC} %s\n" "$1"
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        linux|darwin)
            echo "$OS"
            ;;
        *)
            error "Unsupported OS: $OS. Use Linux or macOS."
            ;;
    esac
}

# Detect architecture
detect_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH. Use x86_64 or arm64."
            ;;
    esac
}

# Get latest version from GitHub API
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download file
download() {
    URL="$1"
    OUTPUT="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$URL" -o "$OUTPUT"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O "$OUTPUT"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Calculate SHA256 checksum
sha256() {
    FILE="$1"
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$FILE" | awk '{print $1}'
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$FILE" | awk '{print $1}'
    else
        error "Neither sha256sum nor shasum found. Cannot verify checksum."
    fi
}

# Verify checksum
verify_checksum() {
    FILE="$1"
    EXPECTED="$2"
    ACTUAL=$(sha256 "$FILE")

    if [ "$ACTUAL" != "$EXPECTED" ]; then
        error "Checksum mismatch!\n  Expected: $EXPECTED\n  Actual:   $ACTUAL"
    fi
    info "Checksum verified"
}

# Check if sudo is needed
need_sudo() {
    if [ -w "$INSTALL_DIR" ]; then
        echo ""
    else
        echo "sudo"
    fi
}

# Check if PATH contains install directory
check_path() {
    case ":$PATH:" in
        *":$INSTALL_DIR:"*)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

main() {
    # Get version from argument or fetch latest
    VERSION="${1:-}"
    if [ -z "$VERSION" ]; then
        info "Fetching latest version..."
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            error "Failed to get latest version. Please specify version manually."
        fi
    fi

    # Ensure version starts with 'v'
    case "$VERSION" in
        v*)
            ;;
        *)
            VERSION="v$VERSION"
            ;;
    esac

    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected: ${OS}/${ARCH}"
    info "Installing brfit $VERSION"

    # Build download URLs
    VERSION_NUM="${VERSION#v}"
    ARCHIVE_NAME="${ARCHIVE_PREFIX}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download archive
    info "Downloading $ARCHIVE_NAME..."
    download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE_NAME"

    # Download and verify checksum
    info "Downloading checksums..."
    download "$CHECKSUM_URL" "$TMP_DIR/checksums.txt"

    EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "$TMP_DIR/checksums.txt" | awk '{print $1}')
    if [ -z "$EXPECTED_CHECKSUM" ]; then
        error "Checksum not found for $ARCHIVE_NAME"
    fi

    info "Verifying checksum..."
    verify_checksum "$TMP_DIR/$ARCHIVE_NAME" "$EXPECTED_CHECKSUM"

    # Extract archive
    info "Extracting..."
    tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"

    # Install binary
    SUDO=$(need_sudo)
    if [ -n "$SUDO" ]; then
        info "Installing to $INSTALL_DIR (requires sudo)..."
    else
        info "Installing to $INSTALL_DIR..."
    fi

    $SUDO mkdir -p "$INSTALL_DIR"
    $SUDO mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    $SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # Success message
    info "brfit $VERSION installed successfully!"

    # macOS quarantine warning
    if [ "$OS" = "darwin" ]; then
        echo ""
        warn "macOS users: If 'brfit' is blocked, run:"
        echo "    xattr -d com.apple.quarantine $INSTALL_DIR/brfit"
    fi

    # Check PATH
    if ! check_path; then
        echo ""
        warn "$INSTALL_DIR is not in your PATH."
        echo ""
        echo "Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
        echo "Then restart your terminal or run:"
        echo "  source ~/.bashrc  # or ~/.zshrc"
    fi

    echo ""
    echo "Run 'brfit --help' to get started."
}

main "$@"
