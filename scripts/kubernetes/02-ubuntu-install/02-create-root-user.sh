#!/bin/bash

set -e

create_root_user() {
    local username="$1"
    local password="$2"
    local theme="${3:-fletcherm}"
    local homedir="$4"
    local currentUser="$(whoami)"
    local defaultUserName="kube" # change as per the os

    if [ "$currentUser" = "root" ]; then
        currentUser=$defaultUserName
    fi

    local srcPathOfSSHDir="/home/$currentUser/.ssh"
    local srcPathOfSSH="$srcPathOfSSHDir/authorized_keys"
    local srcZshrcFile="$PWD/.zshrc"
    local removeTextFromZshrc="plugins=(git)"

    mkdir -p "$srcPathOfSSHDir"
    echo "dir created : $srcPathOfSSHDir"

    if [ -z "$username" ]; then
        echo "Username is required."
        exit 1
    fi

    if [ -z "$homedir" ]; then
        homedir="/home/$username"
    fi

    local destZshrcFilePath="$homedir/.zshrc"
    echo "zshrc src file : $srcZshrcFile"
    echo "zshrc dst file : $destZshrcFilePath"

    useradd -m -d "$homedir" -s /usr/bin/zsh "$username"

    if [ -n "$password" ]; then
        echo "$username:$password" | chpasswd
    else
        passwd "$username"
    fi

    # Add the user to the sudoers file with root privileges
    {
        echo "$username ALL=(ALL) NOPASSWD:ALL"
    } >> /etc/sudoers

    sudo usermod -aG sudo $username && \
    echo "$username ALL=(ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/$username && \
    sudo chmod 0440 /etc/sudoers.d/"$username" && \
    sudo echo "groups of $(groups $username)"

    # Run multiple commands as the new user
    sudo -u "$username" zsh <<EOF
        cd $homedir
        mkdir -p "$homedir/scripts"
        mkdir -p "$homedir/gitlab"
        mkdir -p "$homedir/github"
        mkdir -p "$homedir/.ssh"
        touch "$homedir/.ssh/authorized_keys"
        sudo chmod 600 "$homedir/.ssh/authorized_keys"
        sudo cp $srcPathOfSSH "$homedir/.ssh/"
        sudo chmod 600 "$homedir/.ssh/authorized_keys"
        sudo chown -R $username:$username "$homedir"
        sudo chsh -s "$(which zsh)" "$username"
        wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh
        sudo sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="'$theme'"/' "$destZshrcFilePath"
        sudo sed -i "/$removeTextFromZshrc/d" "$destZshrcFilePath"
        sudo cat "$srcZshrcFile" >> "$destZshrcFilePath"
        git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
        sudo rm -f /home/kube/.zcompdump-kube-*
        source $homedir/.zshrc
EOF

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
