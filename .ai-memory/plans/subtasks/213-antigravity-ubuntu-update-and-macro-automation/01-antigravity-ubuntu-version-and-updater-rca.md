# Subtask 213.1: Antigravity Ubuntu Version & Updater RCA

- **Parent Plan:** [83-antigravity-ubuntu-update-and-macro-automation.md](../../83-antigravity-ubuntu-update-and-macro-automation.md)
- **Spec Reference:** [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md](../../../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md)
- **Status:** Ready
- **Target Node:** Ubuntu U1 (`192.168.1.22`, user `a`)
- **Target Area:** `/home/a/.local/share/antigravity-ide/`, updater diagnostics, process state

---

## 1. Objective

Investigate and document a comprehensive 4-part Root Cause Analysis (RCA) explaining why the Ubuntu node `u1` is pinned at version **2.13.0 (Build 6362815968182272)** while the Windows development host operates on **2.19.1 (Build 6046815158665216)**. Analyze Linux updater mechanics, process lifecycle locks, and Electron SUID sandbox requirements to design a foolproof upgrade pipeline.

---

## 2. 4-Part Root Cause Analysis (RCA)

### Part 1: Symptom Analysis
- **Version Discrepancy:** Remote execution of `/usr/local/bin/antigravity --version` on node `u1` returns `2.13.0 (Build 6362815968182272)`. Meanwhile, Windows 11 host runs `2.19.1 (Build 6046815158665216)`.
- **Silent Update Stagnation:** In-app background update notifications do not apply automatically on headless or unattended Linux instances.
- **Process Persistence:** Long-running `antigravity` background processes, language servers, and IPC sockets remain active, preventing in-place binary replacement without forced termination.

### Part 2: Root Cause Identification
1. **Unpackaged Tarball Distribution on Linux:** Unlike Windows (installed via NSIS installer with registry and background auto-update hooks) or macOS (managed `.app` bundles), the Linux release of Antigravity is distributed as a raw tarball (`Antigravity.tar.gz`) extracted directly into `/home/a/.local/share/antigravity-ide/`.
2. **Lack of Elevated Background Updater Daemon:** Electron's built-in `autoUpdater` module on Linux cannot elevate privileges non-interactively. Because the Chromium sandbox (`chrome-sandbox`) requires `root:root` ownership and `4755` SUID permissions, an unprivileged user process cannot update or re-sandbox the binary in the background.
3. **Hardcoded GitMap Installation Constants:** In `cli/cmdinstall/installantigravity_types.go`, GitMap's installation manifest had default constants pinned to:
   - `AntigravityDefaultVersion = "2.13.0"`
   - `AntigravityDefaultBuildID = "6362815968182272"`
   Any automated provisioning invoking `gitmap install antigravity` defaulted back to 2.13.0.

### Part 3: Blast Radius & Impact
- **Agent Capability Drift:** Newer subagent orchestration APIs, enhanced tools, and agent memory management features present in 2.19.1 are missing on 2.13.0.
- **Cross-OS Workspace Desynchronization:** Workspace state and transcript schemas generated on 2.19.1 can cause parsing errors or degraded performance when accessed on node `u1`.
- **Security & Stability:** Running an outdated Chromium/Electron base exposes the node to known engine vulnerabilities and lacks Wayland compositing stability fixes.

### Part 4: Resolution & Prevention Strategy
- **Orchestrated SSH Upgrade Pipeline (Subtask 213-02):** Execute a deterministic bash script over SSH that cleanly kills running processes, downloads the verified 2.19.1 tarball, performs atomic directory replacement, configures `chrome-sandbox` with `4755` root permissions, and validates both CLI output and graphical launch.
- **GitMap Constant Synchronization (Subtask 213-03):** Update `AntigravityDefaultVersion` to `"2.19.1"` and `AntigravityDefaultBuildID` to `"6046815158665216"` in `cli/cmdinstall/installantigravity_types.go`.
- **Node Telemetry & Version Guard (Subtask 213-05):** Include Antigravity version checks in routine GitMap node telemetry (`gitmap node list -v`) to detect version drift proactively.

---

## 3. Verification Checklist

1. [ ] Inspect active processes on node `u1`: `ssh u1 "pgrep -a antigravity || true"`.
2. [ ] Check current binary version: `ssh u1 "/home/a/.local/share/antigravity-ide/antigravity --version"`.
3. [ ] Verify ownership and permissions of current sandbox helper: `ssh u1 "ls -l /home/a/.local/share/antigravity-ide/chrome-sandbox"`.
4. [ ] Validate target artifact availability at `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz`.
