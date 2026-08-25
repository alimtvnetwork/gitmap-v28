#!/bin/bash

# Exit on error
set -e

# Function to display help
display_help() {
    echo "Md Alim Ul Karim"
    echo "https://alimkarim.com | https://riseup-pro.com"
    echo
    echo "Usage: $0 {theme-name} [append-zshrc]"
    echo
    echo "Change the ZSH theme and optionally append .zshrc content."
    echo
    echo "Arguments:"
    echo "  theme-name       The name of the ZSH theme to apply."
    echo "  append-zshrc     Optional. Whether to append .zshrc content if not already appended. Default is 'false'."
    echo
    echo "Examples:"
    echo "  $0 my-theme-name true"
    echo "  $0 my-theme-name"
    echo
    echo "## Themes"
    echo
    echo "For more themes, visit: https://github.com/ohmyzsh/ohmyzsh/wiki/Themes"
    echo
    echo "Themes Names:"
    echo "  1. amuse"
    echo "  2. robbyrussell"
    echo "  3. af-magic"
    echo "  4. afowler"
    echo "  5. agnoster"
    echo "  6. alanpeabody"
    echo "  7. apple"
    echo "  8. arrow"
    echo "  9. aussiegeek"
    echo "  10. avit"
    echo "  11. awesomepanda"
    echo "  12. bira"
    echo "  13. blinks"
    echo "  14. bureau"
    echo "  15. candy"
    echo "  16. clean"
    echo "  17. cloud"
    echo "  18. crcandy"
    echo "  19. crunch"
    echo "  20. cypher"
    echo "  21. dallas"
    echo "  22. darkblood"
    echo "  23. daveverwer"
    echo "  24. dpoggi"
    echo "  25. dst"
    echo "  26. duellj"
    echo "  27. edvardm"
    echo "  28. fino-time"
    echo "  29. fino"
    echo "  30. fishy"
    echo "  31. flettcherm"
    echo "  32. fox"
    echo "  33. frisk"
    echo "  34. frontcube"
    echo "  35. funky"
    echo "  36. fwalch"
    echo "  37. gnzh"
    echo "  38. godzilla"
    echo "  39. half-life"
    echo "  40. intheloop"
    echo "  41. itchy"
}

prompt_for_theme() {
    local themeName
    read -p "Please enter the theme name (default is 'robbyrussell'): " themeName
    echo "${themeName:-robbyrussell}"
}

# Change ZSH Theme
change_zsh_theme() {
    themeName="$1"
    local destZshrcFilePath="$HOME/.zshrc"
    
    echo "Applying the theme: $themeName"

    # Check if theme name is provided as an argument
    if [ -z "$themeName" ]; then
        themeName=$(prompt_for_theme)
    fi

    echo "  => Final theme: '$themeName' ✅ ."

    if grep -q '^ZSH_THEME=' "$destZshrcFilePath"; then
        sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="'$themeName'"/' "$destZshrcFilePath"
    else
        echo 'ZSH_THEME="'$themeName'"' >> "$destZshrcFilePath"
    fi
}

install_omyzsh() {
    local zshFirst="${ZSH:-$HOME/.oh-my-zsh}"
    local zshRoot=$(eval echo "$zshFirst")
    local savedPWD="$PWD"

    if [ ! -d "$zshRoot" ]; then
        echo ""
        echo "📁 ZSH Root Dir: $zshRoot"
        echo "📁 cd -> $HOME"
        cd "$HOME"
    fi

    # Check if Oh My Zsh is installed
    if [ ! -d "$zshRoot" ]; then
        echo "Oh My Zsh is not installed. Installing now..."
        # Use the --unattended flag to avoid changing the shell during installation
        # sh -c "$(wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O -)" --unattended
        chsh -s "$(which zsh)"
        wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh
        git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions

        echo "Oh My Zsh is installed 👍 ."
        echo ""        
        echo "📁 Expanding dir: $zshRoot"
        
        ls -la "$zshRoot"

        apply_chmod
    else
        echo "Oh My Zsh is already installed 🆗 ."
    fi

    if [ ! -d "$zshRoot" ]; then
        cd "$savedPWD"
        echo "📁 cd -> $savedPWD"
    fi
}

apply_chmod() {
    local zshRoot="${ZSH:-$HOME/.oh-my-zsh}"
    local zshScript="$zshRoot/oh-my-zsh.sh"

    # Apply chmod +x to the oh-my-zsh.sh script
    if [ -f "$zshScript" ]; then
        chmod +x "$zshScript"
        echo "Permissions set for oh-my-zsh.sh 🆗 ."
    else
        echo "oh-my-zsh.sh script not found ❌ ."
    fi
}

# Remove Specific Text
remove_text_from_zshrc() {
    local textToRemove="plugins=(git)"
    local destZshrcFilePath="$HOME/.zshrc"
    
    # echo "disabling default ($textToRemove) and removing it"

    sed -i "/$textToRemove/d" "$destZshrcFilePath"
}

copy_zshrc_if_missing() {
    local srcZshrcFile="$PWD/.zshrc-base"
    local destZshrcFilePath="$HOME/.zshrc"

    if [ ! -f "$destZshrcFilePath" ]; then
        cp -f "$srcZshrcFile" "$destZshrcFilePath"
        chown "$USER:$USER" "$destZshrcFilePath"
        echo "✔ 🆗 .zshrc file was missing and has been created from .zshrc-base."
    else
        echo "✔ .zshrc file already exists. No action taken."
    fi
}

# Append ZSHRC
append_zshrc() {
    local srcZshrcFile="$PWD/.zshrc"
    local destZshrcFilePath="$HOME/.zshrc"
    local isLineByLineDebugLog=true
    echo "Appending from $srcZshrcFile to $destZshrcFilePath - started!"
    echo ""

    # Read the source file line by line
    while IFS= read -r line; do
        # Escape special characters for grep
        # escaped_line=$(printf '%s\n' "$line" | sed 's/[]\/$*.^|[]/\\&/g')
        # Escape '>' character for grep
        if [[ "$line" == *"\${fpath"* ]]; then
            escaped_line=$(printf '%s\n' "$line" | sed 's/>/\\>/g')
        else
            escaped_line="$line"
        fi

        # Check if the line is already in the destination file
        if ! grep -Fxq "$escaped_line" "$destZshrcFilePath"; then
            echo "$line" >> "$destZshrcFilePath"
            if [ $isLineByLineDebugLog = true ]; then
                echo " + $line - ✅"
            fi
        fi

    done < "$srcZshrcFile"

    echo ""
    echo "Append operation - completed."
}

# Append ZSHRC
append_zshrc_or_copy_on_missing() {
    local destZshrcFilePath="$HOME/.zshrc"
    echo ""
    echo "Appending zshrc or copying the missing zshrc - Function Invoked!"

    if [ ! -f "$destZshrcFilePath" ]; then
        echo "Copying the missing zshrc - Function Invoked..."
        echo ""
        copy_zshrc_if_missing
    else

        echo "Appending the zshrc - Function Invoked..."
        echo ""
        append_zshrc
    fi

    echo ""
}

# Source ZSHRC
source_zshrc() {
    local destZshrcFilePath="$HOME/.zshrc"
    sudo zsh -c "source $destZshrcFilePath"
    zsh
    echo "ZSH sourced source $destZshrcFilePath ✅ ."
}

print_zshrc_file_presence() {
    # Check if ~/.zshrc file exists
    if [ -f "$HOME/.zshrc" ]; then
        echo "~/.zshrc is present ✅ ."
    else
        echo "~/.zshrc is missing ❌ !"
    fi
}

# Example Usage:
# main "agnoster" "true"
main() {
    local theme="$1"
    local is_append_zshrc="${2:-true}"

    install_omyzsh
    
    echo "Appending ~/.zshrc flag ❓: $is_append_zshrc"

    if [[ $is_append_zshrc = "true" || $is_append_zshrc = true || $is_append_zshrc = yes || $is_append_zshrc = y ]]; then
        echo "Applying zshrc append or copy ♻..."
        append_zshrc_or_copy_on_missing
    fi

    print_zshrc_file_presence
    change_zsh_theme "$theme"
    remove_text_from_zshrc 
  
    source_zshrc

    echo "ZSH theme changed to $theme and .zshrc content updated ✅ ."
}

# Check for help argument
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    display_help
    exit 0
fi

# Call the function with arguments
main "$1" "${2:-false}"
