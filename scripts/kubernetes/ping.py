#!/usr/bin/env python3
"""Get host name and current IP (cross-platform)."""

import socket

def get_current_host_name() -> str:
    return socket.gethostname()

def get_current_ip() -> str:
    try:
        # Connect to an external address to get preferred outbound IP
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
            s.connect(("8.8.8.8", 80))
            return s.getsockname()[0]
    except Exception:
        return "127.0.0.1"

if __name__ == "__main__":
    print("Getting the host name and ip\n")
    print(f"Hostname:\n{get_current_host_name()}\n")
    print(f"Current Ip:\n{get_current_ip()}\n")
