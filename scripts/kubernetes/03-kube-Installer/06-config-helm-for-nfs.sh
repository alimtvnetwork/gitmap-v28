#!/bin/bash

source "$PWD/../01-base-shell-scripts/00-import-all.sh"

main() {
    local ip=$(get_control_ip)
    log_msg_ip "Received ip : $ip"

    helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/ # "nfs-subdir-external-provisioner"
    kubectl delete storageclass nfs-client
    # helm install nfs-suibdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner --set nfs-server="$ip" --set nfs.path=/nfsexport # not the best practice
    # helm install nfs-suibdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner --set nfs-server="$(hostname -I | awk '{print $1}')" --set nfs.path=/nfsexport # not the best practice
    helm uninstall nfs-provisioner
    helm install nfs-provisioner \
        nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
        --set nfs.server="$(hostname -I | awk '{print $1}')" \
        --set nfs.path=/nfsexport

    kubectl get pods
}

main
