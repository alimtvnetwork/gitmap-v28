#!/bin/bash

log_message() {
    local message="$1"
    echo "[LOG] $(date '+%Y-%m-%d %H:%M:%S') - $message"
}


log_msg_ip() {
    local message="$1"
    local hostname=$(hostname)
    local ip=$(hostname -I | awk '{print $1}')
    echo "[$hostname @ $ip] $(date '+%Y-%m-%d %H:%M:%S') - $message"
}
