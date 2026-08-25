#!/bin/bash

# Exit on error
set -e

create_root_user() {
    local username="$1"
    local password="$2"
    local theme="${3:-fletcherm}"
    local homedir="$4"
    local currentUser="$(whoami)"
    local authorizedKeyFile="$PWD/authorized_keys"

    if [ "$currentUser" = "root" ]; then
        currentUser=$defaultUserName
    fi

    local srcPathOfSSHDir="/home/$currentUser/.ssh"
    local srcPathOfSSH="$srcPathOfSSHDir/authorized_keys"

    mkdir -p "$srcPathOfSSHDir"
    echo "dir created : $srcPathOfSSHDir"

    # Check if username is provided
    if [ -z "$username" ]; then
        echo "Username is required."
        exit 1
    fi

    # Use default home directory if not provided
    if [ -z "$homedir" ]; then
        homedir="/home/$username"
    fi

    # Create the user with the specified home directory
    useradd -m -d "$homedir" -s /usr/bin/zsh "$username"

    # Set the user password if provided, otherwise prompt for it
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
        sudo cp -f "$authorizedKeyFile" "$homedir/.ssh/"
        sudo chmod 600 "$homedir/.ssh/authorized_keys"
        sudo chown -R $username:$username "$homedir"
        sudo chsh -s "$(which zsh)" "$username"
EOF

    echo "User $username created with home directory $homedir, theme $theme, and root privileges."
    echo "------> switch user      : su $username && whoami"
    echo "------> sudo login user  : su $username"
    echo "------> home directory   : $homedir && ls -la $homedir"
    echo '------> change terminal  : chsh -s "$(which zsh)"'
}

# Main script
username="$1"
password="$2"
theme="$3"
homedir="$4"

create_root_user "$username" "$password" "$theme" "$homedir"
