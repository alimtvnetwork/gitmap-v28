#!/bin/bash

show_help() {
    echo "Md Alim Ul Karim"
    echo "https://alimkarim.com | https://riseup-pro.com"
    echo
    echo "Usage: ./09-repo-permissions.sh [permissions] [username]"
    echo "Examples:"
    echo "  ./09-repo-permissions.sh 777             # Applies chmod 777 -R"
    echo "  ./09-repo-permissions.sh                 # Applies default 777"
    echo "  ./09-repo-permissions.sh 777 username    # Applies 777 for the given user"
    echo "  ./09-repo-permissions.sh -h / --help     # Shows this help message"
}

update_permissions() {
    local permissions="${1:-777}"
    local user="${2:-$(whoami)}"
    local repoRoot=$(git rev-parse --show-toplevel)

    cd "$repoRoot"
    sudo chmod "$permissions" -R "$repoRoot"
    sudo chown -R "$user:$user" "$repoRoot"
    sudo git config --global --add safe.directory "$repoRoot"
    sudo cd ls -la "$repoRoot"

    echo "Permissions updated to $permissions for user $user on the git repository root directory: $repoRoot"
}

# Main
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_help
else
    update_permissions "$@"
fi
