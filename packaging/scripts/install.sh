#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="robertusnegoro/k8ctl"
BINARY_NAME="k8ctl"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.k8ctl"

# Detect OS and Architecture
detect_os_arch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        armv7l|armv6l)
            ARCH="arm"
            ;;
    esac
    
    echo "$OS-$ARCH"
}

# Get latest version
get_latest_version() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install
install_k8ctl() {
    echo -e "${GREEN}Installing k8ctl...${NC}"
    
    OS_ARCH=$(detect_os_arch)
    VERSION=$(get_latest_version)
    
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Failed to get latest version${NC}"
        exit 1
    fi
    
    echo -e "Detected: ${YELLOW}$OS_ARCH${NC}"
    echo -e "Version: ${YELLOW}$VERSION${NC}"
    
    # Download
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/k8ctl_${OS_ARCH}.tar.gz"
    echo -e "Downloading from: ${YELLOW}$DOWNLOAD_URL${NC}"
    
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    curl -L -o k8ctl.tar.gz "$DOWNLOAD_URL" || {
        echo -e "${RED}Failed to download k8ctl${NC}"
        exit 1
    }
    
    # Extract
    tar -xzf k8ctl.tar.gz
    
    # Install
    if [ -w "$INSTALL_DIR" ]; then
        cp k8ctl "$INSTALL_DIR/"
    else
        sudo cp k8ctl "$INSTALL_DIR/"
    fi
    
    # Make executable
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # Create config directory
    mkdir -p "$CONFIG_DIR"
    
    # Install shell completion
    install_completion
    
    # Cleanup
    cd -
    rm -rf "$TEMP_DIR"
    
    echo -e "${GREEN}Installation complete!${NC}"
    echo -e "Run '${YELLOW}$BINARY_NAME version${NC}' to verify installation."
}

# Install shell completion
install_completion() {
    SHELL=$(basename "$SHELL")
    
    case $SHELL in
        bash)
            COMPLETION_DIR="/etc/bash_completion.d"
            if [ ! -w "$COMPLETION_DIR" ]; then
                COMPLETION_DIR="$HOME/.bash_completion.d"
                mkdir -p "$COMPLETION_DIR"
            fi
            "$INSTALL_DIR/$BINARY_NAME" completion bash > "$COMPLETION_DIR/$BINARY_NAME" || true
            ;;
        zsh)
            COMPLETION_DIR="${fpath[1]}"
            if [ -z "$COMPLETION_DIR" ]; then
                COMPLETION_DIR="$HOME/.zsh/completions"
                mkdir -p "$COMPLETION_DIR"
            fi
            "$INSTALL_DIR/$BINARY_NAME" completion zsh > "$COMPLETION_DIR/_$BINARY_NAME" || true
            ;;
        fish)
            COMPLETION_DIR="$HOME/.config/fish/completions"
            mkdir -p "$COMPLETION_DIR"
            "$INSTALL_DIR/$BINARY_NAME" completion fish > "$COMPLETION_DIR/$BINARY_NAME.fish" || true
            ;;
    esac
}

# Main
main() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${YELLOW}$BINARY_NAME is already installed${NC}"
        read -p "Do you want to reinstall? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 0
        fi
    fi
    
    install_k8ctl
}

main "$@"
