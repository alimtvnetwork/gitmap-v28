# Subtask Plan 01: Cursor Ubuntu Fleet Setup, Wrapper Engine, and repo-secrets Tracking

- **Subtask Slug:** `01-cursor-ubuntu-fleet-setup-and-repo-secrets`
- **Parent Task:** `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search`
- **Target Files:**
  - `repo-secrets/05-scripts/setup-cursor-ubuntu.py`
  - `repo-secrets/05-scripts/setup-cursor-ubuntu.sh`
  - `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`
  - `cli/cmdcursor/cursor_install.go`

---

## 1. Overview & Objectives

On the GitMap Ubuntu fleet (specifically worker node `u1`), modern AI development requires the Cursor IDE. However, installing Cursor on headless or developer Ubuntu instances faces multiple challenges:
1. Hardcoded download URLs become obsolete quickly as new stable versions release.
2. AppImages on Ubuntu 24.04 and hardened kernels require `libfuse2` and unprivileged user namespace flags (`--no-sandbox`) to prevent Electron sandbox launch crashes.
3. Fleet tracking is absent unless status is recorded in a centralized ledger in `repo-secrets/`.

This subtask implements:
1. Dynamic upstream querying against `https://www.cursor.com/api/download?platform=linux-x64&releaseTrack=stable` to resolve direct download URLs and version metadata.
2. Idempotent provisioning scripts `repo-secrets/05-scripts/setup-cursor-ubuntu.py` and `setup-cursor-ubuntu.sh` that install the AppImage to `/opt/cursor/Cursor.AppImage` (permissions 0755).
3. Creation of universal `--no-sandbox` wrapper scripts at `/usr/local/bin/cursor` and `$HOME/.local/bin/cursor`.
4. Installation of the XDG desktop entry at `/usr/share/applications/cursor.desktop`.
5. Injection of Dracula Dark theme and whitespace formatting invariants into `$HOME/.config/Cursor/User/settings.json`.
6. Automatic persistence of installation status, version, SHA256 checksum, and health check timestamp in `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`.
7. Updating `cli/cmdcursor/cursor_install.go` to delegate `--node <alias>` commands to this setup script.

---

## 2. Step-by-Step Implementation Instructions

### Step 1: Create `repo-secrets/05-scripts/setup-cursor-ubuntu.py`

Create the idempotent Python 3 provisioning engine:
```python
#!/usr/bin/env python3
"""
setup-cursor-ubuntu.py
Automated Cursor IDE Provisioning, Sandboxing Wrapper & Fleet Ledger for Ubuntu Fleet.
"""

import argparse
import datetime
import hashlib
import json
import os
import pathlib
import shutil
import subprocess
import sys
import urllib.request

CURSOR_API_URL = "https://www.cursor.com/api/download?platform=linux-x64&releaseTrack=stable"
TARGET_DIR = pathlib.Path("/opt/cursor")
APPIMAGE_PATH = TARGET_DIR / "Cursor.AppImage"
SYSTEM_WRAPPER = pathlib.Path("/usr/local/bin/cursor")
USER_WRAPPER = pathlib.Path.home() / ".local" / "bin" / "cursor"
DESKTOP_ENTRY = pathlib.Path("/usr/share/applications/cursor.desktop")
STATUS_FILE = pathlib.Path("repo-secrets/04-ubuntu-migration/cursor-fleet-status.json")
CURSOR_SETTINGS = pathlib.Path.home() / ".config" / "Cursor" / "User" / "settings.json"

DRACULA_SETTINGS = {
    "workbench.colorTheme": "Dracula Theme",
    "editor.fontFamily": "'JetBrains Mono', 'Fira Code', Consolas, monospace",
    "editor.fontSize": 14,
    "editor.lineHeight": 22,
    "editor.tabSize": 4,
    "editor.insertSpaces": True,
    "files.autoSave": "afterDelay",
    "files.autoSaveDelay": 1000,
    "files.eol": "\n",
    "files.insertFinalNewline": True,
    "files.trimTrailingWhitespace": True,
    "editor.renderWhitespace": "selection",
    "telemetry.telemetryLevel": "off"
}

def resolve_download_url():
    req = urllib.request.Request(CURSOR_API_URL, headers={"User-Agent": "GitMap-Fleet/2.0"})
    with urllib.request.urlopen(req, timeout=30) as resp:
        return resp.geturl()

def install_dependencies():
    packages = ["libfuse2", "libnss3", "libasound2", "libgbm1", "libxss1", "curl", "ca-certificates"]
    if shutil.which("apt-get"):
        subprocess.run(["sudo", "apt-get", "update", "-y"], check=False)
        subprocess.run(["sudo", "apt-get", "install", "-y"] + packages, check=False)

def download_appimage(download_url):
    TARGET_DIR.mkdir(parents=True, exist_ok=True)
    tmp_path = APPIMAGE_PATH.with_suffix(".tmp")
    urllib.request.urlretrieve(download_url, tmp_path)
    tmp_path.chmod(0o755)
    tmp_path.rename(APPIMAGE_PATH)

def create_wrappers():
    wrapper_content = "#!/bin/sh\nexport ELECTRON_ENABLE_LOGGING=0\nexec /opt/cursor/Cursor.AppImage --no-sandbox \"$@\"\n"
    USER_WRAPPER.parent.mkdir(parents=True, exist_ok=True)
    USER_WRAPPER.write_text(wrapper_content)
    USER_WRAPPER.chmod(0o755)
    try:
        SYSTEM_WRAPPER.write_text(wrapper_content)
        SYSTEM_WRAPPER.chmod(0o755)
    except PermissionError:
        pass

def inject_dracula_theme():
    CURSOR_SETTINGS.parent.mkdir(parents=True, exist_ok=True)
    existing = {}
    if CURSOR_SETTINGS.exists():
        try:
            existing = json.loads(CURSOR_SETTINGS.read_text(encoding="utf-8"))
        except Exception:
            existing = {}
    existing.update(DRACULA_SETTINGS)
    CURSOR_SETTINGS.write_text(json.dumps(existing, indent=2), encoding="utf-8")

def update_fleet_ledger(node_alias, version, sha256_val, status):
    STATUS_FILE.parent.mkdir(parents=True, exist_ok=True)
    ledger = {"lastUpdated": datetime.datetime.now(datetime.timezone.utc).isoformat(), "nodes": {}}
    if STATUS_FILE.exists():
        try:
            ledger = json.loads(STATUS_FILE.read_text(encoding="utf-8"))
        except Exception:
            pass
    ledger.setdefault("nodes", {})[node_alias] = {
        "nodeAlias": node_alias,
        "cursorInstalled": True,
        "cursorVersion": version,
        "appImagePath": str(APPIMAGE_PATH),
        "wrapperPath": str(USER_WRAPPER),
        "sha256": sha256_val,
        "status": status,
        "lastVerified": datetime.datetime.now(datetime.timezone.utc).isoformat()
    }
    STATUS_FILE.write_text(json.dumps(ledger, indent=2), encoding="utf-8")
```

### Step 2: Create `repo-secrets/05-scripts/setup-cursor-ubuntu.sh`

Create a POSIX shell bootstrap script that calls the Python engine or provides shell fallback:
```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_NODE="${1:-u1}"

echo "● Launching Cursor Ubuntu Fleet Setup for node: ${TARGET_NODE}"
if command -v python3 >/dev/null 2>&1; then
    python3 "${SCRIPT_DIR}/setup-cursor-ubuntu.py" --node "${TARGET_NODE}"
else
    echo "⚠ python3 not found, attempting apt bootstrap..."
    sudo apt-get update -y && sudo apt-get install -y python3 curl
    python3 "${SCRIPT_DIR}/setup-cursor-ubuntu.py" --node "${TARGET_NODE}"
fi
```

### Step 3: Initialize `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`

Create the baseline schema ledger in `repo-secrets/04-ubuntu-migration/`:
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "CursorFleetStatus",
  "version": "1.0.0",
  "lastUpdated": "2026-10-05T00:00:00Z",
  "nodes": {}
}
```

### Step 4: Refactor `cli/cmdcursor/cursor_install.go`

Update `installOnLinux` and `delegateRemoteInstall`:
1. Point local Linux setup to execute `repo-secrets/05-scripts/setup-cursor-ubuntu.py`.
2. Point remote SSH execution in `delegateRemoteInstall` to invoke the setup script on the remote host via `cmdssh.RunSSHExec`:
```go
func delegateRemoteInstall(opts cursorInstallOptions) error {
    fmt.Printf("\n%s● Delegating Cursor installation to remote node:%s %s\n", constants.ColorCyan, constants.ColorReset, opts.targetNode)
    if opts.isDryRun {
        fmt.Printf("  [DryRun] Would execute remote provisioning on node '%s'\n", opts.targetNode)
        return nil
    }
    fmt.Printf("  Executing automated fleet setup script on node '%s'...\n", opts.targetNode)
    remoteCmd := "python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py --node " + opts.targetNode + " || bash repo-secrets/05-scripts/setup-cursor-ubuntu.sh " + opts.targetNode
    if opts.isForce {
        remoteCmd = "FORCE=true " + remoteCmd
    }
    if err := cmdssh.RunSSHExec([]string{opts.targetNode, remoteCmd}); err != nil {
        fmt.Printf("%s⚠ Note: Remote setup execution notice: %v%s\n", constants.ColorYellow, err, constants.ColorReset)
        return nil
    }
    fmt.Printf("%s✔ Remote setup completed on %s.%s\n", constants.ColorGreen, opts.targetNode, constants.ColorReset)
    return nil
}
```

---

## 3. Verification & Acceptance Checklist

- [ ] `repo-secrets/05-scripts/setup-cursor-ubuntu.py` is created with executable permissions.
- [ ] `repo-secrets/05-scripts/setup-cursor-ubuntu.sh` is created with executable permissions.
- [ ] `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json` is initialized with valid JSON.
- [ ] AppImage permissions are set to `0755` at `/opt/cursor/Cursor.AppImage`.
- [ ] `/usr/local/bin/cursor` and `$HOME/.local/bin/cursor` include `exec /opt/cursor/Cursor.AppImage --no-sandbox "$@"`.
- [ ] `$HOME/.config/Cursor/User/settings.json` contains `"workbench.colorTheme": "Dracula Theme"`.
- [ ] `cli/cmdcursor/cursor_install.go` delegates remote installation to `setup-cursor-ubuntu.py`.
- [ ] Code strictly follows relative path discipline and `< 100` lines file sizing standards.
