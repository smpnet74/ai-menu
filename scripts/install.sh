#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="smpnet74/ai-menu"
INSTALL_DIR="${INSTALL_DIR:-.}"
BINARY_NAME="ai-menu"

echo -e "${GREEN}ai-menu Installer${NC}"
echo "===================="

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    Linux)
        GOOS="linux"
        ;;
    Darwin)
        echo -e "${RED}Error: ai-menu is designed to run in a Linux devcontainer${NC}"
        exit 1
        ;;
    *)
        echo -e "${RED}Error: Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        GOARCH="amd64"
        ;;
    aarch64|arm64)
        GOARCH="arm64"
        ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Detected: $GOOS/$GOARCH${NC}"
echo ""

# Get latest release version
echo -e "${YELLOW}Fetching latest release...${NC}"
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Could not determine latest release${NC}"
    exit 1
fi

echo -e "${GREEN}Latest version: $VERSION${NC}"
echo ""

# Download binary
BINARY_FILENAME="${BINARY_NAME}-${GOOS}-${GOARCH}"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_FILENAME"
TEMP_FILE="/tmp/${BINARY_FILENAME}.tmp"

echo -e "${YELLOW}Downloading from: $DOWNLOAD_URL${NC}"

if ! curl -fsSL -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    echo -e "${RED}Error: Failed to download binary${NC}"
    exit 1
fi

if [ ! -f "$TEMP_FILE" ] || [ ! -s "$TEMP_FILE" ]; then
    echo -e "${RED}Error: Binary file is empty or missing${NC}"
    exit 1
fi

echo -e "${GREEN}Download successful${NC}"

# Verify checksum (optional, non-fatal)
CHECKSUM_FILE="/tmp/${BINARY_FILENAME}.sha256"
CHECKSUM_URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY_FILENAME}.sha256"

echo -e "${YELLOW}Verifying checksum...${NC}"

if curl -fsSL -o "$CHECKSUM_FILE" "$CHECKSUM_URL" 2>/dev/null; then
    if [ -s "$CHECKSUM_FILE" ]; then
        # Copy binary to temp location for checksum verification
        VERIFY_DIR=$(mktemp -d)
        cp "$TEMP_FILE" "$VERIFY_DIR/$BINARY_FILENAME"
        cd "$VERIFY_DIR"

        if sha256sum -c "$CHECKSUM_FILE" > /dev/null 2>&1; then
            echo -e "${GREEN}Checksum verified${NC}"
        else
            echo -e "${YELLOW}Warning: Checksum verification inconclusive, proceeding anyway${NC}"
        fi

        cd - > /dev/null
        rm -rf "$VERIFY_DIR"
    else
        echo -e "${YELLOW}Warning: Checksum file empty, skipping verification${NC}"
    fi
else
    echo -e "${YELLOW}Warning: Could not download checksum file, skipping verification${NC}"
fi

rm -f "$CHECKSUM_FILE"

# Install binary
echo ""
echo -e "${YELLOW}Installing...${NC}"

mkdir -p "$INSTALL_DIR"

if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${RED}Error: No write permission to $INSTALL_DIR${NC}"
    exit 1
fi

DEST="$INSTALL_DIR/$BINARY_NAME"
mv "$TEMP_FILE" "$DEST"
chmod +x "$DEST"

echo ""
echo -e "${GREEN}Installation complete!${NC}"
echo -e "${GREEN}ai-menu installed to: $DEST${NC}"
echo -e "${GREEN}Run: $DEST${NC}"
