#!/bin/bash

prompt_for_theme() {
    local themeName
    read -p "Please enter the theme name (default is 'robbyrussell'): " themeName
    echo "${themeName:-robbyrussell}"
}

print_zshrc_file_presence() {
    # Check if ~/.zshrc file exists
    if [ -f "$HOME/.zshrc" ]; then
        echo "zshrc file is present"
    else
        echo "zshrc file is not present - vim ~/.zshrc"
    fi
}

clean_up_and_install(){
  local theme="$1"
  local zshRootPath=${ZSH:-$HOME/.oh-my-zsh}
  local savedPwd="$PWD"
  # Remove Oh My Zsh and apply the theme

  echo "cleaning up and removing oh-my-zsh"
  echo ""

  echo "Final theme name: $theme"

  rm -rf $ZSH && \
  rm -rf ~/.oh-my-zsh && \
  rm -rf ~/.zshrc && \
  sudo chmod +x "$PWD/01-zsh-theme-change-v2.sh" && \
  ./01-zsh-theme-change-v2.sh "$theme" true && \
  echo "" && \
  echo "expanding path : $zshRootPath" && \
  ls -la "$zshRootPath" && \
  cd $savedPwd

  print_zshrc_file_presence
}

main() {
  local myTheme="$1"
  # Check if theme name is provided as an argument
  if [ -z "$myTheme" ]; then
      myTheme=$(prompt_for_theme)
  fi

  echo "Applying theme name: $myTheme"

  clean_up_and_install $myTheme
}

main "$1"