# scripts-fixer Comprehensive Script Catalog & Chrome Installation Audit

**Date:** 2026-09-08
**Repository:** `D:\work\scripts-fixer`
**Version:** v1.34.0
**Scope:** Windows (`scripts/`), Linux/macOS (`scripts-linux/`), and Ubuntu (`scripts/os/ubuntu/` & `scripts/run.sh`)

---

## 1. Executive Summary & Repository Architecture

`scripts-fixer` provides an automated cross-platform developer environment bootstrapping engine. The codebase is organized into three distinct execution layers:

1. **Windows Engine (`scripts/` & `run.ps1`):**
   - 51+ modular PowerShell scripts (`01-install-vscode` through `75-install-yarn`).
   - Utilizes Chocolatey (`choco`), Winget, Scoop, official silent installers, and Windows Registry injection.
   - Orchestrated via root `run.ps1` and `scripts/12-install-all-dev-tools/run.ps1`.
   - Tracks install state in `.installed/<tool>.json` and execution logs in `.logs/`.

2. **Cross-Distro Linux & macOS Engine (`scripts-linux/` & `scripts-linux/run.sh`):**
   - Unified multi-distribution bash architecture registered via `registry.yaml` (syncs to `scripts-linux/registry.json`).
   - Scripts numbered `01` through `109`.
   - Uses APT, Snap, Flatpak, Homebrew (`brew`), and direct curl/tarball/binary extraction.

3. **Ubuntu Native Engine (`scripts/os/ubuntu/` & `scripts/run.sh`):**
   - High-performance, tailored Ubuntu bash installers invoked via root `run`, `run.sh`, or `./scripts/run.sh install <keywords|IDs|profiles>`.
   - Directly executes targeted shell scripts (`dep-*.sh`, `install-*.sh`, `profile-ubuntu-*.sh`).

---

## 2. In-Depth Analysis: Ubuntu Chrome Installation Failure & Working Fix

### 2.1 Root Causes of Failure in Standard Approaches

When automating Google Chrome installation on Ubuntu / Debian systems, naive implementations consistently fail due to five fundamental issues:

1. **Absence from Default Canonical Repositories:**
   Google Chrome is proprietary freeware and is never packaged in Ubuntu's default `main`, `universe`, or `multiverse` repositories. Attempting `apt-get install -y google-chrome-stable` immediately fails with `E: Unable to locate package google-chrome-stable`.

2. **The `dpkg -i` Dependency Deadlock:**
   `dpkg` is a low-level unpacker without network repository resolution logic. Chrome depends on over 20 libraries (`libasound2`, `libgbm1`, `libnss3`, `fonts-liberation`, `libu2f-udev`, `libvulkan1`). On fresh systems, `dpkg -i` crashes with code 1, halting scripts running with `set -e`.

3. **Deprecated GPG Key Management (`apt-key add`):**
   On modern Ubuntu releases (22.04 LTS and 24.04 LTS), `apt-key` is deprecated. Manually creating `/etc/apt/sources.list.d/google-chrome.list` causes duplicate repository warnings because Chrome's `.deb` automatically installs its own cron and repository list.

4. **Minimal Base Image & Stale Package Index:**
   Headless cloud VMs, WSL images, and Docker containers frequently lack `wget` and `curl`. Running `wget` without verifying its presence leads to command not found, and attempting to install packages without `apt-get update` throws 404 Not Found against outdated mirrors.

5. **Filesystem Pollution:**
   Downloading the ~115 MB `.deb` into project directories leaves dirty files and causes permission conflicts.

---

### 2.2 The Exact Working Fix in scripts-fixer

In `scripts-fixer`, the Ubuntu Chrome installation is defined in `scripts/os/ubuntu/install-chrome.sh`:

```bash
#!/bin/bash
# Install Google Chrome on Ubuntu

echo -e "\e[1;36mℹ Installing Google Chrome\e[0m"

sudo apt-get update
sudo apt-get install -y wget curl

# Download and install Google Chrome Stable
echo -e "\e[1;33m[  ..  ] Downloading Google Chrome (.deb)...\e[0m"
wget -q -O /tmp/google-chrome-stable_current_amd64.deb https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

echo -e "\e[1;33m[  ..  ] Installing package...\e[0m"
sudo apt-get install -y /tmp/google-chrome-stable_current_amd64.deb

# Clean up
rm -f /tmp/google-chrome-stable_current_amd64.deb

echo -e "\e[1;32m✔ Google Chrome installation complete.\e[0m"
google-chrome --version
```

### 2.3 Why This Fix Succeeds in GitMap:
1. `sudo apt-get update` guarantees local cache is fresh.
2. `sudo apt-get install -y wget curl` guarantees fetch utilities are present.
3. Fetching directly to `/tmp/google-chrome-stable_current_amd64.deb` isolates temporary files.
4. `sudo apt-get install -y /tmp/...` activates APT's local-deb mode, which automatically downloads and configures all missing dependencies in a single atomic transaction.
5. Immediate deletion of `/tmp/*.deb` prevents disk leakage.
6. Verification via `google-chrome --version` confirms successful installation.
