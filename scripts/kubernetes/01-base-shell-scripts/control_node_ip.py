#!/usr/bin/env python3
"""Get control plane node IP and hostname (cross-platform)."""

import json
import socket
import subprocess

def get_control_ip() -> str:
    try:
        out = subprocess.check_output(
            ["kubectl", "get", "nodes", "-o", "json"],
            stderr=subprocess.DEVNULL,
            text=True,
        )
        data = json.loads(out)
        for item in data.get("items", []):
            addresses = item.get("status", {}).get("addresses", [])
            for addr in addresses:
                if addr.get("type") == "InternalIP":
                    return addr.get("address", "")
    except Exception:
        pass
    return ""

def get_current_host_name() -> str:
    return socket.gethostname()

def get_current_ip() -> str:
    try:
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
            s.connect(("8.8.8.8", 80))
            return s.getsockname()[0]
    except Exception:
        return "127.0.0.1"

if __name__ == "__main__":
    print(f"Control Node IP: {get_control_ip()}")
    print(f"Current Host:    {get_current_host_name()}")
    print(f"Current IP:      {get_current_ip()}")
