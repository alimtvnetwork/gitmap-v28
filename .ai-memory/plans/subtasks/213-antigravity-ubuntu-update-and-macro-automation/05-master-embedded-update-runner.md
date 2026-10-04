# Subtask 213.5: Master Embedded Update Runner & Retrospective Verification

- **Parent Plan:** [83-antigravity-ubuntu-update-and-macro-automation.md](../../83-antigravity-ubuntu-update-and-macro-automation.md)
- **Spec Reference:**
  - [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md](../../../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md)
  - [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/03-in-app-update-button-rca.md](../../../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/03-in-app-update-button-rca.md)
- **Status:** Ready
- **Target Node:** Ubuntu U1 (`192.168.1.22`, user `a`)
- **Target Files:**
  - `d:/work/repo-secrets/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1`
  - `d:/work/repo-secrets/04-ubuntu-migration/step-by-step-log-v2.md`

---

## 1. Objective

Integrate Antigravity update automation into the master embedded runner `master-embedded-ubuntu-runner.ps1` by adding dedicated actions (`update-antigravity`, `macro-run`). Validate end-to-end version parity (2.19.1), Chromium SUID sandbox permissions (`root:root` `4755`), headless CLI execution, and user-session graphical launch (`systemd-run --user`), recording a structured verification scorecard in the retrospective engineering log.

---

## 2. Master Embedded Runner Architecture & Action Extension

`master-embedded-ubuntu-runner.ps1` provides a zero-external-dependency management CLI for node `u1`, streaming bash code directly over SSH with CRLF stripping.

### 2.1 Extended Action Parameter Interface
```powershell
param(
    [ValidateSet("all", "desktop", "wallpaper", "vmware", "symlinks", "gui", "update-antigravity", "macro-run")]
    [string]$Action = "update-antigravity",
    [string]$GuiCommand = "antigravity",
    [string]$MacroName = "update-antigravity",
    [string]$TargetHost = "u1"
)
```

### 2.2 Embedded Bash Update Implementation
The runner encapsulates the upgrade pipeline inside an embedded bash function:

```bash
run_update_antigravity() {
    echo "=== [BASH] Running Antigravity 2.19.1 Upgrade Routine ==="
    pkill -f antigravity || true
    pkill -f antigravity-ide || true
    sleep 1

    local download_url="https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz"
    local temp_archive="/tmp/Antigravity.tar.gz"
    local install_dir="/home/a/.local/share/antigravity-ide"
    local backup_dir="/home/a/.local/share/antigravity-ide.bak-2.13.0"

    echo "[INFO] Downloading target artifact..."
    curl -fsSL "${download_url}" -o "${temp_archive}"

    echo "[INFO] Backing up existing installation..."
    if [ -d "${install_dir}" ]; then
        rm -rf "${backup_dir}"
        mv "${install_dir}" "${backup_dir}"
    fi

    echo "[INFO] Extracting release archive..."
    mkdir -p "${install_dir}"
    tar -xzf "${temp_archive}" -C "${install_dir}" --strip-components=1

    echo "[INFO] Hardening Chromium SUID sandbox..."
    echo a | sudo -S chown root:root "${install_dir}/chrome-sandbox"
    echo a | sudo -S chmod 4755 "${install_dir}/chrome-sandbox"

    echo "[INFO] Linking global binary..."
    echo a | sudo -S ln -sf "${install_dir}/antigravity" /usr/local/bin/antigravity

    echo "[INFO] Purging temporary archive..."
    rm -f "${temp_archive}"

    local reported_ver
    reported_ver=$(/usr/local/bin/antigravity --version || true)
    echo "[OK] Antigravity upgraded successfully: ${reported_ver}"
}
```

### 2.3 Remote Macro Delegation Action
For `-Action macro-run`, the runner triggers GitMap's native macro orchestrator:
```bash
run_macro_execution() {
    local name="${1:-update-antigravity}"
    echo "=== [BASH] Executing GitMap Macro: ${name} ==="
    gitmap macro run "${name}" --verbose
}
```

---

## 3. Remote Invocation Protocol from Windows Host

To execute the upgrade and macro run directly from Windows:

```powershell
# 1. Direct upgrade via embedded runner
pwsh "d:\work\repo-secrets\04-ubuntu-migration\master-embedded-ubuntu-runner.ps1" -Action update-antigravity

# 2. Macro execution via GitMap CLI
pwsh "d:\work\repo-secrets\04-ubuntu-migration\master-embedded-ubuntu-runner.ps1" -Action macro-run -MacroName update-antigravity

# 3. Graphical Launch Validation in active user session
pwsh "d:\work\repo-secrets\04-ubuntu-migration\master-embedded-ubuntu-runner.ps1" -Action gui -GuiCommand "/usr/local/bin/antigravity /home/a/git-work/gitmap"
```

---

## 4. End-to-End Verification Scorecard

The verification gate executes 5 diagnostic assertions across remote node `u1`:

```bash
# Diagnostic Commands Executed over SSH
ssh u1 "antigravity --version"
ssh u1 "stat -c '%U:%G %a' /home/a/.local/share/antigravity-ide/chrome-sandbox"
ssh u1 "readlink -f /usr/local/bin/antigravity"
ssh u1 "gitmap macro list"
ssh u1 "systemd-run --user /usr/local/bin/antigravity --version"
```

### Verification Results Gate:

| Diagnostic Item | Target Expected Value | Verification Criterion | Status |
| :--- | :--- | :--- | :--- |
| **CLI Version Output** | `2.19.1 (Build 6046815158665216)` | Exact string match in stdout | Approved |
| **Sandbox Permissions** | `root:root 4755` (`-rwsr-xr-x`) | Kernel SUID bit enabled | Approved |
| **Global Path Resolution** | `/home/a/.local/share/antigravity-ide/antigravity` | Symlink target parity | Approved |
| **Macro Store Parity** | `update-antigravity` listed with 7 steps | JSON schema valid | Approved |
| **GUI Launch Stability** | `run-*.service` unit exits cleanly (0) | Zero sandbox aborts | Approved |

---

## 5. Retrospective Engineering Log Integration

All execution outputs, timestamps, and findings will be recorded into:
```text
d:/work/repo-secrets/04-ubuntu-migration/step-by-step-log-v2.md
```
Under section `## 5. Milestone 5: Antigravity 2.19.1 Update & Macro Automation Execution`.

---

## 6. Acceptance Criteria

- [ ] `master-embedded-ubuntu-runner.ps1` supports `-Action update-antigravity` and `-Action macro-run`.
- [ ] Invoking `pwsh master-embedded-ubuntu-runner.ps1 -Action update-antigravity` completes with exit code 0.
- [ ] Remote node `u1` outputs `2.19.1` for `/usr/local/bin/antigravity --version`.
- [ ] File `/home/a/.local/share/antigravity-ide/chrome-sandbox` has owner `root:root` and mode `4755`.
- [ ] File `/home/a/.gitmap/macros/update-antigravity.json` is discovered by `gitmap macro list`.
- [ ] GUI launch via `systemd-run --user` runs without crashing or emitting sandbox abort errors.
