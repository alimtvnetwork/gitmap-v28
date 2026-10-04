# Subtask 04: Antigravity Upgrade & Deep Brain Migration

> **Task Reference:** `214-ubuntu-fleet-automation-and-workstation-governance`  
> **Parent Plan:** [.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md](./.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `ubuntu-fleet-01`)  

---

## 1. Objective & Requirements
- Upgrade Antigravity IDE on Ubuntu from `2.13.0` to the latest `2.19.1`.
- Secure Chromium SUID root sandbox permissions (`chmod 4755 chrome-sandbox`, owned by `root:root`).
- Document Root Cause Analysis (RCA) explaining why the in-app "Check for Updates" GUI button fails on Linux.
- Migrate deep brain conversation transcripts and SQLite summaries from Windows to Ubuntu with relative path preservation (`$WORKSPACE_DIR/` -> `$HOME/git-work/`).

---

## 2. Implementation & Commands
- Implemented in `run_antigravity_update()` and `sync-antigravity-deep.ps1`.
- RCA documented in `02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/03-in-app-update-button-rca.md`.
- GitMap macro staged in `update-antigravity.json`.
- Execution command:
  ```powershell
  pwsh -File $SECRETS_DIR/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action update-antigravity
  pwsh -File $SECRETS_DIR/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action sync-brain
  ```

---

## 3. Verification & Live Evidence
- Remote Node.js ASAR parser confirms version `2.19.1`.
- `chrome-sandbox` verified as `root:root` mode `4755`.
- 98 conversation summaries and 98 brain transcript folders migrated.
- 2,321 relative path references normalized from Windows to Linux format.
- Status: **DONE**
