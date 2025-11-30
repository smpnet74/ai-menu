#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="anthropics/ai-menu"
INSTALL_DIR="${INSTALL_DIR:-.}"
BINARY_NAME="ai-menu"

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s)
    local arch=$(uname -m)

    case "$os" in
        Linux)
            OS="linux"
            ;;
        Darwin)
            echo -e "${RED}Error: ai-menu is designed to run in a Linux devcontainer${NC}"
            echo "Please run this script inside your devcontainer, not on macOS directly."
            exit 1
            ;;
        *)
            echo -e "${RED}Error: Unsupported OS: $os${NC}"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64)
            ARCH="arm64"
            ;;
        arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Error: Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
}

# Get the latest release
get_latest_release() {
    local url="https://api.github.com/repos/$REPO/releases/latest"

    if ! command -v curl &> /dev/null; then
        echo -e "${RED}Error: curl is required but not installed${NC}"
        exit 1
    fi

    curl -s "$url" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4
}

# Download binary from release
download_binary() {
    local version="$1"
    local filename="$BINARY_NAME-$OS-$ARCH"
    local download_url="https://github.com/$REPO/releases/download/$version/$filename"
    local temp_file="/tmp/$filename.tmp"

    echo -e "${YELLOW}Downloading ai-menu $version for $OS/$ARCH...${NC}"

    if ! curl -fsSL -o "$temp_file" "$download_url"; then
        echo -e "${RED}Error: Failed to download binary from $download_url${NC}"
        exit 1
    fi

    chmod +x "$temp_file"
    echo "$temp_file"
}

# Verify checksum (optional)
verify_checksum() {
    local binary="$1"
    local version="$2"
    local filename="$BINARY_NAME-$OS-$ARCH"
    local checksum_url="https://github.com/$REPO/releases/download/$version/$filename.sha256"
    local temp_checksum="/tmp/$filename.sha256"

    echo -e "${YELLOW}Verifying checksum...${NC}"

    if ! curl -fsSL -o "$temp_checksum" "$checksum_url"; then
        echo -e "${YELLOW}Warning: Could not download checksum file, skipping verification${NC}"
        return 0
    fi

    if command -v sha256sum &> /dev/null; then
        if sha256sum -c "$temp_checksum" --ignore-missing > /dev/null 2>&1; then
            echo -e "${GREEN}Checksum verified${NC}"
            rm -f "$temp_checksum"
            return 0
        else
            echo -e "${RED}Error: Checksum verification failed${NC}"
            rm -f "$temp_checksum"
            exit 1
        fi
    else
        echo -e "${YELLOW}Warning: sha256sum not available, skipping verification${NC}"
        return 0
    fi
}

# Install binary
install_binary() {
    local binary="$1"
    local dest="$INSTALL_DIR/$BINARY_NAME"

    mkdir -p "$INSTALL_DIR"

    if [ ! -w "$INSTALL_DIR" ]; then
        echo -e "${RED}Error: No write permission to $INSTALL_DIR${NC}"
        exit 1
    fi

    mv "$binary" "$dest"
    chmod +x "$dest"

    echo -e "${GREEN}Successfully installed ai-menu to $dest${NC}"
    echo -e "${GREEN}You can now run: $dest${NC}"
}

# Main installation flow
main() {
    echo -e "${GREEN}ai-menu Installer${NC}"
    echo "===================="

    detect_platform
    echo -e "${GREEN}Detected: $OS/$ARCH${NC}"

    echo ""
    VERSION=$(get_latest_release)
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Error: Could not determine latest release${NC}"
        exit 1
    fi

    echo -e "${GREEN}Latest version: $VERSION${NC}"
    echo ""

    BINARY=$(download_binary "$VERSION")
    verify_checksum "$BINARY" "$VERSION"
    install_binary "$BINARY"

    echo ""
    echo -e "${GREEN}Installation complete!${NC}"
}

main "$@"
