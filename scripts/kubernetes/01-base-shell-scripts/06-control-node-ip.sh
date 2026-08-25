#!/bin/bash

get_control_ip() {
    local nodes=($(kubectl get nodes -o jsonpath='{range .items[*]}{.status.addresses[?(@.type=="InternalIP")].address}{"\n"}{end}'))
    local control_plane_ip=${nodes[0]}
    
    echo "$control_plane_ip"
}

get_current_host_name() {
     echo $(hostname)
}

get_current_ip() {
     echo $(hostname -I | awk '{print $1}')
}