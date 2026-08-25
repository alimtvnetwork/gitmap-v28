#!/bin/bash

# Source the log function from the external file
source "$PWD/../01-base-shell-scripts/00-import-all.sh"

main() {
    local nodes=($(kubectl get nodes -o jsonpath='{range .items[*]}{.status.addresses[?(@.type=="InternalIP")].address}{"\n"}{end}'))
    local controlPlane=${nodes[0]}
    local currentIp=$(hostname -I | awk '{print $1}')
    log_message "control ip : $controlPlane, $currentIp"

    install_apt_no_msg nfs-server

    mkdir -p /nfsexport 
    sudo sh -c 'echo "/nfsexport *(rw,no_root_squash)" >> /etc/exports'
    sudo systemctl restart nfs-server 
    echo
    showmount -e $controlPlane
}

main