#!/bin/bash

# Exit on error
set -e

# Function to remove a user
remove_user() {
    local username="$1"

    # Check if username is provided
    if [ -z "$username" ]; then
        echo "Username is required."
        exit 1
    fi

    # Remove the user from the sudoers file if it exists
    if [ -f /etc/sudoers.d/"$username" ]; then
        sudo rm -f /etc/sudoers.d/"$username"
        echo "Removed $username from /etc/sudoers.d/"
    fi

    # Remove the user from the no-password section in the sudoers file
    sudo sed -i "/^$username ALL=(ALL) NOPASSWD:ALL/d" /etc/sudoers

    # Remove the user and their home directory
    sudo deluser --remove-home "$username"
    sudo gpasswd -d "$username" sudo

    echo "User $username has been removed from the sudoers file and their home directory has been deleted."
}

# Loop through all provided usernames
for username in "$@"; do
    remove_user "$username"
done
