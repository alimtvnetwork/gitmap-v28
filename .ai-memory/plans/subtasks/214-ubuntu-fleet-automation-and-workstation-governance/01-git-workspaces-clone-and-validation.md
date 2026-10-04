# Subtask 01: Git Workspaces Clone & Integrity Audit

> **Task Reference:** `214-ubuntu-fleet-automation-and-workstation-governance`  
> **Parent Plan:** [.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md](file:///d:/work/gitmap/.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `192.168.1.22`)  
> **Target Directory:** `/home/a/git-work/`  

---

## 1. Objective & Requirements
- Enforce the automated loading/cloning of all 71+ repositories from the Windows work directory to Ubuntu `/home/a/git-work/` using GitMap and SSH streaming.
- Normalize legacy Windows backslash directories (`02-prompts\prompt-architect`, `Antigravity-Manager`, `movie-cli-v8`).
- Strictly avoid duplicating existing repositories by inspecting `[ -d "$target/.git" ]`.
- Report total verified repositories.

---

## 2. Implementation & Commands
- Implemented inside `run_clone_repos()` in `master-embedded-ubuntu-runner.ps1` with embedded manifest.
- Standalone execution command:
  ```powershell
  pwsh -File d:/work/repo-secrets/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action clone
  ```

---

## 3. Verification & Live Evidence
- Total active git repositories found on node `u1`: **74**.
- Zero clone failures, zero skipped directories unhandled.
- Status: **DONE**
