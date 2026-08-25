#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 7 -- Remote Command Executor
#  Execute commands on cluster nodes via SSH.
#  Usage:  ./run-cmd.sh all|control|workers|worker-1 "<command>"
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

# -- Help ------------------------------------------------------------------
show_help() {
    echo ""
    echo "  Kubernetes Remote Command Executor"
    echo ""
    echo "  Usage:"
    echo "    $0 <target> \"<command>\""
    echo ""
    echo "  Targets:"
    echo "    all        Run on master + all workers"
    echo "    control    Run on master only"
    echo "    workers    Run on all workers"
    echo "    worker-N   Run on a specific worker (e.g., worker-1)"
    echo ""
    echo "  Examples:"
    echo "    $0 all \"hostname && hostname -I\""
    echo "    $0 workers \"kubectl get nodes\""
    echo "    $0 control \"kubeadm token create --print-join-command\""
    echo ""
    echo "  Config file: ../config.json (copy from config-sample.json)"
    echo ""
}

# -- Read JSON config -------------------------------------------------------
read_json() {
    local json_file="$1"
    local key="$2"
    jq -r "$key" "$json_file"
}

# -- Execute on nodes -------------------------------------------------------
execute_on_nodes() {
    local username="$1"
    local password="$2"
    local command="$3"
    local label="$4"
    shift 4
    local nodes=("$@")

    for node in "${nodes[@]}"; do
        echo ""
        log_message "Executing on $label [$node]: \"$command\"" "info"
        sshpass -p "$password" ssh -o StrictHostKeyChecking=no \
            "$username@$node" "echo $password | sudo -S bash -c '$command' 2>/dev/null" || \
            log_message "Failed to execute on $node" "error"
    done
}

# -- Main -------------------------------------------------------------------
main() {
    if [[ "$1" == "-h" || "$1" == "--help" || -z "$1" ]]; then
        show_help
        exit 0
    fi

    local config_file="$SCRIPT_DIR/../config.json"
    if [[ ! -f "$config_file" ]]; then
        log_message "Config file not found: $config_file" "error"
        echo "  Copy config-sample.json to config.json and edit it."
        exit 1
    fi

    install_apt_quiet sshpass jq

    local username
    local password
    local control_node
    username=$(read_json "$config_file" '.user.name')
    password=$(read_json "$config_file" '.user.password')
    control_node=$(read_json "$config_file" '.control.master')

    local worker_ips
    mapfile -t worker_ips < <(jq -r '.nodes | to_entries[] | .value' "$config_file")

    local target="$1"
    local command="$2"

    log_message "Target: $target | Command: \"$command\"" "info"

    case "$target" in
        all)
            execute_on_nodes "$username" "$password" "$command" "Master" "$control_node"
            execute_on_nodes "$username" "$password" "$command" "Worker" "${worker_ips[@]}"
            ;;
        control)
            execute_on_nodes "$username" "$password" "$command" "Master" "$control_node"
            ;;
        workers)
            execute_on_nodes "$username" "$password" "$command" "Worker" "${worker_ips[@]}"
            ;;
        worker-*)
            local node_ip
            node_ip=$(jq -r ".nodes.\"$target\"" "$config_file")
            if [[ "$node_ip" == "null" ]]; then
                log_message "Unknown node: $target" "error"
                exit 1
            fi
            execute_on_nodes "$username" "$password" "$command" "$target" "$node_ip"
            ;;
        *)
            show_help
            exit 1
            ;;
    esac

    echo ""
    log_message "Command execution complete for target: $target" "success"
}

main "$@"
