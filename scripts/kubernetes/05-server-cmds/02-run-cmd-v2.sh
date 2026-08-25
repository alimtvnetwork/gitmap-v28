#!/bin/bash

source "$PWD/../01-base-shell-scripts/00-import-all.sh"

show_help() {
    echo "Md Alim Ul Karim"
    echo "https://alimkarim.com | https://riseup-pro.com"
    echo
    echo "Usage:"
    echo "  $0 all|control|workers|worker-1|worker-2|worker-3 <shell_command> [yes|true|no|false]"
    echo
    echo "eg:"
    echo "  $0 all \"<shell_command>\" [yes|true]"
    echo "  $0 control \"<shell_command>\" [yes|true]"
    echo "  $0 workers \"<shell_command>\" [yes|true]"
    echo "  $0 worker-1 \"<shell_command>\" [yes|true]"
    echo "  $0 worker-2 \"<shell_command>\" [yes|true]"
    echo "  $0 worker-3 \"<shell_command>\" [yes|true]"
    echo
    echo "This script executes a shell command on specified nodes."
    echo
    echo "Options:"
    echo "  -h, --help    Show this help message and exit"
}

read_json_file() {
    local json_file="$1"
    local key="$2"
    jq -r "$key" "$json_file"
}

install_sshpass() {
    install_apt_no_msg sshpass jq
}

map_worker_nodes() {
    local json_file="$1"
    declare -A worker_nodes
    while IFS="=" read -r key value; do
        worker_nodes["$key"]="$value"
    done < <(jq -r '.nodes | to_entries | map("\(.key)=\(.value)") | .[]' "$json_file")
    echo "$(declare -p worker_nodes)"
}

print_logs() {
    local json_file="$1"
    local username="$2"
    local password="$3"
    local control_node="$4"
    shift 4
    local worker_nodes=("$@")
    log_message "JSON file: $json_file"
    log_message "Username: $username"
    log_message "Control Node: $control_node"
    
    log_message "Worker Nodes Found:"
    for worker in "${worker_nodes[@]}"; do
        local worker_name=$(jq -r --arg worker "$worker" '.nodes | to_entries | map(select(.value == $worker))[0].key' "$json_file")
        log_message "   $worker_name : $worker"
    done
}

create_temp_script_remote() {
    local username="$1"
    local password="$2"
    local node="$3"
    local shell_command="$4"
    local script_dir="/tmp/on-the-fly-cmd"
    local script_file="$script_dir/script-$(date +%Y%m%d%H%M%S).sh"
    sshpass -p "$password" ssh -o StrictHostKeyChecking=no "$username@$node" "mkdir -p $script_dir && echo '$shell_command' > $script_file && chmod +x $script_file"
    echo "$script_file"
}

execute_command() {
    local username="$1"
    local password="$2"
    local shell_command="$3"
    local node_label="$4"
    local use_script="$5"
    shift 5
    local nodes=("$@")

    for node in "${nodes[@]}"; do
        echo
        log_message "Executing on $node_label - $node - \"$shell_command\""
        echo
        if [ "$use_script" == "yes" ] || [ "$use_script" == "true" ]; then
            local script_file=$(create_temp_script_remote "$username" "$password" "$node" "$shell_command")
            log_message "Using script file $node_label - $node - \"$shell_command\""
            sshpass -p "$password" ssh -o StrictHostKeyChecking=no "$username@$node" "echo $password | sudo -S bash $script_file 2>/dev/null"
            # sshpass -p "$password" ssh -o StrictHostKeyChecking=no "$username@$node" "echo $password | sudo -S bash $script_file"
        else
            sshpass -p "$password" ssh -o StrictHostKeyChecking=no "$username@$node" "echo $password | sudo -S bash -c '$shell_command' 2>/dev/null"
            # sshpass -p "$password" ssh -o StrictHostKeyChecking=no "$username@$node" "echo $password | sudo -S bash -c '$shell_command'" 2>/dev/null"
        fi
    done
}

main() {
    local json_file="$PWD/01-config.json"
    if [ ! -f "$json_file" ]; then
        log_message "JSON file not found: $json_file"
        exit 1
    fi

    local username=$(read_json_file "$json_file" '.user.name')
    local password=$(read_json_file "$json_file" '.user.password')
    local control_node=$(read_json_file "$json_file" '.control.master')
    eval "$(map_worker_nodes "$json_file")"

    # Pass associative array as a string representation
    # local worker_nodes_str="$(declare -p worker_nodes)"
    print_logs "$json_file" "$username" "$password" "$control_node" "${worker_nodes[@]}"

    if [[ $1 == "-h" || $1 == "--help" ]]; then
        show_help
        exit 0
    fi

    install_sshpass

    local target="$1"
    local shell_command="$2"
    local use_script="${3:-no}"

    log_message "Starting command execution for target: $target"

    if [ "$target" == "all" ]; then
        execute_command "$username" "$password" "$shell_command" "Control" "$use_script" "$control_node"
        execute_command "$username" "$password" "$shell_command" "Worker" "$use_script" "${worker_nodes[@]}"
    elif [ "$target" == "control" ]; then
        execute_command "$username" "$password" "$shell_command" "Control" "$use_script" "$control_node"
    elif [ "$target" == "workers" ]; then
        execute_command "$username" "$password" "$shell_command" "Worker" "$use_script" "${worker_nodes[@]}"
    elif [[ "$target" == worker-* ]]; then
        local node="${worker_nodes[$target]}"
        execute_command "$username" "$password" "$shell_command" "$target" "$use_script" "$node"
    else
        show_help
        exit 1
    fi

    echo
    log_message "Finished command execution for target: $target"
}

# Usage: ./script.sh all|control|workers|worker-1|worker-2|worker-3 "<shell_command>" [yes|true|no|false]
main "$@"
