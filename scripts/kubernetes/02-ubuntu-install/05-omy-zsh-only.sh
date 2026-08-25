#!/bin/bash

# Exit on error
set -e

apply_theme() {
    local destZshrcFilePath="$1"
    local theme="$2"

    if grep -q '^ZSH_THEME=' "$destZshrcFilePath"; then
        sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="'$theme'"/' "$destZshrcFilePath"
    else
        echo 'ZSH_THEME="'$theme'"' >> "$destZshrcFilePath"
    fi
}

# Example Usage:
# install_oh_my_zsh_on_missing
install_oh_my_zsh_on_missing() {
    local homeDir="$HOME"

    if [ ! -d "$homeDir/.oh-my-zsh" ]; then
        wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh
    fi
}

# Example Usage:
# is_append_zshrc_file "yes" "/home/user/source/.zshrc" "/home/user/destination/.zshrc"
is_append_zshrc_file() {
    is_append_zshrc="$1"          # "yes" or "no"
    srcZshrcFile="$2"          # e.g., "/path/to/src/.zshrc"
    destZshrcFilePath="$3"     # e.g., "/path/to/dest/.zshrc"

    if [ "$is_append_zshrc" = "yes" ] && ! grep -qF "$(cat "$srcZshrcFile")" "$destZshrcFilePath"; then
        cat "$srcZshrcFile" >> "$destZshrcFilePath"
    fi
}

# Example Usage:
# install_oh_my_zsh_on_missing
install_oh_my_zsh_on_missing() {
    local homeDir="$HOME"

    if [ ! -d "$homeDir/.oh-my-zsh" ]; then
        wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh
    fi
}

install_oh_my_zsh_root() {
    local theme="$1"
    local is_append_zshrc="${2:-yes}"
    local srcZshrcFile="$PWD/.zshrc"
    local destZshrcFilePath="$HOME/.zshrc"
    local authorizedKeyFile="$PWD/authorized_keys"

    # Install ZSH if not already installed
    if ! command -v zsh &> /dev/null; then
        sudo apt update -y && \
        sudo apt install -y zsh git
    fi

    cd $HOME
    mkdir -p "$HOME/scripts"
    mkdir -p "$HOME/gitlab"
    mkdir -p "$HOME/github"
    mkdir -p "$HOME/.ssh"
    touch "$HOME/.ssh/authorized_keys"
    sudo chmod 600 "$HOME/.ssh/authorized_keys"
    sudo cp -f "$authorizedKeyFile" "$HOME/.ssh/"
    sudo chmod 600 "$HOME/.ssh/authorized_keys"

    # Set ZSH as the default shell
    chsh -s "$(which zsh)"

    # Install Oh My Zsh
    install_oh_my_zsh_on_missing

    # Change ZSH theme dynamically
    apply_theme "$destZshrcFilePath" "$theme"

    # Append text from src .zshrc if not already appended and is_append_zshrc is yes
    is_append_zshrc_file "$is_append_zshrc" "$srcZshrcFile" "$destZshrcFilePath"

    # Install zsh-autosuggestions plugin
    if [ ! -d "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/plugins/zsh-autosuggestions" ]; then
        git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
    fi

    # Source the updated .zshrc
    zsh -c "source $destZshrcFilePath"

    echo "Oh My Zsh installed, theme set to $theme, and .zshrc content updated."
}

# Check for help argument
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    echo "Usage: $0 {theme} {isAppendZshConfig:yes/no}"
    echo
    echo "Install Oh My Zsh, set the ZSH theme, and optionally append .zshrc content."
    echo
    echo "Arguments:"
    echo "  theme             The name of the ZSH theme to apply."
    echo "  isAppendZshConfig Optional. Whether to append .zshrc content if not already appended. Default is 'yes'."
    echo
    echo "Examples:"
    echo "  $0 my-theme-name yes"
    echo "  $0 my-theme-name"
    exit 0
fi

# Call the function with arguments
install_oh_my_zsh_root "$1" "${2:-yes}"
