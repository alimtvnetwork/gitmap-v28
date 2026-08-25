#!/bin/bash

show_help() {
    echo "Md Alim Ul Karim"
    echo "https://alimkarim.com | https://riseup-pro.com"
    echo
    echo "Usage: ./10-git-pull.sh [permissions] [username]"
    echo "Examples:"
    echo "  ./10-git-pull.sh 777             # Applies chmod 777 -R"
    echo "  ./10-git-pull.sh                 # Applies default 777"
    echo "  ./10-git-pull.sh 777 username    # Applies 777 for the given user"
    echo "  ./10-git-pull.sh -h / --help     # Shows this help message"
}

update_permissions() {
    local permissions="${1:-777}"
    local user="${2:-$(whoami)}"
    local curDir="$PWD"
    local repoRoot

    repoRoot=$(git rev-parse --show-toplevel)
    echo "repo root: $repoRoot"
    echo ""

    cd "$repoRoot" || exit
    sudo chmod -R "$permissions" "$repoRoot"
    sudo chown -R "$user:$user" "$repoRoot"
    sudo git config --global --add safe.directory "$repoRoot"
    sudo git reset --hard
    sudo git pull --ff-only
    sudo chown -R "$user:$user" "$repoRoot"
    echo ""
    
    ls -la

    sudo chown -R "$user:$user" "$curDir"
    sudo chmod -R "$permissions" "$curDir"

    echo ""
    echo "Permissions updated to $permissions for user $user on the git repository root directory: $repoRoot"
}

# Main
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_help
else
    update_permissions "$@"
fi
