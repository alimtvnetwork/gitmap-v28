#!/bin/bash

set -e
read_default_username() {
    if [ -f /tmp/default-script-user.txt ]; then
        cat /tmp/default-script-user.txt
    else
        read -p "Enter default username: " defaultUserName
        echo "$defaultUserName" > /tmp/default-script-user.txt
        echo "$defaultUserName"
    fi
}

get_default_username() {
    local username="$1"
    if [ -z "$username" ]; then
        read_default_username
    else
        echo "$username"
    fi
}

verify_username() {
    local username="$1"
    if [ -z "$username" ]; then
        echo "Username is required."
        exit 1
    fi
}

verify_password() {
    local username="$1"
    local password="$2"
    if [ -n "$password" ]; then
        echo "$username:$password" | chpasswd
    else
        passwd "$username"
    fi
}

setup_user_zsh() {
    local username="$1"
    local homedir="$2"
    local theme="$3"
    local srcZshrcFile="$4"
    local removeTextFromZshrc="$5"
    local authorizedKeyFile="$PWD/authorized_keys"
    local sshLocalDir="$PWD/.ssh"

    sudo -u "$username" zsh <<EOF
        cd $homedir
        mkdir -p "scripts" "gitlab" "github" ".ssh"
        touch "$homedir/.ssh/authorized_keys"
        sudo chmod 700 "$homedir/.ssh"
        sudo chmod 600 "$homedir/.ssh/authorized_keys"
        sudo chmod 600 "$authorizedKeyFile"
        sudo cp -f "$authorizedKeyFile" "$homedir/.ssh/"
        sudo cp -rf "$sshLocalDir" "$homedir/"
        sudo chmod 600 "$homedir/.ssh/authorized_keys"
        sudo chown -R $username:$username "$homedir"
        sudo chsh -s "$(which zsh)" "$username"
        wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh
        sudo sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="'$theme'"/' "$homedir/.zshrc"
        sudo sed -i "/$removeTextFromZshrc/d" "$homedir/.zshrc"
        sudo cat "$srcZshrcFile" >> "$homedir/.zshrc"
        git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
        source $homedir/.zshrc
EOF
}

create_root_user() {
    local username="$1"
    local password="$2"
    local theme="${3:-fletcherm}"
    local homedir="$4"
    local currentUser="$(whoami)"
    local defaultUserName

    defaultUserName=$(get_default_username)

    if [ "$currentUser" = "root" ]; then
        currentUser=$defaultUserName
    fi

    local srcPathOfSSHDir="/home/$currentUser/.ssh"
    local srcPathOfSSH="$srcPathOfSSHDir/authorized_keys"
    local srcZshrcFile="$PWD/.zshrc"
    local removeTextFromZshrc="plugins=(git)"

    sudo rm -f "/home/$defaultUserName/.zcompdump-kube-*"
    mkdir -p "$srcPathOfSSHDir"
    echo "dir created : $srcPathOfSSHDir"

    verify_username "$username"

    if [ -z "$homedir" ]; then
        homedir="/home/$username"
    fi

    local destZshrcFilePath="$homedir/.zshrc"
    echo "zshrc src file : $srcZshrcFile"
    echo "zshrc dst file : $destZshrcFilePath"

    useradd -m -d "$homedir" -s /usr/bin/zsh "$username"

    verify_password "$username" "$password"

    # Add the user to the sudoers file with root privileges
    {
        echo "$username ALL=(ALL) NOPASSWD:ALL"
    } >> /etc/sudoers

    sudo usermod -aG sudo $username && \
    echo "$username ALL=(ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/$username && \
    sudo chmod 0440 /etc/sudoers.d/"$username" && \
    sudo echo "groups of $(groups $username)"

    # Setup user-specific zsh configuration
    setup_user_zsh "$username" "$homedir" "$theme" "$srcZshrcFile" "$removeTextFromZshrc"

    echo "User $username created with home directory $homedir, theme $theme, and root privileges."
    echo "------> switch user      : su $username && whoami"
    echo "------> sudo login user  : su - $username"
    echo "------> fix ohmyzsh      : source ~/.zshrc"
    echo "------> home directory   : $homedir && ls -la $homedir"
    echo '------> change terminal  : chsh -s "$(which zsh)"'
}

# Main script
username="$1"
password="$2"
theme="$3"
homedir="$4"

create_root_user "$username" "$password" "$theme" "$homedir"
