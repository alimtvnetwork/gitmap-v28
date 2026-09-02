#!/usr/bin/env python3
"""Step 2 — Ubuntu Prerequisites for Kubernetes setup (cross-platform / Linux)."""

import os
import subprocess
import sys
from pathlib import Path

# Add 01-base-helpers to sys.path
helpers_dir = Path(__file__).resolve().parent.parent / "01-base-helpers"
if str(helpers_dir) not in sys.path:
    sys.path.insert(0, str(helpers_dir))

try:
    from logger import log_message, assert_root
except ImportError:
    def log_message(msg: str, level: str = "info") -> None:
        print(f"[{level.upper()}] {msg}")
    def assert_root() -> None:
        if os.name != "nt" and os.geteuid() != 0:
            print("Must run as root", file=sys.stderr)
            sys.exit(1)

def main() -> int:
    assert_root()
    log_message("=== Ubuntu Prerequisites ===", "info")

    # Disable swap
    log_message("Disabling swap...", "info")
    subprocess.run(["swapoff", "-a"], check=False)
    fstab = Path("/etc/fstab")
    if fstab.is_file():
        try:
            with open(fstab, "r") as f:
                lines = f.readlines()
            new_lines = [("# " + l if " swap " in l and not l.startswith("#") else l) for l in lines]
            with open(fstab, "w") as f:
                f.writelines(new_lines)
            log_message("Swap disabled.", "success")
        except OSError:
            pass

    # Load kernel modules
    log_message("Loading kernel modules...", "info")
    k8s_modules = Path("/etc/modules-load.d/k8s.conf")
    try:
        k8s_modules.parent.mkdir(parents=True, exist_ok=True)
        k8s_modules.write_text("overlay\nbr_netfilter\n")
        subprocess.run(["modprobe", "overlay"], check=False)
        subprocess.run(["modprobe", "br_netfilter"], check=False)
        log_message("Kernel modules loaded.", "success")
    except OSError:
        pass

    # Sysctl
    log_message("Configuring sysctl for Kubernetes networking...", "info")
    k8s_sysctl = Path("/etc/sysctl.d/k8s.conf")
    try:
        k8s_sysctl.parent.mkdir(parents=True, exist_ok=True)
        k8s_sysctl.write_text("net.bridge.bridge-nf-call-iptables = 1\nnet.bridge.bridge-nf-call-ip6tables = 1\nnet.ipv4.ip_forward = 1\n")
        subprocess.run(["sysctl", "--system"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=False)
        log_message("Sysctl configured.", "success")
    except OSError:
        pass

    # Base packages
    log_message("Installing base packages...", "info")
    pkgs = ["curl", "apt-transport-https", "ca-certificates", "gnupg", "lsb-release", "wget", "nano", "vim", "git", "sshpass", "jq", "software-properties-common"]
    subprocess.run(["apt-get", "update", "-y"], check=False)
    subprocess.run(["apt-get", "install", "-y"] + pkgs, check=False)
    log_message("Ubuntu prerequisites complete.", "success")
    return 0

if __name__ == "__main__":
    sys.exit(main())
