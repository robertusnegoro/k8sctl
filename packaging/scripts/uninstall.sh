#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BINARY_NAME="k8ctl"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.k8ctl"

uninstall_k8ctl() {
    echo -e "${YELLOW}Uninstalling k8ctl...${NC}"
    
    # Remove binary
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        if [ -w "$INSTALL_DIR" ]; then
            rm "$INSTALL_DIR/$BINARY_NAME"
        else
            sudo rm "$INSTALL_DIR/$BINARY_NAME"
        fi
        echo -e "${GREEN}Removed binary${NC}"
    fi
    
    # Remove shell completion
    SHELL=$(basename "$SHELL")
    case $SHELL in
        bash)
            rm -f "/etc/bash_completion.d/$BINARY_NAME"
            rm -f "$HOME/.bash_completion.d/$BINARY_NAME"
            ;;
        zsh)
            rm -f "${fpath[1]}/_$BINARY_NAME"
            rm -f "$HOME/.zsh/completions/_$BINARY_NAME"
            ;;
        fish)
            rm -f "$HOME/.config/fish/completions/$BINARY_NAME.fish"
            ;;
    esac
    
    # Remove config (optional)
    read -p "Remove config directory ($CONFIG_DIR)? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        echo -e "${GREEN}Removed config directory${NC}"
    fi
    
    echo -e "${GREEN}Uninstallation complete!${NC}"
}

uninstall_k8ctl
