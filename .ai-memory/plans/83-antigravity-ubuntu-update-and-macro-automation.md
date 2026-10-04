# Plan 83: Antigravity Ubuntu Update & Macro Automation

> **Plan Status:** Active  
> **Traceability IDs:** Subtask 213-01 .. Subtask 213-05  
> **Spec Reference:** [02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md](../../02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md)  
> **Target Node:** Ubuntu U1 (`192.168.1.22`, user `a`)  
> **Source Host:** Windows 11 (`desktop-corei9-direct`)  
> **Execution Location:** `cli/cmdinstall/`, remote node `u1`, `.ai-memory/plans/subtasks/213-antigravity-ubuntu-update-and-macro-automation/`  

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable / Subtask Spec | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 213-01** | Version & Updater RCA | `.ai-memory/plans/subtasks/213-antigravity-ubuntu-update-and-macro-automation/01-antigravity-ubuntu-version-and-updater-rca.md` | Ready |
| **Subtask 213-02** | Remote SSH Upgrade Pipeline | `.ai-memory/plans/subtasks/213-antigravity-ubuntu-update-and-macro-automation/02-remote-ssh-update-to-2-19-pipeline.md` | Ready |
| **Subtask 213-03** | GitMap Install Types Sync | `cli/cmdinstall/installantigravity_types.go` | Queued |
| **Subtask 213-04** | Macro Automation & Verification | Remote macro execution on U1, validation runner | Queued |
| **Subtask 213-05** | Retrospective Log & Verification | Engineering log and end-to-end verification gate | Queued |

---

## 2. 5-Subtask Detailed Breakdown

### Subtask 213-01: Antigravity Ubuntu Version & Updater RCA
- Conduct a 4-part Root Cause Analysis (RCA) on why Ubuntu node `u1` was pinned at version 2.13.0 (Build 6362815968182272) while host nodes operate on 2.19.1 (Build 6046815158665216).
- Document updater mechanics on Linux: unpackaged Electron standalone tarball, lack of background update service, and SUID sandbox permission loss when unprivileged users extract tar archives.
- Establish resolution strategy and preventative validation gates.

### Subtask 213-02: Remote SSH Update to 2.19 Pipeline
- Construct and execute a remote SSH upgrade pipeline to transition node `u1` to 2.19.1.
- Pipeline steps:
  1. Process termination (`pkill -f antigravity` / `pkill -f antigravity-ide`).
  2. Download `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.19.1-6046815158665216/linux-x64/Antigravity.tar.gz` to `/tmp/Antigravity.tar.gz`.
  3. Backup old installation to `/home/a/.local/share/antigravity-ide.bak-2.13.0` and extract tarball to `/home/a/.local/share/antigravity-ide`.
  4. Configure SUID sandbox (`sudo chown root:root chrome-sandbox && sudo chmod 4755 chrome-sandbox`).
  5. Preserve global symlink `/usr/local/bin/antigravity`.
  6. Verify version output (`antigravity --version` -> `2.19.1`) and launch via `systemd-run --user`.
  7. Clean up temporary archive `/tmp/Antigravity.tar.gz`.

### Subtask 213-03: GitMap Antigravity Install Types & Version Constant Sync
- Synchronize default constants in `cli/cmdinstall/installantigravity_types.go`:
  - Update `AntigravityDefaultVersion` to `"2.19.1"`.
  - Update `AntigravityDefaultBuildID` to `"6046815158665216"`.
- Validate that URL resolution correctly yields the official download artifact.
- Ensure all Go tests and linters in `cli/cmdinstall` pass.

### Subtask 213-04: Macro Automation Remote Execution & Test Verification
- Author or execute macro automation tasks targeting remote node `u1` over SSH.
- Test automated application launch and CLI headless execution.
- Validate that macro execution gracefully handles remote environment nuances (Wayland/XWayland display variables, D-Bus session bus).

### Subtask 213-05: End-to-End Verification & Retrospective Engineering Log
- Execute end-to-end verification checklist across all deliverables.
- Confirm node `u1` runs version 2.19.1 with zero sandbox warnings.
- Record comprehensive execution log, including timestamps, exit codes, and operational findings.

---

## 3. Requirements Traceability Matrix

| Requirement / Prompt Item | Canonical Spec File | Plan / Subtask File | Target Components |
| :--- | :--- | :--- | :--- |
| Version & Updater RCA | `02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md#1-system-overview--problem-statement` | `subtasks/213-antigravity-ubuntu-update-and-macro-automation/01-antigravity-ubuntu-version-and-updater-rca.md` | `u1` installation analysis |
| Remote SSH Upgrade Pipeline | `02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md#2-remote-upgrade-architecture--pipeline` | `subtasks/213-antigravity-ubuntu-update-and-macro-automation/02-remote-ssh-update-to-2-19-pipeline.md` | `/home/a/.local/share/antigravity-ide/`, SUID sandbox |
| GitMap Type Constants Sync | `02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md#3-gitmap-integration--version-synchronization` | `83-antigravity-ubuntu-update-and-macro-automation.md` (Subtask 213-03) | `cli/cmdinstall/installantigravity_types.go` |
| Macro Automation Execution | `02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md#1-system-overview--problem-statement` | `83-antigravity-ubuntu-update-and-macro-automation.md` (Subtask 213-04) | Macro engine, remote runner |
| Verification & Retrospective | `02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md#6-acceptance-criteria` | `83-antigravity-ubuntu-update-and-macro-automation.md` (Subtask 213-05) | Telemetry, version check |

---

## 4. Acceptance Criteria & Quality Gates

1. **Version Parity:** Remote node `u1` reports `2.19.1` when invoking `/usr/local/bin/antigravity --version`.
2. **SUID Sandbox Compliance:** `/home/a/.local/share/antigravity-ide/chrome-sandbox` is owned by `root:root` with file mode `4755`.
3. **Clean Execution:** Antigravity launches without crashing or throwing SUID sandbox abort errors.
4. **GitMap Source Sync:** `cli/cmdinstall/installantigravity_types.go` has default version constants set to `2.19.1` and `6046815158665216`.
5. **Zero Artifact Leaks:** Ephemeral download archive `/tmp/Antigravity.tar.gz` is completely deleted.
