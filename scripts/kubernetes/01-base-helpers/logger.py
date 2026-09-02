#!/usr/bin/env python3
"""Colorful timestamped logging for Kubernetes setup scripts (cross-platform)."""

from __future__ import annotations

import datetime
import os
import socket
import sys

_CLR_RESET = "\033[0m"
_CLR_GREEN = "\033[0;32m"
_CLR_YELLOW = "\033[0;33m"
_CLR_RED = "\033[0;31m"
_CLR_CYAN = "\033[0;36m"
_CLR_GRAY = "\033[0;90m"


def log_message(message: str, level: str = "info") -> None:
    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    if level == "success":
        print(f"{_CLR_GREEN}[OK]{_CLR_RESET}  {_CLR_GRAY}{now}{_CLR_RESET} - {message}")
    elif level == "warn":
        print(f"{_CLR_YELLOW}[!!]{_CLR_RESET}  {_CLR_GRAY}{now}{_CLR_RESET} - {message}")
    elif level == "error":
        print(f"{_CLR_RED}[ERR]{_CLR_RESET} {_CLR_GRAY}{now}{_CLR_RESET} - {message}", file=sys.stderr)
    else:
        print(f"{_CLR_CYAN}[>>]{_CLR_RESET}  {_CLR_GRAY}{now}{_CLR_RESET} - {message}")


def log_msg_ip(message: str, level: str = "info") -> None:
    host = socket.gethostname()
    try:
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
            s.connect(("8.8.8.8", 80))
            ip = s.getsockname()[0]
    except Exception:
        ip = "127.0.0.1"
    log_message(f"[{host} @ {ip}] {message}", level)


def assert_root() -> None:
    if os.name != "nt" and os.geteuid() != 0:
        log_message("This script must be run as root (sudo).", "error")
        sys.exit(1)
